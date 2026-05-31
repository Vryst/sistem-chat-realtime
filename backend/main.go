package main

import (
	"fmt"
	"net/http"
	"sistem-chat-realtime/models"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 4096 // 4 KB

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, s *models.Server) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrade:", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)

	client := &models.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	s.Register <- client

	go ReadPump(client, s)
	go WritePump(client)
}

func main() {
	server := newServer()
	go server.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, server)
	})

	fmt.Println("Server jalan di port 8080")
	http.ListenAndServe(":8080", nil)
}
