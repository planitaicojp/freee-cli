package cmdutil

import (
	"fmt"
	"testing"
)

func TestFiscalYearRange(t *testing.T) {
	tests := []struct {
		year         int
		closingMonth int
		wantFrom     string
		wantTo       string
	}{
		{2025, 12, "2025-01-01", "2025-12-31"},
		{2025, 3, "2024-04-01", "2025-03-31"},
		{2025, 9, "2024-10-01", "2025-09-30"},
		{2026, 1, "2025-02-01", "2026-01-31"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("year=%d,month=%d", tt.year, tt.closingMonth), func(t *testing.T) {
			from, to := FiscalYearRange(tt.year, tt.closingMonth)
			if from != tt.wantFrom {
				t.Errorf("from = %q, want %q", from, tt.wantFrom)
			}
			if to != tt.wantTo {
				t.Errorf("to = %q, want %q", to, tt.wantTo)
			}
		})
	}
}

func TestClosingMonthFromEndDate(t *testing.T) {
	tests := []struct {
		endDate string
		want    int
	}{
		{"2025-03-31", 3},
		{"2025-12-31", 12},
		{"2026-09-30", 9},
	}
	for _, tt := range tests {
		t.Run(tt.endDate, func(t *testing.T) {
			got, err := closingMonthFromEndDate(tt.endDate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
