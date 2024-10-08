package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func newServer(serverAddress string) *Server {
	return &Server{
		serverAddress:  serverAddress,
		clients:        make(map[*Client]bool),
		messageHistory: make([]Message, 0),
		lobbies:        make(map[string]*Lobby),
	}
}

func (s *Server) Run() {
	// Serve the WebSocket endpoint
	http.HandleFunc("/ws", s.handleConnections)

	// Start the server on port 8080
	fmt.Println("Server started on http:/localhost:", s.serverAddress)
	err := http.ListenAndServe(":"+s.serverAddress, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *Server) notifyClients(msg Message, lockCheck bool, recieverUsername string) {
	if lockCheck {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	for c := range s.clients {
		// fmt.Println(recieverUsername, " and the c.username is ", c.username)
		if recieverUsername == "" || c.username == recieverUsername {
			// fmt.Println("In the loop: ", c.username)
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Println("Error writing message:", err)
				c.conn.Close()
				delete(s.clients, c)
			}
		}
	}
}

// Notify clients about the connected users
func (s *Server) notifyUserList() {
	s.mu.Lock()
	defer s.mu.Unlock()
	connectedUsernames := make([]string, 0, len(s.clients))
	for c := range s.clients {
		connectedUsernames = append(connectedUsernames, c.username)
	}
	// fmt.Println("Here are connected users: ", connectedUsernames)

	userListMsg := Message{
		Username:    "Server",
		Content:     strings.Join(connectedUsernames, ", "),
		MessageType: "user_list", // Set the message type
	}

	s.notifyClients(userListMsg, false, "")
}

func (s *Server) ListLobbies() {
	s.mu.Lock()
	defer s.mu.Unlock()
	lobbyList := make([]string, 0, len(s.lobbies))
	for _, lobby := range s.lobbies {
		lobbyList = append(lobbyList, fmt.Sprintf("%s:%s", lobby.name, lobby.id))
	}

	lobbyListMsg := Message{
		Username:    "Server",
		Content:     strings.Join(lobbyList, ", "),
		MessageType: "lobby_list",
		Reciever:    "",
	}

	s.notifyClients(lobbyListMsg, false, "")
}

func (s *Server) ListLobbiesToClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lobbyList := make([]string, 0, len(s.lobbies))
	for _, lobby := range s.lobbies {
		lobbyList = append(lobbyList, fmt.Sprintf("%s:%s", lobby.name, lobby.id))
	}

	lobbyListMsg := Message{
		Username:    "Server",
		Content:     strings.Join(lobbyList, ", "),
		MessageType: "lobby_list",
		Reciever:    "",
	}

	if err := client.conn.WriteJSON(lobbyListMsg); err != nil {
		log.Println("Error writing message:", err)
		client.conn.Close()
		delete(s.clients, client)
	}
}
