package user

import (
	"strings"
	"testing"

	userEntity "github.com/Myakun/personal-secretary-user-api/internal/entity/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateName(t *testing.T) {
	validator := userEntity.GetUserValidator()

	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"too short", "A", true},
		{"too long", strings.Repeat("A", userEntity.NameMaxLength+1), true},
		{"valid", "Valid Name", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateName(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail_InvalidFormat(t *testing.T) {
	// Arrange
	validator := userEntity.GetUserValidator()
	require.NotNil(t, validator, "User validator should be initialized")

	email := "invalid-email" // Invalid email format

	// Act
	err := validator.ValidateEmail(email, nil)

	// Assert
	assert.Error(t, err)

	// Check if it's a validation error
	var validationErr *validator.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorInvalidEmail.Error(), validationErr.Error())
}

func TestValidateEmail_AlreadyExists(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	validator := userEntity.GetUserValidator()
	require.NotNil(t, service, "User service should be initialized")
	require.NotNil(t, validator, "User validator should be initialized")

	email := gofakeit.Email()
	name := "Test User"
	password := "password123"

	// Save first user
	firstUser := userEntity.NewUser(email, "", name, password)
	createdFirstUser, err := service.CreateUser(firstUser)
	require.NoError(t, err)
	require.NotNil(t, createdFirstUser)

	// Try to validate with same email for a different user
	secondUser := userEntity.NewUser(email, "", "Another User", "anotherpassword")

	// Act
	err = validator.ValidateEmail(email, secondUser)

	// Assert
	assert.Error(t, err)

	// Check if it's a validation error
	var validationErr *validator.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorEmailAlreadyExists.Error(), validationErr.Error())
}

func TestValidate_Success(t *testing.T) {
	// Arrange
	validator := userEntity.GetUserValidator()
	require.NotNil(t, validator, "User validator should be initialized")

	email := gofakeit.Email()
	name := gofakeit.Name()
	password := "password123"
	user := userEntity.NewUser(email, "", name, password)

	// Act
	err := validator.Validate(user)

	// Assert
	assert.NoError(t, err)
}
