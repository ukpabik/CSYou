package gsi

import (
	"fmt"
	"log"
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

// InitializeEventHandlers initializes all event handlers for every captured event.
func InitializeEventHandlers() {
	cs2gsi.RegisterGlobalHandler(func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		fmt.Printf("(%v) %v %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			events.EnumToEventName[gameEvent.EventType],
			gsiEvent.GetOriginalRequestFlat(),
		)

		if shared.PlayerID == "" {
			shared.PlayerID = gsiEvent.Player.Steamid
		}

		if gsiEvent.Round == nil || gsiEvent.Player == nil {
			return
		}

		eventLog := &model.Log{
			EventType: events.EnumToEventName[gameEvent.EventType],
			Time:      time.Now().Format("2006-01-02 15:04:05.000"),
		}

		// Send log to frontend
		if err := api.LogSender(*eventLog, shared.ADDRESS, shared.FRONTEND_PORT); err != nil {
			fmt.Printf("failed to send log to frontend: %v", err)
		}

		if gsiEvent.CSMap.Name != shared.LastMap ||
			(gsiEvent.CSMap.Round == 1 && shared.LastRound >= 1) {
			shared.CurrentMatchID = uuid.New().String()
			shared.LastMap = gsiEvent.CSMap.Name
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
			if err := api.LogSender(*killEventLog, shared.ADDRESS, shared.FRONTEND_PORT); err != nil {
				fmt.Printf("failed to send log to frontend: %v", err)
			}
			if err := kafka_io.WriteKillEvent(ke, gsiEvent.Player.Steamid); err != nil {
				log.Printf("failed to write kill event to kafka: %v", err)
			}
		}

		// Track last round
		shared.LastRound = gsiEvent.CSMap.Round
	})

	cs2gsi.RegisterNonEventHandler(func(gsiEvent *structs.GSIEvent) {
		fmt.Printf("(%v) #N/A %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			gsiEvent.GetOriginalRequestFlat(),
		)
	})
}

// Listen starts up the GSI server to listen for POST requests (with event data).
func Listen() {
	if err := cs2gsi.StartupAndServe(fmt.Sprintf(":%d", GSI_PORT)); err != nil {
		panic("FAILED TO START SERVER")
	}
}
