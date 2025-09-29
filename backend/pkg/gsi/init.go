package gsi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/google/uuid"
	"github.com/ukpabik/CSYou/pkg/api"
	"github.com/ukpabik/CSYou/pkg/api/model"
	"github.com/ukpabik/CSYou/pkg/kafka_io"
	"github.com/ukpabik/CSYou/pkg/player_events"
	"github.com/ukpabik/CSYou/pkg/shared"
)

const GSI_PORT = 3000

// Config struct for reading config.json
type Config struct {
	SteamID string `json:"steam_id"`
}

var STEAM_ID string

// LoadConfig loads the steam_id from config.json
func LoadConfig() {
	file, err := os.Open("../config.json")
	if err != nil {
		log.Fatalf("failed to open config.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("failed to parse config.json: %v", err)
	}

	if config.SteamID == "" {
		log.Fatalf("steam_id must be set in config.json")
	}

	STEAM_ID = config.SteamID
	log.Printf("Loaded SteamID from config: %s", STEAM_ID)
}

// InitializeEventHandlers initializes all event handlers for every captured event.
//
// NOTE: Make sure you update your `config.json` with your actual SteamID
// before running this, otherwise no events will be tracked.
func InitializeEventHandlers() {
	if STEAM_ID == "" {
		LoadConfig()
	}

	cs2gsi.RegisterGlobalHandler(func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		if shared.PlayerID == "" {
			shared.PlayerID = STEAM_ID
		}

		if gsiEvent.Player == nil || gsiEvent.Player.Steamid != shared.PlayerID || gsiEvent.Round == nil {
			return
		}

		eventLog := &model.Log{
			EventType: events.EnumToEventName[gameEvent.EventType],
			Time:      time.Now().Format("2006-01-02 15:04:05.000"),
		}

		// Send log to frontend
		api.PushLog(*eventLog)

		if gsiEvent.CSMap.Name != shared.LastMap || gsiEvent.Round.Phase == "gameover" {
			shared.CurrentMatchID = uuid.New().String()
			shared.LastMap = gsiEvent.CSMap.Name
			shared.LastRound = 0
			log.Printf("New match started: %s on map %s", shared.CurrentMatchID, shared.LastMap)
		}

		playerEvent := shared.BundlePlayerEvent(gsiEvent, &gameEvent)

		// Publish player event to Kafka
		if err := kafka_io.WritePlayerEvent(playerEvent, gsiEvent.Player.Steamid); err != nil {
			log.Printf("failed to write player event to kafka: %v", err)
		}

		killEvents := player_events.DetectKillEvents(shared.CurrentMatchID, gsiEvent)
		for _, ke := range killEvents {
			if ke.ActiveGun.Type == "C4" {
				continue
			}

			killEventLog := &model.Log{
				EventType: "Player Kill",
				Time:      time.Now().Format("2006-01-02 15:04:05.000"),
			}

			// Send log to frontend
			api.PushLog(*killEventLog)
			if err := kafka_io.WriteKillEvent(ke, gsiEvent.Player.Steamid); err != nil {
				log.Printf("failed to write kill event to kafka: %v", err)
			}
		}

		// Track last round
		shared.LastRound = gsiEvent.CSMap.Round
	})

	cs2gsi.RegisterNonEventHandler(func(gsiEvent *structs.GSIEvent) {
		// fmt.Printf("(%v) #N/A %v\n",
		// 	time.Now().Format("2006-01-02 15:04:05.000"),
		// 	gsiEvent.GetOriginalRequestFlat(),
		// )
	})
}

// Listen starts up the GSI server to listen for POST requests (with event data).
func Listen() {
	if err := cs2gsi.StartupAndServe(fmt.Sprintf(":%d", GSI_PORT)); err != nil {
		log.Fatalf("failed to start GSI server: %v", err)
	}
}
