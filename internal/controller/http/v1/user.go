package v1

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	store  sync.Map
}

func NewUserRoutes(handler *gin.RouterGroup,
	l *slog.Logger,
	a usecases.IAuthUC,
	u usecases.IUserUC,
	r usecases.IRoleUC,
) {
	ur := &userRoutes{l, a, u, r, sync.Map{}}

	h := handler.Group("/user")
	{
		h.GET("/change-password", ur.loadChangePasswordPage)
		h.POST("/change-password", ur.changePasswordHandler)
		h.POST("/update-user-profile", ur.updateUserProfileHandler)
		h.POST("/upload-profile-picture", ur.uploadProfilePictureHandler)
		h.GET("/progress-of-upload-profile-picture", ur.getProgressOfUploadProfilePictureHandler)
		h.POST("/add-social-account", ur.addSocialAccountHandler)
		h.GET("/cancel-add-social-account", ur.cancelAddSocialAccountHandler)
		h.GET("/remove-social-account", ur.removeSocialAccountHandler)
		h.GET("/reset-upload-profile-picture-progress-section", ur.resetUploadProfilePictureProgressSectionHandler)
		h.DELETE("/delete-profile-picture", ur.deleteProfilePictureHandler)
	}
}

func (ur *userRoutes) loadChangePasswordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "change_password.html", gin.H{})
}

func (ur *userRoutes) changePasswordHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "change_password.html", gin.H{})
}

// @Summary Update User Profile
// @Description Updates the profile information for the currently authenticated user.
// @Tags User
// @Security JWT
// @Consumes json
// @Produce html
// @Router /v1/user/update-user-profile [POST]
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

// @Summary Get Profile Picture Upload Progress
// @Description Retrieves the upload progress for the currently ongoing profile picture upload (if any).
// @Tags User
// @Security JWT
// @Produce html
// @Router /v1/user/progress-of-upload-profile-picture [GET]
func (ur *userRoutes) getProgressOfUploadProfilePictureHandler(c *gin.Context) {
	progress, _ := ur.store.Load(c.GetString("username") + "upload-profile-picture-progress")
	if progress == nil {
		// Loop until progress is not nil
		for progress == nil {
			// Implement waiting logic here
			time.Sleep(time.Millisecond * 50)
			progress, _ = ur.store.Load(c.GetString("username") + "upload-profile-picture-progress") // Re-check the progress value
		}
	}
	progressMap := progress.(map[string]interface{})
	c.HTML(http.StatusOK, "upload-profile-picture-progress-section", gin.H{
		"uuid":                uuid.New().String(),
		"fileFormat":          progressMap["fileFormat"].(string),
		"fileName":            progressMap["fileName"].(string),
		"currentPercent":      progressMap["currentPercent"].(int),
		"currentUploadedSize": progressMap["currentUploadedSize"].(string),
		"totalSize":           progressMap["totalSize"].(string),
		"uploading":           progressMap["uploading"].(bool),
		"finish":              progressMap["uploading"].(bool),
	})
	if !progressMap["uploading"].(bool) {
		ur.store.Delete(c.GetString("username") + "upload-profile-picture-progress")
	}
}

