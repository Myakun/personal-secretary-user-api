package user

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestFindOneByEmail_Success(t *testing.T) {
	// Arrange
	service := userEntity.GetUserService()
	require.NotNil(t, service, "User service should be initialized")

	const amount = 10
	users := make([]*userEntity.User, amount)
	for i := 0; i < amount; i++ {
		email := gofakeit.Email()
		name := gofakeit.Name()
		password := "password123"
		user := userEntity.NewUser(email, "", name, password)
		user, err := service.CreateUser(user)
		assert.NoError(t, err)
		users[i] = user
	}

	for _, user := range users {
		foundUser, err := service.FindOneByEmail(user.GetEmail())
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.GetId(), foundUser.GetId())
		assert.True(t, foundUser.IsInserted())
	}
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
