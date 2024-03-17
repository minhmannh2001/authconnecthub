package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecase"
)

const (
	AccessTokenHeader  string = "Authorization"
	RefreshTokenHeader string = "Refresh"
)

func IsHtmxRequest(c *gin.Context) {
	cfg := helper.GetConfig(c)
	swaggerInfo, err := helper.GetSwaggerInfo(cfg.App.SwaggerPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error loading swagger information"})
		return
	}

	if !helper.IsPathMethodInSwagger(c.Request.URL.Path, c.Request.Method, swaggerInfo) {
		c.Next()
		return
	}

	hxRequest := c.GetHeader("HX-Request")
	hxReload := c.GetHeader("HX-Reload")
	if hxRequest == "" && hxReload == "" {
		// Request doesn't have HX-Request header, likely from browser interaction
		// Trigger a response to simulate request with htmx
		triggerHtmxReload(c)
		return
	}

	// Request has HX-Request header, proceed with processing
	c.Next()
}

func triggerHtmxReload(c *gin.Context) {
	c.HTML(http.StatusOK, "reload.html", gin.H{
		"method":       c.Request.Method,
		"path":         c.Request.URL.Path + "?" + c.Request.URL.RawQuery,
		"body":         nil, // The request body if method is PUT or PATCH,
		"reloadHeader": true,
	})
	c.Abort()
}

func IsAuthorized(auth usecase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := helper.GetConfig(c)
		swaggerInfo, err := helper.GetSwaggerInfo(cfg.App.SwaggerPath)
		if err != nil {
			helper.HandleInternalError(c, err)
			return
		}

		isPrivateRoute, err := helper.HasSecurityKeyForPathAndMethod(c.Request.URL.Path, c.Request.Method, swaggerInfo)
		if err != nil {
			helper.HandleInternalError(c, err)
			return
		}
		if !isPrivateRoute {
			c.Next()
			return
		}

		accessToken := extractHeaderToken(c, AccessTokenHeader)

		if accessToken == "" {
			toastMessage := "login-is-required-for-this-action.-sign-in-or-create-an-account-to-continue."
			redirectToLogin(c, toastMessage)
			return
		}

		username, err := auth.ValidateToken(accessToken)
		if err != nil {
			log.Printf("Internal error: %v\n", err)
			toastMessage := "your-session-has-expired.-please-log-in-to-continue."
			redirectToLogin(c, toastMessage)
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

func redirectToLogin(c *gin.Context, message string) {
	toastData := map[string]interface{}{
		"toast-message": message,
		"toast-type":    dto.ToastTypeDanger,
	}

	hashValue, err := helper.HashMap(toastData)
	if err != nil {
		helper.HandleInternalError(c, err)
		return
	}

	redirectURL := fmt.Sprintf("/v1/auth/login?toast-message=%s&toast-type=%s&hash-value=%s", message, dto.ToastTypeDanger, hashValue)
	c.Header("HX-Redirect", redirectURL)
	c.Abort()
}

func IsLoggedIn(auth usecase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractHeaderToken(c, AccessTokenHeader)

		// Validate token if present
		if accessToken != "" {
			if _, err := auth.ValidateToken(accessToken); err != nil {
				if helper.IsTokenExpired(err) {
					refreshToken := extractHeaderToken(c, RefreshTokenHeader)
					newAccessToken, newRerefreshToken, err := auth.CheckAndRefreshTokens(accessToken, refreshToken, helper.GetConfig(c))
					if err != nil {
						goto SESSION_EXPIRE
					}
					HXTriggerEvents, err := helper.MapToJSONString(map[string]interface{}{
						"saveToken": map[string]interface{}{
							"accessToken":  newAccessToken,
							"refreshToken": newRerefreshToken,
						},
					})
					if err != nil {
						goto SESSION_EXPIRE
					}
					c.Header("HX-Trigger", HXTriggerEvents)
					modifyAuthorizationHeaders(c, newAccessToken, newRerefreshToken)
					goto CONTINUE
				}
			SESSION_EXPIRE:
				toastMessage := "your-session-has-expired.-please-log-in-to-continue."
				c.Header("HX-Trigger", "deleteToken")
				redirectToLogin(c, toastMessage)
				return
			}

		CONTINUE:
			// Redirect to home page if logged in user accesses login or register pages
			if shouldRedirectToHome(c.Request.URL.Path) {
				c.Header("HX-Redirect", "/")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// Extract token from request headers
func extractHeaderToken(c *gin.Context, tokenType string) string {
	headerToken := c.GetHeader(tokenType)
	if headerToken != "" && strings.HasPrefix(headerToken, "Bearer ") {
		parts := strings.Split(headerToken, " ")
		if len(parts) == 2 && parts[1] != "null" {
			return parts[1]
		}
	}
	return ""
}

// Determines if a path should be redirected to the home page
func shouldRedirectToHome(path string) bool {
	return path == "/v1/auth/login" || path == "/v1/auth/register"
}

func modifyAuthorizationHeaders(c *gin.Context, newAccessToken string, newRefreshToken string) {
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", newAccessToken))
	c.Request.Header.Set("Refresh", fmt.Sprintf("Bearer %s", newRefreshToken))
}
