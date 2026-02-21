package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
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
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Parse UserID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid user_id format", err)
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
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create account", err)
		return
	}

	// Response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	publicID, err := uuid.Parse(accountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	acc, err := h.service.GetAccount(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc)
}
