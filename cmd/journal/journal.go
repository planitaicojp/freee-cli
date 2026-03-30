package journal

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the journal command group.
var Cmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage journals (仕訳帳)",
}

func init() {
	Cmd.AddCommand(listCmd)

	listCmd.Flags().String("download-type", "csv", "download type: csv or pdf")
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List/download journals",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.ListJournals(client.CompanyID, "", &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}
