package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

func (c *Client) handleSendMessage(payloadBytes []byte) {
	var payload SendMessagePayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		log.Printf("Invalid payload for send_message: %v", err)
		return
	}

	chatID, err := uuid.Parse(payload.ChatID)
	if err != nil {
		log.Printf("Invalid chat_id format: %v", err)
		return
	}

	senderID, _ := uuid.Parse(c.UserID)

	input := domain.CreateMessageInput{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  &payload.Content,
	}

	savedMsg, err := c.Hub.chatService.SendMessage(context.Background(), input)
	if err != nil {
		log.Printf("Failed to save message to DB: %v", err)
		return
	}

	log.Printf("Message saved successfully! ID: %d", savedMsg.ID)

	responseBytes, _ := json.Marshal(map[string]interface{}{
		"action":  "new_message",
		"message": savedMsg,
	})

	c.Hub.BroadcastToChat(context.Background(), chatID, responseBytes)

	c.Hub.SendToUser(c.UserID, responseBytes)
}

