package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ukpabik/CSYou/pkg/api/handlers"
)

var httpClient *http.Client

func InitializeAPIServer(addr, port string) *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Route("/redis", func(r chi.Router) {
		r.Get("/player-events", handlers.GetAllPlayerEventsHandler)
		r.Get("/kill-events", handlers.GetAllKillEventsHandler)
	})

	httpClient = &http.Client{}
	// Add more routes
	return chiRouter
}
