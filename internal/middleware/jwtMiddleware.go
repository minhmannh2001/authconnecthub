package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecase"
)

func IsAuthorized(auth usecase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {

		privateURLs := []string{"/private"}
		if !isPrivateRoute(c.Request.URL.Path, privateURLs) {
			c.Next()
			return
		}

		tokenString, err := parseAuthToken(c)
		if tokenString == "" && err == nil {
			c.Abort()
			return
		}
		if err != nil {
			toastMessage := "login-is-required-for-this-action.-sign-in-or-create-an-account-to-continue."
			hashValue, err := helper.HashMap(map[string]interface{}{
				"toast-message": toastMessage,
				"toast-type":    dto.ToastTypeDanger,
			})
			if err != nil {
				// Implement internal error, return to home page
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}

			c.Redirect(http.StatusFound, fmt.Sprintf("/v1/auth/login?toast-message=%s&toast-type=%s&hash-value=%s", toastMessage, dto.ToastTypeDanger, hashValue))
			c.Abort()
			return
		}

		username, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

func isPrivateRoute(url string, privateURLs []string) bool {
	for _, privateURL := range privateURLs {
		if url == privateURL {
			return true
		}
	}
	return false
}

func parseAuthToken(c *gin.Context) (string, error) {
	headerToken := c.GetHeader("Authorization")
	if headerToken != "" && strings.HasPrefix(headerToken, "Bearer ") {
		parts := strings.Split(headerToken, " ")
		if len(parts) == 2 && parts[1] != "null" {
			return parts[1], nil
		}
	}

	_, err := c.Cookie("alreadyResend")
	if err != nil && err.Error() == "http: named cookie not present" {
		c.SetCookie("alreadyResend", "true", 60, "/", "localhost", false, true)
		c.HTML(http.StatusOK, "reload.html", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"body":   nil,
		})
		return "", nil
	}

	return "", errors.New("unauthorized")
}
