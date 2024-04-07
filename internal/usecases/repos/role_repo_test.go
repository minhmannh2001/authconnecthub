package repos_test

import (
	"context"
	"log"
	"testing"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"github.com/minhmannh2001/authconnecthub/tests/testhelpers"
	"github.com/stretchr/testify/suite"
)

type RoleRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	pg          *postgres.Postgres
	roleRepo    *repos.RoleRepo
	ctx         context.Context
}

func (suite *RoleRepoTestSuite) SetupSuite() {
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
	suite.roleRepo = repos.NewRoleRepo(pg)
}

func (suite *RoleRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *RoleRepoTestSuite) TestGetRoleIDByName_Success() {
	// Name of the role to search
	name := "admin"

	// Call GetRoleIDByName
	roleID, err := suite.roleRepo.GetRoleIDByName(name)

	suite.Equal(uint(1), roleID)
	suite.Nil(err)
}

func (suite *RoleRepoTestSuite) TestGetRoleIDByName_NotFound() {
	// Name of the role to search
	name := "abc"

	// Call GetRoleIDByName
	roleID, err := suite.roleRepo.GetRoleIDByName(name)

	suite.Equal(uint(0), roleID)
	suite.Equal(&entity.RoleNotFoundError{Name: name}, err)
}

func TestRoleRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepoTestSuite))
}
