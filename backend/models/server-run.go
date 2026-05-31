package models

import "encoding/json"

func (s *Server) Run() {
	for {
		select {
		case client := <-s.Register:
			s.Clients[client] = true
			// Kirim history ke client baru
			for _, msg := range s.History {
				select {
				case client.Send <- msg:
				default:
				}
			}

		case client := <-s.Unregister:
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.Send)

				// Broadcast notif leave
				if client.Username != "" {
					notif := Message{
						Type:     "leave",
						Username: client.Username,
						Content:  client.Username + " meninggalkan chat",
					}
					out, _ := json.Marshal(notif)
					for c := range s.Clients {
						select {
						case c.Send <- out:
						default:
							delete(s.Clients, c)
							close(c.Send)
						}
					}
				}
			}

		case message := <-s.Broadcast:
			s.History = append(s.History, message)
			if len(s.History) > maxHistory {
				s.History = s.History[len(s.History)-maxHistory:]
			}
			for client := range s.Clients {
				select {
				case client.Send <- message:
				default:
					delete(s.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}
