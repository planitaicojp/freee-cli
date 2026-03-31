package manualjournal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
	"github.com/planitaicojp/freee-cli/internal/resolve"
)

// Cmd is the manual-journal command group.
var Cmd = &cobra.Command{
	Use:   "manual-journal",
	Short: "Manage manual journals (振替伝票)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)

	// update flags
	updateCmd.Flags().String("date", "", "issue date YYYY-MM-DD")
	updateCmd.Flags().Int64("debit-account-id", 0, "debit account item ID (simple mode: replaces all details)")
	updateCmd.Flags().String("debit-account-name", "", "debit account item name")
	updateCmd.Flags().Int64("credit-account-id", 0, "credit account item ID (simple mode: replaces all details)")
	updateCmd.Flags().String("credit-account-name", "", "credit account item name")
	updateCmd.Flags().Int64("amount", 0, "amount (simple mode)")
	updateCmd.Flags().Int64("tax-code", 0, "tax code (simple mode)")
	updateCmd.Flags().String("description", "", "description/memo (simple mode)")
	updateCmd.Flags().Bool("adjustment", false, "mark as closing/adjustment entry")
	updateCmd.Flags().String("details-json", "", "details as JSON string (replaces all details)")
	updateCmd.Flags().String("details-file", "", "details from JSON file path (replaces all details)")

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

	// create flags — simple mode
	createCmd.Flags().String("date", "", "issue date YYYY-MM-DD (required)")
	createCmd.Flags().Int64("debit-account-id", 0, "debit account item ID")
	createCmd.Flags().String("debit-account-name", "", "debit account item name (resolves to ID)")
	createCmd.Flags().Int64("credit-account-id", 0, "credit account item ID")
	createCmd.Flags().String("credit-account-name", "", "credit account item name (resolves to ID)")
	createCmd.Flags().Int64("amount", 0, "amount for simple 1:1 entry")
	createCmd.Flags().Int64("tax-code", 0, "tax code (required in simple mode, applied to both debit/credit)")
	createCmd.Flags().String("description", "", "description/memo")
	createCmd.Flags().Bool("adjustment", false, "mark as closing/adjustment entry")
	// create flags — JSON mode
	createCmd.Flags().String("details-json", "", "details as JSON string")
	createCmd.Flags().String("details-file", "", "details from JSON file path")
	_ = createCmd.MarkFlagRequired("date")
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

// isSimpleMode returns true if any simple-mode flag was provided.
func isSimpleMode(cmd *cobra.Command) bool {
	for _, f := range []string{"debit-account-id", "debit-account-name", "credit-account-id", "credit-account-name", "amount"} {
		if cmd.Flags().Changed(f) {
			return true
		}
	}
	return false
}

// isJSONMode returns true if any JSON-mode flag was provided.
func isJSONMode(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("details-json") || cmd.Flags().Changed("details-file")
}

// parseJSONDetails parses details from --details-json or --details-file flag.
func parseJSONDetails(cmd *cobra.Command) ([]map[string]any, error) {
	var raw string
	if cmd.Flags().Changed("details-json") {
		raw, _ = cmd.Flags().GetString("details-json")
	} else {
		path, _ := cmd.Flags().GetString("details-file")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("cannot read details file: %w", err)
		}
		raw = string(data)
	}
	var details []map[string]any
	if err := json.Unmarshal([]byte(raw), &details); err != nil {
		return nil, fmt.Errorf("invalid details JSON: %w", err)
	}
	return details, nil
}

