package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

// @Summary      Снятие наличных
// @Description  Списывает указанную сумму со счета или карты. Учитывает кредитный лимит.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        Idempotency-Key header string true "Ключ идемпотентности"
// @Param        body body      domain.WithdrawRequest true  "Данные для снятия"
// @Success      200  {object}  domain.WithdrawResponse "Успешное снятие"
// @Failure      400  {object}  map[string]string "Неверный запрос"
// @Failure      404  {object}  map[string]string "Счет не найден"
// @Failure      500  {object}  map[string]string "Внутренняя ошибка"
// @Router       /api/withdrawals [post]
func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_HEADER", "Idempotency-Key header is required", nil)
		return
	}

	var req domain.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	if req.SourceType == "" || req.SourceID == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "source_type and source_id are required", nil)
		return
	}

	input := domain.ServiceWithdrawInput{
		SourceType:     req.SourceType,
		SourceValue:    req.SourceID,
		AmountStr:      req.Amount,
		IdempotencyKey: idempotencyKey,
	}

	result, err := h.service.Withdraw(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidFormat) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid source ID format", err)
			return
		}
		if errors.Is(err, domain.ErrAccountNotFound) || errors.Is(err, domain.ErrCardNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Source not found", err)
			return
		}
		if errors.Is(err, domain.ErrAccountInactive) || errors.Is(err, domain.ErrCardBlocked) {
			respondWithError(w, http.StatusBadRequest, "INACTIVE", "Source is blocked or inactive", err)
			return
		}
		if errors.Is(err, domain.ErrInsufficientFunds) {
			respondWithError(w, http.StatusBadRequest, "INSUFFICIENT_FUNDS", "Not enough money on balance", err)
			return
		}
		if errors.Is(err, domain.ErrDuplicateTransaction) {
			respondWithError(w, http.StatusConflict, "DUPLICATE_REQUEST", "Transaction with this Idempotency-Key already exists", err)
			return
		}
		if errors.Is(err, domain.ErrInvalidAmountFormat) || errors.Is(err, domain.ErrInvalidWithdrawAmount) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid withdrawal amount", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process withdrawal", err)
		return
	}

	resp := domain.WithdrawResponse{
		TransactionID: result.TransactionID,
		NewBalance:    result.NewBalance,
		Currency:      result.Currency,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
