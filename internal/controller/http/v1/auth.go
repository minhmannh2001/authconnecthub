package v1

import (
	"fmt"
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
		h.GET("/login", ar.getLogin)
		h.POST("/login", ar.postLogin)

		h.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"title": "Personal Hub",
				"toastSettings": map[string]interface{}{
					"hidden": true,
				},
			})
		})
		h.POST("/register", ar.register)

		h.GET("/logout", ar.LogoutHandler)
	}
}

// @Summary Login Page
// @Description This endpoint renders the login page and displays a toast notification if provided query parameters are valid.
// @Tags Authen
// @Produce html
// @Param toast-message query string false "The message to display in the toast notification.""
// @Param toast-type query string false "The type of the toast notification (e.g., success, error).""
// @Param hash-value query string false "A hash value used for validation."
// @router /v1/auth/login [GET]
func (ar *authRoutes) getLogin(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	toastMessage := helper.ExtractQueryParam(queryParams, "toast-message", "")
	toastType := helper.ExtractQueryParam(queryParams, "toast-type", "")
	hashValue := helper.ExtractQueryParam(queryParams, "hash-value", "")

	isValid := helper.IsMapValid(map[string]interface{}{
		"toast-message": toastMessage,
		"toast-type":    toastType,
	}, hashValue)

	toastSettings := map[string]interface{}{
		"hidden":  !isValid, // Toggle based on validity
		"type":    toastType,
		"message": helper.FormatToastMessage(toastMessage),
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":         "Personal Hub",
		"toastSettings": toastSettings,
		"reload":        c.GetHeader("HX-Reload"),
	})
}

