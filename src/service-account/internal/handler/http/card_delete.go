package http

import (
	"errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Удалить карту
// @Description  Полностью удаляет данные карты из Card Vault и помечает её как удаленную в ядре.
// @Tags         Cards
// @Param        card_id  path      string  true  "ID карты (UUID)" Format(uuid)
// @Success      200      {string}  string  "Card deleted successfully"
// @Router       /api/cards/{card_id} [delete]
func (h *Handler) deleteCard(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "card_id")
	cardUUID, err := uuid.Parse(cardIDStr)
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCard(r.Context(), cardUUID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCardNotFound):
			http.Error(w, "card not found", http.StatusNotFound)
		default:
			fmt.Printf("ERROR: HandleDeleteCard failed: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "card deleted successfully"})
}