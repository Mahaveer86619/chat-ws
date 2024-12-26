package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer ws.Close()
    for {
        var msg string
        err := ws.ReadJSON(&msg)
        if err != nil {
            fmt.Println(err)
            break
        }
        fmt.Printf("Received: %s\n", msg)
    }
}

func main() {
    http.HandleFunc("/ws", handleConnections)
    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}