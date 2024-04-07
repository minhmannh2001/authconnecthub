package helper_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/stretchr/testify/assert"
)

// Define a custom error type for testing
type CustomError struct{}

func (e *CustomError) Error() string {
	return "custom error"
}

func TestIsErrOfType(t *testing.T) {
	t.Run("matches custom error", func(t *testing.T) {
		customErr := &CustomError{}
		anotherCustomErr := &CustomError{}
		assert.True(t, helper.IsErrOfType(anotherCustomErr, customErr))
	})

	t.Run("does not match different error types", func(t *testing.T) {
		customErr := &CustomError{}
		err := errors.New("different error")
		assert.False(t, helper.IsErrOfType(err, customErr))
	})
}

func TestHandleInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode) // Disable logging output for testing

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := errors.New("internal error")
	helper.HandleInternalError(c, err)

	// Assertions
	assert.Equal(t, "/500.html", c.Writer.Header().Get("HX-Redirect"))
	assert.Equal(t, 200, c.Writer.Status())
}
