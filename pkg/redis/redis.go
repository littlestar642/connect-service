package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Init(addr string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Redis")

	return redisClient, nil
}
