package transfer

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

// Cmd is the transfer command group.
var Cmd = &cobra.Command{
	Use:   "transfer",
	Short: "Manage transfers (口座振替)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)

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
