package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	UserID       string `json:"user_id"`
	TypeID       int32  `json:"type_id"`
	CurrencyCode string `json:"currency_code"`
	Name         string `json:"name"`
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	// Read JSON
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse UserID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	// Calling service
	acc, err := h.service.CreateAccount(
		r.Context(),
		userID,
		req.TypeID,
		req.CurrencyCode,
		req.Name,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "getting balance for account " + accountID,
		"status":  "not implemented yet",
	})
}
