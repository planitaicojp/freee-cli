package model

type DealsResponse struct {
	Deals []Deal `json:"deals"`
}

type DealResponse struct {
	Deal Deal `json:"deal"`
}

type Deal struct {
	ID          int64        `json:"id"`
	CompanyID   int64        `json:"company_id"`
	IssueDate   string       `json:"issue_date"`
	DueDate     string       `json:"due_date"`
	Type        string       `json:"type"`
	PartnerID   int64        `json:"partner_id"`
	PartnerCode string       `json:"partner_code"`
	RefNumber   string       `json:"ref_number"`
	Status      string       `json:"status"`
	Amount      int64        `json:"amount"`
	DueAmount   int64        `json:"due_amount"`
	Details     []DealDetail `json:"details"`
	Payments    []DealPayment `json:"payments"`
}

type DealDetail struct {
	ID            int64  `json:"id"`
	AccountItemID int64  `json:"account_item_id"`
	TaxCode       int    `json:"tax_code"`
	Amount        int64  `json:"amount"`
	Vat           int64  `json:"vat"`
	Description   string `json:"description"`
	EntrySide     string `json:"entry_side"`
}

type DealPayment struct {
	ID                  int64  `json:"id"`
	Date                string `json:"date"`
	FromWalletableType  string `json:"from_walletable_type"`
	FromWalletableID    int64  `json:"from_walletable_id"`
	Amount              int64  `json:"amount"`
}

type DealRow struct {
	ID     int64  `json:"id"`
	Date   string `json:"date"`
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
}
