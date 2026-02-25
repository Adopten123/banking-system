package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
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