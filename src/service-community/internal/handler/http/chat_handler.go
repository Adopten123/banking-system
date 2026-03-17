package http

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	chatService domain.ChatService
}

func NewChatHandler(chatService domain.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/chats", h.createChat)
}
