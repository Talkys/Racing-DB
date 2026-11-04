package routes

import (
	"backend/utils/apperror"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// internalErrorBody is written as-is when even a response cannot be encoded.
const internalErrorBody = `{"status":500,"error":"internal server error"}` + "\n"

// WriteJSONResponse encodes data before touching the response, so an encoding failure can still become a 500.
func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		slog.Error("encoding response", "status", status, "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(internalErrorBody))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

// WriteError answers client errors (apperror.Error) with their own status and message.
// Anything else is logged and answered with a generic 500, so database details never reach the client.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		slog.Warn("client error", "status", appErr.Status, "err", appErr.Message)
		WriteJSONResponse(w, appErr.Status, ErrorResponse{Status: appErr.Status, Error: appErr.Message})
		return
	}

	slog.Error("internal error", "err", err)
	WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Status: http.StatusInternalServerError, Error: "internal server error"})
}

// NotFound answers any route that is not registered.
func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, apperror.NotFound("resource not found: %s %s", r.Method, r.URL.Path))
}
