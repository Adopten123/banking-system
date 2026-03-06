package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Блокировать счет
// @Description Переводит счет в статус "blocked" (3)
// @Tags 		accounts
// @Produce 	json
// @Param 		id path string true "Public ID счета"
// @Success 	200 {object} map[string]string "Счет заблокирован"
// @Failure 	400 {object} map[string]string "Неверный формат ID"
// @Failure 	404 {object} map[string]string "Счет не найден"
// @Failure 	500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id}/block [post]
func (h *Handler) blockAccount(w http.ResponseWriter, r *http.Request) {
	accIDParam := chi.URLParam(r, "id")
	publicID, err := uuid.Parse(accIDParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	err = h.service.BlockAccount(r.Context(), publicID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to block account", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Account blocked successfully",
	})
}

// @Summary Закрыть счет
// @Description Переводит счет в статус "closed" (4). Требует строго нулевого баланса!
// @Tags accounts
// @Produce json
// @Param id path string true "Public ID счета"
// @Success 200 {object} map[string]string "Счет закрыт"
// @Failure 400 {object} map[string]string "Неверный ID или ненулевой баланс"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id}/close [post]
func (h *Handler) closeAccount(w http.ResponseWriter, r *http.Request) {
	accountIDParam := chi.URLParam(r, "id")
	publicID, err := uuid.Parse(accountIDParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	err = h.service.CloseAccount(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountHasBalance) {
			respondWithError(w, http.StatusBadRequest, "NON_ZERO_BALANCE", "Account balance must be exactly zero to close", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to close account", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Account closed successfully",
	})
}
