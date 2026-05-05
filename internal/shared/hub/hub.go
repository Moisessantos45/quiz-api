package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	RoomKey string
}

type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*Client]bool
}

var Global = &Hub{
	rooms: make(map[string]map[*Client]bool),
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[client.RoomKey] == nil {
		h.rooms[client.RoomKey] = make(map[*Client]bool)
	}
	h.rooms[client.RoomKey][client] = true
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.rooms[client.RoomKey]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.rooms, client.RoomKey)
		}
	}
}

func (h *Hub) Broadcast(roomKey string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.rooms[roomKey] {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
