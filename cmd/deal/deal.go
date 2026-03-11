package deal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the deal command group.
var Cmd = &cobra.Command{
	Use:   "deal",
	Short: "Manage deals (取引)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	listCmd.Flags().String("type", "", "filter by type: income or expense")
	listCmd.Flags().String("partner", "", "filter by partner name")
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
	listCmd.Flags().String("status", "", "filter by status: settled or unsettled")
	listCmd.Flags().Int("limit", 50, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")

	createCmd.Flags().String("type", "", "deal type: income or expense (required)")
	createCmd.Flags().String("date", "", "issue date YYYY-MM-DD (required)")
	createCmd.Flags().Int64("partner-id", 0, "partner ID")
	createCmd.Flags().Int64("account-item-id", 0, "account item ID (required)")
	createCmd.Flags().Int64("amount", 0, "amount (required)")
	createCmd.Flags().Int64("tax-code", 0, "tax code")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deals",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		fetchAll := cmdutil.IsAll(cmd)

		if fetchAll {
			limit, _ := cmd.Flags().GetInt("limit")
			if limit <= 0 {
				limit = 100
			}
			var allDeals []model.Deal
			for offset := 0; ; offset += limit {
				cmd.Flags().Set("offset", fmt.Sprintf("%d", offset))
				cmd.Flags().Set("limit", fmt.Sprintf("%d", limit))
				params := buildListParams(cmd)
				var resp model.DealsResponse
				if err := freeeAPI.ListDeals(client.CompanyID, params, &resp); err != nil {
					return err
				}
				allDeals = append(allDeals, resp.Deals...)
				if len(resp.Deals) < limit {
					break
				}
			}
			if format != "" && format != "table" {
				return output.New(format).Format(os.Stdout, map[string]any{"deals": allDeals})
			}
			rows := make([]model.DealRow, len(allDeals))
			for i, d := range allDeals {
				rows[i] = model.DealRow{ID: d.ID, Date: d.IssueDate, Type: d.Type, Amount: d.Amount, Status: d.Status}
			}
			return output.New("table").Format(os.Stdout, rows)
		}

		params := buildListParams(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListDeals(client.CompanyID, params, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.DealsResponse
		if err := freeeAPI.ListDeals(client.CompanyID, params, &resp); err != nil {
			return err
		}
		rows := make([]model.DealRow, len(resp.Deals))
		for i, d := range resp.Deals {
			rows[i] = model.DealRow{ID: d.ID, Date: d.IssueDate, Type: d.Type, Amount: d.Amount, Status: d.Status}
		}
		return output.New("table").Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show deal details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		var dealID int64
		fmt.Sscanf(args[0], "%d", &dealID)

		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetDeal(client.CompanyID, dealID, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.DealResponse
		if err := freeeAPI.GetDeal(client.CompanyID, dealID, &resp); err != nil {
			return err
		}
		d := resp.Deal
		fmt.Printf("ID:        %d\n", d.ID)
		fmt.Printf("Type:      %s\n", d.Type)
		fmt.Printf("Date:      %s\n", d.IssueDate)
		if d.DueDate != "" {
			fmt.Printf("Due Date:  %s\n", d.DueDate)
		}
		fmt.Printf("Amount:    %d\n", d.Amount)
		fmt.Printf("Status:    %s\n", d.Status)
		if d.PartnerID != 0 {
			fmt.Printf("Partner:   %d\n", d.PartnerID)
		}
		if d.RefNumber != "" {
			fmt.Printf("Ref:       %s\n", d.RefNumber)
		}
		if len(d.Details) > 0 {
			fmt.Println("Details:")
			for i, dt := range d.Details {
				fmt.Printf("  [%d] Account: %d, Amount: %d, Tax: %d, VAT: %d\n", i+1, dt.AccountItemID, dt.Amount, dt.TaxCode, dt.Vat)
			}
		}
		if len(d.Payments) > 0 {
			fmt.Println("Payments:")
			for i, p := range d.Payments {
				fmt.Printf("  [%d] Date: %s, Amount: %d, From: %s/%d\n", i+1, p.Date, p.Amount, p.FromWalletableType, p.FromWalletableID)
			}
		}
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deal",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		dealType, _ := cmd.Flags().GetString("type")
		date, _ := cmd.Flags().GetString("date")
		accountItemID, _ := cmd.Flags().GetInt64("account-item-id")
		amount, _ := cmd.Flags().GetInt64("amount")
		taxCode, _ := cmd.Flags().GetInt64("tax-code")
		partnerID, _ := cmd.Flags().GetInt64("partner-id")

		body := map[string]any{
			"company_id": client.CompanyID,
			"type":       dealType,
			"issue_date": date,
			"details": []map[string]any{
				{
					"account_item_id": accountItemID,
					"amount":          amount,
					"tax_code":        taxCode,
				},
			},
		}
		if partnerID != 0 {
			body["partner_id"] = partnerID
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/deals")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateDeal(body, &resp); err != nil {
			return err
		}

		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a deal",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a deal",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var dealID int64
		fmt.Sscanf(args[0], "%d", &dealID)

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/deals/%d\n", dealID)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteDeal(client.CompanyID, dealID)
	},
}

func buildListParams(cmd *cobra.Command) string {
	params := ""
	add := func(key, value string) {
		if value != "" {
			if params != "" {
				params += "&"
			}
			params += key + "=" + value
		}
	}
	t, _ := cmd.Flags().GetString("type")
	add("type", t)
	partner, _ := cmd.Flags().GetString("partner")
	add("partner_code", partner)
	from, _ := cmd.Flags().GetString("from")
	add("start_issue_date", from)
	to, _ := cmd.Flags().GetString("to")
	add("end_issue_date", to)
	status, _ := cmd.Flags().GetString("status")
	add("status", status)
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
