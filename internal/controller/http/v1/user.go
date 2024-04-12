package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
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
		h.POST("/update-user-profile", ur.updateUserProfileHandler)
		h.POST("/upload-profile-picture", ur.uploadProfilePictureHandler)
		h.POST("/add-social-account", ur.addSocialAccountHandler)
		h.GET("/cancel-add-social-account", ur.cancelAddSocialAccountHandler)
		h.GET("/remove-social-account", ur.removeSocialAccountHandler)
	}
}

func (ur *userRoutes) updateUserProfileHandler(c *gin.Context) {
	var userProfile dto.UpdateUserProfile

	if err := c.ShouldBind(&userProfile); err != nil {
		c.Header("HX-Reswap", "none")
		c.HTML(http.StatusOK, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		return
	}

	username := c.GetString("username")

	user, err := ur.userUC.FindByUsernameOrEmail(username, "")
	if err != nil {
		c.Header("HX-Reswap", "none")
		c.HTML(http.StatusOK, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		return
	}

	user.UserProfile.FirstName = userProfile.FirstName
	user.UserProfile.LastName = userProfile.LastName
	user.UserProfile.Gender = userProfile.Gender
	user.UserProfile.Country = userProfile.Country
	user.UserProfile.City = userProfile.City
	user.UserProfile.Address = userProfile.Address
	user.Email = userProfile.Email
	user.UserProfile.PhoneNumber = userProfile.PhoneNumber
	user.UserProfile.Birthday = userProfile.Birthday
	user.UserProfile.Company = userProfile.Company
	user.UserProfile.Role = userProfile.Role

	err = ur.userUC.Update(user)
	if err != nil {
		c.Header("HX-Reswap", "none")
		c.HTML(http.StatusOK, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		return
	}

	userInfo := prepareUserData(user) // because update succeeded, the user info is up to date, we could use it here. Don't need to query db again

	c.HTML(http.StatusOK, "new-toast-section", gin.H{
		"type":    "success",
		"message": "Your profile information has been saved!",
	})
	c.HTML(http.StatusOK, "update-profile-form", gin.H{
		"userInfo": userInfo,
	})
}

func (ur *userRoutes) uploadProfilePictureHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "upload profile picture",
	})
}

func (ur *userRoutes) addSocialAccountHandler(c *gin.Context) {
	// Define valid account types
	validAccountTypes := []string{"facebook", "twitter", "github", "youtube"}

	// Get account type and link from request
	accountType := c.Query("type")
	accountLink := c.PostForm(accountType + "-link")

	// Validate account type
	if !helper.ContainsString(validAccountTypes, accountType) {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Invalid account type.",
		})
		c.HTML(http.StatusBadRequest, "account-section", entity.SocialAccount{
			AccountType: accountType,
			AccountLink: accountLink,
			ButtonState: "save",
		})
		return
	}

	// Validate account link format (using a helper function)
	if !helper.IsValidUrl(accountLink) {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Invalid account link format.",
		})
		c.HTML(http.StatusBadRequest, "account-section", entity.SocialAccount{
			AccountType: accountType,
			AccountLink: accountLink,
			ButtonState: "save",
		})
		return
	}

	// Call user service to add social accounts
	_, err := ur.userUC.AddUserSocialAccounts(c.GetString("username"), map[string]string{accountType: accountLink})
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Failed to add social account.",
		})
		c.HTML(http.StatusBadRequest, "account-section", entity.SocialAccount{
			AccountType: accountType,
			AccountLink: accountLink,
			ButtonState: "save",
		})
		return
	}

	c.HTML(http.StatusOK, "new-toast-section", gin.H{
		"type":    "success",
		"message": "Social account added successfully.",
	})
	c.HTML(http.StatusBadRequest, "account-section", entity.SocialAccount{
		AccountType: accountType,
		AccountLink: accountLink,
		ButtonState: "save-success",
	})
}

func (ur *userRoutes) cancelAddSocialAccountHandler(c *gin.Context) {
	accountType := c.Query("type")
	c.HTML(http.StatusOK, "account-section", entity.SocialAccount{
		AccountType: accountType,
		AccountLink: "",
		ButtonState: "cancel",
	})
}

func (ur *userRoutes) removeSocialAccountHandler(c *gin.Context) {
	// Get username and account type from request parameters
	username := c.GetString("username")
	accountType := c.Query("type")

	// Call user service to remove the social account
	_, err := ur.userUC.RemoveUserSocialAccount(username, accountType)
	if err != nil {
		c.Header("HX-Reswap", "none")
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Failed to remove social account.",
		})
		return
	}

	c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
		"type":    "success",
		"message": "Social account removed successfully.",
	})
	c.HTML(http.StatusOK, "account-section", entity.SocialAccount{
		AccountType: accountType,
		AccountLink: "",
		ButtonState: "cancel",
	})
}
