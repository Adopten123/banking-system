package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

// @Summary      Проверить реквизиты карты (Внутреннее API)
// @Description  Эндпоинт для платежного шлюза. Принимает сырые реквизиты, проверяет их криптографически через Vault и возвращает привязанный счет для списания средств. Если карта заблокирована, удалена или данные неверны — возвращает is_valid: false.
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Param        request  body      domain.VerifyCardRequest  true  "Сырые реквизиты карты с терминала/сайта"
// @Success      200      {object}  domain.VerifyCardResult "Результат проверки (is_valid и, если успешно, ID счета)"
// @Failure      400      {string}  string  "Неверный формат JSON или отсутствуют обязательные поля"
// @Failure      500      {string}  string  "Внутренняя ошибка сервера (нет связи с Сейфом или БД)"
// @Router       /api/internal/payments/verify-card [post]
func (h *Handler) verifyCard(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Pan == "" || req.Cvv == "" || req.ExpiryMonth == 0 || req.ExpiryYear == 0 {
		http.Error(w, "missing required fields: pan, cvv, expiry_month, expiry_year are mandatory", http.StatusBadRequest)
		return
	}

	input := domain.VerifyCardInput{
		PAN:         req.Pan,
		CVV:         req.Cvv,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
	}

	result, err := h.service.VerifyCardForPayment(r.Context(), input)
	if err != nil {
		fmt.Printf("ERROR: HandleVerifyCard critical failure: %v\n", err)
		http.Error(w, "internal verification error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
