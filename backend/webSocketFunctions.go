package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade to WebSocket: %v", err)
	}
	defer ws.Close()

	username := r.URL.Query().Get("username")

	client := &Client{conn: ws, username: username}
	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	fmt.Printf("New client connected: %s\n", username)

	// Notify all clients about the new connection
	connectMessage := Message{
		Username:    "Server",
		Content:     fmt.Sprintf("%s has connected.", username),
		MessageType: "server_announcement",
	}
	s.notifyClients(connectMessage, true, "")
	s.notifyUserList()

	// Send the message history to the connected client
	s.historyMu.Lock()
	s.ListLobbiesToClient(client)

	for _, msg := range s.messageHistory {
		if msg.Reciever == username || msg.Reciever == "" {
			if err := ws.WriteJSON(msg); err != nil {
				log.Println("Error sending message history:", err)
			}
		}
	}
	s.historyMu.Unlock()

	// Define the client's current lobby
	var currentLobby *Lobby

	// Main message handling loop
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		switch msg.MessageType {

		case "typing":
			s.notifyClients(msg, true, msg.Reciever)

		case "create_lobby":
			lobbyID := msg.LobbyID
			lobbyName := msg.Content
			s.CreateLobby(lobbyID, lobbyName, username)
			fmt.Println("Lobby created: ", lobbyID, " with name: ", lobbyName)
			s.ListLobbies()

		case "join_lobby":
			lobbyID := msg.LobbyID
			s.JoinLobby(client, lobbyID)
			currentLobby = s.lobbies[lobbyID]

		case "leave_lobby":
			if currentLobby != nil {
				s.LeaveLobby(client, currentLobby.id)
				currentLobby = nil
			}

		case "remove_lobby":
			lobbyID := msg.LobbyID
			lobby, exists := s.lobbies[lobbyID]
			if !exists {
				fmt.Println("No lobby exists with lobbyID", lobbyID)
				continue
			}
			if lobby.creatorID != username {
				fmt.Println("User is not the creator of the lobby")
				continue
			}
			delete(s.lobbies, lobbyID)
			fmt.Printf("Lobby %s removed by %s\n", lobbyID, username)
			s.ListLobbies()

		case "lobby_message":
			if currentLobby != nil {
				msg.Username = fmt.Sprintf("%s (Lobby: %s)", username, currentLobby.name)
				currentLobby.notifyClients(msg)
			}

		default:
			msg.Username = username
			s.historyMu.Lock()
			s.messageHistory = append(s.messageHistory, msg)
			s.historyMu.Unlock()
			s.notifyClients(msg, true, msg.Reciever)
		}
	}

	// Client disconnected
	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
	fmt.Printf("Client disconnected: %s\n", username)

	// Notify all clients about the disconnection
	disconnectMsg := Message{
		Username:    "Server",
		Content:     fmt.Sprintf("%s has disconnected.", username),
		MessageType: "server_announcement",
	}
	s.notifyClients(disconnectMsg, true, "")
	s.notifyUserList()
}
