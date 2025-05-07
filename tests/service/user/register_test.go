package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"personal-secretary-user-ap/internal/common/entity"
	userEntity "personal-secretary-user-ap/internal/entity/user"
	userService "personal-secretary-user-ap/internal/service/user"
	"testing"
)

func TestConvertRegisterUserResponseErrorToDto_InvalidEmail(t *testing.T) {
	validationErr := entity.NewValidationError(userEntity.ValidationErrorInvalidEmail)
	responseErr, err := userService.ConvertRegisterUserResponseErrorToDto(validationErr)
	assert.NoError(t, err)
	assert.NotNil(t, responseErr)
	assert.Equal(t, "invalid_email", responseErr.Error)
}

func TestConvertRegisterUserResponseErrorToDto_EmailExists(t *testing.T) {
	validationErr := entity.NewValidationError(userEntity.ValidationErrorEmailAlreadyExists)
	responseErr, err := userService.ConvertRegisterUserResponseErrorToDto(validationErr)
	assert.NoError(t, err)
	assert.NotNil(t, responseErr)
	assert.Equal(t, "email_exists", responseErr.Error)
}

func TestConvertRegisterUserResponseErrorToDto_InvalidPassword(t *testing.T) {
	validationErr := entity.NewValidationError(userEntity.ValidationErrorInvalidPassword)
	responseErr, err := userService.ConvertRegisterUserResponseErrorToDto(validationErr)
	assert.NoError(t, err)
	assert.NotNil(t, responseErr)
	assert.Equal(t, "invalid_password", responseErr.Error)
}

func TestConvertRegisterUserResponseErrorToDto_NilError(t *testing.T) {
	responseErr, err := userService.ConvertRegisterUserResponseErrorToDto(nil)
	assert.Error(t, err)
	assert.Nil(t, responseErr)
	assert.Equal(t, "nil error", err.Error())
}

func TestConvertRegisterUserResponseErrorToDto_UnknownError(t *testing.T) {
	unknownErr := errors.New("some unknown error")
	responseErr, err := userService.ConvertRegisterUserResponseErrorToDto(unknownErr)
	assert.Error(t, err)
	assert.Nil(t, responseErr)
	assert.Contains(t, err.Error(), "unknown validation error")
}
