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
	ValidationErrorInvalidName        = errors.New("invalid_name")
	ValidationErrorInvalidPassword    = errors.New("invalid_password")
)

type userValidator struct {
	loggerService *logger.Logger
}

func (validator *userValidator) Validate(user *User) error {
	err := validator.ValidateName(user.GetName())
	if nil != err {
		return err
	}

	err = validator.ValidatePassword(user.GetPassword())
	if nil != err {
		return err
	}

	err = validator.ValidateEmail(user.GetEmail(), user)
	if nil != err {
		return err
	}

	return nil
}

func (validator *userValidator) ValidateEmail(email string, user *User) error {
	if !commonValidator.ValidateEmail(email) {
		return entity.NewValidationError(ValidationErrorInvalidEmail)
	}

	existingUser, err := GetUserService().FindOneByEmail(email)
	if nil != err {
		return fmt.Errorf("failed to validate user email: %w", err)
	}

	if nil != existingUser && (nil == user || user.GetId() != existingUser.GetId()) {
		return entity.NewValidationError(ValidationErrorEmailAlreadyExists)
	}

	return nil
}

func (validator *userValidator) ValidateName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return entity.NewValidationError(ValidationErrorInvalidName)
	}

	return nil
}

func (validator *userValidator) ValidatePassword(password string) error {
	// TODO: Implement strong password validation logic
	if len(password) < 6 {
		return entity.NewValidationError(ValidationErrorInvalidPassword)
	}

	return nil
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
