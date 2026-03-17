package service

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

type chatService struct {
	repo domain.ChatRepository
}

func NewChatService(repo domain.ChatRepository) domain.ChatService {
	return &chatService{repo: repo}
}