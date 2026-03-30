package invoice

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
	Use:   "invoice",
	Short: "Manage invoices (請求書)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	listCmd.Flags().String("status", "", "filter by status")
	listCmd.Flags().String("partner", "", "filter by partner name")
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
	listCmd.Flags().Int("limit", 50, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")

	createCmd.Flags().Int64("partner-id", 0, "partner ID (required)")
	createCmd.Flags().String("partner-name", "", "partner name (resolves to partner ID)")
	createCmd.Flags().String("date", "", "issue date YYYY-MM-DD (required)")
	_ = createCmd.MarkFlagRequired("date")
	createCmd.Flags().String("due-date", "", "due date YYYY-MM-DD")
	createCmd.Flags().String("title", "", "invoice title")
	createCmd.Flags().String("description", "", "description")

	updateCmd.Flags().Int64("partner-id", 0, "partner ID")
	updateCmd.Flags().String("partner-name", "", "partner name (resolves to partner ID)")
	updateCmd.Flags().String("date", "", "issue date YYYY-MM-DD")
	updateCmd.Flags().String("due-date", "", "due date YYYY-MM-DD")
	updateCmd.Flags().String("title", "", "invoice title")
	updateCmd.Flags().String("description", "", "description")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		fiscalFrom, fiscalTo, err := cmdutil.ResolveFiscalYear(cmd, freeeAPI, client.CompanyID)
		if err != nil {
			return err
		}
		if fiscalFrom != "" {
			_ = cmd.Flags().Set("from", fiscalFrom)
			_ = cmd.Flags().Set("to", fiscalTo)
		}

		params := buildInvoiceListParams(cmd)
		format := cmdutil.GetFormat(cmd)
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListInvoices(client.CompanyID, params, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.InvoicesResponse
		if err := freeeAPI.ListInvoices(client.CompanyID, params, &resp); err != nil {
			return err
		}
		rows := make([]model.InvoiceRow, len(resp.Invoices))
		for i, inv := range resp.Invoices {
			rows[i] = model.InvoiceRow{
				ID:        inv.ID,
				Number:    inv.InvoiceNumber,
				Partner:   inv.PartnerName,
				Amount:    inv.TotalAmount,
				Status:    inv.InvoiceStatus,
				IssueDate: inv.IssueDate,
			}
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

func buildInvoiceListParams(cmd *cobra.Command) string {
	params := ""
	add := func(key, value string) {
		if value != "" {
			if params != "" {
				params += "&"
			}
			params += key + "=" + value
		}
	}
	status, _ := cmd.Flags().GetString("status")
	add("invoice_status", status)
	partner, _ := cmd.Flags().GetString("partner")
	add("partner_code", partner)
	from, _ := cmd.Flags().GetString("from")
	add("start_issue_date", from)
	to, _ := cmd.Flags().GetString("to")
	add("end_issue_date", to)
	limit, _ := cmd.Flags().GetInt("limit")
	if limit > 0 {
		add("limit", fmt.Sprintf("%d", limit))
	}
	offset, _ := cmd.Flags().GetInt("offset")
	if offset > 0 {
		add("offset", fmt.Sprintf("%d", offset))
	}
	return params
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show invoice details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		invoiceID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid invoice ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetInvoice(client.CompanyID, invoiceID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.InvoiceResponse
		if err := freeeAPI.GetInvoice(client.CompanyID, invoiceID, &resp); err != nil {
			return err
		}
		inv := resp.Invoice
		fmt.Printf("ID:        %d\n", inv.ID)
		if inv.InvoiceNumber != "" {
			fmt.Printf("Number:    %s\n", inv.InvoiceNumber)
		}
		if inv.Title != "" {
			fmt.Printf("Title:     %s\n", inv.Title)
		}
		if inv.PartnerName != "" {
			fmt.Printf("Partner:   %s\n", inv.PartnerName)
		}
		fmt.Printf("Amount:    %d\n", inv.TotalAmount)
		fmt.Printf("SubTotal:  %d\n", inv.SubTotal)
		fmt.Printf("VAT:       %d\n", inv.TotalVat)
		fmt.Printf("Status:    %s\n", inv.InvoiceStatus)
		fmt.Printf("Issue:     %s\n", inv.IssueDate)
		if inv.DueDate != "" {
			fmt.Printf("Due:       %s\n", inv.DueDate)
		}
		if inv.Description != "" {
			fmt.Printf("Note:      %s\n", inv.Description)
		}
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invoice",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		partnerID, err := resolve.PartnerID(cmd, freeeAPI, client.CompanyID)
		if err != nil {
			return err
		}
		date, _ := cmd.Flags().GetString("date")

		body := map[string]any{
			"company_id": client.CompanyID,
			"partner_id": partnerID,
			"issue_date": date,
		}
		if v, _ := cmd.Flags().GetString("due-date"); v != "" {
			body["due_date"] = v
		}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			body["title"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/invoices")
			return output.New("json").Format(os.Stdout, body)
		}

		var resp any
		if err := freeeAPI.CreateInvoice(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an invoice",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		invoiceID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid invoice ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		body := map[string]any{
			"company_id": client.CompanyID,
		}
		partnerID, err := resolve.PartnerID(cmd, freeeAPI, client.CompanyID)
		if err != nil {
			return err
		}
		if partnerID != 0 {
			body["partner_id"] = partnerID
		}
		if v, _ := cmd.Flags().GetString("date"); v != "" {
			body["issue_date"] = v
		}
		if v, _ := cmd.Flags().GetString("due-date"); v != "" {
			body["due_date"] = v
		}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			body["title"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/invoices/%d\n", invoiceID)
			return output.New("json").Format(os.Stdout, body)
		}

		var resp any
		if err := freeeAPI.UpdateInvoice(invoiceID, body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an invoice",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		invoiceID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid invoice ID: %s", args[0])
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/invoices/%d\n", invoiceID)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteInvoice(client.CompanyID, invoiceID)
	},
}
