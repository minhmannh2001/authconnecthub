package postgres_test

import (
	"context"
	"log"
	"testing"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"github.com/minhmannh2001/authconnecthub/tests/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repo        *postgres.Postgres
	ctx         context.Context
}

func (suite *PostgresRepoTestSuite) SetupSuite() {
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
	suite.repo = pg
}

func (suite *PostgresRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *PostgresRepoTestSuite) TestPostgresNew_Success() {
	var count int64
	suite.repo.Conn.Model(&entity.User{}).Where("username = ?", "admin").Count(&count)
	suite.True(count > 0)

	suite.repo.Conn.Model(&entity.Role{}).Count(&count)
	suite.True(count == 3)

	suite.repo.Conn.Model(&entity.Role{}).Where("name = ?", "admin").Count(&count)
	suite.True(count == 1)

	suite.repo.Conn.Model(&entity.Role{}).Where("name = ?", "customer").Count(&count)
	suite.True(count == 1)

	suite.repo.Conn.Model(&entity.Role{}).Where("name = ?", "anonymous").Count(&count)
	suite.True(count == 1)
}

func TestPostgresNew_Fail(t *testing.T) {
	pg, err := postgres.New(&config.Config{
		PG: config.PG{
			Host:     "localhost",
			Port:     "1206",
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

	assert.Nil(t, pg)
	assert.Error(t, err)
}

func TestPostgresRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoTestSuite))
}
