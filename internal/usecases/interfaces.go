package usecases

import (
	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

type (
	IAuthUC interface {
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

	IUserUC interface {
		Create(entity.User) (entity.User, error)
		Update(*entity.User) error
		FindByUsernameOrEmail(string, string) (*entity.User, error)
		GetUserSocialAccounts(string) (map[string]entity.SocialAccount, error)
		AddUserSocialAccounts(string, map[string]string) (bool, error)
		RemoveUserSocialAccount(string, string) (bool, error)
	}

	IRoleUC interface {
		GetRoleIDByName(string) (uint, error)
	}
)

// mockery --dir=./internal/usecases --output=./internal/usecases/mocks --outpkg=mocks --all
