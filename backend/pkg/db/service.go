package db

import (
	"context"
	"fmt"

	"github.com/ukpabik/CSYou/pkg/api/model"
	"github.com/ukpabik/CSYou/pkg/shared"
)

const (
	killEventTableName   = "cs2_kill_events"
	playerEventTableName = "cs2_player_events"
)

func GetAllKillEvents() ([]model.ClickHouseKillEvent, error) {
	if ClickHouseClient == nil {
		return nil, fmt.Errorf("clickhouse client is not initialized")
	}
	context := context.Background()
	var events []model.ClickHouseKillEvent
	query := fmt.Sprintf("SELECT * FROM %s", killEventTableName)

	err := ClickHouseClient.Select(context, &events, query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return events, nil
}

func GetAllPlayerEvents() ([]model.ClickHousePlayerEvent, error) {
	if ClickHouseClient == nil {
		return nil, fmt.Errorf("clickhouse client is not initialized")
	}
	context := context.Background()
	var events []model.ClickHousePlayerEvent
	query := fmt.Sprintf("SELECT * FROM %s", playerEventTableName)
	err := ClickHouseClient.Select(context, &events, query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return events, nil
}

// InsertKillEvents inserts multiple kill events using batch operation
func InsertKillEvents(killEvents []shared.RedisKillEvent) error {
	if ClickHouseClient == nil {
		return fmt.Errorf("clickhouse client is not initialized")
	}

	if len(killEvents) == 0 {
		return nil
	}

	ctx := context.Background()

	// Prepare batch insert
	batch, err := ClickHouseClient.PrepareBatch(ctx, fmt.Sprintf(`
        INSERT INTO %s (
            match_id, round, map, team, steamid, name, mode,
            weapon_name, weapon_type, weapon_ammo, weapon_reserve, weapon_skin, weapon_headshot,
            timestamp
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, killEventTableName))

	if err != nil {
		return fmt.Errorf("unable to prepare batch statement: %v", err)
	}

	// Append all kill events to batch
	for _, event := range killEvents {
		err = batch.Append(
			event.MatchID,
			event.Round,
			event.Map,
			event.Team,
			event.SteamID,
			event.Name,
			event.Mode,
			event.ActiveGun.Name,
			event.ActiveGun.Type,
			event.ActiveGun.Ammo,
			event.ActiveGun.Reserve,
			event.ActiveGun.Skin,
			event.ActiveGun.Headshot,
			event.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to append kill event to batch: %v", err)
		}
	}

	// Execute batch
	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to execute batch insert for kill events: %v", err)
	}

	return nil
}

// InsertKillEvent inserts a single kill event
func InsertKillEvent(killEvent *shared.RedisKillEvent) error {
	return InsertKillEvents([]shared.RedisKillEvent{*killEvent})
}

// InsertPlayerEvents inserts multiple player events using batch operation
func InsertPlayerEvents(playerEvents []shared.RedisPlayerEvent) error {
	if ClickHouseClient == nil {
		return fmt.Errorf("clickhouse client is not initialized")
	}

	if len(playerEvents) == 0 {
		return nil
	}

	ctx := context.Background()

	// Prepare batch insert
	batch, err := ClickHouseClient.PrepareBatch(ctx, fmt.Sprintf(`
        INSERT INTO %s (
            match_id, round, map, team, steamid, name, mode,
            health, armor, helmet, money, equip_value,
            round_kills, round_killhs,
            kills, assists, deaths, mvps, score,
            event_timestamp, win_team
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, playerEventTableName))

	if err != nil {
		return fmt.Errorf("unable to prepare batch statement: %v", err)
	}

	// Append all player events to batch
	for _, event := range playerEvents {
		err = batch.Append(
			event.MatchID,
			event.Round,
			event.Map,
			event.Team,
			event.SteamID,
			event.Name,
			event.Mode,
			event.Health,
			event.Armor,
			event.Helmet,
			event.Money,
			event.EquipValue,
			event.RoundKills,
			event.RoundKillHS,
			event.Kills,
			event.Assists,
			event.Deaths,
			event.MVPs,
			event.Score,
			event.EventTS,
			event.WinTeam,
		)
		if err != nil {
			return fmt.Errorf("failed to append player event to batch: %v", err)
		}
	}

	// Execute batch
	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to execute batch insert for player events: %v", err)
	}

	return nil
}

// InsertPlayerEvent inserts a single player event
func InsertPlayerEvent(playerEvent *shared.RedisPlayerEvent) error {
	return InsertPlayerEvents([]shared.RedisPlayerEvent{*playerEvent})
}

// CreateTables creates the necessary ClickHouse tables
func CreateTables() error {
	if ClickHouseClient == nil {
		return fmt.Errorf("clickhouse client is not initialized")
	}

	ctx := context.Background()

	// Create kill events table
	killEventSchema := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            match_id String,
            round UInt32,
            map String,
            team String,
            steamid String,
            name String,
            mode String,
            weapon_name String,
            weapon_type String,
            weapon_ammo UInt32,
            weapon_reserve UInt32,
            weapon_skin String,
            weapon_headshot Bool,
            timestamp Int64
        ) ENGINE = MergeTree()
        ORDER BY (match_id, timestamp)
        PARTITION BY toDate(fromUnixTimestamp(timestamp))
    `, killEventTableName)

	if err := ClickHouseClient.Exec(ctx, killEventSchema); err != nil {
		return fmt.Errorf("failed to create kill events table: %v", err)
	}

	// Create player events table
	playerEventSchema := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            match_id String,
            round UInt32,
            map String,
            team String,
            steamid String,
            name String,
            mode String,
            health UInt32,
            armor UInt32,
            helmet Bool,
            money UInt32,
            equip_value UInt32,
            round_kills UInt32,
            round_killhs UInt32,
            kills UInt32,
            assists UInt32,
            deaths UInt32,
            mvps UInt32,
            score UInt32,
            event_timestamp Int64,
            win_team String
        ) ENGINE = MergeTree()
        ORDER BY (match_id, event_timestamp)
        PARTITION BY toDate(fromUnixTimestamp(event_timestamp))
    `, playerEventTableName)

	if err := ClickHouseClient.Exec(ctx, playerEventSchema); err != nil {
		return fmt.Errorf("failed to create player events table: %v", err)
	}

	return nil
}
