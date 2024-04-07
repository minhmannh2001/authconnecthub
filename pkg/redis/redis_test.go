package redis_test

import (
	"context"
	"log"
	"testing"

	"github.com/minhmannh2001/authconnecthub/config"
	redis_pkg "github.com/minhmannh2001/authconnecthub/pkg/redis"
	"github.com/minhmannh2001/authconnecthub/tests/testhelpers"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	redisContainer *testhelpers.RedisContainer
	redis          *redis_pkg.Redis
	ctx            context.Context
}

func (suite *RedisTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	redisContainer, err := testhelpers.CreateRedisContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.redisContainer = redisContainer
	host, err := redisContainer.ExtractHost(redisContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	port, err := redisContainer.ExtractPort(redisContainer.ConnectionString)
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
}

func (suite *RedisTestSuite) TearDownSuite() {
	if err := suite.redisContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating redis container: %s", err)
	}
}

func (suite *RedisTestSuite) TestRedisNew() {
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
