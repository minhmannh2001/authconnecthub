package usecase

import (
	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

type (
	Auth interface {
		Register()
		CreateAccessToken(entity.User, int) (string, error)
		CreateRefreshToken(entity.User, string, int) (string, error)
		ValidateToken(string) (string, error)
	}
	AuthRepo interface {
	}

	User interface {
		Create(entity.User) (entity.User, error)
		FindByUsernameOrEmail(string, string) (*entity.User, error)
	}
	UserRepo interface {
		Create(entity.User) (entity.User, error)
		RetrieveByID(uint) (entity.User, error)
		Update(entity.User) (entity.User, error)
		Delete(entity.User) (entity.User, error)
		FindByUsernameOrEmail(string, string) (*entity.User, error)
	}

	Role interface {
		GetRoleIDByName(string) (uint, error)
	}
	RoleRepo interface {
		GetRoleIDByName(string) (uint, error)
	}
)
