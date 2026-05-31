package models

import "github.com/gorilla/websocket"

type Client struct {
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}
type Message struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
}
