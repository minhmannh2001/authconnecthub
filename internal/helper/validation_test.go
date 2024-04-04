package helper_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/undefinedlabs/go-mpatch"
)

type mockFieldError struct {
	validator.FieldError
	tag   string
	field string
}

func (e mockFieldError) Tag() string { return e.tag }

func (e mockFieldError) Field() string { return e.field }

func TestGenerateValidationResponse_SingleError(t *testing.T) {
	// Create a mock validation error
	mockValidationError := validator.ValidationErrors{
		mockFieldError{tag: "required", field: "name"},
	}

	// Call the function
	response := helper.GenerateValidationResponse(mockValidationError)

	// Assert expected response
	expectedResponse := dto.ValidationResponse{
		Success: false,
		Validations: []dto.Validation{
			{Field: "name", Message: "Field 'name' is 'required'."},
		},
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
	}
}

func TestGenerateValidationResponse_MultipleErrors(t *testing.T) {
	// Create mock validation errors
	mockValidationError := validator.ValidationErrors{
		mockFieldError{tag: "required", field: "name"},
		mockFieldError{tag: "email", field: "email"},
	}

	// Call the function
	response := helper.GenerateValidationResponse(mockValidationError)

	// Assert expected response
	expectedResponse := dto.ValidationResponse{
		Success: false,
		Validations: []dto.Validation{
			{Field: "name", Message: "Field 'name' is 'required'."},
			{Field: "email", Message: "Field 'email' is not valid."},
		},
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
	}
}

func TestGenerateValidationResponse_NilError(t *testing.T) {
	// Call the function with a nil error
	response := helper.GenerateValidationResponse(nil)

	// Assert expected response
	expectedResponse := dto.ValidationResponse{
		Success: false,
		Validations: []dto.Validation{
			{Field: "general", Message: "Unexpected validation error"}, // Customizable
		},
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
	}
}

func TestGenerateValidationResponse_InvalidErrorType(t *testing.T) {
	// Call the function with an error of a different type
	err := errors.New("some other error") // Use a different error type
	response := helper.GenerateValidationResponse(err)

	// Assert expected response (might depend on your error handling)
	expectedResponse := dto.ValidationResponse{
		Success: false,
		Validations: []dto.Validation{
			{Field: "general", Message: "Unexpected validation error"}, // Customizable
		},
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
	}
}

func TestGenerateValidationMap_SingleError(t *testing.T) {
	// Create a mock validation error
	mockValidationError := validator.ValidationErrors{
		mockFieldError{tag: "required", field: "name"},
	}

	// Call the function
	fieldErrors := helper.GenerateValidationMap(mockValidationError)

	// Assert expected errors
	expectedErrors := map[string]string{"name": "Field 'name' is 'required'."}
	if !reflect.DeepEqual(fieldErrors, expectedErrors) {
		t.Errorf("Expected errors: %v, got: %v", expectedErrors, fieldErrors)
	}
}

func TestGenerateValidationMap_MultipleErrorsSameField(t *testing.T) {
	// Create mock validation errors
	mockValidationError := validator.ValidationErrors{
		mockFieldError{tag: "required", field: "name"},
		mockFieldError{tag: "email", field: "email"},
	}

	// Call the function
	fieldErrors := helper.GenerateValidationMap(mockValidationError)

	// Assert expected errors
	expectedErrors := map[string]string{
		"name":  "Field 'name' is 'required'.",
		"email": "Field 'email' is not valid.",
	}
	if !reflect.DeepEqual(fieldErrors, expectedErrors) {
		t.Errorf("Expected errors: %v, got: %v", expectedErrors, fieldErrors)
	}
}

func TestGenerateValidationMap_NilError(t *testing.T) {
	// Call the function with a nil error
	fieldErrors := helper.GenerateValidationMap(nil)

	// Assert expected response
	expectedErrors := map[string]string{}
	if !reflect.DeepEqual(fieldErrors, expectedErrors) {
		t.Errorf("Expected response: %v, got: %v", expectedErrors, fieldErrors)
	}
}

