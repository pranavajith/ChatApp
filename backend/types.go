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

type Lobby struct {
	creatorID string
	clients   map[*Client]bool
	mu        sync.Mutex
	id        string
	name      string
}

// Message structure for chat messages
type Message struct {
	Username    string `json:"username"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Reciever    string `json:"reciever"`
	LobbyID     string `json:"lobbyId"`
}

type Server struct {
	serverAddress  string
	clients        map[*Client]bool
	lobbies        map[string]*Lobby
	mu             sync.Mutex
	messageHistory []Message
	historyMu      sync.Mutex
}
