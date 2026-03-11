package expense

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
	listCmd.Flags().Int("limit", 50, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")
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

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListExpenseApplications(client.CompanyID, "", &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.ExpenseApplicationsResponse
		if err := freeeAPI.ListExpenseApplications(client.CompanyID, "", &resp); err != nil {
			return err
		}
		rows := make([]model.ExpenseApplicationRow, len(resp.ExpenseApplications))
		for i, e := range resp.ExpenseApplications {
			rows[i] = model.ExpenseApplicationRow{
				ID:     e.ID,
				Title:  e.Title,
				Amount: e.TotalAmount,
				Status: e.Status,
				Date:   e.IssueDate,
			}
		}
		return output.New("table").Format(os.Stdout, rows)
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

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetExpenseApplication(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.ExpenseApplicationResponse
		if err := freeeAPI.GetExpenseApplication(client.CompanyID, id, &resp); err != nil {
			return err
		}
		e := resp.ExpenseApplication
		fmt.Printf("ID:        %d\n", e.ID)
		if e.Title != "" {
			fmt.Printf("Title:     %s\n", e.Title)
		}
		fmt.Printf("Amount:    %d\n", e.TotalAmount)
		fmt.Printf("Status:    %s\n", e.Status)
		fmt.Printf("Date:      %s\n", e.IssueDate)
		if e.Description != "" {
			fmt.Printf("Note:      %s\n", e.Description)
		}
		return nil
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
		var id int64
		fmt.Sscanf(args[0], "%d", &id)

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/expense_applications/%d\n", id)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteExpenseApplication(client.CompanyID, id)
	},
}
