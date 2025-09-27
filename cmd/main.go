package main

import (
	"fmt"

	"github.com/ukpabik/CSYou/pkg/gsi"
	"github.com/ukpabik/CSYou/pkg/redis"
)

const REDIS_PORT = 6379

func main() {

	// Initialize Redis client for hot queries
	redis.InitializeRedisClient(fmt.Sprintf("localhost:%d", REDIS_PORT))

	// Listen for events from CS2 GSI
	gsi.InitializeEventHandlers()
	gsi.Listen()
}
