package company

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/config"
	"github.com/planitaicojp/freee-cli/internal/model"
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

		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, rows)
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
			parsed, parseErr := strconv.ParseInt(args[0], 10, 64)
			if parseErr != nil {
				return fmt.Errorf("invalid company ID: %s", args[0])
			}
			companyID = parsed
		}

		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			// JSON/YAML/CSV: output raw API response
			var resp any
			if err := freeeAPI.GetCompany(companyID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		// Table (default): key-value output
		var resp model.CompanyResponse
		if err := freeeAPI.GetCompany(companyID, &resp); err != nil {
			return err
		}

		c := resp.Company
		fmt.Printf("ID:          %d\n", c.ID)
		fmt.Printf("Name:        %s\n", c.Name)
		fmt.Printf("Display:     %s\n", c.DisplayName)
		fmt.Printf("Name Kana:   %s\n", c.NameKana)
		fmt.Printf("Role:        %s\n", c.Role)
		if c.Phone1 != "" {
			fmt.Printf("Phone 1:     %s\n", c.Phone1)
		}
		if c.Phone2 != "" {
			fmt.Printf("Phone 2:     %s\n", c.Phone2)
		}
		if c.Zipcode != "" {
			fmt.Printf("Zipcode:     %s\n", c.Zipcode)
		}
		addr := c.StreetName1
		if c.StreetName2 != "" {
			addr += " " + c.StreetName2
		}
		if addr != "" {
			fmt.Printf("Address:     %s\n", addr)
		}
		if c.CompanyNumber != "" {
			fmt.Printf("Company No:  %s\n", c.CompanyNumber)
		}
		fmt.Printf("Layout:      %s\n", c.InvoiceLayout)
		fmt.Printf("Workflow:     %s\n", c.WorkflowSetting)
		if len(c.FiscalYears) > 0 {
			fy := c.FiscalYears[0]
			fmt.Printf("Fiscal Year: %s ~ %s\n", fy.StartDate, fy.EndDate)
		}
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <id>",
	Short: "Switch default company",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		companyID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid company ID: %s", args[0])
		}

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
