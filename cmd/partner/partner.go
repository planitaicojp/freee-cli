package partner

import (
	"fmt"
	"os"

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
		var id int64
		fmt.Sscanf(args[0], "%d", &id)
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
		return fmt.Errorf("not yet implemented")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a partner",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a partner",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		var id int64
		fmt.Sscanf(args[0], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeletePartner(client.CompanyID, id)
	},
}
