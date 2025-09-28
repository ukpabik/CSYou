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
	KAFKA_PORT  = 9092
	EVENT_TOPIC = "cs2_events"
)

var (
	EventWriter *kafka.Writer
	EventReader *kafka.Reader
)

func InitializeReaderAndWriter() {
	EventWriter = &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("localhost:%d", KAFKA_PORT)),
		Topic:    EVENT_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	EventReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   EVENT_TOPIC,
	})
}

func CloseWriter() {
	if err := EventReader.Close(); err != nil {
		log.Printf("failed to close writer: %v", err)
	}
}

func SetupGracefulShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("shutting down Kafka writer...")
		CloseWriter()
		os.Exit(0)
	}()
}
