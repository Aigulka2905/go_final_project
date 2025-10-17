package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError sends a JSON error response with the specified message and status code
func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
