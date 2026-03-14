package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Изменить ПИН-код
// @Description  Устанавливает или меняет 4-значный ПИН-код для карты. ПИН не сохраняется в ядре и хешируется в Card Vault.
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card_id  path      string         true  "ID карты (UUID)" Format(uuid)
// @Param        request  body      domain.SetPinRequest  true  "ПИН-код (4 цифры)"
// @Success      200      {string}  string         "PIN set successfully"
// @Failure      400      {string}  string         "Неверный формат ПИН-кода"
// @Failure      403      {string}  string         "Карта заблокирована"
// @Failure      404      {string}  string         "Карта не найдена"
// @Failure      500      {string}  string         "Внутренняя ошибка сервера"
// @Router       /api/cards/{card_id}/pin [post]
func (h *Handler) setCardPin(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "id")
	cardUUID, err := uuid.Parse(cardIDStr)
	if err != nil {
		http.Error(w, "invalid card ID format", http.StatusBadRequest)
		return
	}

	var req domain.SetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.SetCardPin(r.Context(), cardUUID, req.Pin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidPINFormat):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrCardNotFound):
			http.Error(w, "card not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrCardBlocked):
			http.Error(w, "cannot set PIN for a blocked card", http.StatusForbidden)
		default:
			fmt.Printf("ERROR: HandleSetCardPin failed: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "PIN set successfully"})
}

// @Summary      Проверить ПИН-код
// @Description  Проверяет правильность ПИН-кода (например, перед подтверждением перевода).
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card_id  path      string            true  "ID карты (UUID)" Format(uuid)
// @Param        request  body      domain.VerifyPinRequest  true  "ПИН-код для проверки"
// @Success      200      {object}  map[string]bool   "Результат проверки: is_valid"
// @Failure      400      {string}  string            "Неверный формат запроса"
// @Router       /api/cards/{card_id}/verify-pin [post]
func (h *Handler) verifyCardPin(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "card_id")
	cardUUID, err := uuid.Parse(cardIDStr)
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	var req domain.VerifyPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	isValid, err := h.service.VerifyCardPin(r.Context(), cardUUID, req.Pin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidPINFormat):
			http.Error(w, "invalid PIN format: must be exactly 4 digits", http.StatusBadRequest)
		case errors.Is(err, domain.ErrCardNotFound):
			http.Error(w, "card not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrCardBlocked):
			http.Error(w, "cannot verify PIN for a blocked card", http.StatusForbidden)
		default:
			fmt.Printf("ERROR: HandleVerifyCardPin failed: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"is_valid": isValid})
}
