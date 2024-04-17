package v1

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2MB

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
		h.GET("/reset-upload-profile-picture-progress-section", ur.resetUploadProfilePictureProgressSectionHandler)
		h.DELETE("/delete-profile-picture", ur.deleteProfilePictureHandler)
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
	file, fileHeader, err := c.Request.FormFile("upload-profile-picture-file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": "default",
			"fileName":   "uploaded file",
			"failReason": "Internal server error",
		})
		return
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": "default",
			"fileName":   fileHeader.Filename,
			"failReason": "Internal server error",
		})
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
		c.HTML(http.StatusOK, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "The provided file format is not allowed. Please upload a JPEG or PNG image.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": "default",
			"fileName":   fileHeader.Filename,
			"failReason": "Invalid file format",
		})
		return
	}

	// Set maximum upload size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)

	// Parse multipart form
	err = c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "The uploaded file is too big. Please choose an file that's less than 1MB in size",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": filetype[6:],
			"fileName":   fileHeader.Filename,
			"failReason": "File is too large (max. 2MB)",
		})
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": filetype[6:],
			"fileName":   fileHeader.Filename,
			"failReason": "Internal server error",
		})
		return
	}

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./static/images/uploads/profile_pictures", os.ModePerm)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": filetype[6:],
			"fileName":   fileHeader.Filename,
			"failReason": "Internal server error",
		})
		return
	}

	filePath := fmt.Sprintf("./static/images/uploads/profile_pictures/%s%s", c.GetString("username"), filepath.Ext(fileHeader.Filename))

	successMessage := "Profile picture uploaded successfully!"

	// Create a new file in the uploads directory
	dst, err := os.Create(filePath)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": filetype[6:],
			"fileName":   fileHeader.Filename,
			"failReason": "Internal server error",
		})
		return
	}

	defer dst.Close()

	pr := &helper.Progress{
		TotalSize: fileHeader.Size,
	}

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, io.TeeReader(file, pr))
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
			"fileFormat": filetype[6:],
			"fileName":   fileHeader.Filename,
			"failReason": "Internal server error",
		})
		return
	}

	user, _ := ur.userUC.FindByUsernameOrEmail(c.GetString("username"), "")
	if user.UserProfile.ProfilePicture != "" && user.UserProfile.ProfilePicture != c.GetString("username")+filepath.Ext(fileHeader.Filename) {
		successMessage = "Your profile picture has been updated."
		err = os.Remove("./static/images/uploads/profile_pictures/" + user.UserProfile.ProfilePicture)
		if err != nil {
			c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
				"type":    "danger",
				"message": "Oops! An error occurred. Please try again later.",
			})
			c.HTML(http.StatusBadRequest, "upload-profile-picture-progress-section", gin.H{
				"fileFormat": filetype[6:],
				"fileName":   fileHeader.Filename,
				"failReason": "Internal server error",
			})
			return
		}
	}
	user.UserProfile.ProfilePicture = c.GetString("username") + filepath.Ext(fileHeader.Filename)
	ur.userUC.Update(user)

	c.HTML(http.StatusOK, "user-profile-picture-in-dashboard-navbar", gin.H{
		"profilePictureURL": "/static/images/uploads/profile_pictures/" + c.GetString("username") + filepath.Ext(fileHeader.Filename),
	})
	c.HTML(http.StatusOK, "user-profile-picture", gin.H{
		"profilePictureURL": "/static/images/uploads/profile_pictures/" + c.GetString("username") + filepath.Ext(fileHeader.Filename),
	})
	c.HTML(http.StatusOK, "new-toast-section", gin.H{
		"type":    "success",
		"message": successMessage,
	})
	c.HTML(http.StatusOK, "upload-profile-picture-progress-section", gin.H{
		"fileFormat": filetype[6:],
		"fileName":   fileHeader.Filename,
		"failReason": "",
		"fileSize":   helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
	})
	c.HTML(http.StatusOK, "upload-profile-picture-modal-done-button", gin.H{})
}

func (ur *userRoutes) resetUploadProfilePictureProgressSectionHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload-profile-picture-progress-section-wrapper", gin.H{})
}

func (ur *userRoutes) deleteProfilePictureHandler(c *gin.Context) {
	c.Header("HX-Reswap", "none")
	user, _ := ur.userUC.FindByUsernameOrEmail(c.GetString("username"), "")
	err := os.Remove("./static/images/uploads/profile_pictures/" + user.UserProfile.ProfilePicture)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		return
	}
	user.UserProfile.ProfilePicture = ""
	err = ur.userUC.Update(user)
	if err != nil {
		c.HTML(http.StatusBadRequest, "new-toast-section", gin.H{
			"type":    "danger",
			"message": "Oops! An error occurred. Please try again later.",
		})
		return
	}
	c.HTML(http.StatusOK, "user-profile-picture-in-dashboard-navbar", gin.H{
		"profilePictureURL": "/static/images/uploads/profile_pictures/default.jpg",
	})
	c.HTML(http.StatusOK, "user-profile-picture", gin.H{
		"profilePictureURL": "/static/images/uploads/profile_pictures/default.jpg",
	})
	c.HTML(http.StatusOK, "new-toast-section", gin.H{
		"type":    "success",
		"message": "Profile picture deleted successfully!",
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
