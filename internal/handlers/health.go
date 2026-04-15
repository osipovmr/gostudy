package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthCheck(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	healthResponse := map[string]string{"status": "ok"}
	if err := json.NewEncoder(responseWriter).Encode(healthResponse); err != nil {
		http.Error(responseWriter, "internal server error", http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}
