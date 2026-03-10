package walletable

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

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListWalletables(client.CompanyID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.WalletablesResponse
		if err := freeeAPI.ListWalletables(client.CompanyID, &resp); err != nil {
			return err
		}
		rows := make([]model.WalletableRow, len(resp.Walletables))
		for i, w := range resp.Walletables {
			rows[i] = model.WalletableRow{ID: w.ID, Name: w.Name, Type: w.Type}
		}
		return output.New("table").Format(os.Stdout, rows)
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
		fmt.Sscanf(args[1], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetWalletable(client.CompanyID, walletableType, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.WalletableResponse
		if err := freeeAPI.GetWalletable(client.CompanyID, walletableType, id, &resp); err != nil {
			return err
		}
		w := resp.Walletable
		fmt.Printf("ID:       %d\n", w.ID)
		fmt.Printf("Name:     %s\n", w.Name)
		fmt.Printf("Type:     %s\n", w.Type)
		fmt.Printf("Balance:  %d\n", w.WalletableBalance)
		fmt.Printf("Updated:  %s\n", w.UpdateDate)
		return nil
	},
}
