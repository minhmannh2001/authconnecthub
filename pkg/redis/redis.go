package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg config.Config, opts ...Option) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	err := ping(client)

	if err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) Set(key string, value interface{}, ttl int) error {
	ctx := context.Background()
	// Set the key-value pair with the specified TTL
	err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return errors.New("failed to set key-value pair: " + err.Error())
	}

	return nil
}

func ping(client *redis.Client) error {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
