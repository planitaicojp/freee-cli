package resolve

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/model"
)

// PartnerID resolves a partner ID from --partner-id or --partner-name flags.
// Returns (0, nil) if neither flag is set.
func PartnerID(cmd *cobra.Command, freeeAPI *api.FreeeAPI, companyID int64) (int64, error) {
	idChanged := cmd.Flags().Changed("partner-id")
	nameChanged := cmd.Flags().Changed("partner-name")

	if idChanged && nameChanged {
		return 0, &cerrors.ValidationError{
			Message: "--partner-id and --partner-name are mutually exclusive\nhint: use one or the other",
		}
	}

	if idChanged {
		id, _ := cmd.Flags().GetInt64("partner-id")
		return id, nil
	}

	if !nameChanged {
		return 0, nil
	}

	name, _ := cmd.Flags().GetString("partner-name")

	// Fetch all partners (paginated)
	var allPartners []model.Partner
	limit := 100
	for offset := 0; ; offset += limit {
		params := fmt.Sprintf("limit=%d&offset=%d", limit, offset)
		var resp model.PartnersResponse
		if err := freeeAPI.ListPartners(companyID, params, &resp); err != nil {
			return 0, err
		}
		allPartners = append(allPartners, resp.Partners...)
		if len(resp.Partners) < limit {
			break
		}
	}

	return matchByName(name, allPartners, "partner")
}

// matchByName finds a partner by name. Exact match first (case-insensitive),
// then partial match (contains, case-insensitive) as fallback.
func matchByName(name string, partners []model.Partner, resource string) (int64, error) {
	nameLower := strings.ToLower(name)

	// Exact match (case-insensitive)
	var exactMatches []model.Partner
	for _, p := range partners {
		if strings.EqualFold(p.Name, name) {
			exactMatches = append(exactMatches, p)
		}
	}
	if len(exactMatches) == 1 {
		return exactMatches[0].ID, nil
	}
	if len(exactMatches) > 1 {
		return 0, multipleMatchError(resource, name, exactMatches, true)
	}

	// Partial match (contains, case-insensitive)
	var partialMatches []model.Partner
	for _, p := range partners {
		if strings.Contains(strings.ToLower(p.Name), nameLower) {
			partialMatches = append(partialMatches, p)
		}
	}
	if len(partialMatches) == 1 {
		return partialMatches[0].ID, nil
	}
	if len(partialMatches) > 1 {
		return 0, multipleMatchError(resource, name, partialMatches, false)
	}

	return 0, &cerrors.NotFoundError{
		Resource: resource,
		ID:       fmt.Sprintf("no %s found matching %q\nhint: run 'freee %s list' to see available %ss", resource, name, resource, resource),
	}
}

// AccountItemID resolves an account item ID from --account-item-id or --account-name flags.
// Also accepts --account-item-name as an alias for --account-name.
// Returns (0, nil) if no flag is set.
func AccountItemID(cmd *cobra.Command, freeeAPI *api.FreeeAPI, companyID int64) (int64, error) {
	idChanged := cmd.Flags().Changed("account-item-id")
	nameChanged := cmd.Flags().Changed("account-name")
	// Check alias too
	if !nameChanged && cmd.Flags().Lookup("account-item-name") != nil {
		nameChanged = cmd.Flags().Changed("account-item-name")
	}

	if idChanged && nameChanged {
		return 0, &cerrors.ValidationError{
			Message: "--account-item-id and --account-name are mutually exclusive\nhint: use one or the other",
		}
	}

	if idChanged {
		id, _ := cmd.Flags().GetInt64("account-item-id")
		return id, nil
	}

	if !nameChanged {
		return 0, nil
	}

	name, _ := cmd.Flags().GetString("account-name")
	if name == "" {
		// Try alias
		name, _ = cmd.Flags().GetString("account-item-name")
	}

	// Fetch all account items (no pagination needed — single call)
	var resp model.AccountItemsResponse
	if err := freeeAPI.ListAccountItems(companyID, &resp); err != nil {
		return 0, err
	}

	return matchAccountItemByName(name, resp.AccountItems, "account-item")
}

// matchAccountItemByName finds an account item by name.
func matchAccountItemByName(name string, items []model.AccountItem, resource string) (int64, error) {
	nameLower := strings.ToLower(name)

	// Exact match
	var exactMatches []model.AccountItem
	for _, item := range items {
		if strings.EqualFold(item.Name, name) {
			exactMatches = append(exactMatches, item)
		}
	}
	if len(exactMatches) == 1 {
		return exactMatches[0].ID, nil
	}
	if len(exactMatches) > 1 {
		return 0, multipleAccountMatchError(resource, name, exactMatches, true)
	}

	// Partial match
	var partialMatches []model.AccountItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), nameLower) {
			partialMatches = append(partialMatches, item)
		}
	}
	if len(partialMatches) == 1 {
		return partialMatches[0].ID, nil
	}
	if len(partialMatches) > 1 {
		return 0, multipleAccountMatchError(resource, name, partialMatches, false)
	}

	return 0, &cerrors.NotFoundError{
		Resource: resource,
		ID:       fmt.Sprintf("no %s found matching %q\nhint: run 'freee account list' to see available account items", resource, name),
	}
}

func multipleAccountMatchError(resource, name string, matches []model.AccountItem, exact bool) error {
	matchType := "partially"
	if exact {
		matchType = "exactly"
	}
	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("  - %s (id: %d)", m.Name, m.ID))
	}
	hint := "hint: use --account-item-id to specify, or use the full name"
	if exact {
		hint = "hint: use --account-item-id to specify"
	}
	msg := fmt.Sprintf("multiple %ss %s match %q:\n%s\n%s", resource, matchType, name, strings.Join(lines, "\n"), hint)
	return &cerrors.ValidationError{Message: msg}
}

func multipleMatchError(resource, name string, matches []model.Partner, exact bool) error {
	matchType := "partially"
	if exact {
		matchType = "exactly"
	}
	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("  - %s (id: %d)", m.Name, m.ID))
	}
	hint := "hint: use --partner-id to specify, or use the full name"
	if exact {
		hint = "hint: use --partner-id to specify"
	}
	msg := fmt.Sprintf("multiple %ss %s match %q:\n%s\n%s", resource, matchType, name, strings.Join(lines, "\n"), hint)
	return &cerrors.ValidationError{Message: msg}
}
