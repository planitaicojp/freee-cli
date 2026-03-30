package cmdutil

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/model"
)

// FiscalYearRange calculates the date range for a fiscal year given the closing month.
func FiscalYearRange(year int, closingMonth int) (from string, to string) {
	// Go's time.Date normalizes out-of-range values: month 13 day 0 → Dec 31.
	endDate := time.Date(year, time.Month(closingMonth+1), 0, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(-1, 0, 1)
	return startDate.Format("2006-01-02"), endDate.Format("2006-01-02")
}

// ResolveFiscalYear resolves --fiscal-year flag to from/to date strings.
// Returns ("", "", nil) if --fiscal-year is not set.
func ResolveFiscalYear(cmd *cobra.Command, freeeAPI *api.FreeeAPI, companyID int64) (string, string, error) {
	year, _ := cmd.Flags().GetInt("fiscal-year")
	if year == 0 {
		return "", "", nil
	}

	fromSet := cmd.Flags().Changed("from")
	toSet := cmd.Flags().Changed("to")
	if fromSet || toSet {
		return "", "", &cerrors.ValidationError{
			Message: "--fiscal-year and --from/--to are mutually exclusive\nhint: use --fiscal-year for full year, or --from/--to for custom range",
		}
	}

	var resp model.CompanyResponse
	if err := freeeAPI.GetCompany(companyID, &resp); err != nil {
		return "", "", err
	}

	closingMonth, err := findClosingMonth(resp.Company.FiscalYears, year)
	if err != nil {
		return "", "", err
	}

	from, to := FiscalYearRange(year, closingMonth)
	return from, to, nil
}

func findClosingMonth(fiscalYears []model.FiscalYear, year int) (int, error) {
	if len(fiscalYears) == 0 {
		return 0, &cerrors.ValidationError{
			Message: "no fiscal year data found for this company\nhint: run 'freee company show' to check company settings",
		}
	}

	for _, fy := range fiscalYears {
		month, err := closingMonthFromEndDate(fy.EndDate)
		if err != nil {
			continue
		}
		endYear, err := yearFromDate(fy.EndDate)
		if err != nil {
			continue
		}
		if endYear == year {
			return month, nil
		}
	}

	// Fallback: find the most recent fiscal year by EndDate (API order not guaranteed)
	var mostRecentFY model.FiscalYear
	var mostRecentEnd time.Time
	found := false
	for _, fy := range fiscalYears {
		end, err := time.Parse("2006-01-02", fy.EndDate)
		if err != nil {
			continue
		}
		if !found || end.After(mostRecentEnd) {
			mostRecentEnd = end
			mostRecentFY = fy
			found = true
		}
	}
	if !found {
		return 0, fmt.Errorf("no valid fiscal years found to determine closing month")
	}
	return closingMonthFromEndDate(mostRecentFY.EndDate)
}

func closingMonthFromEndDate(endDate string) (int, error) {
	t, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, err
	}
	return int(t.Month()), nil
}

func yearFromDate(date string) (int, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, err
	}
	return t.Year(), nil
}
