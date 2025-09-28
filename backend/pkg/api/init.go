package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var httpClient *http.Client

func InitializeAPIServer(addr, port string) *chi.Mux {
	chiRouter := chi.NewRouter()

	httpClient = &http.Client{}
	// Add more routes
	return chiRouter
}
