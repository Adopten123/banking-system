package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

// @Summary      Пополнение счета/карты
// @Description  Зачисляет указанную сумму на счет или карту
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        Idempotency-Key header    string          true "Ключ идемпотентности"
// @Param        body            body      domain.DepositRequest  true "Данные пополнения"
// @Success      201  {object}  domain.DepositResponse "Успешное пополнение"
// @Failure      400  {object}  map[string]string   "Неверный запрос или сумма"
// @Failure      404  {object}  map[string]string   "Счет/карта не найдены"
// @Failure      409  {object}  map[string]string   "Дубликат транзакции"
// @Failure      500  {object}  map[string]string   "Внутренняя ошибка"
// @Router       /api/deposits [post]
func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	// Больше не берем ID из chi.URLParam(r, "id")

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_HEADER", "Idempotency-Key header is required", nil)
		return
	}

	var req domain.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	result, err := h.service.Deposit(r.Context(), domain.ServiceDepositInput{
		DestinationType:  req.DestinationType,
		DestinationValue: req.DestinationID,
		AmountStr:        req.Amount,
		IdempotencyKey:   idempotencyKey,
	})

	if err != nil {
		// Расширенный блок ошибок, как в трансферах
		if errors.Is(err, domain.ErrInvalidFormat) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid destination ID format", err)
			return
		}
		if errors.Is(err, domain.ErrAccountNotFound) || errors.Is(err, domain.ErrCardNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Destination not found", err)
			return
		}
		if errors.Is(err, domain.ErrAccountInactive) || errors.Is(err, domain.ErrCardBlocked) {
			respondWithError(w, http.StatusForbidden, "FORBIDDEN", "Destination is not active", err)
			return
		}
		if errors.Is(err, domain.ErrInvalidDepositAmount) || errors.Is(err, domain.ErrInvalidAmountFormat) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid deposit amount", err)
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
