package repo

import (
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"github.com/minhmannh2001/authconnecthub/pkg/redis"
)

type AuthRepo struct {
	*postgres.Postgres
	*redis.Redis
}

func NewAuthRepo(pg *postgres.Postgres, redis *redis.Redis) *AuthRepo {
	return &AuthRepo{pg, redis}
}
