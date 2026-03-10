package account

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
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
		var resp any
		if err := freeeAPI.ListAccountItems(client.CompanyID, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
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
		var id int64
		fmt.Sscanf(args[0], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.GetAccountItem(client.CompanyID, id, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}
