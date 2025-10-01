package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ukpabik/CSYou/pkg/db"
)

func GetAllKillEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := db.GetAllKillEvents()
	if err != nil {
		http.Error(w, "Failed to get kill events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetAllPlayerEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := db.GetAllPlayerEvents()
	if err != nil {
		http.Error(w, "Failed to get player events", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
