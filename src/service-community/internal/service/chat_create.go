package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

func (s *chatService) CreateChat(ctx context.Context, input domain.CreateChatInput) (*domain.Chat, error) {
	chat, err := s.repo.CreateChat(ctx, input.TypeID, input.Title, input.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	for _, memberID := range input.MemberIDs {
		err := s.repo.AddChatMember(ctx, chat.ID, memberID, "member")
		if err != nil {
			log.Printf("Failed to add user %s to chat %s: %v", memberID, chat.ID, err)
		}
	}

	return chat, nil
}