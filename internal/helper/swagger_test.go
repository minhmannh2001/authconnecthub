package helper_test

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestIsPathMethodInSwagger(t *testing.T) {
	cases := []struct {
		name        string
		path        string
		method      string
		swaggerInfo *entity.SwaggerInfo
		expected    bool
	}{
		{
			name:        "Path and method get exist",
			path:        "/users",
			method:      http.MethodGet,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expected:    true,
		},
		{
			name:        "Path and method post exist",
			path:        "/users",
			method:      http.MethodPost,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Post: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expected:    true,
		},
		{
			name:        "Path and method put exist",
			path:        "/users",
			method:      http.MethodPut,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Put: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expected:    true,
		},
		{
			name:        "Path and method delete exist",
			path:        "/users",
			method:      http.MethodDelete,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Delete: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expected:    true,
		},
		{
			name:        "Path and method patch exist",
			path:        "/users",
			method:      http.MethodPatch,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Patch: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expected:    true,
		},
		{
			name:        "Path exists but method doesn't",
			path:        "/users",
			method:      http.MethodPatch,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}},
			expected:    false,
		},
		{
			name:        "Path doesn't exist",
			path:        "/nonexistent",
			method:      http.MethodGet,
			swaggerInfo: &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}},
			expected:    false,
		},
		{
			name:        "SwaggerInfo is nil",
			path:        "/any",
			method:      http.MethodGet,
			swaggerInfo: nil,
			expected:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := helper.IsPathMethodInSwagger(tc.path, tc.method, tc.swaggerInfo)
			if actual != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestHasSecurityKeyForPathAndMethod(t *testing.T) {
	cases := []struct {
		name          string
		path          string
		method        string
		swaggerInfo   *entity.SwaggerInfo
		expectedHas   bool
		expectedError error
	}{
		{
			name:          "Path and method get exist, security keys present",
			path:          "/users",
			method:        http.MethodGet,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expectedHas:   true,
			expectedError: nil,
		},
		{
			name:          "Path and method post exist, security keys present",
			path:          "/users",
			method:        http.MethodPost,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Post: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expectedHas:   true,
			expectedError: nil,
		},
		{
			name:          "Path and method putexist, security keys present",
			path:          "/users",
			method:        http.MethodPut,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Put: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expectedHas:   true,
			expectedError: nil,
		},
		{
			name:          "Path and method delete exist, security keys present",
			path:          "/users",
			method:        http.MethodDelete,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Delete: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}},
			expectedHas:   true,
			expectedError: nil,
		},
		{
			name:          "Path exists, method doesn't, no error",
			path:          "/users",
			method:        http.MethodPatch,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}},
			expectedHas:   false,
			expectedError: nil,
		},
		{
			name:          "Path doesn't exist",
			path:          "/nonexistent",
			method:        http.MethodGet,
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}},
			expectedHas:   false,
			expectedError: nil,
		},
		{
			name:          "SwaggerInfo is nil",
			path:          "/any",
			method:        http.MethodGet,
			swaggerInfo:   nil,
			expectedHas:   false,
			expectedError: nil,
		},
		{
			name:          "Unsupported HTTP method",
			path:          "/users",
			method:        "INVALID",
			swaggerInfo:   &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}},
			expectedHas:   false,
			expectedError: fmt.Errorf("unsupported HTTP method: %s", "INVALID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			has, err := helper.HasSecurityKeyForPathAndMethod(tc.path, tc.method, tc.swaggerInfo)
			if has != tc.expectedHas {
				t.Errorf("Expected has: %v; got has: %v", tc.expectedHas, has)
			}
			if tc.expectedError != nil {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetSwaggerInfo_Success(t *testing.T) {
	fileContents := []byte(`{"paths": {"/users": {"get": {}}}}`)
	expectedInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{}}}}

	// Create a temporary file for the test
	tempFile, err := os.CreateTemp("", "swagger-test")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write file contents
	_, err = tempFile.Write(fileContents)
	if err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}

	// Read the file
	info, err := helper.GetSwaggerInfo(tempFile.Name())
	assert.Nil(t, err)
	if !reflect.DeepEqual(info, expectedInfo) {
		t.Errorf("Expected info: %v, got: %v", expectedInfo, info)
	}
}
