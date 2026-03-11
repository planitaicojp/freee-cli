package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type testRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNew(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{"json", "*output.JSONFormatter"},
		{"yaml", "*output.YAMLFormatter"},
		{"csv", "*output.CSVFormatter"},
		{"table", "*output.TableFormatter"},
		{"", "*output.TableFormatter"},
		{"unknown", "*output.TableFormatter"},
	}
	for _, tt := range tests {
		name := tt.format
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			f := New(tt.format)
			if f == nil {
				t.Fatal("New returned nil")
			}
			got := fmt.Sprintf("%T", f)
			if got != tt.want {
				t.Errorf("New(%q) = %s, want %s", tt.format, got, tt.want)
			}
		})
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	data := []testRow{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}

	if err := f.Format(&buf, data); err != nil {
		t.Fatalf("Format error: %v", err)
	}

	var result []testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if len(result) != 2 || result[0].Name != "Alice" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestJSONFormatter_SingleObject(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	data := map[string]string{"key": "value"}

	if err := f.Format(&buf, data); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), `"key"`) {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestYAMLFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &YAMLFormatter{}
	data := []testRow{{ID: 1, Name: "Test"}}

	if err := f.Format(&buf, data); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "id: 1") || !strings.Contains(out, "name: Test") {
		t.Errorf("unexpected YAML output: %s", out)
	}
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &CSVFormatter{}
	data := []testRow{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}

	if err := f.Format(&buf, data); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if lines[0] != "id,name" {
		t.Errorf("unexpected header: %q", lines[0])
	}
	if lines[1] != "1,Alice" {
		t.Errorf("unexpected row 1: %q", lines[1])
	}
}

func TestCSVFormatter_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	f := &CSVFormatter{}
	if err := f.Format(&buf, []testRow{}); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty slice, got %q", buf.String())
	}
}

func TestCSVFormatter_NonSlice(t *testing.T) {
	var buf bytes.Buffer
	f := &CSVFormatter{}
	err := f.Format(&buf, "not a slice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	data := []testRow{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}

	if err := f.Format(&buf, data); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "NAME") {
		t.Errorf("expected uppercase headers, got: %s", out)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("expected data rows, got: %s", out)
	}
}

func TestTableFormatter_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	if err := f.Format(&buf, []testRow{}); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestTableFormatter_NonSlice(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	if err := f.Format(&buf, "hello"); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected string output, got %q", buf.String())
	}
}
