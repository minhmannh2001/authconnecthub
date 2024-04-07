package helper_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestIsTokenExpired_ExpiredError(t *testing.T) {
	expiredErr := errors.New("token is expired")
	assert.True(t, helper.IsTokenExpired(expiredErr))
}

func TestIsTokenExpired_NonExpiredError(t *testing.T) {
	nonExpiredErr := errors.New("some other error")
	assert.False(t, helper.IsTokenExpired(nonExpiredErr))
}

func TestIsTokenExpired_NoError(t *testing.T) {
	assert.False(t, helper.IsTokenExpired(nil))
}

func TestExtractHeaderToken(t *testing.T) {
	gin.SetMode(gin.TestMode) // Set Gin to test mode for all tests

	testCases := []struct {
		name          string
		headerValue   string
		expectedToken string
	}{
		{
			name:          "Valid Bearer Token",
			headerValue:   "Bearer some_valid_token",
			expectedToken: "some_valid_token",
		},
		{
			name:          "Missing Bearer Prefix",
			headerValue:   "some_valid_token",
			expectedToken: "",
		},
		{
			name:          "Empty Token",
			headerValue:   "Bearer ",
			expectedToken: "",
		},
		{
			name:          "Missing Header",
			headerValue:   "",
			expectedToken: "",
		},
		{
			name:          "Null Token",
			headerValue:   "Bearer null",
			expectedToken: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = &http.Request{}
			c.Request.Header = http.Header{}
			if tc.headerValue != "" {
				c.Request.Header.Set(helper.AccessTokenHeader, tc.headerValue)
			}

			token := helper.ExtractHeaderToken(c, helper.AccessTokenHeader)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestDeleteTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Define test cases
	testCases := []struct {
		rememberMe bool
		expected   string
	}{
		{false, `{"deleteToken":{"deleteFrom":"session"}}`},
		{true, `{"deleteToken":{"deleteFrom":"local"}}`},
	}

	for _, tc := range testCases {
		mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		mockContext.Request = &http.Request{}
		mockContext.Request.Header = http.Header{}

		helper.DeleteTokens(mockContext, tc.rememberMe, false)

		assert.Equal(t, tc.expected, mockContext.Writer.Header().Get("HX-Trigger"))
	}
}
