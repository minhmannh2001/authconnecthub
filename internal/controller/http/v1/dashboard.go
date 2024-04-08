package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
)

type dashboardRoutes struct {
	logger *slog.Logger
	authUC usecases.IAuthUC
	userUC usecases.IUserUC
	roleUC usecases.IRoleUC
}

func NewDashboardRoutes(handler *gin.RouterGroup,
	l *slog.Logger,
	a usecases.IAuthUC,
	u usecases.IUserUC,
	r usecases.IRoleUC,
) {
	dr := &dashboardRoutes{l, a, u, r}

	h := handler.Group("/dashboard")
	{
		h.GET("", dr.getDashboardHandler)
		h.GET("/profile", dr.getProfileHandler)
	}
}

// @Summary Get User Profile
// @Description Retrieves the profile information for the currently authenticated user.
// @Tags dashboard
// @Security JWT
// @Param toast-message query string false "The message to display in the toast notification.""
// @Param toast-type query string false "The type of the toast notification (e.g., success, error).""
// @Param hash-value query string false "A hash value used for validation."
// @Produce html
// @Router /v1/dashboard/profile [GET]
func (dr *dashboardRoutes) getDashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
		"subPage": "dashboard",
		"reload":  c.GetHeader("HX-Reload"),
	})
}

func (dr *dashboardRoutes) getProfileHandler(c *gin.Context) {
	username := c.GetString("username")

	user, err := dr.userUC.FindByUsernameOrEmail(username, "")
	if err != nil {
		// Handle user retrieval error
	}

	userInfo := prepareUserData(user)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
		"subPage":  "profile",
		"userInfo": userInfo,
		"reload":   c.GetHeader("HX-Reload"),
	})
}

func prepareUserData(user *entity.User) map[string]interface{} {
	return map[string]interface{}{
		"username":    user.Username,
		"firstName":   user.UserProfile.FirstName,
		"lastName":    user.UserProfile.LastName,
		"gender":      user.UserProfile.Gender,
		"country":     user.UserProfile.Country,
		"city":        user.UserProfile.City,
		"address":     user.UserProfile.Address,
		"email":       user.Email,
		"phoneNumber": user.UserProfile.PhoneNumber,
		"birthday":    user.UserProfile.Birthday,
		"company":     user.UserProfile.Company,
		"role":        user.UserProfile.Role,
	}
}
