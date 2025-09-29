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

func ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := redis.ClearCache(ctx)
	if err != nil {
		http.Error(w, "failed to clear cache", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("cache cleared successfully"))
}

func GetCacheSizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	memoryUsage, err := redis.GetCacheSize(ctx)
	if err != nil || memoryUsage == -1 {
		http.Error(w, "failed to get cache size", http.StatusInternalServerError)
		return
	}

	type MemoryObject struct {
		MemoryValue int64 `json:"memory_value"`
	}

	memoryObject := &MemoryObject{
		MemoryValue: memoryUsage,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memoryObject)
}
