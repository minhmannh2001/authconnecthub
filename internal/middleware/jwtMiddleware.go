package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecase"
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
			// try to get user information if there is access token
			accessToken := helper.ExtractHeaderToken(c, helper.AccessTokenHeader)

			if accessToken != "" {
				username, _ := auth.ValidateToken(accessToken)

				c.Set("username", username)
				c.Next()
				return
			}
			c.Next()
			return
		}

		accessToken := helper.ExtractHeaderToken(c, helper.AccessTokenHeader)

		if accessToken == "" {
			toastMessage := "login-is-required-for-this-action.-sign-in-or-create-an-account-to-continue."
			redirectToLogin(c, toastMessage)
			return
		}

		// check if token is in blacklist or not, if it is then return to login page with error notification
		isInBlackList, _ := auth.IsTokenBlacklisted(accessToken)

		if isInBlackList {
			rememberMe, _ := auth.RetrieveFieldFromJwtToken(accessToken, "remember_me", false)
			toastMessage := "your-token-is-invalid.-please-log-in-to-continue."
			helper.DeleteTokens(c, rememberMe.(bool))
			redirectToLogin(c, toastMessage)
			return
		}

		username, err := auth.ValidateToken(accessToken)
		if err != nil {
			if helper.IsTokenExpired(err) && c.Request.URL.Path == "/v1/auth/logout" {
				c.Next()
				return
			}
			log.Printf("Internal error: %v\n", err)
			rememberMe, _ := auth.RetrieveFieldFromJwtToken(accessToken, "remember_me", true)
			helper.DeleteTokens(c, rememberMe.(bool))
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
		accessToken := helper.ExtractHeaderToken(c, helper.AccessTokenHeader)

		// Validate token if present
		if accessToken != "" {
			if _, err := auth.ValidateToken(accessToken); err != nil {
				if helper.IsTokenExpired(err) {
					refreshToken := helper.ExtractHeaderToken(c, helper.RefreshTokenHeader)
					// check if token is in blacklist or not, if it is then return to login page with error notification
					isInBlackList, _ := auth.IsTokenBlacklisted(refreshToken)

					if isInBlackList {
						rememberMe, _ := auth.RetrieveFieldFromJwtToken(accessToken, "remember_me", false)
						toastMessage := "your-token-is-invalid.-please-log-in-to-continue."
						helper.DeleteTokens(c, rememberMe.(bool))
						redirectToLogin(c, toastMessage)
						return
					}

					newAccessToken, newRerefreshToken, err := auth.CheckAndRefreshTokens(accessToken, refreshToken, helper.GetConfig(c))
					if err != nil {
						goto SESSION_EXPIRE
					}
					// don't change the tokens in header, just keep old tokens
					// don't save new tokens because eventually we will delete them
					if c.Request.URL.Path == "/v1/auth/logout" {
						c.Next()
						return
					}
					rememberMe, _ := auth.RetrieveFieldFromJwtToken(newAccessToken, "remember_me", true)
					var saveTo string
					if rememberMe.(bool) {
						saveTo = "local"
					} else {
						saveTo = "session"
					}
					HXTriggerEvents, err := helper.MapToJSONString(map[string]interface{}{
						"saveToken": map[string]interface{}{
							"saveTo":       saveTo,
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
				rememberMe, err := auth.RetrieveFieldFromJwtToken(accessToken, "remember_me", false)
				if err != nil {
					// delete tokens in both local storage and session storage when it is not valid token
					helper.DeleteTokens(c, true)
					helper.DeleteTokens(c, false)
					helper.HandleInternalError(c, err)
					return
				}
				helper.DeleteTokens(c, rememberMe.(bool))
				toastMessage := "your-session-has-expired.-please-log-in-to-continue."
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

// Determines if a path should be redirected to the home page
func shouldRedirectToHome(path string) bool {
	return path == "/v1/auth/login" || path == "/v1/auth/register"
}

func modifyAuthorizationHeaders(c *gin.Context, newAccessToken string, newRefreshToken string) {
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", newAccessToken))
	c.Request.Header.Set("Refresh", fmt.Sprintf("Bearer %s", newRefreshToken))
}
