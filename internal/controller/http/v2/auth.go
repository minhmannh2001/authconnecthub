package v2

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"golang.org/x/crypto/bcrypt"
)

type authRoutes struct {
	logger *slog.Logger
	authUC usecases.IAuthUC
	userUC usecases.IUserUC
	roleUC usecases.IRoleUC
}

// NewAuthenRoutes creates new authen routes
func NewAuthenRoutes(handler *gin.RouterGroup,
	l *slog.Logger,
	a usecases.IAuthUC,
	u usecases.IUserUC,
	r usecases.IRoleUC,
) {
	ar := &authRoutes{l, a, u, r}

	h := handler.Group("/auth")
	{
		h.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"title": "Personal Hub",
			})
		})
		h.POST("/register", ar.register)
	}
}

func (ar *authRoutes) register(c *gin.Context) {
	var registerRequestBody dto.RegisterRequestBody

	// validate json
	err := c.ShouldBind(&registerRequestBody)
	// validation errors
	if err != nil {
		// generate validation errors response
		response := helper.GenerateValidationResponse(err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(registerRequestBody.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ar.logger.Error("Password hashing failed", slog.Any("err", err))
		response := dto.Response{
			Success: false,
			Data:    nil,
			Message: "Password hashing failed due to an internal error. Please try again later.",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	roleID, err := ar.roleUC.GetRoleIDByName("customer")

	if err != nil {
		ar.logger.Error("Error getting role ID", slog.Any("err", err))
		response := dto.Response{Success: false, Data: nil, Message: "Failed to create user"}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	user := entity.User{
		Username: registerRequestBody.Username,
		Email:    registerRequestBody.Email,
		Password: string(encryptedPassword),
		RoleID:   roleID,
	}

	newUser, err := ar.userUC.Create(user)

	if err != nil {
		response := dto.Response{Success: false, Data: nil}

		if helper.IsErrOfType(err, &entity.ErrDuplicateUser{}) {
			response.Message = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}

		ar.logger.Error("Error creating user", slog.Any("err", err))
		response.Message = "Failed to create user"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	ar.logger.Info("User created", slog.String("username", newUser.Username), slog.Any("roleID", newUser.RoleID), slog.String("email", newUser.Email))
	response := dto.Response{
		Success: true,
		Data:    newUser,
		Message: "User registered successfully",
	}
	c.JSON(http.StatusOK, response)
}
