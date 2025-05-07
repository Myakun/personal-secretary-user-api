package user

import (
	"errors"
	"fmt"
	"personal-secretary-user-ap/internal/common/entity"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
)

const (
	registerErrorCodeEmailExists     = "email_exists"
	registerErrorCodeInvalidEmail    = "invalid_email"
	registerErrorCodeInvalidPassword = "invalid_password"
)

type RegisterUserRequest struct {
	Email    string
	Name     string
	Password string
}

type RegisterUserResponseError struct {
	Error string `json:"error"`
}

func ConvertRegisterUserResponseErrorToDto(err error) (*RegisterUserResponseError, error) {
	if nil == err {
		return nil, errors.New("nil error")
	}

	var validationErr *entity.ValidationError
	if errors.As(err, &validationErr) {
		switch {
		case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidEmail):
			return &RegisterUserResponseError{Error: registerErrorCodeInvalidEmail}, nil
		case errors.Is(validationErr, userEntityPackage.ValidationErrorEmailAlreadyExists):
			return &RegisterUserResponseError{Error: registerErrorCodeEmailExists}, nil
		case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidPassword):
			return &RegisterUserResponseError{Error: registerErrorCodeInvalidPassword}, nil
		}
	}

	msg := fmt.Sprintf("unknown validation error: %v", validationErr)
	logger.GetLoggerService().Critical(msg)

	return nil, errors.New(msg)
}
