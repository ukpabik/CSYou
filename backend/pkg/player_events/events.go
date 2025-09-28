package player_events

import (
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/ukpabik/CSYou/pkg/shared"
)

// Track last known kills per player
var lastKills = make(map[string]int)

// DetectKillEvents checks for kill deltas and emits RedisKillEvent(s).
func DetectKillEvents(matchID string, event *structs.GSIEvent) []*shared.RedisKillEvent {
	killsNow := event.Player.MatchStats.Kills
	steamid := event.Player.Steamid

	// Get last recorded kills (default 0 if not tracked yet)
	prevKills := lastKills[steamid]

	if killsNow <= prevKills {
		return nil
	}

	// Find active gun
	active := shared.ActiveGun{}
	for _, w := range event.Player.Weapons {
		if w.State == "active" {
			active = shared.ActiveGun{
				Name:     w.Name,
				Type:     string(w.Type),
				Skin:     w.Paintkit,
				Headshot: event.Player.State.RoundKillHS > 0,
			}
			if w.AmmoClip != nil {
				active.Ammo = *w.AmmoClip
			}
			if w.AmmoReserve != nil {
				active.Reserve = *w.AmmoReserve
			}
			break
		}
	}

	// Generate kill events for each new kill
	var killEvents []*shared.RedisKillEvent
	for i := prevKills + 1; i <= killsNow; i++ {
		killEvents = append(killEvents, &shared.RedisKillEvent{
			MatchID:   matchID,
			Round:     event.CSMap.Round,
			Map:       event.CSMap.Name,
			Team:      event.Player.Team,
			SteamID:   steamid,
			Name:      event.Player.Name,
			Mode:      event.CSMap.Mode,
			ActiveGun: active,
			Timestamp: time.Now().Unix(),
		})
	}

	// Update lastKills cache
	lastKills[steamid] = killsNow

	return killEvents
}
