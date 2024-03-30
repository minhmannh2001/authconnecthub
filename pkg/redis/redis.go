package redis

import (
	"context"
	"fmt"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(cfg *config.Config, opts ...Option) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	err := ping(client)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{
		Client: client,
	}, nil
}

func ping(client *redis.Client) error {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
