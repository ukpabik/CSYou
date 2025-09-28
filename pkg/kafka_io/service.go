package kafka_io

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/ukpabik/CSYou/pkg/redis"
	"github.com/ukpabik/CSYou/pkg/shared"
)

// WritePlayerEvent writes player event to player_events topic
func WritePlayerEvent(event *shared.RedisPlayerEvent, key string) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return PlayerEventWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: eventBytes,
		},
	)
}

// WriteKillEvent writes kill event to kill_events topic
func WriteKillEvent(event *shared.RedisKillEvent, key string) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return KillEventWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: eventBytes,
		},
	)
}

// ReadPlayerEventLoop reads from player_events topic
func ReadPlayerEventLoop() {
	log.Println("Starting Kafka player event consumer loop...")
	for {
		message, err := PlayerEventReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading kafka player event message: %v", err)
			break
		}

		log.Printf("Received player event from Kafka (key: %s)", string(message.Key))

		var playerEvent shared.RedisPlayerEvent
		if err := json.Unmarshal(message.Value, &playerEvent); err != nil {
			log.Printf("failed to unmarshal player event: %v", err)
			continue
		}
		redis.HandlePlayerEvent(&playerEvent)

		// TODO: Process event into Clickhouse
		log.Printf("Processing player event for match %s, player %s", playerEvent.MatchID, playerEvent.SteamID)
	}
	log.Println("Kafka player event consumer loop ended")
}

// ReadKillEventLoop reads from kill_events topic
func ReadKillEventLoop() {
	log.Println("Starting Kafka kill event consumer loop...")
	for {
		message, err := KillEventReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading kafka kill event message: %v", err)
			break
		}

		log.Printf("Received kill event from Kafka (key: %s)", string(message.Key))

		var killEvent shared.RedisKillEvent
		if err := json.Unmarshal(message.Value, &killEvent); err != nil {
			log.Printf("failed to unmarshal kill event: %v", err)
			continue
		}

		redis.HandleKillEvent(&killEvent)

		// TODO: Process event into clickhouse
		log.Printf("Processing kill event for match %s, player %s with %s", killEvent.MatchID, killEvent.SteamID, killEvent.ActiveGun.Name)
	}
	log.Println("Kafka kill event consumer loop ended")
}
