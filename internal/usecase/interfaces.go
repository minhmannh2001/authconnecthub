package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

type (
	Auth interface {
		Login(*gin.Context, dto.LoginRequestBody) (*dto.JwtTokens, error)
		Register()
		CreateAccessToken(entity.User, int) (string, error)
		CreateRefreshToken(entity.User, string, int) (string, error)
		ValidateToken(string) (string, error)
		RetrieveFieldFromJwtToken(string, string, bool) (interface{}, error)
		IsRefreshTokenValidForAccessToken(string, string) (bool, error)
		CheckAndRefreshTokens(string, string, *config.Config) (string, string, error)
		Logout(c *gin.Context) error
		IsTokenBlacklisted(string) (bool, error)
	}
	AuthRepo interface {
		BlacklistToken(string, int) error
		IsTokenBlacklisted(string) (bool, error)
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
