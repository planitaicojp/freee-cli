package company

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/config"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the company command group.
var Cmd = &cobra.Command{
	Use:   "company",
	Short: "Manage companies (事業所)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(switchCmd)
}

type companyRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List companies",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp struct {
			Companies []struct {
				ID          int64  `json:"id"`
				DisplayName string `json:"display_name"`
				Role        string `json:"role"`
			} `json:"companies"`
		}
		if err := freeeAPI.Client.Get(client.BaseURL()+"/api/1/companies", &resp); err != nil {
			return err
		}

		rows := make([]companyRow, len(resp.Companies))
		for i, c := range resp.Companies {
			rows[i] = companyRow{ID: c.ID, Name: c.DisplayName, Role: c.Role}
		}

		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show company details",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		companyID := client.CompanyID
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &companyID)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.GetCompany(companyID, &resp); err != nil {
			return err
		}

		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <id>",
	Short: "Switch default company",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var companyID int64
		fmt.Sscanf(args[0], "%d", &companyID)

		profileName := "default"
		if p, _ := cmd.Flags().GetString("profile"); p != "" {
			profileName = p
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		profile, ok := cfg.Profiles[profileName]
		if !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}

		profile.CompanyID = companyID
		cfg.Profiles[profileName] = profile
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Switched default company to %d\n", companyID)
		return nil
	},
}
