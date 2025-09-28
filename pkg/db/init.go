package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ClickHouseClient driver.Conn

func InitializeClickHouseClient(addr string, port int) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", addr, port)},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Debugf: func(format string, v ...interface{}) {
			log.Printf(format, v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      time.Second * 30,
		MaxOpenConns:     10,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	})

	if err != nil {
		log.Printf("unable to initialize connection to ClickHouse: %v", err)
		return
	}

	// Test the connection
	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		log.Printf("unable to ping ClickHouse: %v", err)
		return
	}

	log.Println("Connected to ClickHouse client")
	ClickHouseClient = conn
}

func CloseClickHouseConnection() {
	if err := ClickHouseClient.Close(); err != nil {
		log.Printf("unable to close clickhouse client: %v", err)
	}
}
