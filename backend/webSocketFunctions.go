package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	connectMessage := Message{Username: "Server", Content: fmt.Sprintf("%s has connected.", username), MessageType: "server_announcement"}
	s.notifyClients(connectMessage, true, "")
	s.notifyUserList()

	s.historyMu.Lock()
	for _, msg := range s.messageHistory {
		if msg.Reciever == username || msg.Reciever == "" {
			if err := ws.WriteJSON(msg); err != nil {
				log.Println("Error sending message history:", err)
			}
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

		if msg.MessageType == "typing" {
			s.notifyClients(msg, true, msg.Reciever)
			continue
		}

		msg.Username = username
		fmt.Printf("Received message from %s: %s\n", username, msg.Content)

		s.historyMu.Lock()
		s.messageHistory = append(s.messageHistory, msg)
		s.historyMu.Unlock()

		s.notifyClients(msg, true, msg.Reciever)
	}

	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
	fmt.Printf("Client disconnected: %s\n", username)

	disconnectMsg := Message{Username: "Server", Content: fmt.Sprintf("%s has disconnected.", username), MessageType: "server_announcement"}
	s.notifyClients(disconnectMsg, true, "")
	s.notifyUserList()
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
	// fmt.Println("Message sent to clients: ", msg)
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
