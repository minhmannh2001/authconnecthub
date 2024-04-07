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

type UserRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	pg          *postgres.Postgres
	userRepo    *repos.UserRepo
	ctx         context.Context
}

func (suite *UserRepoTestSuite) SetupSuite() {
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
	suite.userRepo = repos.NewUserRepo(pg)
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *UserRepoTestSuite) TestCreate_Success() {
	newUser := entity.User{
		Username: "newuser",
	}

	// Call Create
	createdUser, err := suite.userRepo.Create(newUser)

	suite.Nil(err)
	suite.Equal(newUser.Username, createdUser.Username)
}

func (suite *UserRepoTestSuite) TestCreate_DuplicateUser() {
	newUser := entity.User{
		Username: "admin",
	}
	// Call Create
	createdUser, err := suite.userRepo.Create(newUser)

	suite.Equal(createdUser, entity.User{})
	suite.NotNil(err)
	suite.Equal(&entity.ErrDuplicateUser{Username: newUser.Username, Email: newUser.Email}, err)
}

func (suite *UserRepoTestSuite) TestFindByUsernameOrEmail_ByUsernameSuccess() {
	user, err := suite.userRepo.FindByUsernameOrEmail("admin", "")

	suite.Nil(err)
	suite.Equal("admin", user.Username)
	suite.Equal("admin@localhost", user.Email)
}

func (suite *UserRepoTestSuite) TestFindByUsernameOrEmail_ByEmailSuccess() {
	user, err := suite.userRepo.FindByUsernameOrEmail("", "admin@localhost")

	suite.Nil(err)
	suite.Equal("admin", user.Username)
	suite.Equal("admin@localhost", user.Email)
}

func (suite *UserRepoTestSuite) TestFindByUsernameOrEmail_NotFound() {
	user, err := suite.userRepo.FindByUsernameOrEmail("minhmannh2001", "")

	suite.Nil(user)
	suite.Equal(&entity.InvalidCredentialsError{}, err)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
