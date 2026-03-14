package format

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestOutputJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "test"}
	if err := outputJSON(&buf, data); err != nil {
		t.Fatalf("json output: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected 'test', got %q", result["name"])
	}
}

func TestOutputText(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"Name", "Type"}
	rows := [][]string{{"checking", "bank"}, {"wallet", "cash"}}

	if err := outputText(&buf, headers, rows); err != nil {
		t.Fatalf("text output: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Name: checking") {
		t.Errorf("expected 'Name: checking' in output: %s", out)
	}
	if !strings.Contains(out, "Type: cash") {
		t.Errorf("expected 'Type: cash' in output: %s", out)
	}
}

func TestOutputTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"ID", "Name"}
	rows := [][]string{{"1", "alice"}, {"2", "bob"}}

	if err := outputTable(&buf, headers, rows); err != nil {
		t.Fatalf("table output: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 { // header + separator + 2 rows
		t.Errorf("expected 4 lines, got %d: %s", len(lines), out)
	}
	if !strings.Contains(lines[0], "ID") || !strings.Contains(lines[0], "Name") {
		t.Errorf("expected headers in first line: %s", lines[0])
	}
}

func TestOutputUnknownFormat(t *testing.T) {
	err := Output("xml", nil, nil, nil)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
