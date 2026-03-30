package resolve

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/model"
)

// Named is implemented by types that have an ID and a Name.
type Named interface {
	GetID() int64
	GetName() string
}

// resolveNotFoundError is a not-found error with exit code 3 and clean message.
type resolveNotFoundError struct {
	message string
}

func (e *resolveNotFoundError) Error() string { return e.message }
func (e *resolveNotFoundError) ExitCode() int { return cerrors.ExitNotFound }

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

	return matchByName(name, allPartners, "partner", "--partner-id")
}

// AccountItemID resolves an account item ID from --account-item-id or --account-name flags.
// Also accepts --account-item-name as an alias for --account-name.
// Returns (0, nil) if no flag is set.
func AccountItemID(cmd *cobra.Command, freeeAPI *api.FreeeAPI, companyID int64) (int64, error) {
	idChanged := cmd.Flags().Changed("account-item-id")
	nameChanged := cmd.Flags().Changed("account-name")
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
		name, _ = cmd.Flags().GetString("account-item-name")
	}

	var resp model.AccountItemsResponse
	if err := freeeAPI.ListAccountItems(companyID, &resp); err != nil {
		return 0, err
	}

	return matchByName(name, resp.AccountItems, "account-item", "--account-item-id")
}

// matchByName finds an item by name. Exact match first (case-insensitive),
// then partial match (contains, case-insensitive) as fallback.
func matchByName[T Named](name string, items []T, resource, idFlag string) (int64, error) {
	nameLower := strings.ToLower(name)

	// Exact match (case-insensitive)
	var exactMatches []T
	for _, item := range items {
		if strings.EqualFold(item.GetName(), name) {
			exactMatches = append(exactMatches, item)
		}
	}
	if len(exactMatches) == 1 {
		return exactMatches[0].GetID(), nil
	}
	if len(exactMatches) > 1 {
		return 0, multipleMatchError(resource, name, exactMatches, true, idFlag)
	}

	// Partial match (contains, case-insensitive)
	var partialMatches []T
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.GetName()), nameLower) {
			partialMatches = append(partialMatches, item)
		}
	}
	if len(partialMatches) == 1 {
		return partialMatches[0].GetID(), nil
	}
	if len(partialMatches) > 1 {
		return 0, multipleMatchError(resource, name, partialMatches, false, idFlag)
	}

	listCmd := resource
	if resource == "account-item" {
		listCmd = "account"
	}
	return 0, &resolveNotFoundError{
		message: fmt.Sprintf("no %s found matching %q\nhint: run 'freee %s list' to see available %ss", resource, name, listCmd, resource),
	}
}

func multipleMatchError[T Named](resource, name string, matches []T, exact bool, idFlag string) error {
	matchType := "partially"
	if exact {
		matchType = "exactly"
	}
	var lines []string
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf("  - %s (id: %d)", m.GetName(), m.GetID()))
	}
	hint := fmt.Sprintf("hint: use %s to specify, or use the full name", idFlag)
	if exact {
		hint = fmt.Sprintf("hint: use %s to specify", idFlag)
	}
	msg := fmt.Sprintf("multiple %ss %s match %q:\n%s\n%s", resource, matchType, name, strings.Join(lines, "\n"), hint)
	return &cerrors.ValidationError{Message: msg}
}
