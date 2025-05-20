package storage

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
)

type RedisClient struct {
    *redis.Client
}

func NewRedisClient(addr, password string) *RedisClient {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })
    return &RedisClient{client}
}

func (r *RedisClient) PublishMessage(ctx context.Context, message interface{}) error {
    msgBytes, err := json.Marshal(message)
    if err != nil {
        return err
    }
    return r.Publish(ctx, "messages", msgBytes).Err()
}