package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(c.Hub.cfg.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.Hub.cfg.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(c.Hub.cfg.PongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, []byte{'\n'}, []byte{' '}, -1))

		var event WSEvent
		if err := json.Unmarshal(messageBytes, &event); err != nil {
			log.Printf("Invalid WS message format from user %s: %v", c.UserID, err)
			continue
		}

		switch event.Action {
		case "send_message":
			c.handleSendMessage(event.Payload)
		default:
			log.Printf("Unknown action '%s' from user %s", event.Action, c.UserID)
		}
	}
}

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

	// TODO: c.Hub.BroadcastToChat(chatID, responseBytes),
	c.Hub.SendToUser(c.UserID, responseBytes)
}
