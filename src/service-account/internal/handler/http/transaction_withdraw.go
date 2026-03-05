package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// @Summary      Снятие наличных
// @Description  Списывает указанную сумму со счета. Учитывает кредитный лимит.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string          true  "ID счета (UUID)" Format(uuid)
// @Param        body body      domain.WithdrawRequest true  "Сумма для снятия"
// @Success      200  {object}  domain.WithdrawResponse "Успешное снятие (возвращает чек)"
// @Failure      400  {object}  ErrorResponse "Неверный запрос, нехватка средств или счет неактивен"
// @Failure      404  {object}  ErrorResponse "Счет не найден"
// @Failure      500  {object}  ErrorResponse "Внутренняя ошибка сервера"
// @Router       /api/accounts/{id}/withdraw [post]
func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "id")

	publicID, err := uuid.Parse(accountIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	// Decode JSON
	var req domain.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Parse and validate amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid amount format", err)
		return
	}

	if !amount.IsPositive() {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Withdraw amount must be strictly greater than zero", nil)
		return
	}

	// Calling service
	result, err := h.service.Withdraw(r.Context(), publicID, amount)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		if errors.Is(err, domain.ErrAccountInactive) {
			respondWithError(w, http.StatusBadRequest, "ACCOUNT_INACTIVE", "Account is blocked or inactive", err)
			return
		}
		if errors.Is(err, domain.ErrInsufficientFunds) {
			respondWithError(w, http.StatusBadRequest, "INSUFFICIENT_FUNDS", "Not enough money on balance (including credit limit)", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process withdrawal", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
