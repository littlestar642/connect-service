package repository

import (
	"context"

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

func (r *Repo) IsUniqueRequestId(ctx context.Context, id int) bool {
	return r.RedisClient.SetNX(ctx, "id", id, 0).Val()
}

func (r *Repo) GetRequestCount(ctx context.Context) int {
	val, err := r.RedisClient.Get(ctx, "count").Int()
	if err != nil {
		return 0
	}
	return val
}