// @Summary Upload Profile Picture
// @Description Uploads a new profile picture for the currently authenticated user.
// @Tags User
// @Security JWT
// @Consumes multipart/form-data
// @Param upload-profile-picture-file formData file true "The profile picture file to upload"
// @Produce html
// @Router /v1/user/upload-profile-picture [POST]
func (ur *userRoutes) uploadProfilePictureHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload-profile-picture-modal-done-button", gin.H{})

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
			"uploading":  false,
			"finish":     true,
		})
		return
	}

	defer file.Close()

	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		ur.store.Store(c.GetString("username")+"upload-profile-picture-progress", map[string]interface{}{
			"fileFormat":          "default",
			"fileName":            fileHeader.Filename,
			"currentPercent":      i,
			"currentUploadedSize": "0",
			"totalSize":           helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
			"uploading":           true,
			"finish":              false,
		})
	}

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
			"uploading":  false,
			"finish":     true,
		})
		return
	}

	for i := 11; i <= 20; i++ {
		time.Sleep(10 * time.Millisecond)
		ur.store.Store(c.GetString("username")+"upload-profile-picture-progress", map[string]interface{}{
			"fileFormat":          filetype[6:],
			"fileName":            fileHeader.Filename,
			"currentPercent":      i,
			"currentUploadedSize": "0",
			"totalSize":           helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
			"uploading":           true,
			"finish":              false,
		})
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
			"uploading":  false,
			"finish":     true,
		})
		return
	}

	for i := 21; i <= 30; i++ {
		time.Sleep(10 * time.Millisecond)
		ur.store.Store(c.GetString("username")+"upload-profile-picture-progress", map[string]interface{}{
			"fileFormat":          filetype[6:],
			"fileName":            fileHeader.Filename,
			"currentPercent":      i,
			"currentUploadedSize": "0",
			"totalSize":           helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
			"uploading":           true,
			"finish":              false,
		})
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
			"uploading":  false,
			"finish":     true,
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
			"uploading":  false,
			"finish":     true,
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
			"uploading":  false,
			"finish":     true,
		})
		return
	}

	defer dst.Close()

	pr := &helper.Progress{
		Store:     &ur.store,
		Username:  c.GetString("username"),
		Filename:  fileHeader.Filename,
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
			"uploading":  false,
			"finish":     true,
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
				"uploading":  false,
				"finish":     true,
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
	ur.store.Store(c.GetString("username")+"upload-profile-picture-progress", map[string]interface{}{
		"fileFormat":          filetype[6:],
		"fileName":            fileHeader.Filename,
		"currentPercent":      100,
		"currentUploadedSize": helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
		"totalSize":           helper.FormatFileSize(float64(fileHeader.Size), 1024.0),
		"uploading":           false,
		"finish":              true,
	})
	c.HTML(http.StatusOK, "upload-profile-picture-modal-done-button", gin.H{})
}

// @Summary Reset Upload Profile Picture Progress Section
// @Description Resets the upload profile picture progress section on the frontend.
// @Tags User
// @Security JWT
// @Produce html
// @Response 200 { description: "Upload profile picture progress section reset" }
// @Router /v1/user/reset-upload-profile-picture-progress-section [GET]
func (ur *userRoutes) resetUploadProfilePictureProgressSectionHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload-profile-picture-progress-section-wrapper", gin.H{})
}

// @Summary Delete Profile Picture
// @Description Deletes the currently authenticated user's profile picture.
// @Tags User
// @Security JWT
// @Produce html
// @Response 200 { description: "Profile picture deleted successfully" }
// @Response 400 { description: "Bad Request - Error deleting profile picture" }
// @Router /v1/user/delete-profile-picture [DELETE]
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

// @Summary Add Social Account
// @Description Adds a social media account to the currently authenticated user's profile.
// @Tags User
// @Security JWT
// @Produce html
// @Router /v1/user/add-social-account [POST]
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

// @Summary Cancel Add Social Account
// @Description Cancels the process of adding a social media account to the user's profile.
// @Tags User
// @Security JWT
// @Produce html
// @Router /v1/user/cancel-add-social-account [GET]
func (ur *userRoutes) cancelAddSocialAccountHandler(c *gin.Context) {
	accountType := c.Query("type")
	c.HTML(http.StatusOK, "account-section", entity.SocialAccount{
		AccountType: accountType,
		AccountLink: "",
		ButtonState: "cancel",
	})
}

// @Summary Remove Social Account
// @Description Removes a social media account from the currently authenticated user's profile.
// @Tags User
// @Security JWT
// @Consumes application/x-www-form-urlencoded  // Optional, but can be added for consistency
// @Param username query string true "Username of the authenticated user"
// @Param type query string true "The type of social account to remove (e.g., facebook, twitter)"
// @Produce html
// @Router /v1/user/remove-social-account [GET]
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

// https://celery.school/celery-progress-bars-with-fastapi-htmx
// https://freshman.tech/file-upload-golang/
// https://medium.com/@relia/an-in-depth-guide-for-using-go-sync-map-with-code-sample-e742814e7bce
