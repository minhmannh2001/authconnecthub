package repos_test

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	redis_pkg "github.com/minhmannh2001/authconnecthub/pkg/redis"
	"github.com/minhmannh2001/authconnecthub/tests/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthRepoTestSuite struct {
	suite.Suite
	pgContainer    *testhelpers.PostgresContainer
	pg             *postgres.Postgres
	redisContainer *testhelpers.RedisContainer
	redis          *redis_pkg.Redis
	authRepo       *repos.AuthRepo
	ctx            context.Context
}

func (suite *AuthRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	host, err := pgContainer.ExtractHost(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	port, err := pgContainer.ExtractPort(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	pg, err := postgres.New(&config.Config{
		PG: config.PG{
			Host:     host,
			Port:     port,
			Username: "postgres",
			Password: "postgres",
			Dbname:   "test-db",
			Sslmode:  "disable",
		},
		Authen: config.Authen{
			AdminUsername: "admin",
			AdminPassword: "password",
			AdminEmail:    "admin@localhost",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	suite.pg = pg

	redisContainer, err := testhelpers.CreateRedisContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.redisContainer = redisContainer
	host, err = redisContainer.ExtractHost(redisContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	port, err = redisContainer.ExtractPort(redisContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	redis, err := redis_pkg.New(&config.Config{
		Redis: config.Redis{
			Host:     host,
			Port:     port,
			Password: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	suite.redis = redis

	suite.authRepo = repos.NewAuthRepo(suite.pg, suite.redis)
}

func (suite *AuthRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}

	if err := suite.redisContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating redis container: %s", err)
	}
}

func (suite *AuthRepoTestSuite) TestBlacklistToken_Success() {
	// Token and expiration to test
	token := "testtoken123"
	expiration := 60 // Seconds

	// Call BlacklistToken
	err := suite.authRepo.BlacklistToken(token, expiration)

	suite.Nil(err)
}

func TestBlacklistToken_RedisError(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	host, err := pgContainer.ExtractHost(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	port, err := pgContainer.ExtractPort(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	pg, err := postgres.New(&config.Config{
		PG: config.PG{
			Host:     host,
			Port:     port,
			Username: "postgres",
			Password: "postgres",
			Dbname:   "test-db",
			Sslmode:  "disable",
		},
		Authen: config.Authen{
			AdminUsername: "admin",
			AdminPassword: "password",
			AdminEmail:    "admin@localhost",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	redisContainer, err := testhelpers.CreateRedisContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	host, err = redisContainer.ExtractHost(redisContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	port, err = redisContainer.ExtractPort(redisContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	redis, err := redis_pkg.New(&config.Config{
		Redis: config.Redis{
			Host:     host,
			Port:     port,
			Password: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	authRepo := repos.NewAuthRepo(pg, redis)

	// Token and expiration to test
	token := "testtoken123"
	expiration := 60 // Seconds

	duration := 5 * time.Second
	redisContainer.Stop(ctx, &duration)

	// Call BlacklistToken
	err = authRepo.BlacklistToken(token, expiration)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "failed to blacklist token"))
}

func (suite *AuthRepoTestSuite) TestIsTokenBlacklisted_NotFound() {
	token := "nonexistenttoken"

	isBlacklisted, err := suite.authRepo.IsTokenBlacklisted(token)

	suite.Nil(err)
	suite.False(isBlacklisted)
}

func (suite *AuthRepoTestSuite) TestIsTokenBlacklisted_Blacklisted() {
	token := "blacklistedtoken"
	expiration := 60 // Seconds

	// Call BlacklistToken
	err := suite.authRepo.BlacklistToken(token, expiration)

	suite.Nil(err)

	isBlacklisted, err := suite.authRepo.IsTokenBlacklisted(token)

	suite.Nil(err)
	suite.True(isBlacklisted)
}

func TestAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepoTestSuite))
}
