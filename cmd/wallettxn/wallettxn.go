package wallettxn

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the wallet-txn command group.
var Cmd = &cobra.Command{
	Use:   "wallet-txn",
	Short: "Manage wallet transactions (口座明細)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(deleteCmd)

	// list flags
	listCmd.Flags().String("walletable-type", "", "account type: bank_account, credit_card, wallet (requires --walletable-id)")
	listCmd.Flags().Int64("walletable-id", 0, "account ID (requires --walletable-type)")
	listCmd.Flags().String("entry-side", "", "filter by: income, expense")
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
	listCmd.Flags().Int("limit", 20, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")

	// create flags
	createCmd.Flags().String("walletable-type", "", "account type: bank_account, credit_card, wallet (required)")
	createCmd.Flags().Int64("walletable-id", 0, "account ID (required)")
	createCmd.Flags().String("entry-side", "", "income or expense (required)")
	createCmd.Flags().Int64("amount", 0, "transaction amount (required)")
	createCmd.Flags().String("date", "", "transaction date YYYY-MM-DD (required)")
	createCmd.Flags().String("description", "", "transaction description")
	createCmd.Flags().Int64("balance", 0, "account balance after transaction")
	_ = createCmd.MarkFlagRequired("walletable-type")
	_ = createCmd.MarkFlagRequired("walletable-id")
	_ = createCmd.MarkFlagRequired("entry-side")
	_ = createCmd.MarkFlagRequired("amount")
	_ = createCmd.MarkFlagRequired("date")
}

var validEntrySides = map[string]bool{
	"income":  true,
	"expense": true,
}

func validateEntrySide(value string) error {
	if !validEntrySides[value] {
		return &cerrors.ValidationError{
			Message: fmt.Sprintf("--entry-side must be one of: income, expense (got %q)", value),
		}
	}
	return nil
}

func buildBaseListParams(cmd *cobra.Command) (string, error) {
	params := ""
	add := func(key, value string) {
		if value != "" {
			if params != "" {
				params += "&"
			}
			params += key + "=" + value
		}
	}

	wType, _ := cmd.Flags().GetString("walletable-type")
	wID, _ := cmd.Flags().GetInt64("walletable-id")
	hasType := cmd.Flags().Changed("walletable-type")
	hasID := cmd.Flags().Changed("walletable-id")
	if hasType != hasID {
		return "", &cerrors.ValidationError{
			Message: "both --walletable-type and --walletable-id must be specified together\nhint: run 'freee walletable list' to see available accounts",
		}
	}
	if hasType {
		if err := cmdutil.ValidateWalletableType("walletable-type", wType); err != nil {
			return "", err
		}
		add("walletable_type", wType)
		add("walletable_id", fmt.Sprintf("%d", wID))
	}

	if es, _ := cmd.Flags().GetString("entry-side"); es != "" {
		if err := validateEntrySide(es); err != nil {
			return "", err
		}
		add("entry_side", es)
	}

	from, _ := cmd.Flags().GetString("from")
	add("start_date", from)
	to, _ := cmd.Flags().GetString("to")
	add("end_date", to)

	return params, nil
}

func buildListParams(cmd *cobra.Command) (string, error) {
	params, err := buildBaseListParams(cmd)
	if err != nil {
		return "", err
	}
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
	return params, nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallet transactions",
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
			const fetchLimit = 100
			baseParams, err := buildBaseListParams(cmd)
			if err != nil {
				return err
			}
			var all []model.WalletTxn
			for offset := 0; ; offset += fetchLimit {
				p := baseParams
				if p != "" {
					p += "&"
				}
				p += fmt.Sprintf("limit=%d&offset=%d", fetchLimit, offset)
				var resp model.WalletTxnsResponse
				if err := freeeAPI.ListWalletTxns(client.CompanyID, p, &resp); err != nil {
					return err
				}
				if len(resp.WalletTxns) == 0 {
					break
				}
				all = append(all, resp.WalletTxns...)
				if len(resp.WalletTxns) < fetchLimit {
					break
				}
			}
			if format != "" && format != "table" {
				return output.New(format, opts).Format(os.Stdout, map[string]any{"wallet_txns": all})
			}
			rows := make([]model.WalletTxnRow, len(all))
			for i, wt := range all {
				rows[i] = wt.ToRow()
			}
			return output.New("table", opts).Format(os.Stdout, rows)
		}

		params, err := buildListParams(cmd)
		if err != nil {
			return err
		}
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListWalletTxns(client.CompanyID, params, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.WalletTxnsResponse
		if err := freeeAPI.ListWalletTxns(client.CompanyID, params, &resp); err != nil {
			return err
		}
		rows := make([]model.WalletTxnRow, len(resp.WalletTxns))
		for i, wt := range resp.WalletTxns {
			rows[i] = wt.ToRow()
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show wallet transaction details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid wallet transaction ID: %s", args[0])
		}

		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetWalletTxn(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.WalletTxnResponse
		if err := freeeAPI.GetWalletTxn(client.CompanyID, id, &resp); err != nil {
			return err
		}
		wt := resp.WalletTxn
		fmt.Printf("ID:           %d\n", wt.ID)
		fmt.Printf("Date:         %s\n", wt.Date)
		fmt.Printf("Entry Side:   %s\n", wt.EntrySide)
		fmt.Printf("Amount:       %d\n", wt.Amount)
		fmt.Printf("Due Amount:   %d\n", wt.DueAmount)
		if wt.Balance != nil {
			fmt.Printf("Balance:      %d\n", *wt.Balance)
		}
		fmt.Printf("Walletable:   %s:%d\n", wt.WalletableType, wt.WalletableID)
		fmt.Printf("Status:       %s\n", wt.StatusLabel())
		if wt.Description != "" {
			fmt.Printf("Description:  %s\n", wt.Description)
		}
		fmt.Printf("Rule Matched: %t\n", wt.RuleMatched)
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		wType, _ := cmd.Flags().GetString("walletable-type")
		if err := cmdutil.ValidateWalletableType("walletable-type", wType); err != nil {
			return err
		}
		entrySide, _ := cmd.Flags().GetString("entry-side")
		if err := validateEntrySide(entrySide); err != nil {
			return err
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		wID, _ := cmd.Flags().GetInt64("walletable-id")
		amount, _ := cmd.Flags().GetInt64("amount")
		date, _ := cmd.Flags().GetString("date")
		description, _ := cmd.Flags().GetString("description")

		body := map[string]any{
			"company_id":      client.CompanyID,
			"walletable_type": wType,
			"walletable_id":   wID,
			"entry_side":      entrySide,
			"amount":          amount,
			"date":            date,
		}
		if description != "" {
			body["description"] = description
		}
		if cmd.Flags().Changed("balance") {
			balance, _ := cmd.Flags().GetInt64("balance")
			body["balance"] = balance
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/wallet_txns")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateWalletTxn(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a wallet transaction",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid wallet transaction ID: %s", args[0])
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/wallet_txns/%d\n", id)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		if err := freeeAPI.DeleteWalletTxn(client.CompanyID, id); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Deleted wallet transaction %d\n", id)
		return nil
	},
}
