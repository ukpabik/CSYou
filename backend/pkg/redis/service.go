package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ukpabik/CSYou/pkg/shared"
)

// HandlePlayerEvent parses the game event and stores it in Redis.
func HandlePlayerEvent(event *shared.RedisPlayerEvent) {
	if event == nil {
		return
	}
	ctx := context.Background()
	err := storePlayerEvent(ctx, event)
	if err != nil {
		log.Printf("failed to store player event: %v", err)
	}

}

// HandleKillEvent is the public function to process a kill event and store it.
func HandleKillEvent(event *shared.RedisKillEvent) {
	if event == nil {
		return
	}

	ctx := context.Background()
	err := storeKillEvent(ctx, event)
	if err != nil {
		log.Printf("failed to store kill event: %v", err)
	}
}

// storePlayerEvent is a helper function to store a player event into Redis.
func storePlayerEvent(ctx context.Context, event *shared.RedisPlayerEvent) error {
	// Check if the user is in game
	if event == nil {
		return fmt.Errorf("nil player event")
	}

	key := fmt.Sprintf("matches:%s:round:%d:player:%s:events",
		event.MatchID, event.Round, event.SteamID)

	_, err := RedisClient.JSONSet(ctx, key, ".", event).Result()
	if err != nil {
		return fmt.Errorf("unable to add event to Redis: %v", err)
	}

	return nil
}

func storeKillEvent(ctx context.Context, event *shared.RedisKillEvent) error {
	if event == nil {
		return fmt.Errorf("nil kill event")
	}

	key := fmt.Sprintf("matches:%s:round:%d:player:%s:kills",
		event.MatchID, event.Round, event.SteamID)

	_, err := RedisClient.JSONSet(ctx, key, ".", event).Result()
	if err != nil {
		return fmt.Errorf("unable to add kill event to Redis: %v", err)
	}
	return nil
}

func GetAllPlayerEvents(ctx context.Context) ([]RedisPlayerEvent, error) {
	var playerEvents []RedisPlayerEvent

	pattern := fmt.Sprintf("matches:*:round:*:player:%s:events", shared.PlayerID)
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		val, err := RedisClient.JSONGet(ctx, iter.Val()).Result()
		if err != nil {
			return nil, fmt.Errorf("jsonget failed for key %s: %w", iter.Val(), err)
		}
		if val == "" {
			continue
		}
		var ev RedisPlayerEvent
		if err := json.Unmarshal([]byte(val), &ev); err != nil {
			return nil, fmt.Errorf("unmarshal failed for key %s: %w", iter.Val(), err)
		}
		playerEvents = append(playerEvents, ev)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return playerEvents, nil
}

func GetAllKillEvents(ctx context.Context) ([]RedisKillEvent, error) {
	var killEvents []RedisKillEvent

	pattern := fmt.Sprintf("matches:*:round:*:player:%s:kills", shared.PlayerID)
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		val, err := RedisClient.JSONGet(ctx, iter.Val()).Result()
		if err != nil {
			return nil, fmt.Errorf("jsonget failed for key %s: %w", iter.Val(), err)
		}
		if val == "" {
			continue
		}
		var ev RedisKillEvent
		if err := json.Unmarshal([]byte(val), &ev); err != nil {
			return nil, fmt.Errorf("unmarshal failed for key %s: %w", iter.Val(), err)
		}
		killEvents = append(killEvents, ev)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return killEvents, nil
}
