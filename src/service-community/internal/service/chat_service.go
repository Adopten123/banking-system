package service

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

type chatService struct {
	repo      domain.ChatRepository
	publisher domain.TransferPublisher
}

func NewChatService(repo domain.ChatRepository, publisher domain.TransferPublisher) domain.ChatService {
	return &chatService{repo: repo, publisher: publisher}
}