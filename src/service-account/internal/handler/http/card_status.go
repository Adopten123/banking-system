package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Изменить статус карты
// @Description  Блокирует или разблокирует карту. Синхронизирует статус с защищенным Card Vault.
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card_id  path   string                   true  "ID карты (UUID)" Format(uuid)
// @Param        request  body   domain.UpdateCardStatusRequest  true  "Новый статус"
// @Success      200      {string} string "Status updated successfully"
// @Failure      400      {string} string "Неверный ID карты или статус"
// @Failure      500      {string} string "Внутренняя ошибка сервера"
// @Router       /api/cards/{card_id}/status [patch]
func (h *Handler) updateCardStatus(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "card_id")
	cardUUID, err := uuid.Parse(cardIDStr)
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateCardStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateCardStatus(r.Context(), cardUUID, req.Status)
	if err != nil {
		fmt.Printf("ERROR: HandleUpdateCardStatus failed: %v\n", err)
		http.Error(w, "failed to update card status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "status updated successfully"})
}
