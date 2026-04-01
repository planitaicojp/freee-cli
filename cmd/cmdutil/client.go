package cmdutil

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/config"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

// NewClient creates an API client from the cobra command context.
func NewClient(cmd *cobra.Command) (*api.Client, error) {
	profileName := getProfile(cmd)

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return nil, &cerrors.ConfigError{Message: "profile not found, run 'freee auth login'"}
	}

	creds, err := config.LoadCredentials()
	if err != nil {
		return nil, err
	}

	cred, ok := creds.Get(profileName)
	if !ok {
		return nil, &cerrors.AuthError{Message: fmt.Sprintf("no credentials for profile %q, run 'freee auth login'", profileName)}
	}

	token, err := api.EnsureToken(profileName, cred, cfg)
	if err != nil {
		return nil, err
	}

	// Company ID: flag > env > profile default
	companyID := getCompanyID(cmd, profile)

	return api.NewClient(token, companyID), nil
}

// GetFormat returns the output format from flags.
func GetFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultFormat
	}
	return cfg.Defaults.Format
}

func getProfile(cmd *cobra.Command) string {
	if p, _ := cmd.Flags().GetString("profile"); p != "" {
		return p
	}
	if p := config.EnvOr(config.EnvProfile, ""); p != "" {
		return p
	}
	cfg, _ := config.Load()
	if cfg != nil && cfg.ActiveProfile != "" {
		return cfg.ActiveProfile
	}
	return "default"
}

// IsDryRun returns whether dry-run mode is enabled.
func IsDryRun(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("dry-run")
	return v
}

// IsAll returns whether --all flag is set.
func IsAll(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("all")
	return v
}

// IsNoHeader returns whether --no-header is set.
func IsNoHeader(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("no-header")
	return v
}

// ValidWalletableTypes lists the valid walletable type values.
var ValidWalletableTypes = map[string]bool{
	"bank_account": true,
	"credit_card":  true,
	"wallet":       true,
}

// ValidateWalletableType validates that the given value is a valid walletable type.
func ValidateWalletableType(flagName, value string) error {
	if !ValidWalletableTypes[value] {
		return &cerrors.ValidationError{
			Message: fmt.Sprintf("--%s must be one of: bank_account, credit_card, wallet (got %q)\nhint: run 'freee walletable list' to see available accounts", flagName, value),
		}
	}
	return nil
}

func getCompanyID(cmd *cobra.Command, profile config.Profile) int64 {
	if id, _ := cmd.Flags().GetInt64("company-id"); id != 0 {
		return id
	}
	if s := config.EnvOr(config.EnvCompanyID, ""); s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			return id
		}
	}
	return profile.CompanyID
}
