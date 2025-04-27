package user

import (
	"errors"
	"fmt"
	"personal-secretary-user-ap/internal/common/entity"
	commonValidator "personal-secretary-user-ap/internal/common/validator"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
)

var userValidatorInstance *userValidator
var initUserValidatorOnce sync.Once

var (
	ValidationErrorEmailAlreadyExists = errors.New("email_already_exists")
	ValidationErrorInvalidEmail       = errors.New("invalid_email")
	ValidationErrorInvalidPassword    = errors.New("invalid_password")
)

type userValidator struct {
	loggerService *logger.Logger
}

func (validator *userValidator) Validate(user *User) (bool, error) {
	_, err := validator.ValidateEmail(user.GetEmail(), user)
	if nil != err {
		return false, err
	}

	return true, nil
}

func (validator *userValidator) ValidateEmail(email string, user *User) (bool, error) {
	if !commonValidator.ValidateEmail(email) {
		return false, entity.NewValidationError(ValidationErrorInvalidEmail)
	}

	existingUser, err := GetUserService().FindOneByEmail(email)
	if nil != err {
		return false, fmt.Errorf("failed to validate user email: %w", err)
	}

	if nil != existingUser && (nil == user || user.GetId() != existingUser.GetId()) {
		return false, entity.NewValidationError(ValidationErrorEmailAlreadyExists)
	}

	return true, nil
}

func (validator *userValidator) ValidatePassword(password string) (bool, error) {
	// TODO: Implement strong password validation logic
	if len(password) < 6 {
		return false, ValidationErrorInvalidPassword
	}

	return true, nil
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetUserValidator() *userValidator {
	if nil == userValidatorInstance {
		panic("user validator is not initialized. Use InitUserValidator() to initialize.")
	}

	return userValidatorInstance
}

func InitUserValidator(loggerService *logger.Logger) {
	initUserValidatorOnce.Do(func() {
		userValidatorInstance = &userValidator{
			loggerService: loggerService,
		}
	})
}
