package model

type WalletablesResponse struct {
	Walletables []Walletable `json:"walletables"`
}

type WalletableResponse struct {
	Walletable Walletable `json:"walletable"`
}

type Walletable struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	BankID            int64  `json:"bank_id"`
	LastBalance       int64  `json:"last_balance"`
	WalletableBalance int64  `json:"walletable_balance"`
	UpdateDate        string `json:"update_date"`
}

type WalletableRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
