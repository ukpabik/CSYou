package redis

import "github.com/redis/go-redis/v9"

var RedisClient *redis.Client

func InitializeRedisClient(addr string) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password used
		DB:       0,  // Using Default DB
		Protocol: 2,  // Using the connection Protocol
	})

	RedisClient = client
}
