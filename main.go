package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Global variables
var (
	clients    = make(map[*websocket.Conn]bool) // Connected clients
	broadcast  = make(chan string)             // Broadcast channel
	upgrader   = websocket.Upgrader{           // Upgrader for HTTP to WebSocket
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for testing
		},
	}
	mu sync.Mutex // Mutex to synchronize client access
)

// handleConnections upgrades HTTP requests to WebSocket and manages client connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Add client to the map
	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	log.Println("New client connected, total clients:", len(clients))
	log.Println("Joined Client:", r.RemoteAddr)

	for {
		var msg []byte
		// Read a message from the client
		_, msg, err := ws.ReadMessage()
		// err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		log.Printf("Received message from %s: %s", r.RemoteAddr, msg)

		// Send the message to the broadcast channel
		var stringMag = string(msg)
		broadcast <- stringMag
	}
}

// handleMessages listens for messages on the broadcast channel and sends them to all connected clients
func handleMessages() {
	for {
		// Receive a message from the broadcast channel
		msg := <-broadcast
		log.Printf("Broadcasting message: %s", msg)

		// Send the message to all connected clients
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func main() {
	// HTTP route for WebSocket connections
	http.HandleFunc("/ws", handleConnections)

	// Start the message broadcasting goroutine
	go handleMessages()

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server started on %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
