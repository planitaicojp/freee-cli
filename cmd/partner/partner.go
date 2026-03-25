package partner

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
)

var Cmd = &cobra.Command{
	Use:   "partner",
	Short: "Manage partners (取引先)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	createCmd.Flags().String("name", "", "partner name (required)")
	createCmd.Flags().String("code", "", "partner code")
	createCmd.Flags().String("long-name", "", "long name (正式名称)")
	createCmd.Flags().String("shortcut1", "", "shortcut 1 (カナ)")
	createCmd.Flags().String("shortcut2", "", "shortcut 2")
	createCmd.Flags().String("country-code", "JP", "country code")

	updateCmd.Flags().String("name", "", "partner name")
	updateCmd.Flags().String("code", "", "partner code")
	updateCmd.Flags().String("long-name", "", "long name (正式名称)")
	updateCmd.Flags().String("shortcut1", "", "shortcut 1 (カナ)")
	updateCmd.Flags().String("shortcut2", "", "shortcut 2")
	updateCmd.Flags().String("country-code", "", "country code")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List partners",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListPartners(client.CompanyID, "", &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.PartnersResponse
		if err := freeeAPI.ListPartners(client.CompanyID, "", &resp); err != nil {
			return err
		}
		rows := make([]model.PartnerRow, len(resp.Partners))
		for i, p := range resp.Partners {
			rows[i] = model.PartnerRow{ID: p.ID, Name: p.Name, Code: p.Code}
		}
		return output.New("table").Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show partner details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid partner ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetPartner(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.PartnerResponse
		if err := freeeAPI.GetPartner(client.CompanyID, id, &resp); err != nil {
			return err
		}
		p := resp.Partner
		fmt.Printf("ID:        %d\n", p.ID)
		fmt.Printf("Name:      %s\n", p.Name)
		if p.Code != "" {
			fmt.Printf("Code:      %s\n", p.Code)
		}
		if p.LongName != "" {
			fmt.Printf("Long Name: %s\n", p.LongName)
		}
		if p.Shortcut1 != "" {
			fmt.Printf("Shortcut1: %s\n", p.Shortcut1)
		}
		if p.Shortcut2 != "" {
			fmt.Printf("Shortcut2: %s\n", p.Shortcut2)
		}
		fmt.Printf("Country:   %s\n", p.CountryCode)
		fmt.Printf("Available: %v\n", p.Available)
		fmt.Printf("Updated:   %s\n", p.UpdateDate)
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a partner",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		countryCode, _ := cmd.Flags().GetString("country-code")
		body := map[string]any{
			"company_id":   client.CompanyID,
			"name":         name,
			"country_code": countryCode,
		}
		if v, _ := cmd.Flags().GetString("code"); v != "" {
			body["code"] = v
		}
		if v, _ := cmd.Flags().GetString("long-name"); v != "" {
			body["long_name"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/partners")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreatePartner(body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a partner",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid partner ID: %s", args[0])
		}

		body := map[string]any{
			"company_id": client.CompanyID,
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("code"); v != "" {
			body["code"] = v
		}
		if v, _ := cmd.Flags().GetString("long-name"); v != "" {
			body["long_name"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}
		if v, _ := cmd.Flags().GetString("country-code"); v != "" {
			body["country_code"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/partners/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.UpdatePartner(id, body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a partner",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid partner ID: %s", args[0])
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/partners/%d\n", id)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeletePartner(client.CompanyID, id)
	},
}
