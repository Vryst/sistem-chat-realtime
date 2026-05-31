package models

const maxHistory = 100

type Server struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	History    [][]byte
}
