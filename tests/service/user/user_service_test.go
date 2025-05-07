package user

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	userEntity "personal-secretary-user-ap/internal/entity/user"
	userService "personal-secretary-user-ap/internal/service/user"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := gofakeit.Name()
	password := "password123"
	request := userService.RegisterUserRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

	user, err := service.Register(request)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.GetEmail())
	assert.Equal(t, name, user.GetName())
	assert.NotEqual(t, password, user.GetPassword()) // Password should be hashed
	assert.True(t, user.IsInserted())

	// Verify the user was actually created in the database
	foundUser, err := userEntity.GetUserService().FindOneByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.GetId(), foundUser.GetId())
}

func TestRegister_DuplicateEmail(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := gofakeit.Name()
	password := "password123"
	request := userService.RegisterUserRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

	firstUser, err := service.Register(request)
	require.NoError(t, err)
	require.NotNil(t, firstUser)

	request.Name = "Another User"
	request.Password = "anotherpassword"

	secondUser, err := service.Register(request)

	assert.Error(t, err)
	assert.Nil(t, secondUser)
}

func TestRegister_InvalidEmail(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	request := userService.RegisterUserRequest{
		Email:    "invalid-email",
		Name:     gofakeit.Name(),
		Password: "password123",
	}

	user, err := service.Register(request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestRegister_InvalidPassword(t *testing.T) {
	service := userService.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	request := userService.RegisterUserRequest{
		Email:    gofakeit.Email(),
		Name:     gofakeit.Name(),
		Password: "12345",
	}

	user, err := service.Register(request)

	assert.Error(t, err)
	assert.Nil(t, user)
}
