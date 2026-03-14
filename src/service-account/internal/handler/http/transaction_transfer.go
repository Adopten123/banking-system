package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

// @Summary Перевод средств
// @Description Осуществляет безопасный перевод денег. Источником и получателем могут выступать счет (account) или карта (card).
// @Tags       transactions
// @Accept     json
// @Produce    json
// @Param      Idempotency-Key header string true "Уникальный ключ запроса"
// @Param      request body domain.TransferRequest true "Данные для перевода"
// @Success    200 {object} domain.TransferResponse "Успешный перевод"
// @Failure    400 {object} map[string]string "Неверный запрос"
// @Failure    404 {object} map[string]string "Счет или карта не найдены"
// @Failure    409 {object} map[string]string "Дубликат транзакции"
// @Failure    500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/transfers [post]
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Idempotency-Key header is required", nil)
		return
	}

	var req domain.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	if req.SourceType == "" || req.SourceID == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "source type and source id are required", nil)
	}

	if req.DestinationType == "" || req.DestinationID == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "destination_type and destination_id are required", nil)
		return
	}

	result, err := h.service.Transfer(
		r.Context(),
		domain.TransferInput{
			SourceType:      req.SourceType,
			SourceID:        req.SourceID,
			DestinationType: req.DestinationType,
			DestinationID:   req.DestinationID,
			Amount:          req.Amount,
			Currency:        req.CurrencyCode,
			IdempotencyKey:  idempotencyKey,
			Description:     req.Description,
		},
	)

	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			respondWithError(w, http.StatusBadRequest, "INSUFFICIENT_FUNDS", "Not enough money to complete the transfer", err)
			return
		}

		if errors.Is(err, domain.ErrInvalidFormat) || errors.Is(err, domain.ErrTransferToSelf) {
			respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), err)
			return
		}

		if errors.Is(err, domain.ErrAccountInactive) || errors.Is(err, domain.ErrCardBlocked) {
			respondWithError(w, http.StatusBadRequest, "INACTIVE", "Source or destination is blocked or inactive", err)
			return
		}

		if errors.Is(err, domain.ErrAccountNotFound) || errors.Is(err, domain.ErrCardNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Source or destination not found", err)
			return
		}
		if errors.Is(err, domain.ErrDuplicateTransaction) {
			respondWithError(w, http.StatusConflict, "DUPLICATE_REQUEST", "Transaction with this Idempotency-Key already exists", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Transfer failed", err)
		return
	}

	resp := domain.TransferResponse{
		TransactionID: result.TransactionID,
		NewBalance:    result.SenderNewBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