func (ar *authRoutes) register(c *gin.Context) {
	var registerRequestBody dto.RegisterRequestBody

	// validate json
	err := c.ShouldBind(&registerRequestBody)
	// validation errors
	if err != nil {
		// generate validation errors response
		validationMap := helper.GenerateValidationMap(err)
		_ = validationMap

		c.HTML(http.StatusBadRequest, "toast-section", gin.H{
			"hidden":  false,
			"type":    dto.ToastTypeDanger,
			"message": "Failed to create user.",
		})

		inputData := map[string]string{
			"username":        registerRequestBody.Username,
			"email":           registerRequestBody.Email,
			"password":        registerRequestBody.Password,
			"confirmPassword": registerRequestBody.ConfirmPassword,
		}

		c.HTML(http.StatusOK, "register-form", gin.H{
			"inputData":      inputData,
			"validationFail": true,
			"validationMap":  validationMap,
		})

		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(registerRequestBody.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ar.logger.Error("Password hashing failed", slog.Any("err", err))
		c.HTML(http.StatusBadRequest, "toast-section", gin.H{
			"hidden":  false,
			"type":    dto.ToastTypeDanger,
			"message": "Password hashing failed due to an internal error. Please try again later.",
		})

		inputData := map[string]string{
			"username":        registerRequestBody.Username,
			"email":           registerRequestBody.Email,
			"password":        registerRequestBody.Password,
			"confirmPassword": registerRequestBody.ConfirmPassword,
		}

		c.HTML(http.StatusOK, "register-form", gin.H{
			"inputData":      inputData,
			"validationFail": true,
			"validationMap":  map[string]string{},
		})
		return
	}

	roleID, err := ar.roleUC.GetRoleIDByName("customer")

	if err != nil {
		ar.logger.Error("Error getting role ID", slog.Any("err", err))
		c.HTML(http.StatusBadRequest, "toast-section", gin.H{
			"hidden":  false,
			"type":    dto.ToastTypeDanger,
			"message": "Failed to create user due to an internal error. Please try again later.",
		})

		inputData := map[string]string{
			"username":        registerRequestBody.Username,
			"email":           registerRequestBody.Email,
			"password":        registerRequestBody.Password,
			"confirmPassword": registerRequestBody.ConfirmPassword,
		}

		c.HTML(http.StatusOK, "register-form", gin.H{
			"inputData":      inputData,
			"validationFail": true,
			"validationMap":  map[string]string{},
		})
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

		if helper.IsErrOfType(err, &entity.ErrDuplicateUser{}) {
			c.HTML(http.StatusBadRequest, "toast-section", gin.H{
				"hidden":  false,
				"type":    dto.ToastTypeDanger,
				"message": err.Error(),
			})

			c.HTML(http.StatusOK, "register-form", gin.H{
				"inputData":      map[string]string{},
				"validationFail": false,
				"validationMap":  map[string]string{},
			})
			return
		}

		ar.logger.Error("Error creating user", slog.Any("err", err))
		c.HTML(http.StatusBadRequest, "toast-section", gin.H{
			"hidden":  false,
			"type":    dto.ToastTypeDanger,
			"message": "Failed to create user.",
		})

		inputData := map[string]string{
			"username":        registerRequestBody.Username,
			"email":           registerRequestBody.Email,
			"password":        registerRequestBody.Password,
			"confirmPassword": registerRequestBody.ConfirmPassword,
		}

		c.HTML(http.StatusOK, "register-form", gin.H{
			"inputData":      inputData,
			"validationFail": true,
			"validationMap":  map[string]string{},
		})
		return
	}

	cfg := helper.GetConfig(c)
	accessToken, _ := ar.authUC.CreateAccessToken(user, cfg.Authen.AccessTokenTTL)
	refreshToken, _ := ar.authUC.CreateRefreshToken(user, accessToken, cfg.Authen.RefreshTokenTTL)

	hashValue, err := helper.HashMap(map[string]interface{}{
		"toast-message": "user-registered-successfully",
		"toast-type":    dto.ToastTypeSuccess,
	})
	if err != nil {
		ar.logger.Error("Failed to generate toast message hash", slog.Any("err", err))
	}

	ar.logger.Info("User created", slog.String("username", newUser.Username), slog.Any("roleID", newUser.RoleID), slog.String("email", newUser.Email))
	HXTriggerEvents, _ := helper.MapToJSONString(map[string]interface{}{
		"saveToken": map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
	c.Header("HX-Trigger", HXTriggerEvents)
	c.Header("HX-Redirect", fmt.Sprintf("/?toast-message=user-registered-successfully&toast-type=%s&hash-value=%s", dto.ToastTypeSuccess, hashValue))
}

func (ar *authRoutes) postLogin(c *gin.Context) {
	var loginRequestBody dto.LoginRequestBody

	if err := c.ShouldBind(&loginRequestBody); err != nil {
		ar.logger.Error("Validation error", slog.Any("err", err))
		// generate validation errors response
		validationMap := helper.GenerateValidationMap(err)
		_ = validationMap

		c.HTML(http.StatusBadRequest, "toast-section", gin.H{
			"hidden":  false,
			"type":    dto.ToastTypeDanger,
			"message": "Invalid credentials. Please try again.",
		})

		inputData := map[string]string{
			"username": loginRequestBody.Username,
			"password": loginRequestBody.Password,
		}

		c.HTML(http.StatusOK, "login-form", gin.H{
			"inputData":      inputData,
			"validationFail": true,
			"validationMap":  validationMap,
		})

		return
	}

	jwtTokens, err := ar.authUC.Login(c, loginRequestBody)
	if err != nil {
		inputData := map[string]string{
			"username":    loginRequestBody.Username,
			"password":    loginRequestBody.Password,
			"remember_me": loginRequestBody.RememberMe,
		}

		if helper.IsErrOfType(err, &entity.InvalidCredentialsError{}) {
			c.HTML(http.StatusBadRequest, "toast-section", gin.H{
				"hidden":  false,
				"type":    dto.ToastTypeDanger,
				"message": err.Error(),
			})

			c.HTML(http.StatusOK, "login-form", gin.H{
				"inputData":      inputData,
				"validationFail": true,
				"validationMap":  map[string]string{},
			})
		} else {
			c.HTML(http.StatusBadRequest, "toast-section", gin.H{
				"hidden":  false,
				"type":    dto.ToastTypeDanger,
				"message": "An unexpected error occurred. Please try again later.",
			})

			c.HTML(http.StatusOK, "login-form", gin.H{
				"inputData":      inputData,
				"validationFail": true,
				"validationMap":  map[string]string{},
			})
		}
		return
	}

	saveTo := "session"
	if loginRequestBody.RememberMe == "on" {
		saveTo = "local"
	}

	HXTriggerEvents, err := helper.MapToJSONString(map[string]interface{}{
		"saveToken": map[string]interface{}{
			"saveTo":       saveTo,
			"accessToken":  jwtTokens.AccessToken,
			"refreshToken": jwtTokens.RefreshToken,
		},
	})
	if err != nil {
		ar.logger.Error("Failed to create HX-Trigger events", slog.Any("err", err))
	}

	hashValue, err := helper.HashMap(map[string]interface{}{
		"toast-message": "login-successfully",
		"toast-type":    dto.ToastTypeSuccess,
	})
	if err != nil {
		ar.logger.Error("Failed to generate toast message hash", slog.Any("err", err))
		helper.HandleInternalError(c, err)
		return
	}

	ar.logger.Info("User logged in", slog.String("username", loginRequestBody.Username))
	c.Header("HX-Trigger", HXTriggerEvents)
	c.Header("HX-Redirect", fmt.Sprintf("/?toast-message=login-successfully&toast-type=%s&hash-value=%s", dto.ToastTypeSuccess, hashValue))
}

// @Summary Logout User
// @Description Logs out the currently authenticated user and redirects to the home page with a success toast notification.
// @Tags Authen
// @Security JWT
// @Router /v1/auth/logout [GET]
func (ar *authRoutes) LogoutHandler(c *gin.Context) {
	err := ar.authUC.Logout(c)
	if err != nil {
		ar.logger.Error("Failed to logout user", slog.Any("username", c.MustGet("username")), slog.Any("err", err))
		helper.HandleInternalError(c, err)
		return
	}

	hashValue, err := helper.HashMap(map[string]interface{}{
		"toast-message": "logout-successfully",
		"toast-type":    dto.ToastTypeSuccess,
	})
	if err != nil {
		ar.logger.Error("Failed to generate toast message hash", slog.Any("err", err))
		helper.HandleInternalError(c, err)
		return
	}
	ar.logger.Error("User logged out", slog.Any("username", c.MustGet("username")))
	c.Header("HX-Redirect", fmt.Sprintf("/?toast-message=logout-successfully&toast-type=%s&hash-value=%s", dto.ToastTypeSuccess, hashValue))
}
