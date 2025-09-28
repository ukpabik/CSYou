package gsi

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/ukpabik/CSYou/pkg/kafka_io"
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

		if gsiEvent.Round == nil || gsiEvent.Player == nil {
			return
		}

		// Create a wrapper with both GSI event and game event details
		eventWrapper := &shared.EventWrapper{
			GSIEvent:    gsiEvent,
			GameDetails: &gameEvent,
			Timestamp:   time.Now().Unix(),
		}

		convertedEvent, err := json.Marshal(eventWrapper)
		if err != nil {
			log.Println("unable to marshal event")
			return
		}

		// Write to Kafka with player steamid as key for partitioning
		key := gsiEvent.Player.Steamid
		if err := kafka_io.WriteEvent(convertedEvent, key); err != nil {
			log.Printf("failed to write event to kafka: %v", err)
		}
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
