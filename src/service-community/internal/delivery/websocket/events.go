package websocket

import "encoding/json"

type WSEvent struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}