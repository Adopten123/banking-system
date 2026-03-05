package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Получить баланс счета
// @Description  Возвращает текущий баланс счета по его публичному UUID
// @Tags         accounts
// @Produce      json
// @Param        id   path      string  true  "ID счета (UUID)" Format(uuid)
// @Success      200  {object}  domain.AccountBalanceResponse "Успешное получение баланса"
// @Failure      400  {object}  ErrorResponse "Неверный формат ID (INVALID_REQUEST)"
// @Failure      404  {object}  ErrorResponse "Счет не найден (NOT_FOUND)"
// @Failure      500  {object}  ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR)"
// @Router       /api/accounts/{id}/balance [get]
func (h *Handler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountPathID := chi.URLParam(r, "id")

	publicID, err := uuid.Parse(accountPathID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	balance, err := h.service.GetAccountBalance(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", err)
		return
	}

	resp := domain.AccountBalanceResponse{
		AccountID: publicID,
		Balance:   balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
