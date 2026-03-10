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

// @Summary      Получить реквизиты карты
// @Description  Возвращает полные данные карты (PAN, CVV) для показа в приложении. Данные запрашиваются из защищенного Card Vault.
// @Tags         Cards
// @Produce      json
// @Param        card_id   path      string  true  "ID карты (UUID токена)" Format(uuid)
// @Success      200       {object}  domain.CardDetails
// @Failure      400       {string}  string  "Неверный ID карты"
// @Failure      403       {string}  string  "Карта заблокирована"
// @Failure      404       {string}  string  "Карта не найдена"
// @Router       /api/cards/{card_id}/details [get]
func (h *Handler) getCardDetails(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "card_id")
	cardUUID, err := uuid.Parse(cardIDStr)
	if err != nil {
		http.Error(w, "invalid card ID format", http.StatusBadRequest)
		return
	}

	details, err := h.service.GetCardDetails(r.Context(), cardUUID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCardNotFound):
			http.Error(w, "card not found", http.StatusNotFound)

		case errors.Is(err, domain.ErrCardBlocked):
			http.Error(w, "cannot get details for a blocked card", http.StatusForbidden)

		default:
			fmt.Printf("ERROR: HandleGetCardDetails failed: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(details)
}
