package helper

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
)

func GenerateValidationResponse(err error) (response dto.ValidationResponse) {
	response.Success = false

	var validations []dto.Validation

	// get validation errors
	validationErrors := err.(validator.ValidationErrors)

	for _, value := range validationErrors {
		// get field & rule (tag)
		field, rule := value.Field(), value.Tag()

		// create validation object
		validation := dto.Validation{Field: field, Message: GenerateValidationMessage(field, rule)}

		// add validation object to validations
		validations = append(validations, validation)
	}

	// set Validations response
	response.Validations = validations

	return response
}

func GenerateValidationMap(err error) (fieldValidationErrors map[string]string) {
	fieldValidationErrors = make(map[string]string)

	if err == nil {
		return fieldValidationErrors // Handle no errors case
	}

	// get validation errors
	validationErrors := err.(validator.ValidationErrors)

	for _, value := range validationErrors {
		// Get field name and violated rule
		field, rule := value.Field(), value.Tag()
		message := GenerateValidationMessage(field, rule)

		if _, exists := fieldValidationErrors[field]; !exists {
			fieldValidationErrors[ToLowerFirstChar(field)] = message
		}
	}

	return fieldValidationErrors
}

func HashMap(data map[string]interface{}) (string, error) {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		return "", fmt.Errorf("error encoding data: %s", err.Error())
	}

	if data == nil {
		return "", fmt.Errorf("nil map cannot be hashed")
	}

	data["secretKey"] = cfg.Authen.SecretKey

	// Encode map value to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error encoding data: %s", err.Error())
	}

	// Hash the JSON data
	hash := sha256.Sum256(jsonData)

	// Convert hash to string representation (optional)
	hashString := fmt.Sprintf("%x", hash)

	return hashString, nil
}

func IsMapValid(data map[string]interface{}, expectedHash string) bool {
	actualHash, err := HashMap(data)
	if err != nil {
		log.Printf("Error hashing map: %v", err)
		return false
	}

	return actualHash == expectedHash
}

func ExtractQueryParam(queryParams url.Values, key, defaultValue string) string {
	value, ok := queryParams[key]
	if !ok {
		return defaultValue
	}
	return value[0]
}

func IsTokenExpired(err error) bool {
	return err != nil && strings.Contains(err.Error(), "expired")
}
