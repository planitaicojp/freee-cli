package transfer

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

// Cmd is the transfer command group.
var Cmd = &cobra.Command{
	Use:   "transfer",
	Short: "Manage transfers (口座振替)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	// create flags
	createCmd.Flags().String("date", "", "transfer date YYYY-MM-DD (required)")
	createCmd.Flags().Int64("amount", 0, "transfer amount (required)")
	createCmd.Flags().String("from-type", "", "source account type: bank_account, credit_card, wallet (required)")
	createCmd.Flags().Int64("from-id", 0, "source account ID (required)")
	createCmd.Flags().String("to-type", "", "destination account type: bank_account, credit_card, wallet (required)")
	createCmd.Flags().Int64("to-id", 0, "destination account ID (required)")
	createCmd.Flags().String("description", "", "description/memo")
	_ = createCmd.MarkFlagRequired("date")
	_ = createCmd.MarkFlagRequired("amount")
	_ = createCmd.MarkFlagRequired("from-type")
	_ = createCmd.MarkFlagRequired("from-id")
	_ = createCmd.MarkFlagRequired("to-type")
	_ = createCmd.MarkFlagRequired("to-id")

	// update flags
	updateCmd.Flags().String("date", "", "transfer date YYYY-MM-DD")
	updateCmd.Flags().Int64("amount", 0, "transfer amount")
	updateCmd.Flags().String("from-type", "", "source account type: bank_account, credit_card, wallet")
	updateCmd.Flags().Int64("from-id", 0, "source account ID")
	updateCmd.Flags().String("to-type", "", "destination account type: bank_account, credit_card, wallet")
	updateCmd.Flags().Int64("to-id", 0, "destination account ID")
	updateCmd.Flags().String("description", "", "description/memo")

	// list flags
	listCmd.Flags().String("from", "", "start date (YYYY-MM-DD)")
	listCmd.Flags().String("to", "", "end date (YYYY-MM-DD)")
	listCmd.Flags().Int("limit", 20, "max number of results per page")
	listCmd.Flags().Int("offset", 0, "offset for pagination")
	listCmd.Flags().Bool("all", false, "fetch all pages automatically")
}

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
	add("start_date", from)
	to, _ := cmd.Flags().GetString("to")
	add("end_date", to)
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
	Short: "List transfers",
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
			baseParams := buildBaseListParams(cmd)
			var all []model.Transfer
			for offset := 0; ; offset += fetchLimit {
				p := baseParams
				if p != "" {
					p += "&"
				}
				p += fmt.Sprintf("limit=%d&offset=%d", fetchLimit, offset)
				var resp model.TransfersResponse
				if err := freeeAPI.ListTransfers(client.CompanyID, p, &resp); err != nil {
					return err
				}
				if len(resp.Transfers) == 0 {
					break
				}
				all = append(all, resp.Transfers...)
				if len(resp.Transfers) < fetchLimit {
					break
				}
			}
			if format != "" && format != "table" {
				return output.New(format, opts).Format(os.Stdout, map[string]any{"transfers": all})
			}
			rows := make([]model.TransferRow, len(all))
			for i, tr := range all {
				rows[i] = tr.ToRow()
			}
			return output.New("table", opts).Format(os.Stdout, rows)
		}

		params := buildListParams(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListTransfers(client.CompanyID, params, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.TransfersResponse
		if err := freeeAPI.ListTransfers(client.CompanyID, params, &resp); err != nil {
			return err
		}
		rows := make([]model.TransferRow, len(resp.Transfers))
		for i, tr := range resp.Transfers {
			rows[i] = tr.ToRow()
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

var validWalletableTypes = map[string]bool{
	"bank_account": true,
	"credit_card":  true,
	"wallet":       true,
}

func validateWalletableType(flagName, value string) error {
	if !validWalletableTypes[value] {
		return &cerrors.ValidationError{
			Message: fmt.Sprintf("--%s must be one of: bank_account, credit_card, wallet (got %q)\nhint: run 'freee walletable list' to see available accounts", flagName, value),
		}
	}
	return nil
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromType, _ := cmd.Flags().GetString("from-type")
		if err := validateWalletableType("from-type", fromType); err != nil {
			return err
		}
		toType, _ := cmd.Flags().GetString("to-type")
		if err := validateWalletableType("to-type", toType); err != nil {
			return err
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		amount, _ := cmd.Flags().GetInt64("amount")
		fromID, _ := cmd.Flags().GetInt64("from-id")
		toID, _ := cmd.Flags().GetInt64("to-id")
		description, _ := cmd.Flags().GetString("description")

		body := map[string]any{
			"company_id":           client.CompanyID,
			"date":                 date,
			"amount":               amount,
			"from_walletable_type": fromType,
			"from_walletable_id":   fromID,
			"to_walletable_type":   toType,
			"to_walletable_id":     toID,
			"description":          description,
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/transfers")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateTransfer(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a transfer",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transfer ID: %s", args[0])
		}

		// Validate walletable types if provided
		if cmd.Flags().Changed("from-type") {
			fromType, _ := cmd.Flags().GetString("from-type")
			if err := validateWalletableType("from-type", fromType); err != nil {
				return err
			}
		}
		if cmd.Flags().Changed("to-type") {
			toType, _ := cmd.Flags().GetString("to-type")
			if err := validateWalletableType("to-type", toType); err != nil {
				return err
			}
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		// GET current state for merge
		var current model.TransferResponse
		if err := freeeAPI.GetTransfer(client.CompanyID, id, &current); err != nil {
			return err
		}

		// Start with current values
		tr := current.Transfer
		date := tr.Date
		amount := tr.Amount
		fromType := tr.FromWalletableType
		fromID := tr.FromWalletableID
		toType := tr.ToWalletableType
		toID := tr.ToWalletableID
		description := tr.Description

		// Merge changed fields
		if cmd.Flags().Changed("date") {
			date, _ = cmd.Flags().GetString("date")
		}
		if cmd.Flags().Changed("amount") {
			amount, _ = cmd.Flags().GetInt64("amount")
		}
		if cmd.Flags().Changed("from-type") {
			fromType, _ = cmd.Flags().GetString("from-type")
		}
		if cmd.Flags().Changed("from-id") {
			fromID, _ = cmd.Flags().GetInt64("from-id")
		}
		if cmd.Flags().Changed("to-type") {
			toType, _ = cmd.Flags().GetString("to-type")
		}
		if cmd.Flags().Changed("to-id") {
			toID, _ = cmd.Flags().GetInt64("to-id")
		}
		if cmd.Flags().Changed("description") {
			description, _ = cmd.Flags().GetString("description")
		}

		body := map[string]any{
			"company_id":           client.CompanyID,
			"date":                 date,
			"amount":               amount,
			"from_walletable_type": fromType,
			"from_walletable_id":   fromID,
			"to_walletable_type":   toType,
			"to_walletable_id":     toID,
			"description":          description,
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/transfers/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		var resp any
		if err := freeeAPI.UpdateTransfer(id, body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a transfer",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transfer ID: %s", args[0])
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] DELETE /api/1/transfers/%d\n", id)
			return nil
		}

		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteTransfer(client.CompanyID, id)
	},
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show transfer details",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transfer ID: %s", args[0])
		}

		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.GetTransfer(client.CompanyID, id, &resp); err != nil {
				return err
			}
			return output.New(format).Format(os.Stdout, resp)
		}

		var resp model.TransferResponse
		if err := freeeAPI.GetTransfer(client.CompanyID, id, &resp); err != nil {
			return err
		}
		tr := resp.Transfer
		fmt.Printf("ID:          %d\n", tr.ID)
		fmt.Printf("Date:        %s\n", tr.Date)
		fmt.Printf("Amount:      %d\n", tr.Amount)
		fmt.Printf("From:        %s:%d\n", tr.FromWalletableType, tr.FromWalletableID)
		fmt.Printf("To:          %s:%d\n", tr.ToWalletableType, tr.ToWalletableID)
		if tr.Description != "" {
			fmt.Printf("Description: %s\n", tr.Description)
		}
		return nil
	},
}
