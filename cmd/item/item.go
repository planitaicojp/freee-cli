package item

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
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
		var resp any
		if err := freeeAPI.ListItems(client.CompanyID, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an item",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an item",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
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
		var id int64
		fmt.Sscanf(args[0], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteItem(client.CompanyID, id)
	},
}
