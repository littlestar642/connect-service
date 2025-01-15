package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientI interface {
	SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	IsNil(err error) bool
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
	Close()
}

type client struct {
	RedisClient *redis.Client
}

func Init(addr string) (ClientI, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("connected to Redis")

	return &client{
		RedisClient: redisClient,
	}, nil
}

func (c *client) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return c.RedisClient.SIsMember(ctx, key, member)
}

func (c *client) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.RedisClient.Incr(ctx, key)
}

func (c *client) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.RedisClient.SAdd(ctx, key, members...)
}

func (c *client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.RedisClient.Get(ctx, key)
}

func (c *client) IsNil(err error) bool {
	return err == redis.Nil
}

func (c *client) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	return c.RedisClient.Expire(ctx, key, ttl)
}

func (c *client) Close() {
	err := c.RedisClient.Close()
	if err != nil {
		log.Println("failed to close redis connection: ", err.Error())
	}
}
