package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID string
	Send   chan []byte
}

type GetHistoryPayload struct {
	ChatID string `json:"chat_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}