package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
)

type userRoutes struct {
	logger *slog.Logger
	authUC usecases.IAuthUC
	userUC usecases.IUserUC
	roleUC usecases.IRoleUC
}

func NewUserRoutes(handler *gin.RouterGroup,
	l *slog.Logger,
	a usecases.IAuthUC,
	u usecases.IUserUC,
	r usecases.IRoleUC,
) {
	ur := &userRoutes{l, a, u, r}

	h := handler.Group("/user")
	{
		h.POST("/save-user-profile", ur.saveUserProfileHandler)
		h.POST("/upload-profile-picture", ur.uploadProfilePictureHandler)
	}
}

func (ur *userRoutes) saveUserProfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "upload profile picture",
	})
}

func (ur *userRoutes) uploadProfilePictureHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "upload profile picture",
	})
}
