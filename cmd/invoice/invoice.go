package invoice

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
	listCmd.Flags().Int("limit", 50, "max number of results")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
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

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListInvoices(client.CompanyID, "", &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.InvoicesResponse
		if err := freeeAPI.ListInvoices(client.CompanyID, "", &resp); err != nil {
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
		return output.New("table").Format(os.Stdout, rows)
	},
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
		var invoiceID int64
		fmt.Sscanf(args[0], "%d", &invoiceID)
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
		return fmt.Errorf("not yet implemented")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an invoice",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an invoice",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		var invoiceID int64
		fmt.Sscanf(args[0], "%d", &invoiceID)
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteInvoice(client.CompanyID, invoiceID)
	},
}
