package routes

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serve(t *testing.T, h http.HandlerFunc, req *http.Request) (*httptest.ResponseRecorder, []map[string]any) {
	t.Helper()
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))
	rec := httptest.NewRecorder()
	Middleware(log, h).ServeHTTP(rec, req)

	var entries []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var e map[string]any
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("log line is not JSON: %q", line)
		}
		entries = append(entries, e)
	}
	return rec, entries
}

func TestMiddlewareAccessLog(t *testing.T) {
	var seenID string
	h := func(w http.ResponseWriter, r *http.Request) {
		seenID = RequestID(r.Context())
		w.WriteHeader(http.StatusTeapot)
	}
	rec, entries := serve(t, h, httptest.NewRequest("GET", "/api/1/races?page=2", nil))

	id := rec.Header().Get(RequestIDHeader)
	if id == "" || id != seenID {
		t.Errorf("response id %q, handler saw %q", id, seenID)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d log lines, want 1: %v", len(entries), entries)
	}
	e := entries[0]
	if e["method"] != "GET" || e["path"] != "/api/1/races" || e["status"] != float64(http.StatusTeapot) || e["id"] != id {
		t.Errorf("unexpected access log: %v", e)
	}
}

func TestMiddlewareImplicitOK(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("hi")) }
	_, entries := serve(t, h, httptest.NewRequest("GET", "/", nil))
	if entries[0]["status"] != float64(http.StatusOK) {
		t.Errorf("status = %v, want 200", entries[0]["status"])
	}
}

func TestMiddlewareRequestIDHeader(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(RequestIDHeader, "abc-123")
	rec, _ := serve(t, h, req)
	if got := rec.Header().Get(RequestIDHeader); got != "abc-123" {
		t.Errorf("valid incoming id not reused: %q", got)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set(RequestIDHeader, "bad id\nwith newline")
	rec, _ = serve(t, h, req)
	if got := rec.Header().Get(RequestIDHeader); got == "" || strings.Contains(got, " ") {
		t.Errorf("invalid incoming id not replaced: %q", got)
	}
}

func TestMiddlewareRecoversPanic(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) { panic("boom") }
	rec, entries := serve(t, h, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Errorf("body leaks the panic: %q", rec.Body.String())
	}
	if len(entries) != 2 || entries[0]["msg"] != "panic" || entries[1]["status"] != float64(http.StatusInternalServerError) {
		t.Errorf("want a panic line then an access line with 500, got %v", entries)
	}
}
