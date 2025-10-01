package model

type Log struct {
	EventType string `json:"event_type"`
	Time      string `json:"time"`
}

type ClickHouseKillEvent struct {
	MatchId string `ch:"match_id"`
	Round   uint32 `ch:"round"`
	Map     string `ch:"map"`
	Team    string `ch:"team"`
	SteamID string `ch:"steamid"`
	Name    string `ch:"name"`
	Mode    string `ch:"mode"`

	WeaponName     string `ch:"weapon_name"`
	WeaponType     string `ch:"weapon_type"`
	WeaponAmmo     uint32 `ch:"weapon_ammo"`
	WeaponReserve  uint32 `ch:"weapon_reserve"`
	WeaponSkin     string `ch:"weapon_skin"`
	WeaponHeadshot bool   `ch:"weapon_headshot"`

	Timestamp int64 `ch:"timestamp"`
}

type ClickHousePlayerEvent struct {
	MatchId string `ch:"match_id"`
	Round   uint32 `ch:"round"`
	Map     string `ch:"map"`
	Team    string `ch:"team"`
	SteamID string `ch:"steamid"`
	Name    string `ch:"name"`
	Mode    string `ch:"mode"`

	Health      uint32 `ch:"health"`
	Armor       uint32 `ch:"armor"`
	Helmet      bool   `ch:"helmet"`
	Money       uint32 `ch:"money"`
	EquipValue  uint32 `ch:"equip_value"`
	RoundKills  uint32 `ch:"round_kills"`
	RoundKillHS uint32 `ch:"round_killhs"`
	Kills       uint32 `ch:"kills"`
	Assists     uint32 `ch:"assists"`
	Deaths      uint32 `ch:"deaths"`
	MVPs        uint32 `ch:"mvps"`
	Score       uint32 `ch:"score"`

	EventTS int64  `ch:"event_timestamp"`
	WinTeam string `ch:"win_team"`
}
