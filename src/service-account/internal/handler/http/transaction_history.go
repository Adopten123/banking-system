package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Получить историю транзакций
// @Description Возвращает список транзакций по счету с поддержкой пагинации
// @Tags accounts
// @Produce json
// @Param id path string true "Public ID счета"
// @Param limit query int false "Количество записей (по умолчанию 20)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Param start_date query string false "Начальная дата (RFC3339, например: 2026-02-01T00:00:00Z)"
// @Param end_date query string false "Конечная дата (RFC3339)"
// @Success 200 {array} domain.TransactionHistory"
// @Failure 400 {object} map[string]string "Неверный запрос"
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

	query := r.URL.Query()
	limit := int32(20)
	offset := int32(0)

	if l := query.Get("limit"); l != "" {
		if parsedLimit, err := strconv.ParseInt(l, 10, 32); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	if o := query.Get("offset"); o != "" {
		if parsedOffset, err := strconv.ParseInt(o, 10, 32); err == nil && parsedOffset >= 0 {
			offset = int32(parsedOffset)
		}
	}

	var startDate, endDate *time.Time

	if sd := query.Get("start_date"); sd != "" {
		if parsedDate, err := time.Parse(time.RFC3339, sd); err == nil {
			startDate = &parsedDate
		} else {
			respondWithError(w, http.StatusBadRequest, "INVALID_DATE", "Invalid start_date format. Use RFC3339", err)
			return
		}
	}

	if ed := query.Get("end_date"); ed != "" {
		if parsedDate, err := time.Parse(time.RFC3339, ed); err == nil {
			endDate = &parsedDate
		} else {
			respondWithError(w, http.StatusBadRequest, "INVALID_DATE", "Invalid end_date format. Use RFC3339", err)
			return
		}
	}

	history, err := h.service.GetAccountTransactions(r.Context(), publicID, limit, offset, startDate, endDate)
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
