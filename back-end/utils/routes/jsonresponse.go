package routes

import (
	"backend/modules/betterlogger"
	"encoding/json"
	"net/http"
)

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
