package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	UserID       string `json:"user_id"`
	TypeID       int32  `json:"type_id"`
	CurrencyCode string `json:"currency_code"`
	Name         string `json:"name"`
}

// @Summary Создание нового счета
// @Description Создает новый банковский счет для пользователя с нулевым балансом
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Данные для создания счета"
// @Success 201 {object} domain.Account "Счет успешно создан"
// @Failure 400 {object} map[string]string "Неверный запрос (ошибка валидации)"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts [post]
func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	// Read JSON
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	// Parse UserID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid user_id format", err)
		return
	}

	// Calling service
	acc, err := h.service.CreateAccount(
		r.Context(),
		userID,
		req.TypeID,
		req.CurrencyCode,
		req.Name,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create account", err)
		return
	}

	// Response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

// @Summary Получение информации о счете
// @Description Возвращает текущий баланс, валюту и статус счета по его публичному ID
// @Tags accounts
// @Produce json
// @Param id path string true "Public ID счета (UUID)"
// @Success 200 {object} domain.Account "Информация о счете"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Счет не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/accounts/{id} [get]
func (h *Handler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	publicID, err := uuid.Parse(accountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid account ID format", err)
		return
	}

	acc, err := h.service.GetAccount(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Account not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc)
}
