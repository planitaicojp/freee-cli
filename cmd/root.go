package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/account"
	"github.com/planitaicojp/freee-cli/cmd/auth"
	"github.com/planitaicojp/freee-cli/cmd/company"
	cmdconfig "github.com/planitaicojp/freee-cli/cmd/config"
	"github.com/planitaicojp/freee-cli/cmd/deal"
	"github.com/planitaicojp/freee-cli/cmd/expense"
	"github.com/planitaicojp/freee-cli/cmd/invoice"
	"github.com/planitaicojp/freee-cli/cmd/item"
	"github.com/planitaicojp/freee-cli/cmd/journal"
	"github.com/planitaicojp/freee-cli/cmd/manualjournal"
	"github.com/planitaicojp/freee-cli/cmd/partner"
	"github.com/planitaicojp/freee-cli/cmd/schema"
	"github.com/planitaicojp/freee-cli/cmd/section"
	"github.com/planitaicojp/freee-cli/cmd/skill"
	"github.com/planitaicojp/freee-cli/cmd/tag"
	"github.com/planitaicojp/freee-cli/cmd/transfer"
	"github.com/planitaicojp/freee-cli/cmd/walletable"
	"github.com/planitaicojp/freee-cli/cmd/wallettxn"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/config"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

var (
	version = "dev"

	flagProfile    string
	flagFormat     string
	flagCompanyID  int64
	flagNoInput    bool
	flagQuiet      bool
	flagVerbose    bool
	flagNoColor    bool
	flagDryRun     bool
	flagNoHeader   bool
	flagFiscalYear int
)

// rootCmd is the base command.
var rootCmd = &cobra.Command{
	Use:           "freee",
	Short:         "freee API CLI",
	Long:          "Command-line interface for freee public API (Accounting, HR, Invoice, and more)",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api.UserAgent = "planitaicojp/freee-cli/" + version
		if flagVerbose {
			api.SetDebugLevel(api.DebugVerbose)
		}
		if flagNoInput {
			_ = os.Setenv(config.EnvNoInput, "1")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagProfile, "profile", "", "config profile to use")
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "", "output format: table, json, yaml, csv")
	rootCmd.PersistentFlags().Int64Var(&flagCompanyID, "company-id", 0, "freee company ID override")
	rootCmd.PersistentFlags().BoolVar(&flagNoInput, "no-input", false, "disable interactive prompts")
	rootCmd.PersistentFlags().BoolVar(&flagQuiet, "quiet", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "preview request without executing (mutating commands only)")
	rootCmd.PersistentFlags().BoolVar(&flagNoHeader, "no-header", false, "suppress table/CSV headers")
	rootCmd.PersistentFlags().IntVar(&flagFiscalYear, "fiscal-year", 0, "fiscal year (auto-sets --from/--to based on company closing month)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(cmdconfig.Cmd)
	rootCmd.AddCommand(company.Cmd)
	rootCmd.AddCommand(deal.Cmd)
	rootCmd.AddCommand(invoice.Cmd)
	rootCmd.AddCommand(partner.Cmd)
	rootCmd.AddCommand(account.Cmd)
	rootCmd.AddCommand(section.Cmd)
	rootCmd.AddCommand(tag.Cmd)
	rootCmd.AddCommand(item.Cmd)
	rootCmd.AddCommand(journal.Cmd)
	rootCmd.AddCommand(manualjournal.Cmd)
	rootCmd.AddCommand(expense.Cmd)
	rootCmd.AddCommand(transfer.Cmd)
	rootCmd.AddCommand(walletable.Cmd)
	rootCmd.AddCommand(wallettxn.Cmd)
	rootCmd.AddCommand(schema.NewCmd(rootCmd))
	rootCmd.AddCommand(skill.Cmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}

// GetProfile returns the active profile name.
func GetProfile() string {
	if flagProfile != "" {
		return flagProfile
	}
	if p := config.EnvOr(config.EnvProfile, ""); p != "" {
		return p
	}
	cfg, err := config.Load()
	if err != nil {
		return "default"
	}
	if cfg.ActiveProfile != "" {
		return cfg.ActiveProfile
	}
	return "default"
}

// GetFormat returns the output format.
func GetFormat() string {
	if flagFormat != "" {
		return flagFormat
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

// GetCompanyID returns the company ID override.
func GetCompanyID() int64 {
	return flagCompanyID
}

// IsQuiet returns whether quiet mode is enabled.
func IsQuiet() bool {
	return flagQuiet
}

// IsVerbose returns whether verbose mode is enabled.
func IsVerbose() bool {
	return flagVerbose
}

// IsDryRun returns whether dry-run mode is enabled.
func IsDryRun() bool {
	return flagDryRun
}
