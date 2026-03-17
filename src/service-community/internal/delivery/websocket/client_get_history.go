package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (c *Client) handleGetHistory(payloadBytes []byte) {
	var payload GetHistoryPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return
	}

	chatID, err := uuid.Parse(payload.ChatID)
	if err != nil {
		return
	}

	userID, _ := uuid.Parse(c.UserID)
	limit := payload.Limit
	if limit == 0 {
		limit = 50
	}

	messages, err := c.Hub.chatService.GetHistory(context.Background(), chatID, userID, limit, payload.Offset)
	if err != nil {
		log.Printf("Failed to get history: %v", err)
		return
	}

	responseBytes, _ := json.Marshal(map[string]interface{}{
		"action":   "chat_history",
		"chat_id":  chatID.String(),
		"messages": messages,
	})

	c.Hub.SendToUser(c.UserID, responseBytes)
}
