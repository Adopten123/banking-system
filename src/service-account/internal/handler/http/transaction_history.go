package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary История транзакций (выписка)
// @Description Возвращает список всех операций по указанному счету, отсортированный от новых к старым
// @Tags accounts
// @Produce json
// @Param id path string true "Public ID счета (UUID)"
// @Success 200 {array} domain.TransactionHistory "Список транзакций"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id}/transactions [get]
func (h *Handler) getTransactions(w http.ResponseWriter, r *http.Request) {
	accountIDParam := chi.URLParam(r, "id")
	publicID, err := uuid.Parse(accountIDParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	history, err := h.service.GetAccountTransactions(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch transactions", err)
		return
	}

	if history == nil {
		history = []domain.TransactionHistory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
