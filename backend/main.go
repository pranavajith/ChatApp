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

// Client structure to hold the connection and any other necessary information
type Client struct {
	conn *websocket.Conn
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

	// Create a new client and add it to the global clients map
	client := &Client{conn: ws}
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	fmt.Println("New client connected")

	// Infinite loop to keep the connection open
	for {
		// Read in a new message
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received message: %s\n", p)

		// Broadcast the message to all clients
		mu.Lock()
		for c := range clients {
			if err := c.conn.WriteMessage(messageType, p); err != nil {
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
	fmt.Println("Client disconnected")
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
