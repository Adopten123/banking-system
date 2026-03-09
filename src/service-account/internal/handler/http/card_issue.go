package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Выпуск новой карты
// @Description  Создает новую физическую или виртуальную карту для указанного счета. Внутри обращается к защищенному Card Vault для генерации PAN и CVV.
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        account_id   path      string            true  "Public ID счета (UUID)" Format(uuid)
// @Param        request      body      IssueCardRequest  true  "Параметры выпуска карты"
// @Success      201          {object}  domain.Card       "Карта успешно выпущена"
// @Failure      400          {string}  string            "Неверный ID счета или формат запроса"
// @Failure      500          {string}  string            "Внутренняя ошибка сервера"
// @Router       /api/accounts/{account_id}/cards [post]

func (h *Handler) issueCard(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "account_id")
	accountUUID, err := uuid.Parse(accountIDStr)
	if err != nil {
		http.Error(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	var req domain.IssueCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	card, err := h.service.IssueCard(r.Context(), domain.IssueCardInput{
		AccountPublicID: accountUUID,
		PaymentSystem:   req.PaymentSystem,
		IsVirtual:       req.IsVirtual,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
