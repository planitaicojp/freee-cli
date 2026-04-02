# Wallet-Txn Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `freee wallet-txn` command with list/show/create/delete subcommands for managing wallet transactions (口座明細).

**Architecture:** Single-file command package (`cmd/wallettxn/`) following the transfer pattern. Model structs in `internal/model/wallet_txn.go`. API methods appended to `internal/api/freee.go`. Extract `validateWalletableType` from transfer to `cmd/cmdutil/` for sharing.

**Tech Stack:** Go, cobra, existing internal packages (api, model, output, errors, cmdutil)

**Spec:** `docs/superpowers/specs/2026-04-01-wallet-txn-design.md`

---

### Task 1: Extract validateWalletableType to cmdutil

**Files:**
- Modify: `cmd/cmdutil/client.go` (add validation function)
- Modify: `cmd/transfer/transfer.go` (use shared function)

- [ ] **Step 1: Add ValidateWalletableType to cmdutil**

Add to `cmd/cmdutil/client.go`:

```go
import (
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

// ValidateWalletableType validates that the given value is a valid walletable type.
var ValidWalletableTypes = map[string]bool{
	"bank_account": true,
	"credit_card":  true,
	"wallet":       true,
}

func ValidateWalletableType(flagName, value string) error {
	if !ValidWalletableTypes[value] {
		return &cerrors.ValidationError{
			Message: fmt.Sprintf("--%s must be one of: bank_account, credit_card, wallet (got %q)\nhint: run 'freee walletable list' to see available accounts", flagName, value),
		}
	}
	return nil
}
```

- [ ] **Step 2: Update transfer.go to use shared function**

Replace `validWalletableTypes` map and `validateWalletableType` function in `cmd/transfer/transfer.go` with calls to `cmdutil.ValidateWalletableType`.

Remove:
```go
var validWalletableTypes = map[string]bool{...}
func validateWalletableType(flagName, value string) error {...}
```

Replace all `validateWalletableType(` with `cmdutil.ValidateWalletableType(`.

- [ ] **Step 3: Build and verify**

Run: `cd /root/dev/planitai/planitai-freee-cli && go build ./...`
Expected: clean build, no errors

- [ ] **Step 4: Commit**

```bash
git add cmd/cmdutil/client.go cmd/transfer/transfer.go
git commit -m "refactor: extract ValidateWalletableType to cmdutil for sharing"
```

---

### Task 2: Add WalletTxn model

**Files:**
- Create: `internal/model/wallet_txn.go`

- [ ] **Step 1: Create wallet_txn.go**

```go
package model

import "fmt"

// WalletTxnsResponse wraps the wallet_txns list API response.
type WalletTxnsResponse struct {
	WalletTxns []WalletTxn `json:"wallet_txns"`
}

// WalletTxnResponse wraps a single wallet_txn API response.
type WalletTxnResponse struct {
	WalletTxn WalletTxn `json:"wallet_txn"`
}

// WalletTxn represents a wallet transaction (口座明細).
type WalletTxn struct {
	ID             int64  `json:"id"`
	CompanyID      int64  `json:"company_id"`
	Date           string `json:"date"`
	Amount         int64  `json:"amount"`
	DueAmount      int64  `json:"due_amount"`
	Balance        *int64 `json:"balance"`
	EntrySide      string `json:"entry_side"`
	WalletableType string `json:"walletable_type"`
	WalletableID   int64  `json:"walletable_id"`
	Description    string `json:"description"`
	Status         int    `json:"status"`
	RuleMatched    bool   `json:"rule_matched"`
}

// WalletTxnRow is the table output row.
type WalletTxnRow struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"`
	EntrySide   string `json:"entry_side"`
	Amount      int64  `json:"amount"`
	Walletable  string `json:"walletable"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// StatusLabel maps numeric status codes to human-readable labels.
var walletTxnStatusLabels = map[int]string{
	1: "waiting",
	2: "settled",
	3: "ignored",
	4: "settling",
	6: "excluded",
}

// StatusLabel returns the human-readable label for the status code.
func (w WalletTxn) StatusLabel() string {
	if label, ok := walletTxnStatusLabels[w.Status]; ok {
		return label
	}
	return fmt.Sprintf("unknown(%d)", w.Status)
}

// ToRow converts a WalletTxn to a table row.
func (w WalletTxn) ToRow() WalletTxnRow {
	return WalletTxnRow{
		ID:          w.ID,
		Date:        w.Date,
		EntrySide:   w.EntrySide,
		Amount:      w.Amount,
		Walletable:  fmt.Sprintf("%s:%d", w.WalletableType, w.WalletableID),
		Status:      w.StatusLabel(),
		Description: w.Description,
	}
}
```

