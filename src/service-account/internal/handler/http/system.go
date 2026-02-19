package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	status := h.service.CheckHealth()

	response := map[string]string{
		"service": "account",
		"status":  status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
