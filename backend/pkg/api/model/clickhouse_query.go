package model

import "fmt"

// Column Name, Value
type QueryField[T any] struct {
	Value *T
	Name  string
}

func (q QueryField[T]) IsSet() bool {
	return q.Value != nil
}

func (q QueryField[T]) String() string {
	if q.Value == nil {
		return ""
	}

	switch v := any(*q.Value).(type) {
	case string:
		if v == "" {
			return ""
		}
		return fmt.Sprintf("%s = '%v'", q.Name, v)
	case bool:
		// Only include if true, skip false
		if !v {
			return ""
		}
		return fmt.Sprintf("%s = %v", q.Name, v)
	default:
		return fmt.Sprintf("%s = %v", q.Name, v)
	}
}

// Configs

type ClickHouseEventQueryConfig struct {
	Round   QueryField[int]
	MatchID QueryField[string]
}

type ClickHouseKillEventQueryConfig struct {
	ClickHouseEventQueryConfig
	WeaponName     QueryField[string]
	WeaponHeadshot QueryField[bool]
}

// Constructors

func NewClickHouseEventQueryConfig(options []QueryOption) *ClickHouseEventQueryConfig {
	config := &ClickHouseEventQueryConfig{
		Round:   QueryField[int]{Name: "round"},
		MatchID: QueryField[string]{Name: "match_id"},
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

func NewClickHouseKillEventQueryConfig(playerOptions []QueryOption, killOptions []KillQueryOption) *ClickHouseKillEventQueryConfig {
	config := &ClickHouseKillEventQueryConfig{
		ClickHouseEventQueryConfig: ClickHouseEventQueryConfig{
			Round:   QueryField[int]{Name: "round"},
			MatchID: QueryField[string]{Name: "match_id"},
		},
		WeaponName:     QueryField[string]{Name: "weapon_name"},
		WeaponHeadshot: QueryField[bool]{Name: "weapon_headshot"},
	}

	for _, opt := range playerOptions {
		opt(&config.ClickHouseEventQueryConfig)
	}
	for _, opt := range killOptions {
		opt(config)
	}
	return config
}

// Options

type QueryOption func(*ClickHouseEventQueryConfig)
type KillQueryOption func(*ClickHouseKillEventQueryConfig)

func WithRound(round int) QueryOption {
	return func(c *ClickHouseEventQueryConfig) {
		c.Round.Value = &round
	}
}

func WithMatchID(matchID string) QueryOption {
	return func(c *ClickHouseEventQueryConfig) {
		c.MatchID.Value = &matchID
	}
}

func WithWeaponName(weaponName string) KillQueryOption {
	return func(c *ClickHouseKillEventQueryConfig) {
		c.WeaponName.Value = &weaponName
	}
}

func WithWeaponHeadshot(weaponHeadshot bool) KillQueryOption {
	return func(c *ClickHouseKillEventQueryConfig) {
		c.WeaponHeadshot.Value = &weaponHeadshot
	}
}
