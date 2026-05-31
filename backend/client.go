package main

import (
	"encoding/json"
	"sistem-chat-realtime/models"

	"github.com/gorilla/websocket"
)

func ReadPump(c *models.Client, s *models.Server) {
	defer func() {
		s.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			break
		}

		if msg.Type == "username" {
			c.Username = msg.Content
			// Broadcast notif join ke semua client
			notif := models.Message{
				Type:     "join",
				Username: c.Username,
				Content:  c.Username + " bergabung ke chat",
			}
			out, _ := json.Marshal(notif)
			s.Broadcast <- out
		} else {
			msg.Username = c.Username
			out, _ := json.Marshal(msg)
			s.Broadcast <- out
		}
	}
}

func WritePump(c *models.Client) {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
