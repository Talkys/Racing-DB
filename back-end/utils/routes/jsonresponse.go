package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log this error in production
		log.Println("Error processing response:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
