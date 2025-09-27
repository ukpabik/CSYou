package gsi

import (
	"fmt"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/ukpabik/CSYou/pkg/redis"
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
		redis.HandlePlayerEvent(gsiEvent, &gameEvent)
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
