package user

import (
	userService "github.com/Myakun/personal-secretary-user-api/internal/service/user"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	// Initialize JWT config for testing
	userService.InitJWTConfig("test-secret", 60)

	// First register a user
	email := gofakeit.Email()
	name := gofakeit.Name()
	password := "password123"
	registerRequest := userService.RegisterUserRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

	registeredUser, err := service.RegisterUser(registerRequest)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)

	// Now try to login with the registered user
	loginRequest := userService.LoginUserRequest{
		Email:    email,
		Password: password,
	}

	loginResult, err := service.LoginUser(loginRequest)

	assert.NoError(t, err)
	assert.NotNil(t, loginResult)
	assert.NotEmpty(t, loginResult.Token)
	assert.NotEmpty(t, loginResult.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	// Initialize JWT config for testing if not already initialized
	userService.InitJWTConfig("test-secret", 60)

	// Try to login with a non-existent user
	email := "nonexistent_" + gofakeit.Email()
	loginRequest := userService.LoginUserRequest{
		Email:    email,
		Password: "password123",
	}

	loginResult, err := service.LoginUser(loginRequest)

	assert.Error(t, err)
	assert.Nil(t, loginResult)

	// Check that the error is a LoginError with LoginErrorUserNotFound
	var loginErr *userService.LoginError
	assert.True(t, userService.LoginErrorUserNotFound.Error() == "user_not_found")
	assert.ErrorAs(t, err, &loginErr)
	assert.ErrorIs(t, loginErr, userService.LoginErrorUserNotFound)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	// Initialize JWT config for testing if not already initialized
	userService.InitJWTConfig("test-secret", 60)

	// First register a user
	email := gofakeit.Email()
	name := gofakeit.Name()
	password := "password123"
	registerRequest := userService.RegisterUserRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

	registeredUser, err := service.RegisterUser(registerRequest)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)

	// Now try to login with the wrong password
	loginRequest := userService.LoginUserRequest{
		Email:    email,
		Password: "wrongpassword",
	}

	loginResult, err := service.LoginUser(loginRequest)

	assert.Error(t, err)
	assert.Nil(t, loginResult)

	// Check that the error is a LoginError with LoginErrorInvalidCredentials
	var loginErr *userService.LoginError
	assert.True(t, userService.LoginErrorInvalidCredentials.Error() == "invalid_credentials")
	assert.ErrorAs(t, err, &loginErr)
	assert.ErrorIs(t, loginErr, userService.LoginErrorInvalidCredentials)
}
