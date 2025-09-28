package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukpabik/CSYou/pkg/api"
	"github.com/ukpabik/CSYou/pkg/db"
	"github.com/ukpabik/CSYou/pkg/gsi"
	"github.com/ukpabik/CSYou/pkg/kafka_io"
	"github.com/ukpabik/CSYou/pkg/redis"
	"github.com/ukpabik/CSYou/pkg/shared"
)

func setupGracefulShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")

		// Close Kafka connections
		kafka_io.CloseReaderAndWriters()

		// Close ClickHouse connection
		db.CloseClickHouseConnection()

		log.Println("Shutdown complete")
		os.Exit(0)
	}()
}

func main() {

	// Initialize Redis client for hot queries
	redis.InitializeRedisClient(fmt.Sprintf("%s:%d", shared.ADDRESS, shared.REDIS_PORT))

	// Initialize ClickHouse Client
	db.InitializeClickHouseClient(shared.ADDRESS, shared.CLICKHOUSE_PORT)

	if err := db.CreateTables(); err != nil {
		log.Fatalf("unable to create clickhouse tables: %v", err)
	}

	// Initialize Kafka Reader and Writer, and ensure graceful shutdown
	kafka_io.InitializeReaderAndWriter(shared.ADDRESS, shared.KAFKA_PORT)
	setupGracefulShutdown()
	// Run kafka reading in a goroutine
	go kafka_io.ReadPlayerEventLoop()
	go kafka_io.ReadKillEventLoop()

	// Listen for events from CS2 GSI
	gsi.InitializeEventHandlers()
	gsi.Listen()

	// Initialize API server
	http.ListenAndServe(shared.API_PORT, api.InitializeAPIServer(shared.ADDRESS, shared.API_PORT))
}
