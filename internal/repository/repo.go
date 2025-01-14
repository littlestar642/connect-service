package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	RedisClient *redis.Client
}

func New(redisClient *redis.Client) *Repo {
	return &Repo{
		RedisClient: redisClient,
	}
}

var countKey = "unique_request_count:%s"
var requestSetKey = "request_set:%s"

func (r *Repo) IsUniqueRequestId(ctx context.Context, id int) bool {
	timestamp := time.Now().Unix() / 60
	key := fmt.Sprintf(requestSetKey, timestamp)
	return r.RedisClient.SIsMember(ctx, key, id).Val()
}

func (r *Repo) IncrementRequestCount(ctx context.Context, id int) error {
	timestamp := time.Now().Unix() / 60
	countTsKey := fmt.Sprintf(countKey, timestamp)
	setKey := fmt.Sprintf(requestSetKey, timestamp)
	defer expireKeys(ctx, r.RedisClient, countTsKey, setKey)

	err := r.RedisClient.Incr(ctx, countTsKey).Err()
	if err != nil {
		return err
	}

	err = r.RedisClient.SAdd(ctx, setKey, id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetRequestCount(ctx context.Context) int {
	timestamp := time.Now().Unix() / 60
	countTsKey := fmt.Sprintf(countKey, timestamp)
	val, err := r.RedisClient.Get(ctx, countTsKey).Int()
	if err != nil {
		return 0
	}
	return val
}

func expireKeys(ctx context.Context, client *redis.Client, keys ...string) {
	for _, key := range keys {
		client.Expire(ctx, key, time.Minute)
	}
}
