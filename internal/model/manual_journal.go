package model

import "fmt"

type ManualJournalsResponse struct {
	ManualJournals []ManualJournal `json:"manual_journals"`
}

type ManualJournalResponse struct {
	ManualJournal ManualJournal `json:"manual_journal"`
}

type ManualJournal struct {
	ID         int64                 `json:"id"`
	CompanyID  int64                 `json:"company_id"`
	IssueDate  string                `json:"issue_date"`
	Adjustment bool                  `json:"adjustment"`
	TxnNumber  string                `json:"txn_number"`
	Details    []ManualJournalDetail `json:"details"`
	ReceiptIDs []int64               `json:"receipt_ids"`
}

type ManualJournalDetail struct {
	ID            int64   `json:"id"`
	EntrySide     string  `json:"entry_side"`
	AccountItemID int64   `json:"account_item_id"`
	Amount        int64   `json:"amount"`
	Vat           int64   `json:"vat"`
	TaxCode       int64   `json:"tax_code"`
	PartnerID     int64   `json:"partner_id"`
	ItemID        int64   `json:"item_id"`
	SectionID     int64   `json:"section_id"`
	TagIDs        []int64 `json:"tag_ids"`
	Description   string  `json:"description"`
}

type ManualJournalRow struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"`
	Adjustment  string `json:"adjustment"`
	Entries     string `json:"entries"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// ToRow converts a ManualJournal to a ManualJournalRow for table display.
func (mj ManualJournal) ToRow() ManualJournalRow {
	adj := "No"
	if mj.Adjustment {
		adj = "Yes"
	}

	var debitCount, creditCount int
	var debitSum int64
	var desc string
	for _, d := range mj.Details {
		switch d.EntrySide {
		case "debit":
			debitCount++
			debitSum += d.Amount
		case "credit":
			creditCount++
		}
		if desc == "" && d.Description != "" {
			desc = d.Description
		}
	}

	return ManualJournalRow{
		ID:          mj.ID,
		Date:        mj.IssueDate,
		Adjustment:  adj,
		Entries:     fmt.Sprintf("%d debit / %d credit", debitCount, creditCount),
		Amount:      debitSum,
		Description: desc,
	}
}
