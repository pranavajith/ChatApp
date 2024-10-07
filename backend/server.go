package main

import (
	"fmt"
	"log"
	"net/http"
)

func newServer(serverAddress string) *Server {
	return &Server{
		serverAddress:  serverAddress,
		clients:        make(map[*Client]bool),
		messageHistory: make([]Message, 0),
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
