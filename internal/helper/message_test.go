package helper_test

import (
	"testing"

	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestGenerateValidationMessage_Required(t *testing.T) {
	message := helper.GenerateValidationMessage("Email", "required")
	assert.Equal(t, "Field 'Email' is 'required'.", message)
}

func TestGenerateValidationMessage_Min_Password(t *testing.T) {
	message := helper.GenerateValidationMessage("Password", "min")
	assert.Equal(t, "'Password' must be at least 8 characters long.", message)
}

func TestGenerateValidationMessage_Eqfield_ConfirmPasswordMismatch(t *testing.T) {
	message := helper.GenerateValidationMessage("ConfirmPassword", "eqfield")
	assert.Equal(t, "Passwords do not match. Please try again.", message)
}

func TestGenerateValidationMessage_Default(t *testing.T) {
	message := helper.GenerateValidationMessage("Username", "unknown_rule")
	assert.Equal(t, "Field 'Username' is not valid.", message)
}
