package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	redis_pkg "github.com/minhmannh2001/authconnecthub/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type AuthRepo struct {
	*postgres.Postgres
	*redis_pkg.Redis
}

func NewAuthRepo(pg *postgres.Postgres, redis *redis_pkg.Redis) *AuthRepo {
	return &AuthRepo{pg, redis}
}

func (a *AuthRepo) BlacklistToken(token string, expiration int) error {
	ctx := context.Background()
	err := a.Client.Set(ctx, "blacklist:"+token, "", time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (a *AuthRepo) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	val, err := a.Client.Get(ctx, "blacklist:"+token).Result()
	if err != nil {
		if err == redis.Nil { // Key not found, not necessarily an error
			return false, nil
		}
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return val == "", nil
}
