package kafka_io

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

const (
	KAFKA_PORT         = 9092
	PLAYER_EVENT_TOPIC = "player_events"
	KILL_EVENT_TOPIC   = "kill_events"
)

var (
	PlayerEventWriter *kafka.Writer
	KillEventWriter   *kafka.Writer
	PlayerEventReader *kafka.Reader
	KillEventReader   *kafka.Reader
)

func InitializeReaderAndWriter() {
	// Writers
	PlayerEventWriter = &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("localhost:%d", KAFKA_PORT)),
		Topic:    PLAYER_EVENT_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	KillEventWriter = &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("localhost:%d", KAFKA_PORT)),
		Topic:    KILL_EVENT_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	// Readers
	PlayerEventReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   PLAYER_EVENT_TOPIC,
		GroupID: "cs2-player-processor",
	})

	KillEventReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
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

func SetupGracefulShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("shutting down Kafka clients...")
		CloseReaderAndWriters()
		os.Exit(0)
	}()
}
