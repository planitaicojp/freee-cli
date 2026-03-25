package account

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

// Cmd is the account command group.
var Cmd = &cobra.Command{
	Use:   "account",
	Short: "Manage account items (勘定科目)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List account items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListAccountItems(client.CompanyID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.AccountItemsResponse
		if err := freeeAPI.ListAccountItems(client.CompanyID, &resp); err != nil {
			return err
		}

		rows := make([]model.AccountItemRow, len(resp.AccountItems))
		for i, a := range resp.AccountItems {
			rows[i] = model.AccountItemRow{
				ID:       a.ID,
				Name:     a.Name,
				Category: a.AccountCategory,
				TaxCode:  a.TaxCode,
				Shortcut: a.ShortcutNum,
			}
		}
		return output.New("table").Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show account item details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid account item ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetAccountItem(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.AccountItemResponse
		if err := freeeAPI.GetAccountItem(client.CompanyID, id, &resp); err != nil {
			return err
		}
		a := resp.AccountItem
		fmt.Printf("ID:        %d\n", a.ID)
		fmt.Printf("Name:      %s\n", a.Name)
		fmt.Printf("Category:  %s\n", a.AccountCategory)
		if a.GroupName != "" {
			fmt.Printf("Group:     %s\n", a.GroupName)
		}
		fmt.Printf("Tax Code:  %d\n", a.TaxCode)
		if a.Shortcut != "" {
			fmt.Printf("Shortcut:  %s\n", a.Shortcut)
		}
		if a.ShortcutNum != "" {
			fmt.Printf("Number:    %s\n", a.ShortcutNum)
		}
		fmt.Printf("Available: %v\n", a.Available)
		return nil
	},
}
