package repos

import "github.com/minhmannh2001/authconnecthub/internal/entity"

type (
	IAuthRepo interface {
		BlacklistToken(string, int) error
		IsTokenBlacklisted(string) (bool, error)
	}

	IUserRepo interface {
		Create(entity.User) (entity.User, error)
		RetrieveByID(uint) (entity.User, error)
		Update(entity.User) (entity.User, error)
		Delete(entity.User) (entity.User, error)
		FindByUsernameOrEmail(string, string) (*entity.User, error)
	}

	IRoleRepo interface {
		GetRoleIDByName(string) (uint, error)
	}
)
