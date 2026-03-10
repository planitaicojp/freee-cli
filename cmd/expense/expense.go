package expense

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the expense command group.
var Cmd = &cobra.Command{
	Use:   "expense",
	Short: "Manage expense applications (経費申請)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	listCmd.Flags().String("status", "", "filter by status")
	listCmd.Flags().Int("limit", 50, "max number of results")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List expense applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.ListExpenseApplications(client.CompanyID, "", &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show expense application details",
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
		if err := freeeAPI.GetExpenseApplication(client.CompanyID, id, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an expense application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an expense application",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an expense application",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		var id int64
		fmt.Sscanf(args[0], "%d", &id)
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteExpenseApplication(client.CompanyID, id)
	},
}
