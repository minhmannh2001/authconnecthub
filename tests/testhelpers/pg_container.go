package testhelpers

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.5-bullseye"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func (con PostgresContainer) ExtractHost(connString string) (string, error) {
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

func (con PostgresContainer) ExtractPort(connString string) (string, error) {
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