func TestHashMap_Success(t *testing.T) {
	// Mock successful config retrieval
	mockCfg := &config.Config{Authen: config.Authen{SecretKey: "my_secret_key"}}
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return mockCfg, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Sample data
	data := map[string]interface{}{"key1": "value1", "key2": 10}

	// Call the function
	hashString, err := helper.HashMap(data)

	// Assert expected outcome
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate hash format (optional, adjust based on your needs)
	if len(hashString) != 64 { // Assuming SHA256 hash string length
		t.Errorf("Expected hash string length of 64, got: %d", len(hashString))
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashMap_ConfigError(t *testing.T) {
	// Mock config retrieval error
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return nil, errors.New("config error")
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sample data (unused in this test)
	data := map[string]interface{}{}

	// Call the function
	hashString, err := helper.HashMap(data)

	// Assert expected error
	expectedErrMsg := "error encoding data: config error"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected error: %s, got: %v", expectedErrMsg, err)
	}
	if hashString != "" {
		t.Errorf("Expected empty hash string on error, got: %s", hashString)
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashMap_NilMap(t *testing.T) {
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Call the function with nil map
	hashString, err := helper.HashMap(nil)

	// Assert expected error
	expectedErrMsg := "nil map cannot be hashed"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected error: %s, got: %v", expectedErrMsg, err)
	}
	if hashString != "" {
		t.Errorf("Expected empty hash string on error, got: %s", hashString)
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashMap_MarshalError(t *testing.T) {
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mock data that causes marshaling error (e.g., unsupported type)
	data := map[string]interface{}{"key1": func() {}}

	// Call the function
	hashString, err := helper.HashMap(data)

	// Assert expected error
	expectedErrMsg := "error encoding data: json: unsupported type: func()"
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Expected error: %s, got: %v", expectedErrMsg, err)
	}
	if hashString != "" {
		t.Errorf("Expected empty hash string on error, got: %s", hashString)
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsMapValid_Success(t *testing.T) {
	// Mock successful config retrieval
	mockCfg := &config.Config{Authen: config.Authen{SecretKey: "my_secret_key"}}
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return mockCfg, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sample data
	data := map[string]interface{}{"key1": "value1", "key2": 10}

	expectedHash, _ := helper.HashMap(data)

	// Call the function
	isValid := helper.IsMapValid(data, expectedHash)

	// Assert expected outcome
	if !isValid {
		t.Errorf("Expected map to be valid, but got false")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsMapValid_HashingError(t *testing.T) {
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return nil, errors.New("config error")
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sample data
	data := map[string]interface{}{"key1": "value1", "key2": 10}
	isValid := helper.IsMapValid(data, "abc")

	if isValid {
		t.Errorf("Expected map to be invalid due to hashing error, but got true")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsMapValid_MapMismatch(t *testing.T) {
	// Mock successful config retrieval
	mockCfg := &config.Config{Authen: config.Authen{SecretKey: "my_secret_key"}}
	patch, err := mpatch.PatchMethod(config.NewConfig, func() (*config.Config, error) {
		return mockCfg, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sample data and different expected hash
	data := map[string]interface{}{"key1": "value1"}
	expectedHash := "different_hash"

	// Call the function
	isValid := helper.IsMapValid(data, expectedHash)

	// Assert expected outcome
	if isValid {
		t.Errorf("Expected map to be invalid due to hash mismatch, but got true")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtractQueryParam_KeyPresent(t *testing.T) {
	queryParams := url.Values{"key1": {"value1"}, "key2": {"value2"}}
	key := "key1"
	defaultValue := "default" // Unused in this case

	extractedValue := helper.ExtractQueryParam(queryParams, key, defaultValue)
	expectedValue := "value1"

	if extractedValue != expectedValue {
		t.Errorf("Expected value for key %s: %s, got: %s", key, expectedValue, extractedValue)
	}
}

func TestExtractQueryParam_KeyMissing(t *testing.T) {
	queryParams := url.Values{"key2": {"value2"}}
	key := "key1"
	defaultValue := "default_value"

	extractedValue := helper.ExtractQueryParam(queryParams, key, defaultValue)

	if extractedValue != defaultValue {
		t.Errorf("Expected default value: %s, got: %s", defaultValue, extractedValue)
	}
}

func TestExtractQueryParam_EmptyQueryParams(t *testing.T) {
	queryParams := url.Values{}
	key := "any_key"
	defaultValue := "default_value"

	extractedValue := helper.ExtractQueryParam(queryParams, key, defaultValue)

	if extractedValue != defaultValue {
		t.Errorf("Expected default value: %s, got: %s", defaultValue, extractedValue)
	}
}

// https://github.com/undefinedlabs/go-mpatch
