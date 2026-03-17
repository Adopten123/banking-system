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