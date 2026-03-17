package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

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
		case "get_history":
			c.handleGetHistory(event.Payload)
		default:
			log.Printf("Unknown action '%s' from user %s", event.Action, c.UserID)
		}
	}
}
