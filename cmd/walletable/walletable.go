package walletable

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the walletable command group.
var Cmd = &cobra.Command{
	Use:   "walletable",
	Short: "Manage walletables (口座)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List walletables",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.ListWalletables(client.CompanyID, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <type> <id>",
	Short: "Show walletable details (type: bank_account, credit_card, wallet)",
	Args:  cmdutil.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		walletableType := args[0]
		var id int64
		_, _ = fmt.Sscanf(args[1], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.GetWalletable(client.CompanyID, walletableType, id, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}
