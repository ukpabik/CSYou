package shared

import (
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

// Variables for current match information.
var CurrentMatchID string
var LastMap string
var LastRound int

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
