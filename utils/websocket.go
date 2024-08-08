package utils

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan map[string]interface{}
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan map[string]interface{}),
	}
}

func (hub *WebSocketHub) Run() {
	for {
		data := <-hub.Broadcast
		for client := range hub.Clients {
			err := client.WriteJSON(data)
			if err != nil {
				log.Printf("WebSocket error: %v", err)
				client.Close()
				delete(hub.Clients, client)
			}
		}
	}
}

func (hub *WebSocketHub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade error: %v", err)
		return
	}
	hub.Clients[conn] = true
}
