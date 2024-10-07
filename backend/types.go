package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client structure to hold the connection, username, and any other necessary information
type Client struct {
	conn     *websocket.Conn
	username string
}

// Message structure for chat messages
type Message struct {
	Username    string `json:"username"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
}

type Server struct {
	serverAddress  string
	clients        map[*Client]bool
	mu             sync.Mutex
	messageHistory []Message
	historyMu      sync.Mutex
}
