package routes

import (
	"backend/modules/betterlogger"
	"backend/utils/apperror"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	log := betterlogger.GetInstance()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log this error in production
		log.Error("%s", "Error processing response: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteError answers client errors (apperror.Error) with their own status and message.
// Anything else is logged and answered with a generic 500, so database details never reach the client.
func WriteError(w http.ResponseWriter, err error) {
	log := betterlogger.GetInstance()

	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		log.Warn("%s", appErr.Message)
		WriteJSONResponse(w, appErr.Status, ErrorResponse{Status: appErr.Status, Error: appErr.Message})
		return
	}

	log.Error("%s", err.Error())
	WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Status: http.StatusInternalServerError, Error: "internal server error"})
}

// NotFound answers any route that is not registered.
func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, apperror.NotFound("resource not found: %s %s", r.Method, r.URL.Path))
}
