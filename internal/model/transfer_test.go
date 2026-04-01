package model

import (
	"encoding/json"
	"testing"
)

func TestTransferUnmarshal(t *testing.T) {
	raw := `{
		"transfers": [{
			"id": 1,
			"company_id": 100,
			"amount": 5000,
			"date": "2026-03-15",
			"from_walletable_type": "credit_card",
			"from_walletable_id": 101,
			"to_walletable_type": "bank_account",
			"to_walletable_id": 201,
			"description": "test transfer"
		}]
	}`
	var resp TransfersResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(resp.Transfers) != 1 {
		t.Fatalf("expected 1 transfer, got %d", len(resp.Transfers))
	}
	tr := resp.Transfers[0]
	if tr.ID != 1 {
		t.Errorf("ID = %d, want 1", tr.ID)
	}
	if tr.Date != "2026-03-15" {
		t.Errorf("Date = %q, want 2026-03-15", tr.Date)
	}
	if tr.Amount != 5000 {
		t.Errorf("Amount = %d, want 5000", tr.Amount)
	}
	if tr.FromWalletableType != "credit_card" {
		t.Errorf("FromWalletableType = %q, want credit_card", tr.FromWalletableType)
	}
	if tr.FromWalletableID != 101 {
		t.Errorf("FromWalletableID = %d, want 101", tr.FromWalletableID)
	}
	if tr.ToWalletableType != "bank_account" {
		t.Errorf("ToWalletableType = %q, want bank_account", tr.ToWalletableType)
	}
	if tr.ToWalletableID != 201 {
		t.Errorf("ToWalletableID = %d, want 201", tr.ToWalletableID)
	}
	if tr.Description != "test transfer" {
		t.Errorf("Description = %q, want test transfer", tr.Description)
	}
}

func TestTransferSingleUnmarshal(t *testing.T) {
	raw := `{
		"transfer": {
			"id": 1,
			"company_id": 100,
			"amount": 5000,
			"date": "2026-03-15",
			"from_walletable_type": "credit_card",
			"from_walletable_id": 101,
			"to_walletable_type": "bank_account",
			"to_walletable_id": 201,
			"description": "single"
		}
	}`
	var resp TransferResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Transfer.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.Transfer.ID)
	}
}

func TestTransferToRow(t *testing.T) {
	tr := Transfer{
		ID:                 1,
		CompanyID:          100,
		Amount:             5000,
		Date:               "2026-03-15",
		FromWalletableType: "credit_card",
		FromWalletableID:   101,
		ToWalletableType:   "bank_account",
		ToWalletableID:     201,
		Description:        "memo",
	}
	row := tr.ToRow()
	if row.ID != 1 {
		t.Errorf("Row.ID = %d, want 1", row.ID)
	}
	if row.From != "credit_card:101" {
		t.Errorf("Row.From = %q, want credit_card:101", row.From)
	}
	if row.To != "bank_account:201" {
		t.Errorf("Row.To = %q, want bank_account:201", row.To)
	}
	if row.Amount != 5000 {
		t.Errorf("Row.Amount = %d, want 5000", row.Amount)
	}
	if row.Description != "memo" {
		t.Errorf("Row.Description = %q, want memo", row.Description)
	}
}
