package main

import (
	"fmt"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

func main() {
	// Test handlers to ensure events are being sent to the server
	cs2gsi.RegisterGlobalHandler(func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		fmt.Printf("(%v) %v %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			events.EnumToEventName[gameEvent.EventType],
			gsiEvent.GetOriginalRequestFlat(),
		)
	})
	cs2gsi.RegisterNonEventHandler(func(gsiEvent *structs.GSIEvent) {
		fmt.Printf("(%v) #N/A %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			gsiEvent.GetOriginalRequestFlat(),
		)
	})

	err := cs2gsi.StartupAndServe(":3000")
	if err != nil {
		return
	}
}
