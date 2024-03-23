package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenHeader  string = "Authorization"
	RefreshTokenHeader string = "Refresh"
)

func IsTokenExpired(err error) bool {
	return err != nil && strings.Contains(err.Error(), "expired")
}

// Extract token from request headers
func ExtractHeaderToken(c *gin.Context, tokenType string) string {
	headerToken := c.GetHeader(tokenType)
	if headerToken != "" && strings.HasPrefix(headerToken, "Bearer ") {
		parts := strings.Split(headerToken, " ")
		if len(parts) == 2 && parts[1] != "null" {
			return parts[1]
		}
	}
	return ""
}

func DeleteTokens(c *gin.Context, rememberMe bool) {
	deleteFrom := "session"
	if rememberMe {
		deleteFrom = "local"
	}
	HXTriggerEvents, _ := MapToJSONString(map[string]interface{}{
		"deleteToken": map[string]interface{}{
			"deleteFrom": deleteFrom,
		},
	})
	c.Header("HX-Trigger", HXTriggerEvents)
}
