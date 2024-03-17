package helper

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

func IsErrOfType(err, customErrorType error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(customErrorType)
}

func HandleInternalError(c *gin.Context, err error) {
	// Log the error
	log.Printf("Internal error: %v\n", err)

	c.Header("HX-Redirect", "/500.html")
	c.Abort()
}
