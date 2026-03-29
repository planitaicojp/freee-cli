package expense

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
	"github.com/planitaicojp/freee-cli/internal/resolve"
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

	createCmd.Flags().String("title", "", "expense title (required)")
	createCmd.Flags().String("date", "", "issue date YYYY-MM-DD (required)")
	createCmd.Flags().String("description", "", "description")
	createCmd.Flags().Int64("account-item-id", 0, "account item ID")
	createCmd.Flags().String("account-name", "", "account item name (resolves to account item ID)")
	createCmd.Flags().String("account-item-name", "", "alias for --account-name")
	createCmd.Flags().Int64("amount", 0, "amount")

	updateCmd.Flags().String("title", "", "expense title")
	updateCmd.Flags().String("date", "", "issue date YYYY-MM-DD")
	updateCmd.Flags().String("description", "", "description")
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
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid expense application ID: %s", args[0])
		}
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
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		freeeAPI := &api.FreeeAPI{Client: client}

		accountItemID, err := resolve.AccountItemID(cmd, freeeAPI, client.CompanyID)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		date, _ := cmd.Flags().GetString("date")

		body := map[string]any{
			"company_id": client.CompanyID,
			"title":      title,
			"issue_date": date,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}

		// Build expense lines if provided
		amount, _ := cmd.Flags().GetInt64("amount")
		if accountItemID != 0 || amount != 0 {
			body["expense_application_lines"] = []map[string]any{
				{
					"account_item_id":        accountItemID,
					"amount":                 amount,
					"transaction_date":       date,
					"description":            title,
					"expense_application_id": nil,
				},
			}
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/expense_applications")
			return output.New("json").Format(os.Stdout, body)
		}
		var resp any
		if err := freeeAPI.CreateExpenseApplication(body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an expense application",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid expense application ID: %s", args[0])
		}

		body := map[string]any{
			"company_id": client.CompanyID,
		}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			body["title"] = v
		}
		if v, _ := cmd.Flags().GetString("date"); v != "" {
			body["issue_date"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/expense_applications/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.UpdateExpenseApplication(id, body, &resp); err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an expense application",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid expense application ID: %s", args[0])
		}

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
