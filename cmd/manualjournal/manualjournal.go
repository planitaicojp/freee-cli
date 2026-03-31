package manualjournal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the manual-journal command group.
var Cmd = &cobra.Command{
	Use:   "manual-journal",
	Short: "Manage manual journals (振替伝票)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)

	// list flags
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
	listCmd.Flags().String("entry-side", "", "filter by entry side: debit or credit")
	listCmd.Flags().Int64("account-item-id", 0, "account item ID filter")
	listCmd.Flags().Int64("partner-id", 0, "partner ID filter")
	listCmd.Flags().String("adjustment", "", "filter: only or without")
	listCmd.Flags().Int64("min-amount", 0, "minimum amount filter")
	listCmd.Flags().Int64("max-amount", 0, "maximum amount filter")
	listCmd.Flags().Int64("item-id", 0, "item ID filter")
	listCmd.Flags().Int64("section-id", 0, "section ID filter")
	listCmd.Flags().String("txn-number", "", "journal number filter")
	listCmd.Flags().Int("limit", 50, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")
}

// buildBaseListParams builds query params excluding pagination.
func buildBaseListParams(cmd *cobra.Command) string {
	params := ""
	add := func(key, value string) {
		if value != "" {
			if params != "" {
				params += "&"
			}
			params += key + "=" + value
		}
	}
	from, _ := cmd.Flags().GetString("from")
	add("start_issue_date", from)
	to, _ := cmd.Flags().GetString("to")
	add("end_issue_date", to)
	entrySide, _ := cmd.Flags().GetString("entry-side")
	add("entry_side", entrySide)
	if v, _ := cmd.Flags().GetInt64("account-item-id"); v != 0 {
		add("account_item_id", fmt.Sprintf("%d", v))
	}
	if v, _ := cmd.Flags().GetInt64("partner-id"); v != 0 {
		add("partner_id", fmt.Sprintf("%d", v))
	}
	adj, _ := cmd.Flags().GetString("adjustment")
	add("adjustment", adj)
	if v, _ := cmd.Flags().GetInt64("min-amount"); v != 0 {
		add("min_amount", fmt.Sprintf("%d", v))
	}
	if v, _ := cmd.Flags().GetInt64("max-amount"); v != 0 {
		add("max_amount", fmt.Sprintf("%d", v))
	}
	if v, _ := cmd.Flags().GetInt64("item-id"); v != 0 {
		add("item_id", fmt.Sprintf("%d", v))
	}
	if v, _ := cmd.Flags().GetInt64("section-id"); v != 0 {
		add("section_id", fmt.Sprintf("%d", v))
	}
	txn, _ := cmd.Flags().GetString("txn-number")
	add("txn_number", txn)
	return params
}

func buildListParams(cmd *cobra.Command) string {
	params := buildBaseListParams(cmd)
	add := func(key, value string) {
		if value != "" {
			if params != "" {
				params += "&"
			}
			params += key + "=" + value
		}
	}
	if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
		add("limit", fmt.Sprintf("%d", v))
	}
	if v, _ := cmd.Flags().GetInt("offset"); v > 0 {
		add("offset", fmt.Sprintf("%d", v))
	}
	return params
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List manual journals",
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

		format := cmdutil.GetFormat(cmd)
		fetchAll := cmdutil.IsAll(cmd)
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}

		if fetchAll {
			limit, _ := cmd.Flags().GetInt("limit")
			if limit <= 0 {
				limit = 100
			}
			baseParams := buildBaseListParams(cmd)
			var all []model.ManualJournal
			for offset := 0; ; offset += limit {
				p := baseParams
				if p != "" {
					p += "&"
				}
				p += fmt.Sprintf("limit=%d&offset=%d", limit, offset)
				var resp model.ManualJournalsResponse
				if err := freeeAPI.ListManualJournals(client.CompanyID, p, &resp); err != nil {
					return err
				}
				all = append(all, resp.ManualJournals...)
				if len(resp.ManualJournals) < limit {
					break
				}
			}
			if format != "" && format != "table" {
				return output.New(format, opts).Format(os.Stdout, map[string]any{"manual_journals": all})
			}
			rows := make([]model.ManualJournalRow, len(all))
			for i, mj := range all {
				rows[i] = mj.ToRow()
			}
			return output.New("table", opts).Format(os.Stdout, rows)
		}

		params := buildListParams(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListManualJournals(client.CompanyID, params, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.ManualJournalsResponse
		if err := freeeAPI.ListManualJournals(client.CompanyID, params, &resp); err != nil {
			return err
		}
		rows := make([]model.ManualJournalRow, len(resp.ManualJournals))
		for i, mj := range resp.ManualJournals {
			rows[i] = mj.ToRow()
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show manual journal details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid manual journal ID: %s", args[0])
		}

		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetManualJournal(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.ManualJournalResponse
		if err := freeeAPI.GetManualJournal(client.CompanyID, id, &resp); err != nil {
			return err
		}
		mj := resp.ManualJournal
		fmt.Printf("ID:         %d\n", mj.ID)
		fmt.Printf("Date:       %s\n", mj.IssueDate)
		adj := "No"
		if mj.Adjustment {
			adj = "Yes"
		}
		fmt.Printf("Adjustment: %s\n", adj)
		if mj.TxnNumber != "" {
			fmt.Printf("TxnNumber:  %s\n", mj.TxnNumber)
		}
		if len(mj.Details) > 0 {
			fmt.Println("Details:")
			for i, d := range mj.Details {
				fmt.Printf("  [%d] %s Account:%d Amount:%d VAT:%d Tax:%d",
					i+1, d.EntrySide, d.AccountItemID, d.Amount, d.Vat, d.TaxCode)
				if d.PartnerID != 0 {
					fmt.Printf(" Partner:%d", d.PartnerID)
				}
				if d.ItemID != 0 {
					fmt.Printf(" Item:%d", d.ItemID)
				}
				if d.SectionID != 0 {
					fmt.Printf(" Section:%d", d.SectionID)
				}
				if d.Description != "" {
					fmt.Printf(" %q", d.Description)
				}
				fmt.Println()
			}
		}
		return nil
	},
}
