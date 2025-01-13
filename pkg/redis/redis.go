package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Init(addr string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis %s", err)
	}
	log.Println("Connected to Redis")

	return redisClient
}
