package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/google/uuid"
	"github.com/ukpabik/CSYou/pkg/player_events"
	"github.com/ukpabik/CSYou/pkg/shared"
)

// BundleEvent maps an event to Redis JSON structure.
func bundlePlayerEvent(event *structs.GSIEvent, _ *events.GameEventDetails) *RedisPlayerEvent {

	redisEvent := &RedisPlayerEvent{
		// Match Information
		MatchID: shared.CurrentMatchID,
		Round:   event.CSMap.Round,
		Map:     event.CSMap.Name,
		Team:    event.Player.Team,
		SteamID: event.Player.Steamid,
		Name:    event.Player.Name,
		Mode:    event.CSMap.Mode,

		// Player State
		Health:     *event.Player.State.Health,
		Armor:      *event.Player.State.Armor,
		Helmet:     event.Player.State.Helmet,
		Money:      event.Player.State.Money,
		EquipValue: event.Player.State.EquipValue,

		// Per-round stats
		RoundKills:  event.Player.State.RoundKills,
		RoundKillHS: event.Player.State.RoundKillHS,

		// Match stats
		Kills:   event.Player.MatchStats.Kills,
		Assists: event.Player.MatchStats.Assists,
		Deaths:  event.Player.MatchStats.Deaths,
		MVPs:    event.Player.MatchStats.Mvps,
		Score:   event.Player.MatchStats.Score,

		// Match Context
		EventTS: int64(event.Provider.Timestamp),
		WinTeam: event.Round.WinTeam,
	}

	return redisEvent
}

func HandlePlayerEvent(event *structs.GSIEvent, gameDetails *events.GameEventDetails) {
	if event.Round == nil || event.Player == nil {
		return
	}
	ctx := context.Background()
	if event.CSMap.Name != shared.LastMap || event.CSMap.Round == 1 && shared.LastRound > 1 {
		shared.CurrentMatchID = uuid.New().String()
		shared.LastMap = event.CSMap.Name
		log.Printf("New match started: %s on map %s", shared.CurrentMatchID, shared.LastMap)
	}

	err := storePlayerEvent(ctx, event, gameDetails)
	if err != nil {
		log.Printf("%v", err)
	}

	killEvents := player_events.DetectKillEvents(shared.CurrentMatchID, event)
	for _, ke := range killEvents {
		key := fmt.Sprintf("matches:%s:round:%d:kills", ke.MatchID, ke.Round)
		_, err := RedisClient.JSONSet(ctx, key, ".", ke).Result()
		if err != nil {
			log.Printf("failed to store kill event: %v", err)
		}
	}

	shared.LastRound = event.CSMap.Round
}

func storePlayerEvent(ctx context.Context, event *structs.GSIEvent, eventDetails *events.GameEventDetails) error {
	// Check if the user is in game
	if event.Round == nil || event.Player == nil {
		return fmt.Errorf("user not in game")
	}

	redisEvent := bundlePlayerEvent(event, eventDetails)

	key := fmt.Sprintf("matches:%s:round:%d:player:%s",
		shared.CurrentMatchID, redisEvent.Round, redisEvent.SteamID)

	_, err := RedisClient.JSONSet(ctx, key, ".", redisEvent).Result()
	if err != nil {
		return fmt.Errorf("unable to add event to Redis: %v", err)
	}

	return nil
}
