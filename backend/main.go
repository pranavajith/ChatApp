package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

// Client structure to hold the connection, username, and any other necessary information
type Client struct {
	conn     *websocket.Conn
	username string
}

// Message structure for chat messages
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// Global slice to hold connected clients and a mutex for concurrent access
var clients = make(map[*Client]bool)
var mu sync.Mutex

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
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
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	fmt.Printf("New client connected: %s\n", username)

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

		// Broadcast the message to all clients
		mu.Lock()
		for c := range clients {
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Println("Error writing message:", err)
				c.conn.Close()
				delete(clients, c)
			}
		}
		mu.Unlock()
	}

	// Remove the client from the global clients map
	mu.Lock()
	delete(clients, client)
	mu.Unlock()
	fmt.Printf("Client disconnected: %s\n", username)

	// Notify other clients about the disconnection
	disconnectMsg := Message{Username: "Server", Content: fmt.Sprintf("%s has disconnected.", username)}
	mu.Lock()
	for c := range clients {
		if err := c.conn.WriteJSON(disconnectMsg); err != nil {
			log.Println("Error notifying clients:", err)
			c.conn.Close()
			delete(clients, c)
		}
	}
	mu.Unlock()
}

func main() {
	// Serve the WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the server on port 8080
	fmt.Println("WebSocket server started on ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