- [ ] **Step 2: Build and verify**

Run: `go build ./...`
Expected: clean build

- [ ] **Step 3: Commit**

```bash
git add internal/model/wallet_txn.go
git commit -m "feat: add WalletTxn model with status label mapping"
```

---

### Task 3: Add WalletTxn API methods

**Files:**
- Modify: `internal/api/freee.go` (append 4 methods after Transfers section)

- [ ] **Step 1: Add API methods**

Append after the Transfers section in `internal/api/freee.go`:

```go
// --- Wallet Transactions ---

func (a *FreeeAPI) ListWalletTxns(companyID int64, params string, result any) error {
	url := a.url("/wallet_txns?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetWalletTxn(companyID, id int64, result any) error {
	return a.Client.Get(a.url("/wallet_txns/%d?company_id=%d", id, companyID), result)
}

func (a *FreeeAPI) CreateWalletTxn(body, result any) error {
	_, err := a.Client.Post(a.url("/wallet_txns"), body, result)
	return err
}

func (a *FreeeAPI) DeleteWalletTxn(companyID, id int64) error {
	return a.Client.Delete(a.url("/wallet_txns/%d?company_id=%d", id, companyID))
}
```

- [ ] **Step 2: Build and verify**

Run: `go build ./...`
Expected: clean build

- [ ] **Step 3: Commit**

```bash
git add internal/api/freee.go
git commit -m "feat: add WalletTxn API methods (list/get/create/delete)"
```

---

### Task 4: Implement wallet-txn command

**Files:**
- Create: `cmd/wallettxn/wallettxn.go`
- Modify: `cmd/root.go` (register command)

- [ ] **Step 1: Create cmd/wallettxn/wallettxn.go**

Follow `cmd/transfer/transfer.go` pattern. Key differences from transfer:
- No update command
- Additional list filters: `--walletable-type`, `--walletable-id`, `--entry-side`
- Walletable pair validation (both must be specified together)
- Entry-side validation (income/expense only)
- Status label conversion in table output
- Balance field (nullable) in show output
- Create has `--balance` optional flag
- Create supports `--format`

```go
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
```

- [ ] **Step 2: Register in root.go**

Add import:
```go
"github.com/planitaicojp/freee-cli/cmd/wallettxn"
```

Add command registration after `walletable.Cmd` (before `schema.NewCmd`):
```go
rootCmd.AddCommand(wallettxn.Cmd)
```

- [ ] **Step 3: Build and verify**

Run: `go build ./...`
Expected: clean build

- [ ] **Step 4: Smoke test help output**

Run: `go run . wallet-txn --help`
Expected: Shows list, show, create, delete subcommands

Run: `go run . wallet-txn create --help`
Expected: Shows required flags (walletable-type, walletable-id, entry-side, amount, date)

- [ ] **Step 5: Commit**

```bash
git add cmd/wallettxn/wallettxn.go cmd/root.go
git commit -m "feat: add wallet-txn command (list/show/create/delete)"
```

---

### Task 5: Update SPEC.md

**Files:**
- Modify: `SPEC.md`

- [ ] **Step 1: Update Phase 2 table**

Change line 927:
```
| 63–66 | `freee wallet-txn` CRUD | 구좌 명세 (口座明細) | ☐ |
```
To:
```
| 63–65 | `freee wallet-txn` list/show/create | 구좌 명세 (口座明細) | ☐ |
| 66 | `freee wallet-txn` delete | 구좌 명세 삭제 | ☐ |
```

- [ ] **Step 2: Update command count**

Update Phase 2 row in summary table: 19 -> 18 commands.
Update total: 272 -> 271.

- [ ] **Step 3: Commit**

```bash
git add SPEC.md
git commit -m "docs: update SPEC.md wallet-txn command count (no update API)"
```
