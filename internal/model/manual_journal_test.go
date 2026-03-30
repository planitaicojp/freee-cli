package model

import (
	"encoding/json"
	"testing"
)

func TestManualJournalUnmarshal(t *testing.T) {
	raw := `{
		"manual_journals": [{
			"id": 1,
			"company_id": 100,
			"issue_date": "2026-03-15",
			"adjustment": false,
			"txn_number": "MJ-001",
			"details": [
				{"id": 10, "entry_side": "debit", "account_item_id": 200, "amount": 50000, "vat": 5000, "tax_code": 1, "partner_id": 0, "item_id": 0, "section_id": 0, "tag_ids": [1,2], "description": "test debit"},
				{"id": 11, "entry_side": "credit", "account_item_id": 300, "amount": 50000, "vat": 5000, "tax_code": 1, "partner_id": 5, "item_id": 0, "section_id": 0, "tag_ids": null, "description": "test credit"}
			],
			"receipt_ids": [100, 200]
		}]
	}`
	var resp ManualJournalsResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(resp.ManualJournals) != 1 {
		t.Fatalf("expected 1 journal, got %d", len(resp.ManualJournals))
	}
	mj := resp.ManualJournals[0]
	if mj.ID != 1 {
		t.Errorf("ID = %d, want 1", mj.ID)
	}
	if mj.IssueDate != "2026-03-15" {
		t.Errorf("IssueDate = %q, want 2026-03-15", mj.IssueDate)
	}
	if len(mj.Details) != 2 {
		t.Fatalf("details count = %d, want 2", len(mj.Details))
	}
	if mj.Details[0].EntrySide != "debit" {
		t.Errorf("Details[0].EntrySide = %q, want debit", mj.Details[0].EntrySide)
	}
	if mj.Details[0].Amount != 50000 {
		t.Errorf("Details[0].Amount = %d, want 50000", mj.Details[0].Amount)
	}
	if mj.Details[1].PartnerID != 5 {
		t.Errorf("Details[1].PartnerID = %d, want 5", mj.Details[1].PartnerID)
	}
	if len(mj.Details[0].TagIDs) != 2 {
		t.Errorf("Details[0].TagIDs len = %d, want 2", len(mj.Details[0].TagIDs))
	}
	if len(mj.ReceiptIDs) != 2 {
		t.Errorf("ReceiptIDs len = %d, want 2", len(mj.ReceiptIDs))
	}
}

func TestManualJournalSingleUnmarshal(t *testing.T) {
	raw := `{"manual_journal": {"id": 42, "company_id": 1, "issue_date": "2026-01-01", "adjustment": true, "txn_number": "", "details": [], "receipt_ids": []}}`
	var resp ManualJournalResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ManualJournal.ID != 42 {
		t.Errorf("ID = %d, want 42", resp.ManualJournal.ID)
	}
	if !resp.ManualJournal.Adjustment {
		t.Error("Adjustment should be true")
	}
}

func TestManualJournalRowConversion(t *testing.T) {
	mj := ManualJournal{
		ID:         1,
		IssueDate:  "2026-03-15",
		Adjustment: false,
		Details: []ManualJournalDetail{
			{EntrySide: "debit", Amount: 30000, Description: "first line"},
			{EntrySide: "debit", Amount: 20000, Description: ""},
			{EntrySide: "credit", Amount: 50000, Description: ""},
		},
	}
	row := mj.ToRow()
	if row.ID != 1 {
		t.Errorf("ID = %d, want 1", row.ID)
	}
	if row.Date != "2026-03-15" {
		t.Errorf("Date = %q, want 2026-03-15", row.Date)
	}
	if row.Adjustment != "No" {
		t.Errorf("Adjustment = %q, want No", row.Adjustment)
	}
	if row.Entries != "2 debit / 1 credit" {
		t.Errorf("Entries = %q, want '2 debit / 1 credit'", row.Entries)
	}
	if row.Amount != 50000 {
		t.Errorf("Amount = %d, want 50000 (debit sum)", row.Amount)
	}
	if row.Description != "first line" {
		t.Errorf("Description = %q, want 'first line'", row.Description)
	}
}
