package http

import "net/http"

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: распарсить JSON из r.Body, валидировать, передать в h.accountService
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	// TODO: достать ID из URL (chi.URLParam(r, "id")), передать в h.accountService
	w.WriteHeader(http.StatusNotImplemented)
}
