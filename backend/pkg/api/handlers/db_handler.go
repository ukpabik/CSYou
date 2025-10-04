package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ukpabik/CSYou/pkg/api/model"
	"github.com/ukpabik/CSYou/pkg/db"
)

func GetPlayerEventsByParamsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	matchId := queryParams.Get("match_id")
	roundStr := queryParams.Get("round")

	var roundInt int
	if roundStr != "" {
		if val, err := strconv.Atoi(roundStr); err == nil {
			roundInt = val
		}
	}

	playerOptions := []model.QueryOption{
		model.WithMatchID(matchId),
	}

	if roundInt != 0 {
		playerOptions = append(playerOptions, model.WithRound(roundInt))
	}

	// Build config
	paramConfig := model.NewClickHouseEventQueryConfig(playerOptions)

	events, err := db.GetPlayerEventsByParams(*paramConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get player events from db: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetKillEventsByParamsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	matchId := queryParams.Get("match_id")
	headshot := queryParams.Get("headshot")
	weaponName := queryParams.Get("weapon_name")

	playerOptions := []model.QueryOption{
		model.WithMatchID(matchId),
	}
	roundStr := queryParams.Get("round")

	var roundInt int
	if roundStr != "" {
		if val, err := strconv.Atoi(roundStr); err == nil {
			roundInt = val
		}
	}

	if roundInt != 0 {
		playerOptions = append(playerOptions, model.WithRound(roundInt))
	}
	headshotBool, _ := strconv.ParseBool(headshot)
	killOptions := []model.KillQueryOption{
		model.WithWeaponHeadshot(headshotBool),
		model.WithWeaponName(weaponName),
	}

	// Check for all types of params and build config
	paramConfig := model.NewClickHouseKillEventQueryConfig(playerOptions, killOptions)

	events, err := db.GetKillEventsByParams(*paramConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get kill events from db: %v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

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
