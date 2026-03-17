package websocket

import (
	"context"
	"log"

	"github.com/google/uuid"
)

func (h *Hub) BroadcastToChat(ctx context.Context, chatID uuid.UUID, message []byte) {
	memberIDs, err := h.chatService.GetChatMemberIDs(ctx, chatID)
	if err != nil {
		log.Printf("Failed to get chat members for broadcast: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, memberID := range memberIDs {
		userIDStr := memberID.String()
		if client, ok := h.clients[userIDStr]; ok {
			select {
			case client.Send <- message:
			default:
				log.Printf("Buffer full for user %s, dropping message", userIDStr)
			}
		}
	}
}