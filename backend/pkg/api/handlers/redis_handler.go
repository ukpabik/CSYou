package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ukpabik/CSYou/pkg/redis"
)

func GetAllPlayerEventsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := redis.GetAllPlayerEvents(ctx)
	if err != nil {
		http.Error(w, "Failed to get player events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func GetAllKillEventsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := redis.GetAllKillEvents(ctx)
	if err != nil {
		http.Error(w, "Failed to get kill events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
