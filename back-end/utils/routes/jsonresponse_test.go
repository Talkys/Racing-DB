package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/utils/apperror"
)

func TestWriteJSONResponseEncodes(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSONResponse(rec, http.StatusCreated, map[string]int{"a": 1})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	if body := strings.TrimSpace(rec.Body.String()); body != `{"a":1}` {
		t.Errorf("body = %q", body)
	}
}

func TestWriteJSONResponseUnencodable(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSONResponse(rec, http.StatusOK, map[string]any{"ch": make(chan int)})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not JSON: %q", rec.Body.String())
	}
	if body.Status != http.StatusInternalServerError || body.Error != "internal server error" {
		t.Errorf("body = %+v", body)
	}
	if strings.Contains(rec.Body.String(), "chan") {
		t.Errorf("body leaks the encoding error: %q", rec.Body.String())
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, apperror.BadRequest("bad %s", "thing"))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "bad thing") {
		t.Errorf("client error: %d %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	WriteError(rec, errors.New("pq: secret details"))
	if rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), "secret") {
		t.Errorf("internal error: %d %q", rec.Code, rec.Body.String())
	}
}
