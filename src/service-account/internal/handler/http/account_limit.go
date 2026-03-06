package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UpdateCreditLimitRequest struct {
	CreditLimit string `json:"credit_limit"`
}

// @Summary Установить кредитный лимит
// @Description Позволяет задать сумму овердрафта для счета (чтобы баланс мог уходить в минус)
// @Tags 		accounts
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Public ID счета"
// @Param 		request body UpdateCreditLimitRequest true "Данные для обновления"
// @Success 	200 {object} map[string]string "Лимит обновлен"
// @Failure 	400 {object} map[string]string "Неверный запрос"
// @Failure 	404 {object} map[string]string "Счет не найден"
// @Failure 	500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id}/credit-limit [put]
func (h *Handler) updateCreditLimit(w http.ResponseWriter, r *http.Request) {
	accountIDParam := chi.URLParam(r, "id")
	publicID, err := uuid.Parse(accountIDParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	var req UpdateCreditLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	err = h.service.UpdateCreditLimit(r.Context(), publicID, req.CreditLimit)
	if err != nil {
		if errors.Is(err, domain.ErrAccountInactive) {
			respondWithError(w, http.StatusForbidden, "FORBIDDEN", "Account is not active", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update credit limit", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Credit limit updated successfully",
	})
}
