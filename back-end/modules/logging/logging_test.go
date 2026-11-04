package logging

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewJSONRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	log, err := New(&buf, "json", "warn")
	if err != nil {
		t.Fatal(err)
	}

	log.Info("hidden")
	log.Warn("shown", "key", "value")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1: %q", len(lines), buf.String())
	}
	var entry map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if entry["msg"] != "shown" || entry["key"] != "value" {
		t.Errorf("unexpected entry: %v", entry)
	}
}

func TestNewDefaultsToTextInfo(t *testing.T) {
	var buf bytes.Buffer
	log, err := New(&buf, "", "")
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("hidden")
	log.Info("shown")

	if out := buf.String(); strings.Contains(out, "hidden") || !strings.Contains(out, "msg=shown") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestNewRejectsBadConfig(t *testing.T) {
	if _, err := New(&bytes.Buffer{}, "xml", ""); err == nil {
		t.Error("want error for LOG_FORMAT=xml")
	}
	if _, err := New(&bytes.Buffer{}, "", "loud"); err == nil {
		t.Error("want error for LOG_LEVEL=loud")
	}
}
