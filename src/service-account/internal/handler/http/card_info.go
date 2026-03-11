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

// @Summary      Получить список карт счета
// @Description  Возвращает список всех карт (с маскированным PAN), привязанных к указанному счету.
// @Tags         Cards
// @Produce      json
// @Param        account_id  path      string  true  "Public ID счета (UUID)" Format(uuid)
// @Success      200         {array}   domain.Card
// @Failure      400         {string}  string  "Неверный ID счета"
// @Failure      404         {string}  string  "Счет не найден"
// @Failure      500         {string}  string  "Внутренняя ошибка сервера"
// @Router       /api/accounts/{account_id}/cards [get]
func (h *Handler) getAccountCards(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "id")
	accountUUID, err := uuid.Parse(accountIDStr)
	if err != nil {
		http.Error(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetAccountCards(r.Context(), accountUUID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusNotFound)
		default:
			fmt.Printf("ERROR: HandleGetAccountCards failed: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}
