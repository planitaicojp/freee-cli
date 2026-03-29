package resolve

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/internal/api"
	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

// newTestServer returns an httptest server that serves partner list responses.
func newTestServer(partners []map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"partners": partners}) //nolint:errcheck
	}))
}

// newTestCmd creates a cobra command with --partner-id and --partner-name flags.
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	cmd.Flags().Int64("partner-id", 0, "partner ID")
	cmd.Flags().String("partner-name", "", "partner name")
	return cmd
}

// newTestAPI creates a FreeeAPI pointing at the test server.
func newTestAPI(ts *httptest.Server) *api.FreeeAPI {
	client := api.NewClient("test-token", 1)
	client.HTTP = ts.Client()
	client.SetBaseURL(ts.URL)
	return &api.FreeeAPI{Client: client}
}

func TestPartnerID_ByID(t *testing.T) {
	ts := newTestServer(nil)
	defer ts.Close()

	cmd := newTestCmd()
	cmd.Flags().Set("partner-id", "123")

	id, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 123 {
		t.Errorf("got %d, want 123", id)
	}
}

func TestPartnerID_ByName_ExactMatch(t *testing.T) {
	partners := []map[string]any{
		{"id": float64(100), "name": "株式会社A"},
		{"id": float64(200), "name": "株式会社B"},
	}
	ts := newTestServer(partners)
	defer ts.Close()

	cmd := newTestCmd()
	cmd.Flags().Set("partner-name", "株式会社A")

	id, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 100 {
		t.Errorf("got %d, want 100", id)
	}
}

func TestPartnerID_ByName_PartialMatch(t *testing.T) {
	partners := []map[string]any{
		{"id": float64(100), "name": "株式会社Alpha"},
		{"id": float64(200), "name": "合同会社Beta"},
	}
	ts := newTestServer(partners)
	defer ts.Close()

	cmd := newTestCmd()
	cmd.Flags().Set("partner-name", "Alpha")

	id, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 100 {
		t.Errorf("got %d, want 100", id)
	}
}

func TestPartnerID_ByName_MultipleMatch(t *testing.T) {
	partners := []map[string]any{
		{"id": float64(100), "name": "株式会社A"},
		{"id": float64(200), "name": "株式会社AB"},
	}
	ts := newTestServer(partners)
	defer ts.Close()

	cmd2 := newTestCmd()
	cmd2.Flags().Set("partner-name", "株式")

	_, err := PartnerID(cmd2, newTestAPI(ts), 1)
	if err == nil {
		t.Fatal("expected error for multiple matches")
	}
	// Should be a ValidationError (exit code 4)
	if ec, ok := err.(cerrors.ExitCoder); !ok || ec.ExitCode() != 4 {
		t.Errorf("expected exit code 4, got %T: %v", err, err)
	}
}

func TestPartnerID_ByName_NotFound(t *testing.T) {
	partners := []map[string]any{
		{"id": float64(100), "name": "株式会社A"},
	}
	ts := newTestServer(partners)
	defer ts.Close()

	cmd := newTestCmd()
	cmd.Flags().Set("partner-name", "存在しない会社")

	_, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err == nil {
		t.Fatal("expected error for not found")
	}
	// Should be a NotFoundError (exit code 3)
	if ec, ok := err.(cerrors.ExitCoder); !ok || ec.ExitCode() != 3 {
		t.Errorf("expected exit code 3, got %T: %v", err, err)
	}
}

func TestPartnerID_BothFlags(t *testing.T) {
	ts := newTestServer(nil)
	defer ts.Close()

	cmd := newTestCmd()
	cmd.Flags().Set("partner-id", "123")
	cmd.Flags().Set("partner-name", "test")

	_, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
	if ec, ok := err.(cerrors.ExitCoder); !ok || ec.ExitCode() != 4 {
		t.Errorf("expected exit code 4, got %T: %v", err, err)
	}
}

func TestPartnerID_NeitherFlag(t *testing.T) {
	ts := newTestServer(nil)
	defer ts.Close()

	cmd := newTestCmd()

	id, err := PartnerID(cmd, newTestAPI(ts), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 0 {
		t.Errorf("got %d, want 0", id)
	}
}
