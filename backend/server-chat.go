package main

import "sistem-chat-realtime/models"

func newServer() *models.Server {
	return &models.Server{
		Clients:    make(map[*models.Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
	}
}
