package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func respondWithError(w http.ResponseWriter, statusCode int, code, message string, internalErr error) {
	if internalErr != nil {
		log.Printf("[ERROR] %s: %v", code, internalErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}
