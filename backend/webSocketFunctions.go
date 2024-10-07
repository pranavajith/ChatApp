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
	connectMessage := Message{Username: "Server", Content: fmt.Sprintf("%s has connected.", username)}
	s.notifyClients(connectMessage)

	s.historyMu.Lock()
	for _, msg := range s.messageHistory {
		if err := ws.WriteJSON(msg); err != nil {
			log.Println("Error sending message history:", err)
		}
	}
	s.historyMu.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		msg.Username = username
		fmt.Printf("Received message from %s: %s\n", username, msg.Content)

		s.historyMu.Lock()
		s.messageHistory = append(s.messageHistory, msg)
		s.historyMu.Unlock()

		s.notifyClients(msg)
	}

	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
	fmt.Printf("Client disconnected: %s\n", username)

	disconnectMsg := Message{Username: "Server", Content: fmt.Sprintf("%s has disconnected.", username)}
	s.notifyClients(disconnectMsg)
}

func (s *Server) notifyClients(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for c := range s.clients {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println("Error writing message:", err)
			c.conn.Close()
			delete(s.clients, c)
		}
	}
}
