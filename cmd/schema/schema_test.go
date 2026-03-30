package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newTestRoot builds a test cobra tree and returns (root, schemaCmd).
// The schema command is registered under root so persistent flags are inherited.
func newTestRoot() (*cobra.Command, *cobra.Command) {
	root := &cobra.Command{Use: "freee", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().String("format", "", "output format")
	root.PersistentFlags().String("profile", "", "profile")
	root.PersistentFlags().Bool("verbose", false, "verbose")

	deal := &cobra.Command{Use: "deal", Short: "Manage deals (取引)"}
	listCmd := &cobra.Command{Use: "list", Short: "List deals"}
	listCmd.Flags().String("type", "", "filter by type: income or expense")
	listCmd.Flags().Int("limit", 50, "max results per page")
	listCmd.Flags().Bool("all", false, "fetch all pages")
	createCmd := &cobra.Command{Use: "create", Short: "Create a new deal"}
	createCmd.Flags().String("type", "", "deal type: income or expense (required)")
	_ = createCmd.MarkFlagRequired("type")
	createCmd.Flags().String("date", "", "issue date YYYY-MM-DD (required)")
	_ = createCmd.MarkFlagRequired("date")
	createCmd.Flags().Int64("partner-id", 0, "partner ID")
	createCmd.Flags().Int64("amount", 0, "amount (required)")
	_ = createCmd.MarkFlagRequired("amount")
	deal.AddCommand(listCmd, createCmd)

	version := &cobra.Command{Use: "version", Short: "Print version"}
	hidden := &cobra.Command{Use: "completion", Short: "Generate completion", Hidden: true}

	partner := &cobra.Command{Use: "partner", Short: "Manage partners (取引先)"}
	partner.AddCommand(&cobra.Command{Use: "list", Short: "List partners"})

	root.AddCommand(deal, version, hidden, partner)

	schemaCmd := NewCmd(root)
	root.AddCommand(schemaCmd)

	return root, schemaCmd
}

// executeSchema runs the schema command via root to inherit persistent flags.
func executeSchema(args ...string) (*bytes.Buffer, error) {
	root, _ := newTestRoot()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(append([]string{"schema"}, args...))
	err := root.Execute()
	return buf, err
}

func TestFindCommand(t *testing.T) {
	root, _ := newTestRoot()

	tests := []struct {
		name    string
		args    []string
		wantUse string
		wantErr bool
	}{
		{"root", nil, "freee", false},
		{"resource", []string{"deal"}, "deal", false},
		{"action", []string{"deal", "create"}, "create", false},
		{"leaf", []string{"version"}, "version", false},
		{"not found resource", []string{"xyz"}, "", true},
		{"not found action", []string{"deal", "xyz"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := findCommand(root, tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cmd.Use != tt.wantUse {
				t.Errorf("got Use=%q, want %q", cmd.Use, tt.wantUse)
			}
		})
	}
}

func TestListSubcommands(t *testing.T) {
	root, _ := newTestRoot()

	// Root level: should exclude hidden commands and auto-generated help/schema
	items := listSubcommands(root)
	names := make(map[string]bool)
	for _, item := range items {
		names[item.Name] = true
	}
	if names["completion"] {
		t.Error("hidden command 'completion' should be excluded")
	}
	if names["help"] {
		t.Error("auto-generated 'help' should be excluded")
	}
	if names["schema"] {
		t.Error("'schema' itself should be excluded from listing")
	}
	if !names["deal"] {
		t.Error("expected 'deal' in list")
	}
	if !names["version"] {
		t.Error("expected 'version' in list")
	}

	// Deal subcommands
	dealCmd, _ := findCommand(root, []string{"deal"})
	dealSubs := listSubcommands(dealCmd)
	if len(dealSubs) != 2 {
		t.Errorf("expected 2 deal subcommands, got %d", len(dealSubs))
	}
}

func TestExtractFlags(t *testing.T) {
	root, _ := newTestRoot()
	createCmd, _ := findCommand(root, []string{"deal", "create"})

	t.Run("without global", func(t *testing.T) {
		flags := extractFlags(createCmd, false, root)
		names := make(map[string]bool)
		for _, f := range flags {
			names[f.Name] = true
		}
		if !names["type"] {
			t.Error("expected 'type' flag")
		}
		if !names["amount"] {
			t.Error("expected 'amount' flag")
		}
		if names["format"] {
			t.Error("global flag 'format' should be excluded")
		}
		if names["profile"] {
			t.Error("global flag 'profile' should be excluded")
		}
	})

	t.Run("with global", func(t *testing.T) {
		flags := extractFlags(createCmd, true, root)
		names := make(map[string]bool)
		for _, f := range flags {
			names[f.Name] = true
		}
		if !names["format"] {
			t.Error("expected global flag 'format' with --show-global")
		}
	})

	t.Run("required detection", func(t *testing.T) {
		flags := extractFlags(createCmd, false, root)
		for _, f := range flags {
			switch f.Name {
			case "type":
				if !f.Required {
					t.Error("'type' should be required")
				}
			case "partner-id":
				if f.Required {
					t.Error("'partner-id' should not be required")
				}
			}
		}
	})
}

func TestSchemaOutput_JSON_Flags(t *testing.T) {
	buf, err := executeSchema("deal", "create", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var schema CommandSchema
	if err := json.Unmarshal(buf.Bytes(), &schema); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if schema.Command != "freee deal create" {
		t.Errorf("got command=%q, want %q", schema.Command, "freee deal create")
	}
	if len(schema.Flags) == 0 {
		t.Error("expected flags in output")
	}
}

func TestSchemaOutput_JSON_CommandList(t *testing.T) {
	buf, err := executeSchema("deal", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var list CommandList
	if err := json.Unmarshal(buf.Bytes(), &list); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(list.Commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(list.Commands))
	}
}

func TestSchemaOutput_Table(t *testing.T) {
	buf, err := executeSchema("deal", "create")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Command:") {
		t.Error("expected 'Command:' header in table output")
	}
	if !strings.Contains(out, "freee deal create") {
		t.Error("expected command path in output")
	}
	if !strings.Contains(out, "NAME") {
		t.Error("expected column headers in output")
	}
}

func TestLeafCommand(t *testing.T) {
	buf, err := executeSchema("version", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var schema CommandSchema
	if err := json.Unmarshal(buf.Bytes(), &schema); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if schema.Command != "freee version" {
		t.Errorf("got command=%q, want %q", schema.Command, "freee version")
	}
}

func TestSchemaNotFound(t *testing.T) {
	_, err := executeSchema("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}
