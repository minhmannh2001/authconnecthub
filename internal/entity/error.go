package entity

import "fmt"

type CustomErrorType interface {
	Error() string
}

type ErrDuplicateUser struct {
	Username string
	Email    string
}

func (e *ErrDuplicateUser) Error() string {
	errMsg := ""
	switch {
	case e.Username != "" && e.Email != "":
		errMsg = "User with the same username and email already exists"
	case e.Username != "":
		errMsg = fmt.Sprintf("Username '%s' already exists", e.Username)
	case e.Email != "":
		errMsg = fmt.Sprintf("Email '%s' already exists", e.Email)
	}
	return errMsg
}

type RoleNotFoundError struct {
	Name string
}

func (e *RoleNotFoundError) Error() string {
	return fmt.Sprintf("Role with name '%s' not found", e.Name)
}
