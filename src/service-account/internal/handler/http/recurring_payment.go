package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Создать регулярный платеж (подписку)
// @Description  Создает новый автоплатеж по CRON расписанию
// @Tags         recurring_payments
// @Accept       json
// @Produce      json
// @Param        body body CreateRecurringPaymentRequest true "Данные автоплатежа"
// @Success      201  {object}  map[string]interface{}
// @Router       /api/recurring-payments [post]
func (h *Handler) createRecurringPayment(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateRecurringPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Базовая валидация
	if req.SourceType == "" || req.SourceID == "" || req.CronExpression == "" {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "source_type, source_id, and cron_expression are required", nil)
		return
	}

	input := domain.CreateRecurringPaymentInput{
		SourceType:       req.SourceType,
		SourceValue:      req.SourceID,
		DestinationType:  req.DestinationType,
		DestinationValue: req.DestinationID,
		Amount:           req.Amount,
		CurrencyCode:     req.CurrencyCode,
		CategoryID:       req.CategoryID,
		CronExpression:   req.CronExpression,
		Description:      req.Description,
	}

	paymentID, err := h.service.CreateRecurringPayment(r.Context(), input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create recurring payment", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      paymentID,
		"message": "Recurring payment scheduled successfully",
	})
}

// @Summary      Отменить регулярный платеж
// @Description  Отключает (деактивирует) автоплатеж по его ID
// @Tags         recurring_payments
// @Produce      json
// @Param        id path string true "ID подписки"
// @Success      200  {object}  map[string]string
// @Router       /api/recurring-payments/{id} [delete]
func (h *Handler) cancelRecurringPayment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	paymentID, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid payment ID format", err)
		return
	}

	err = h.service.CancelRecurringPayment(r.Context(), paymentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to cancel recurring payment", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Recurring payment cancelled successfully",
	})
}
