package model

import "fmt"

type TransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

type TransferResponse struct {
	Transfer Transfer `json:"transfer"`
}

type Transfer struct {
	ID                 int64  `json:"id"`
	CompanyID          int64  `json:"company_id"`
	Amount             int64  `json:"amount"`
	Date               string `json:"date"`
	FromWalletableType string `json:"from_walletable_type"`
	FromWalletableID   int64  `json:"from_walletable_id"`
	ToWalletableType   string `json:"to_walletable_type"`
	ToWalletableID     int64  `json:"to_walletable_id"`
	Description        string `json:"description"`
}

type TransferRow struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"`
	Amount      int64  `json:"amount"`
	From        string `json:"from"`
	To          string `json:"to"`
	Description string `json:"description"`
}

func (t Transfer) ToRow() TransferRow {
	return TransferRow{
		ID:          t.ID,
		Date:        t.Date,
		Amount:      t.Amount,
		From:        fmt.Sprintf("%s:%d", t.FromWalletableType, t.FromWalletableID),
		To:          fmt.Sprintf("%s:%d", t.ToWalletableType, t.ToWalletableID),
		Description: t.Description,
	}
}
