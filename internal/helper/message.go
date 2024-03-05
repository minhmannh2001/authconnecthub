package helper

import "fmt"

func GenerateValidationMessage(field string, rule string) (message string) {
	switch rule {
	// required rule
	case "required":
		return fmt.Sprintf("Field '%s' is '%s'.", field, rule)
	// TODO: add another validator rule here
	case "min":
		if field == "Password" {
			message = "Password must be at least 8 characters long."
		}
		if field == "ConfirmPassword" {
			message = "Confirm password must be at least 8 characters long."
		}
		return message
	case "eqfield":
		if field == "ConfirmPassword" {
			message = "Passwords do not match. Please try again."
		}
		return message
	default:
		return fmt.Sprintf("Field '%s' is not valid.", field)
	}
}
