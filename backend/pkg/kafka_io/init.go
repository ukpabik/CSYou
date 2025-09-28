package kafka_io

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	PLAYER_EVENT_TOPIC = "player_events"
	KILL_EVENT_TOPIC   = "kill_events"
)

var (
	PlayerEventWriter *kafka.Writer
	KillEventWriter   *kafka.Writer
	PlayerEventReader *kafka.Reader
	KillEventReader   *kafka.Reader
)

func InitializeReaderAndWriter(addr string, port int) {

	location := fmt.Sprintf("%s:%d", addr, port)
	// Writers
	PlayerEventWriter = &kafka.Writer{
		Addr:     kafka.TCP(location),
		Topic:    PLAYER_EVENT_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	KillEventWriter = &kafka.Writer{
		Addr:     kafka.TCP(location),
		Topic:    KILL_EVENT_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	// Readers
	PlayerEventReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{location},
		Topic:   PLAYER_EVENT_TOPIC,
		GroupID: "cs2-player-processor",
	})

	KillEventReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{location},
		Topic:   KILL_EVENT_TOPIC,
		GroupID: "cs2-kill-processor",
	})
}

func CloseReaderAndWriters() {
	writers := []*kafka.Writer{PlayerEventWriter, KillEventWriter}
	for _, writer := range writers {
		if writer != nil {
			if err := writer.Close(); err != nil {
				log.Printf("failed to close writer: %v", err)
			}
		}
	}

	readers := []*kafka.Reader{PlayerEventReader, KillEventReader}
	for _, reader := range readers {
		if reader != nil {
			if err := reader.Close(); err != nil {
				log.Printf("failed to close reader: %v", err)
			}
		}
	}
}