// validateDetails checks balance and line count for details.
func validateDetails(details []map[string]any) error {
	if len(details) > 100 {
		return &cerrors.ValidationError{Message: "details exceed 100 lines\nhint: freee API allows max 100 debit+credit lines combined"}
	}
	var debitSum, creditSum int64
	for i, d := range details {
		tc, ok := d["tax_code"]
		if !ok || tc == nil {
			return &cerrors.ValidationError{Message: fmt.Sprintf("details[%d]: tax_code is required\nhint: each detail entry must include tax_code", i)}
		}
		amt, _ := d["amount"].(float64) // JSON numbers are float64
		side, _ := d["entry_side"].(string)
		switch side {
		case "debit":
			debitSum += int64(amt)
		case "credit":
			creditSum += int64(amt)
		default:
			return &cerrors.ValidationError{Message: fmt.Sprintf("details[%d]: entry_side must be 'debit' or 'credit', got %q", i, side)}
		}
	}
	if debitSum != creditSum {
		return &cerrors.ValidationError{Message: fmt.Sprintf("debit sum (%d) != credit sum (%d)\nhint: debit and credit totals must match", debitSum, creditSum)}
	}
	return nil
}

// resolveAccountID resolves account item ID from --<prefix>-account-id or --<prefix>-account-name.
// Uses resolve.AccountItemIDByName directly (no hidden flag mutation).
func resolveAccountID(cmd *cobra.Command, freeeAPI *api.FreeeAPI, companyID int64, prefix string) (int64, error) {
	idFlag := prefix + "-account-id"
	nameFlag := prefix + "-account-name"
	idChanged := cmd.Flags().Changed(idFlag)
	nameChanged := cmd.Flags().Changed(nameFlag)

	if idChanged && nameChanged {
		return 0, &cerrors.ValidationError{Message: fmt.Sprintf("--%s and --%s are mutually exclusive", idFlag, nameFlag)}
	}
	if idChanged {
		id, _ := cmd.Flags().GetInt64(idFlag)
		return id, nil
	}
	if nameChanged {
		name, _ := cmd.Flags().GetString(nameFlag)
		return resolve.AccountItemIDByName(name, freeeAPI, companyID)
	}
	return 0, nil
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a manual journal",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid manual journal ID: %s", args[0])
		}

		simple := isSimpleMode(cmd)
		jsonMode := isJSONMode(cmd)

		if simple && jsonMode {
			return &cerrors.ValidationError{Message: "simple mode flags and --details-json/--details-file are mutually exclusive"}
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		// GET current state for merge
		var current model.ManualJournalResponse
		if err := freeeAPI.GetManualJournal(client.CompanyID, id, &current); err != nil {
			return err
		}

		// Start with current values
		issueDate := current.ManualJournal.IssueDate
		adjustment := current.ManualJournal.Adjustment

		// Merge changed fields
		if cmd.Flags().Changed("date") {
			issueDate, _ = cmd.Flags().GetString("date")
		}
		if cmd.Flags().Changed("adjustment") {
			adjustment, _ = cmd.Flags().GetBool("adjustment")
		}

		// Build details
		var details any
		if simple {
			// Simple mode: replace all details with new 1:1 entry
			debitAccountID, err := resolveAccountID(cmd, freeeAPI, client.CompanyID, "debit")
			if err != nil {
				return err
			}
			if debitAccountID == 0 {
				return &cerrors.ValidationError{Message: "--debit-account-id or --debit-account-name is required in simple mode"}
			}
			creditAccountID, err := resolveAccountID(cmd, freeeAPI, client.CompanyID, "credit")
			if err != nil {
				return err
			}
			if creditAccountID == 0 {
				return &cerrors.ValidationError{Message: "--credit-account-id or --credit-account-name is required in simple mode"}
			}
			amount, _ := cmd.Flags().GetInt64("amount")
			if amount == 0 {
				return &cerrors.ValidationError{Message: "--amount is required in simple mode"}
			}
			taxCode, _ := cmd.Flags().GetInt64("tax-code")
			if taxCode == 0 {
				return &cerrors.ValidationError{Message: "--tax-code is required in simple mode"}
			}
			desc, _ := cmd.Flags().GetString("description")
			details = []map[string]any{
				{"entry_side": "debit", "account_item_id": debitAccountID, "amount": amount, "tax_code": taxCode, "description": desc},
				{"entry_side": "credit", "account_item_id": creditAccountID, "amount": amount, "tax_code": taxCode, "description": desc},
			}
		} else if jsonMode {
			parsed, err := parseJSONDetails(cmd)
			if err != nil {
				return err
			}
			if err := validateDetails(parsed); err != nil {
				return err
			}
			details = parsed
		} else {
			// No details change — re-send existing details
			details = current.ManualJournal.Details
		}

		body := map[string]any{
			"company_id": client.CompanyID,
			"issue_date": issueDate,
			"adjustment": adjustment,
			"details":    details,
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/manual_journals/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		var resp any
		if err := freeeAPI.UpdateManualJournal(id, body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a manual journal",
	RunE: func(cmd *cobra.Command, args []string) error {
		simple := isSimpleMode(cmd)
		jsonMode := isJSONMode(cmd)

		if simple && jsonMode {
			return &cerrors.ValidationError{Message: "simple mode flags and --details-json/--details-file are mutually exclusive\nhint: use one or the other"}
		}
		if !simple && !jsonMode {
			return &cerrors.ValidationError{Message: "must provide either simple mode flags (--debit-account-id, --amount, etc.) or --details-json/--details-file"}
		}

		date, _ := cmd.Flags().GetString("date")
		adjustment, _ := cmd.Flags().GetBool("adjustment")

		var details []map[string]any

		if simple {
			client, err := cmdutil.NewClient(cmd)
			if err != nil {
				return err
			}
			freeeAPI := &api.FreeeAPI{Client: client}

			debitAccountID, err := resolveAccountID(cmd, freeeAPI, client.CompanyID, "debit")
			if err != nil {
				return err
			}
			if debitAccountID == 0 {
				return &cerrors.ValidationError{Message: "--debit-account-id or --debit-account-name is required in simple mode"}
			}

			creditAccountID, err := resolveAccountID(cmd, freeeAPI, client.CompanyID, "credit")
			if err != nil {
				return err
			}
			if creditAccountID == 0 {
				return &cerrors.ValidationError{Message: "--credit-account-id or --credit-account-name is required in simple mode"}
			}

			amount, _ := cmd.Flags().GetInt64("amount")
			if amount == 0 {
				return &cerrors.ValidationError{Message: "--amount is required in simple mode"}
			}
			taxCode, _ := cmd.Flags().GetInt64("tax-code")
			if taxCode == 0 {
				return &cerrors.ValidationError{Message: "--tax-code is required in simple mode\nhint: freee API requires tax_code for each detail entry"}
			}
			desc, _ := cmd.Flags().GetString("description")

			details = []map[string]any{
				{"entry_side": "debit", "account_item_id": debitAccountID, "amount": amount, "tax_code": taxCode, "description": desc},
				{"entry_side": "credit", "account_item_id": creditAccountID, "amount": amount, "tax_code": taxCode, "description": desc},
			}

			body := map[string]any{
				"company_id": client.CompanyID,
				"issue_date": date,
				"adjustment": adjustment,
				"details":    details,
			}

			if cmdutil.IsDryRun(cmd) {
				fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/manual_journals")
				return output.New("json").Format(os.Stdout, body)
			}

			var resp any
			if err := freeeAPI.CreateManualJournal(body, &resp); err != nil {
				return err
			}
			opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
			return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
		}

		// JSON mode
		var err error
		details, err = parseJSONDetails(cmd)
		if err != nil {
			return err
		}
		if err := validateDetails(details); err != nil {
			return err
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		body := map[string]any{
			"company_id": client.CompanyID,
			"issue_date": date,
			"adjustment": adjustment,
			"details":    details,
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/manual_journals")
			return output.New("json").Format(os.Stdout, body)
		}

		var resp any
		if err := freeeAPI.CreateManualJournal(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}
