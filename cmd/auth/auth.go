package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/config"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/prompt"
)

// Cmd is the auth command group.
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

func init() {
	Cmd.AddCommand(loginCmd)
	Cmd.AddCommand(logoutCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(switchCmd)
	Cmd.AddCommand(tokenCmd)
	Cmd.AddCommand(removeCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via OAuth2 browser flow",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		profileName := getProfileFlag(cmd)

		// Prompt for client credentials
		clientID, _ := cmd.Flags().GetString("client-id")
		if clientID == "" {
			// Check existing credentials
			creds, _ := config.LoadCredentials()
			if c, ok := creds.Get(profileName); ok && c.ClientID != "" {
				clientID = c.ClientID
			} else {
				clientID, err = prompt.String("Client ID (from https://app.secure.freee.co.jp/developers)")
				if err != nil {
					return err
				}
			}
		}

		clientSecret, _ := cmd.Flags().GetString("client-secret")
		if clientSecret == "" {
			creds, _ := config.LoadCredentials()
			if c, ok := creds.Get(profileName); ok && c.ClientSecret != "" {
				clientSecret = c.ClientSecret
			} else {
				clientSecret, err = prompt.Password("Client Secret")
				if err != nil {
					return err
				}
			}
		}

		// Perform OAuth2 login
		result, err := api.Login(clientID, clientSecret)
		if err != nil {
			return err
		}

		// Fetch user info to get email and companies
		client := api.NewClient(result.AccessToken, 0)
		freeeAPI := &api.FreeeAPI{Client: client}

		var userResp struct {
			User struct {
				Email     string `json:"email"`
				Companies []struct {
					ID          int64  `json:"id"`
					DisplayName string `json:"display_name"`
					Role        string `json:"role"`
				} `json:"companies"`
			} `json:"user"`
		}
		if err := freeeAPI.Client.Get(client.BaseURL()+"/api/1/users/me?companies=true", &userResp); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch user info: %v\n", err)
		}

		// Save credentials
		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}
		creds.Set(profileName, config.OAuthCredentials{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    result.ExpiresAt,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		})
		if err := creds.Save(); err != nil {
			return fmt.Errorf("saving credentials: %w", err)
		}

		// Save profile
		if cfg.Profiles == nil {
			cfg.Profiles = map[string]config.Profile{}
		}
		profile := config.Profile{
			Email: userResp.User.Email,
		}
		// Set first company as default if available
		if len(userResp.User.Companies) > 0 {
			profile.CompanyID = userResp.User.Companies[0].ID
			profile.CompanyName = userResp.User.Companies[0].DisplayName
		}
		cfg.Profiles[profileName] = profile
		if cfg.ActiveProfile == "" {
			cfg.ActiveProfile = profileName
		}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		jst := time.FixedZone("JST", 9*60*60)
		fmt.Fprintf(os.Stderr, "Logged in as %s\n", userResp.User.Email)
		if profile.CompanyName != "" {
			fmt.Fprintf(os.Stderr, "Default company: %s (ID: %d)\n", profile.CompanyName, profile.CompanyID)
		}
		fmt.Fprintf(os.Stderr, "Token expires: %s JST\n", result.ExpiresAt.In(jst).Format("2006-01-02 15:04"))
		return nil
	},
}

func init() {
	loginCmd.Flags().String("client-id", "", "OAuth2 client ID")
	loginCmd.Flags().String("client-secret", "", "OAuth2 client secret")
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove token and credentials for the active profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := getProfileFlag(cmd)

		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}
		creds.Delete(profileName)
		if err := creds.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Logged out of profile %q\n", profileName)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		profileName := getProfileFlag(cmd)
		profile, ok := cfg.Profiles[profileName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Profile %q: not configured\n", profileName)
			return &cerrors.ConfigError{Message: fmt.Sprintf("profile %q not found", profileName)}
		}

		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}

		fmt.Printf("Profile:   %s\n", profileName)
		fmt.Printf("Email:     %s\n", profile.Email)
		fmt.Printf("Company:   %s (ID: %d)\n", profile.CompanyName, profile.CompanyID)

		if cred, ok := creds.Get(profileName); ok {
			jst := time.FixedZone("JST", 9*60*60)
			remaining := time.Until(cred.ExpiresAt)
			if remaining > 0 {
				fmt.Printf("Token:     valid (expires in %s, %s JST)\n",
					remaining.Truncate(time.Minute),
					cred.ExpiresAt.In(jst).Format("2006-01-02 15:04"))
			} else {
				fmt.Printf("Token:     expired (%s ago)\n", (-remaining).Truncate(time.Minute))
				if cred.RefreshToken != "" {
					fmt.Printf("Refresh:   available (will auto-refresh on next API call)\n")
				}
			}
		} else {
			fmt.Printf("Token:     none\n")
		}

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}

		if len(cfg.Profiles) == 0 {
			fmt.Fprintln(os.Stderr, "No profiles configured. Run 'freee auth login' to create one.")
			return nil
		}

		for name, profile := range cfg.Profiles {
			marker := " "
			if name == cfg.ActiveProfile {
				marker = "*"
			}
			tokenStatus := "no token"
			if creds.IsTokenValid(name) {
				tokenStatus = "authenticated"
			} else if c, ok := creds.Get(name); ok {
				if c.RefreshToken != "" {
					tokenStatus = "expired (refreshable)"
				} else {
					tokenStatus = "expired"
				}
			}
			fmt.Printf("%s %s\t%s\t%s\t%s\n", marker, name, profile.Email, profile.CompanyName, tokenStatus)
		}
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <profile>",
	Short: "Switch active profile",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if _, ok := cfg.Profiles[name]; !ok {
			return &cerrors.ConfigError{Message: fmt.Sprintf("profile %q not found", name)}
		}

		cfg.ActiveProfile = name
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Switched to profile %q\n", name)
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print current access token to stdout (for scripting)",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := getProfileFlag(cmd)

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}

		cred, ok := creds.Get(profileName)
		if !ok {
			return &cerrors.AuthError{Message: fmt.Sprintf("no credentials for profile %q", profileName)}
		}

		token, err := api.EnsureToken(profileName, cred, cfg)
		if err != nil {
			return err
		}

		fmt.Print(token)
		return nil
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <profile>",
	Short: "Completely remove a profile",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		delete(cfg.Profiles, name)
		if cfg.ActiveProfile == name {
			cfg.ActiveProfile = ""
			for k := range cfg.Profiles {
				cfg.ActiveProfile = k
				break
			}
		}
		if err := cfg.Save(); err != nil {
			return err
		}

		creds, err := config.LoadCredentials()
		if err != nil {
			return err
		}
		creds.Delete(name)
		_ = creds.Save()

		fmt.Fprintf(os.Stderr, "Removed profile %q\n", name)
		return nil
	},
}

func getProfileFlag(cmd *cobra.Command) string {
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
