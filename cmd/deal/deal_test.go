package deal

import (
	"net/url"
	"testing"

	"github.com/spf13/cobra"
)

func TestDealListParamsEncodePartnerCode(t *testing.T) {
	const partner = "東京 & Co.+/#"
	cmd := &cobra.Command{}
	cmd.Flags().String("partner", partner, "")
	cmd.Flags().Int("limit", 20, "")
	cmd.Flags().Int("offset", 40, "")

	for _, tc := range []struct {
		name   string
		query  string
		limit  string
		offset string
	}{
		{"all pages", buildBaseListParams(cmd), "", ""},
		{"single page", buildListParams(cmd), "20", "40"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values, err := url.ParseQuery(tc.query)
			if err != nil {
				t.Fatal(err)
			}
			if got := values.Get("partner_code"); got != partner {
				t.Errorf("partner_code = %q, want %q", got, partner)
			}
			if got := values.Get("limit"); got != tc.limit {
				t.Errorf("limit = %q, want %q", got, tc.limit)
			}
			if got := values.Get("offset"); got != tc.offset {
				t.Errorf("offset = %q, want %q", got, tc.offset)
			}
		})
	}
}
