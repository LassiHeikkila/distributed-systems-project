package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectToRedis(ctx context.Context, addr string, pw string, db int) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "default",
		Password: pw,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	redisClient = rdb
	return nil
}

func SubmitJobToQueue(ctx context.Context, key string, job []byte) error {
	if redisClient == nil {
		return errors.New("not connected to redis")
	}

	err := redisClient.LPush(ctx, key, job).Err()
	if err != nil {
		return fmt.Errorf("failed to submit redis job to queue: %w", err)
	}

	return nil
}

func SubmitToRedisAsJSON(ctx context.Context, key string, job any) error {
	b, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	return SubmitJobToQueue(ctx, key, b)
}
