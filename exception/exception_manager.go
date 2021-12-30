package exception

import (
	"fmt"
)

// UserException ... Represents a user input problem
type UserException struct {
	Message string `json:"message" example:"record not found"`
}

func (e UserException) Error() string {
	return e.Message
}

// CreateUserException ... Create a new exception for a user input problem
func CreateUserException(message string, data ...string) UserException {
	if len(data) > 0 {
		return UserException{
			Message: fmt.Sprintf(message, data),
		}
	}

	return UserException{
		Message: message,
	}
}

// AppException ... Represents a system exception
type AppException struct {
	Message string `json:"message" example:"internal server exception"`
}

func (e AppException) Error() string {
	return e.Message
}

// CreateAppException ... Create a new exception for an internal system problem
func CreateAppException(message string, data ...string) AppException {
	if len(data) > 0 {
		return AppException{
			Message: fmt.Sprintf(message, data),
		}
	}

	return AppException{
		Message: message,
	}
}

// UnauthorizedException ... Represents an unauthorized exception
type UnauthorizedException struct {
	Message string `json:"message" example:"invalid authentication token"`
}

func (e UnauthorizedException) Error() string {
	return e.Message
}

// CreateUnauthorizedException ... Create a new exception for unauthorized users
func CreateUnauthorizedException(message string, data ...string) UnauthorizedException {
	if len(data) > 0 {
		return UnauthorizedException{
			Message: fmt.Sprintf(message, data),
		}
	}

	return UnauthorizedException{
		Message: message,
	}
}

// AccessDeniedException ... Represents access denied exception
type AccessDeniedException struct {
	Message string `json:"message" example:"You are not allowed to access this resource"`
}

func (e AccessDeniedException) Error() string {
	return e.Message
}

// CreateAccessDeniedException ... Create a new exception for when access is denied
func CreateAccessDeniedException(message string, data ...string) AccessDeniedException {
	if len(data) > 0 {
		return AccessDeniedException{
			Message: fmt.Sprintf(message, data),
		}
	}

	return AccessDeniedException{
		Message: message,
	}
}

// ValidationException ... Represents an invalid request
type ValidationException struct {
	Message string `json:"message" example:"Validation failed"`
}

func (r ValidationException) Error() string {
	return r.Message
}

// CreateValidationException ... Create a new exception for invalid requests
func CreateValidationException(msg string, data ...string) ValidationException {
	return ValidationException{
		Message: fmt.Sprintf(msg, data),
	}
}
