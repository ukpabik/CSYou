package kafka_io

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/ukpabik/CSYou/pkg/redis"
	"github.com/ukpabik/CSYou/pkg/shared"
)

func WriteEvent(event []byte, key string) error {
	return EventWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: event,
		},
	)
}

func ReadEventLoop() {
	log.Println("Starting Kafka consumer loop...")
	for {
		message, err := EventReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading kafka message: %v", err)
			break
		}

		log.Printf("Received message from Kafka (key: %s)", string(message.Key))

		// Unmarshal the event wrapper
		var eventWrapper shared.EventWrapper
		if err := json.Unmarshal(message.Value, &eventWrapper); err != nil {
			log.Printf("failed to unmarshal event wrapper: %v", err)
			continue
		}

		// Process the event and store in Redis
		redis.HandlePlayerEvent(eventWrapper.GSIEvent, eventWrapper.GameDetails)
	}
	log.Println("Kafka consumer loop ended")
}
