package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransferRequest struct {
	ToAccountID  string `json:"to_account_id"`
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
	Description  string `json:"description"`
}

// @Summary Перевод средств
// @Description Осуществляет безопасный перевод денег между двумя счетами
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Public ID счета отправителя (UUID)"
// @Param Idempotency-Key header string true "Уникальный ключ запроса"
// @Param request body TransferRequest true "Данные для перевода"
// @Success 200 {object} map[string]string "Успешный перевод"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/accounts/{id}/transfer [post]
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	fromAccID := chi.URLParam(r, "id")
	fromPublicID, err := uuid.Parse(fromAccID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid sender account ID", err)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Idempotency-Key header is required", nil)
		return
	}
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	toPublicID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid receiver account ID", err)
		return
	}

	err = h.service.Transfer(
		r.Context(),
		fromPublicID,
		toPublicID,
		req.Amount,
		req.CurrencyCode,
		idempotencyKey,
		req.Description,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Transfer failed", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
