package main

import (
	"fmt"
	"log"
)

func (s *Server) CreateLobby(lobbyID, lobbyName string, username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lobbies[lobbyID]; exists {
		log.Printf("Lobby with ID %s already exists\n", lobbyID)
		return
	}

	newLobby := &Lobby{
		creatorID: username,
		clients:   make(map[*Client]bool),
		id:        lobbyID,
		name:      lobbyName,
	}

	s.lobbies[lobbyID] = newLobby
	log.Printf("Lobby created with ID: %s and Name: %s\n", lobbyID, lobbyName)
}

// Client joins a lobby
func (s *Server) JoinLobby(client *Client, lobbyID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, exists := s.lobbies[lobbyID]
	if !exists {
		log.Printf("Lobby with ID %s does not exist\n", lobbyID)
		return
	}

	lobby.mu.Lock()
	lobby.clients[client] = true
	lobby.mu.Unlock()

	joinMsg := Message{
		Username:    "Server",
		Content:     fmt.Sprintf("%s has joined the lobby %s", client.username, lobby.name),
		MessageType: "lobby_message",
	}
	lobby.notifyClients(joinMsg)
	log.Printf("Client %s joined lobby %s\n", client.username, lobbyID)
}

// Client leaves a lobby
func (s *Server) LeaveLobby(client *Client, lobbyID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, exists := s.lobbies[lobbyID]
	if !exists {
		log.Printf("Lobby with ID %s does not exist\n", lobbyID)
		return
	}

	lobby.mu.Lock()
	delete(lobby.clients, client)
	lobby.mu.Unlock()

	leaveMsg := Message{
		Username:    "Server",
		Content:     fmt.Sprintf("%s has left the lobby %s", client.username, lobby.name),
		MessageType: "lobby_message",
	}
	lobby.notifyClients(leaveMsg)
	log.Printf("Client %s left lobby %s\n", client.username, lobbyID)
}

func (l *Lobby) notifyClients(msg Message) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for c := range l.clients {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println("Error writing message:", err)
			c.conn.Close()
			delete(l.clients, c)

		}
	}
	fmt.Println("Message sent to all clients in lobby")
}
