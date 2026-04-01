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
