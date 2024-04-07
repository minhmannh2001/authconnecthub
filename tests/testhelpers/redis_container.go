package testhelpers

import (
	"context"
	"errors"
	"log"
	"net/url"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainer struct {
	*redis.RedisContainer
	ConnectionString string
}

func CreateRedisContainer(ctx context.Context) (*RedisContainer, error) {
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("redis:7.2.4"),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		log.Fatalf("failed to start redis container: %s", err)
		return nil, err
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get redis connection string: %s", err)
		return nil, err
	}

	return &RedisContainer{
		RedisContainer:   redisContainer,
		ConnectionString: connStr,
	}, nil
}

func (con RedisContainer) ExtractHost(connString string) (string, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	if host == "" {
		return "", errors.New("missing host in connection string")
	}

	return host, nil
}

func (con RedisContainer) ExtractPort(connString string) (string, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return "", err
	}

	portStr := u.Port()
	if portStr == "" {
		return "", errors.New("missing port in connection string")
	}

	return portStr, nil
}
