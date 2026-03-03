package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // for dev
}

var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	go func() {
		defer func() {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

func Broadcast(v any) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		if err := conn.WriteJSON(v); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}
