package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Получение информации о счете
// @Description Возвращает текущий баланс, валюту и статус счета по его публичному ID
// @Tags accounts
// @Produce json
// @Param id path string true "Public ID счета (UUID)"
// @Success 200 {object} domain.Account "Информация о счете"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id} [get]
func (h *Handler) getAccountInfo(w http.ResponseWriter, r *http.Request) {
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

	resp := domain.AccountInfoResponse{
		ID:           acc.PublicID,
		UserID:       acc.UserID,
		TypeID:       acc.TypeID,
		StatusID:     acc.StatusID,
		CurrencyCode: acc.CurrencyCode,
		Name:         acc.Name,
		CreatedAt:    acc.CreatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
