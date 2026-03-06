package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Пополнение счета
// @Description  Зачисляет указанную сумму на счет
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id              path      string          true "ID счета (UUID)" Format(uuid)
// @Param        Idempotency-Key header    string          true "Ключ идемпотентности"
// @Param        body            body      domain.DepositRequest  true "Сумма пополнения"
// @Success      200  {object}  domain.DepositResponse "Успешное пополнение"
// @Failure      400  {object}  ErrorResponse   "Неверный запрос или сумма"
// @Failure      404  {object}  ErrorResponse   "Счет не найден"
// @Failure      409  {object}  ErrorResponse   "Дубликат транзакции"
// @Failure      500  {object}  ErrorResponse   "Внутренняя ошибка"
// @Router       /api/accounts/{id}/deposit [post]
func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	// Getting UUID from url
	accountIDParam := chi.URLParam(r, "id")
	publicID, err := uuid.Parse(accountIDParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	// Reading idempotency-key from http-header
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_HEADER", "Idempotency-Key header is required", nil)
		return
	}

	// Parsing request
	var req domain.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Sending data into service
	result, err := h.service.Deposit(r.Context(), publicID,
		domain.ServiceDepositInput{
			AmountStr:      req.Amount,
			IdempotencyKey: idempotencyKey,
		},
	)

	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		if errors.Is(err, domain.ErrAccountInactive) {
			respondWithError(w, http.StatusForbidden, "FORBIDDEN", "Account is not active", err)
			return
		}
		if errors.Is(err, domain.ErrInvalidDepositAmount) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid deposit amount", err)
			return
		}
		if errors.Is(err, domain.ErrInvalidAmountFormat) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid amount format", err)
			return
		}
		if errors.Is(err, domain.ErrDuplicateTransaction) {
			respondWithError(w, http.StatusConflict, "DUPLICATE_REQUEST", "Transaction with this Idempotency-Key already exists", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process deposit", err)
		return
	}

	resp := domain.DepositResponse{
		TransactionID: result.TransactionID,
		NewBalance:    result.NewBalance,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
