package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

// Handle WebSocket connections
func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade to WebSocket: %v", err)
	}
	defer ws.Close()

	// Get the username from the query parameters
	username := r.URL.Query().Get("username")

	// Create a new client and add it to the global clients map
	client := &Client{conn: ws, username: username}
	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	fmt.Printf("New client connected: %s\n", username)

	// Send existing messages to the newly connected client
	s.historyMu.Lock()
	for _, msg := range s.messageHistory {
		if err := ws.WriteJSON(msg); err != nil {
			log.Println("Error sending message history:", err)
		}
	}
	s.historyMu.Unlock()

	// Infinite loop to keep the connection open
	for {
		// Read in a new message
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Include the username in the broadcast message
		msg.Username = username
		fmt.Printf("Received message from %s: %s\n", username, msg.Content)

		// Add the message to the history
		s.historyMu.Lock()
		s.messageHistory = append(s.messageHistory, msg)
		s.historyMu.Unlock()

		// Broadcast the message to all clients
		s.mu.Lock()
		for c := range s.clients {
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Println("Error writing message:", err)
				c.conn.Close()
				delete(s.clients, c)
			}
		}
		s.mu.Unlock()
	}

	// Remove the client from the global clients map
	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
	fmt.Printf("Client disconnected: %s\n", username)

	// Notify other clients about the disconnection
	disconnectMsg := Message{Username: "Server", Content: fmt.Sprintf("%s has disconnected.", username)}
	s.mu.Lock()
	for c := range s.clients {
		if err := c.conn.WriteJSON(disconnectMsg); err != nil {
			log.Println("Error notifying clients:", err)
			c.conn.Close()
			delete(s.clients, c)
		}
	}
	s.mu.Unlock()
}
