package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// User struct
type User struct {
	UID  string
	Name string
}

// Global variables
var (
	clients   = make(map[*websocket.Conn]*User) // Map WebSocket connections to users
	users     = make(map[string]*User)          // Map UID to User
	broadcast = make(chan string)               // Broadcast channel
	upgrader  = websocket.Upgrader{             // Upgrader for HTTP to WebSocket
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mu sync.Mutex // Mutex for thread-safe access to maps
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed for UID generation
}

// generateUID generates a unique 4-digit UID
func generateUID() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// registerUser handles user registration via HTTP
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the user name from the request
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Generate a UID and create a user
	uid := generateUID()
	user := &User{UID: uid, Name: req.Name}

	// Save the user
	mu.Lock()
	users[uid] = user
	mu.Unlock()

	// Respond with the UID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"uid": uid})
}

// handleConnections handles WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Get the UID from the query parameters
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "UID is required", http.StatusBadRequest)
		return
	}

	// Validate the UID
	mu.Lock()
	user, exists := users[uid]
	mu.Unlock()
	if !exists {
		http.Error(w, "Invalid UID", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Register the WebSocket connection with the user
	mu.Lock()
	clients[ws] = user
	mu.Unlock()
	log.Printf("User %s (UID: %s) connected", user.Name, user.UID)

	for {
		var msg []byte
		// err := ws.ReadJSON(&msg)
		_, msg, err = ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}

		log.Printf("Received message from %s: %s", user.Name, string(msg))
		broadcast <- fmt.Sprintf("%s (%s): %s", user.Name, user.UID, string(msg))
	}
}

// handleMessages broadcasts messages to all connected clients
func handleMessages() {
	for {
		msg := <-broadcast
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
	// HTTP endpoint for user registration
	http.HandleFunc("/register", registerUser)

	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// dummy endpoint
	http.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

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
