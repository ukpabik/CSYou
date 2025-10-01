package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/ukpabik/CSYou/pkg/api/handlers"
	"github.com/ukpabik/CSYou/pkg/api/model"
)

var httpClient *http.Client

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var clientsMu sync.Mutex

// Broadcast channel for logs
var broadcast = make(chan model.Log, 100)

// WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	// Keep connection alive
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			break
		}
	}
}

// broadcaster sends log entries to all connected clients
func broadcaster() {
	for logEntry := range broadcast {
		data, _ := json.Marshal(logEntry)

		clientsMu.Lock()
		for ws := range clients {
			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				ws.Close()
				delete(clients, ws)
			}
		}
		clientsMu.Unlock()
	}
}

// PushLog pushes a log entry to the broadcast channel
func PushLog(logg model.Log) {
	select {
	case broadcast <- logg:
	default:
		log.Println("broadcast channel full, dropping log")
	}
}

// InitializeAPIServer sets up the API server with routes and middleware
func InitializeAPIServer(addr, port string) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler)

	// Routes
	chiRouter.Route("/redis", func(r chi.Router) {
		r.Get("/player-events", handlers.GetAllRedisPlayerEventsHandler)
		r.Get("/kill-events", handlers.GetAllRedisKillEventsHandler)
		r.Get("/cache-size", handlers.GetCacheSizeHandler)
		r.Delete("/clear", handlers.ClearCacheHandler)
	})

	chiRouter.Route("/db", func(r chi.Router) {
		r.Get("/kill-events", handlers.GetAllKillEventsHandler)
		r.Get("/player-events", handlers.GetAllPlayerEventsHandler)
	})

	// WebSocket endpoint
	chiRouter.Get("/ws", wsHandler)

	// Start broadcaster in background
	go broadcaster()

	httpClient = &http.Client{}
	return chiRouter
}
