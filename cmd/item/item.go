package item

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

// Cmd is the item command group.
var Cmd = &cobra.Command{
	Use:   "item",
	Short: "Manage items (品目)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	createCmd.Flags().String("name", "", "item name (required)")
	createCmd.Flags().String("shortcut1", "", "shortcut 1")
	createCmd.Flags().String("shortcut2", "", "shortcut 2")

	updateCmd.Flags().String("name", "", "item name")
	updateCmd.Flags().String("shortcut1", "", "shortcut 1")
	updateCmd.Flags().String("shortcut2", "", "shortcut 2")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListItems(client.CompanyID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.ItemsResponse
		if err := freeeAPI.ListItems(client.CompanyID, &resp); err != nil {
			return err
		}

		rows := make([]model.ItemRow, len(resp.Items))
		for i, item := range resp.Items {
			rows[i] = model.ItemRow{
				ID:        item.ID,
				Name:      item.Name,
				Available: item.Available,
			}
		}
		return output.New("table").Format(os.Stdout, rows)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an item",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		body := map[string]any{
			"company_id": client.CompanyID,
			"name":       name,
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/items")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateItem(body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an item",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid item ID: %s", args[0])
		}

		body := map[string]any{
			"company_id": client.CompanyID,
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/items/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.UpdateItem(id, body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an item",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid item ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteItem(client.CompanyID, id)
	},
}
