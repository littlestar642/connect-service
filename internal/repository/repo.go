package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RepoI interface {
	IsUniqueRequestId(ctx context.Context, id int) bool
	IncrementRequestCount(ctx context.Context, id int) error
	GetLastMinuteRequestCount(ctx context.Context) (int, error)
	GetCurrentMinuteRequestCount(ctx context.Context) (int, error)
}

type repo struct {
	RedisClient *redis.Client
}

func New(redisClient *redis.Client) RepoI {
	return &repo{
		RedisClient: redisClient,
	}
}

const (
	TIME_FORMAT     = "2006-01-02 15:04"
	COUNT_KEY       = "unique_request_count:%s"
	REQUEST_SET_KEY = "request_set:%s"
)

func (r *repo) IsUniqueRequestId(ctx context.Context, id int) bool {
	timestamp := time.Now().Format(TIME_FORMAT)
	key := fmt.Sprintf(REQUEST_SET_KEY, timestamp)
	return !r.RedisClient.SIsMember(ctx, key, id).Val()
}

func (r *repo) IncrementRequestCount(ctx context.Context, id int) error {
	timestamp := time.Now().Format(TIME_FORMAT)
	countTsKey := fmt.Sprintf(COUNT_KEY, timestamp)
	setKey := fmt.Sprintf(REQUEST_SET_KEY, timestamp)
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

func (r *repo) GetLastMinuteRequestCount(ctx context.Context) (int, error) {
	currentTime := time.Now()
	prevMinute := currentTime.Add(-time.Minute)
	fomattedTime := prevMinute.Format(TIME_FORMAT)
	countTsKey := fmt.Sprintf(COUNT_KEY, fomattedTime)
	val, err := r.RedisClient.Get(ctx, countTsKey).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Println("GetLastMinuteRequestCount: failed to get last minute request count: ", err)
		return 0, err
	}
	
	return val, nil
}

func (r *repo) GetCurrentMinuteRequestCount(ctx context.Context) (int, error) {
	fomattedTime := time.Now().Format(TIME_FORMAT)
	countTsKey := fmt.Sprintf(COUNT_KEY, fomattedTime)
	val, err := r.RedisClient.Get(ctx, countTsKey).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Println("GetCurrentMinuteRequestCount: failed to get current minute request count: ", err)
		return 0, err
	}

	return val, nil
}

func expireKeys(ctx context.Context, client *redis.Client, keys ...string) {
	for _, key := range keys {
		client.Expire(ctx, key, time.Minute)
	}
}
