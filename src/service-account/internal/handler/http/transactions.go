package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// DepositRequest - waiting sum for deposit
type DepositRequest struct {
	Amount string `json:"amount"`
}

// @Summary Пополнение счета
// @Description Зачисляет средства на указанный счет. Требует передачи Idempotency-Key в заголовках.
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Public ID счета (UUID)"
// @Param Idempotency-Key header string true "Уникальный ключ запроса"
// @Param request body DepositRequest true "Сумма пополнения"
// @Success 200 {object} map[string]string "Успешное пополнение"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/accounts/{id}/deposit [post]
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
	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Sending data into service
	err = h.service.Deposit(r.Context(), publicID,
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

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process deposit", err)
		return
	}

	// Making response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Deposit processed successfully",
	})
}
