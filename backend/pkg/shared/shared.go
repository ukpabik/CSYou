package shared

import (
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

// Variables for current match information.
var CurrentMatchID string
var LastMap string
var LastRound int

const (
	REDIS_PORT      = 6379
	CLICKHOUSE_PORT = 9000
	KAFKA_PORT      = 9092
	API_PORT        = "8080"
	FRONTEND_PORT   = 1420
	ADDRESS         = "localhost"
)

// Wrapper for Kafka event handling
type EventWrapper struct {
	GSIEvent    *structs.GSIEvent        `json:"gsi_event"`
	GameDetails *events.GameEventDetails `json:"game_details"`
	Timestamp   int64                    `json:"timestamp"`
}

type RedisPlayerEvent struct {
	MatchID string `json:"match_id"` // UUID you generate
	Round   int    `json:"round"`    // current round
	Map     string `json:"map"`      // map name (e.g., de_dust2)
	Team    string `json:"team"`     // "T" or "CT"
	SteamID string `json:"steamid"`  // player steamid
	Name    string `json:"name"`     // player name
	Mode    string `json:"mode"`     // gamemode

	// Player state
	Health     int  `json:"health"`
	Armor      int  `json:"armor"`
	Helmet     bool `json:"helmet"`
	Money      int  `json:"money"`
	EquipValue int  `json:"equip_value"`

	// Per-round stats
	RoundKills  int `json:"round_kills"`
	RoundKillHS int `json:"round_killhs"`

	// Match stats (cumulative)
	Kills   int `json:"kills"`
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
	MVPs    int `json:"mvps"`
	Score   int `json:"score"`

	// Context
	EventTS int64  `json:"timestamp"` // from provider.timestamp
	WinTeam string `json:"win_team"`  // who won round (T or CT)
}

type RedisKillEvent struct {
	MatchID string `json:"match_id"` // UUID you generate
	Round   int    `json:"round"`    // current round
	Map     string `json:"map"`      // map name (e.g., de_dust2)
	Team    string `json:"team"`     // "T" or "CT"
	SteamID string `json:"steamid"`  // player steamid
	Name    string `json:"name"`     // player name
	Mode    string `json:"mode"`     // gamemode

	ActiveGun ActiveGun `json:"active_gun"` // weapon details

	Timestamp int64 `json:"timestamp"` // provider timestamp
}

type ActiveGun struct {
	Name     string `json:"name"`     // weapon_ak47, weapon_glock, etc.
	Type     string `json:"type"`     // Rifle, Pistol, Knife, C4
	Ammo     int    `json:"ammo"`     // clip at time of kill
	Reserve  int    `json:"reserve"`  // reserve ammo at time of kill
	Skin     string `json:"skin"`     // paintkit / skin
	Headshot bool   `json:"headshot"` // true if the kill was HS
}

// BundleEvent maps an event to Redis JSON structure.
func BundlePlayerEvent(event *structs.GSIEvent, _ *events.GameEventDetails) *RedisPlayerEvent {

	redisEvent := &RedisPlayerEvent{
		// Match Information
		MatchID: CurrentMatchID,
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
