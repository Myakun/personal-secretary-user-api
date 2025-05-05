package user

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"personal-secretary-user-ap/internal/common/entity"
	userEntity "personal-secretary-user-ap/internal/entity/user"
	"testing"
)

func TestConvertUserToDTo(t *testing.T) {
	// Arrange
	email := "test@example.com"
	id := "user123"
	name := "Test User"
	password := "password123"
	user := userEntity.NewUser(email, id, name, password)

	// Act
	dto := userEntity.ConvertUserToDTo(user)

	// Assert
	assert.Equal(t, email, dto.Email)
	assert.Equal(t, id, dto.Id)
	assert.Equal(t, name, dto.Name)
}

func TestCreateUser_Success(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := "Test User"
	password := "password123"
	user := userEntity.NewUser(email, "", name, password)

	// Act
	createdUser, err := service.CreateUser(user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, email, createdUser.GetEmail())
	assert.Equal(t, name, createdUser.GetName())
	assert.NotEmpty(t, createdUser.GetId())
	assert.NotEqual(t, password, createdUser.GetPassword()) // Password should be hashed
	assert.True(t, createdUser.IsInserted())

	// Cleanup - not needed as we're using a test database that will be reset
}

func TestCreateUser_InvalidName_TooShort(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := "Te" // Too short (less than NameMinLength)
	password := "password123"
	user := userEntity.NewUser(email, "", name, password)

	// Act
	createdUser, err := service.CreateUser(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)

	// Check if it's a validation error
	var validationErr *entity.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorInvalidName.Error(), validationErr.Error())
}

func TestCreateUser_InvalidName_TooLong(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := gofakeit.LetterN(uint(userEntity.NameMaxLength + 1))
	password := "password123"
	user := userEntity.NewUser(email, "", name, password)

	// Act
	createdUser, err := service.CreateUser(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)

	// Check if it's a validation error
	var validationErr *entity.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorInvalidName.Error(), validationErr.Error())
}

func TestCreateUser_InvalidPassword(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := "Test User"
	password := "12345" // Too short (less than 6 characters)
	user := userEntity.NewUser(email, "", name, password)

	// Act
	createdUser, err := service.CreateUser(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)

	// Check if it's a validation error
	var validationErr *entity.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorInvalidPassword.Error(), validationErr.Error())
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := "invalid-email" // Invalid email format
	name := "Test User"
	password := "password123"
	user := userEntity.NewUser(email, "", name, password)

	// Act
	createdUser, err := service.CreateUser(user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdUser)

	// Check if it's a validation error
	var validationErr *entity.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorInvalidEmail.Error(), validationErr.Error())
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := "Test User"
	password := "password123"

	// Create first user
	firstUser := userEntity.NewUser(email, "", name, password)
	createdFirstUser, err := service.CreateUser(firstUser)
	require.NoError(t, err)
	require.NotNil(t, createdFirstUser)

	// Try to create second user with same email
	secondUser := userEntity.NewUser(email, "", "Another User", "anotherpassword")

	// Act
	createdSecondUser, err := service.CreateUser(secondUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdSecondUser)

	// Check if it's a validation error
	var validationErr *entity.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, userEntity.ValidationErrorEmailAlreadyExists.Error(), validationErr.Error())
}

func TestFindOneByEmail_Success(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := gofakeit.Email()
	name := "Test User"
	password := "password123"

	// Create user first
	user := userEntity.NewUser(email, "", name, password)
	createdUser, err := service.CreateUser(user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	// Act
	foundUser, err := service.FindOneByEmail(email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, email, foundUser.GetEmail())
	assert.Equal(t, name, foundUser.GetName())
	assert.Equal(t, createdUser.GetId(), foundUser.GetId())
	assert.True(t, foundUser.IsInserted())
}

func TestFindOneByEmail_NotFound(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	email := "nonexistent_" + gofakeit.Email()

	// Act
	foundUser, err := service.FindOneByEmail(email)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, foundUser)
}

func TestHashPassword(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	password := "password123"

	// Act
	hashedPassword, err := service.HashPassword(password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// The hash should be different each time
	hashedPassword2, err := service.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, hashedPassword2)
}
