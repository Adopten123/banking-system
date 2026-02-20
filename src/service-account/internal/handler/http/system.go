package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Берем контекст из запроса и отдаем сервису
	status := h.service.CheckHealth(r.Context())

	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
	})
}
