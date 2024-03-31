package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

var swaggerInfo *entity.SwaggerInfo
var once sync.Once

func IsPathMethodInSwagger(path string, method string, swaggerInfo *entity.SwaggerInfo) bool {
	pathItem, ok := swaggerInfo.Paths[path]
	if !ok {
		return false // we don't need to explicitly check for non-existent paths because Gin automatically handles them
	}

	// Get the operation based on the HTTP method
	var operation *entity.Operation
	switch method {
	case http.MethodGet:
		operation = pathItem.Get
	case http.MethodPost:
		operation = pathItem.Post
	case http.MethodPut:
		operation = pathItem.Put
	case http.MethodDelete:
		operation = pathItem.Delete
	case http.MethodPatch:
		operation = pathItem.Patch
	default:
		return false
	}

	return operation != nil // No operation defined for this path and method (not an error)
}

func HasSecurityKeyForPathAndMethod(path string, method string, swaggerInfo *entity.SwaggerInfo) (bool, error) {
	// Check if the path exists in the Swagger content
	pathItem, ok := swaggerInfo.Paths[path]
	if !ok {
		return false, nil // we don't need to explicitly check for non-existent paths because Gin automatically handles them
	}

	// Get the operation based on the HTTP method
	var operation *entity.Operation
	switch method {
	case http.MethodGet:
		operation = pathItem.Get
	case http.MethodPost:
		operation = pathItem.Post
	case http.MethodPut:
		operation = pathItem.Put
	case http.MethodDelete:
		operation = pathItem.Delete
	case http.MethodPatch:
		operation = pathItem.Patch
	default:
		return false, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Check if the operation exists and has security keys
	if operation == nil {
		return false, nil // No operation defined for this path and method (not an error)
	}
	return len(operation.Security) > 0, nil
}

func readSwaggerFile(filePath string) (*entity.SwaggerInfo, error) {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into the SwaggerInfo struct
	var info entity.SwaggerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &info, nil
}

// GetSwaggerInfo returns the singleton instance of SwaggerInfo, creating it if necessary
func GetSwaggerInfo(filePath string) (*entity.SwaggerInfo, error) {
	once.Do(func() {
		var err error
		swaggerInfo, err = readSwaggerFile(filePath)
		if err != nil {
			panic(err) // Handle error more gracefully in production
		}
	})
	return swaggerInfo, nil
}
