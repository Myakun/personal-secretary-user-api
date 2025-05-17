package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser_ValidInputs(t *testing.T) {
	email := "test@example.com"
	name := "Test User"
	password := "password123"

	user, err := NewUser(email, name, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, password, user.Password)
	assert.Empty(t, user.Id)
	assert.False(t, user.IsInserted())
}

func TestNewUser_EmptyEmail(t *testing.T) {
	email := ""
	name := "Test User"
	password := "password123"

	user, err := NewUser(email, name, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidEmail, err)
	assert.Nil(t, user)
}

func TestNewUser_InvalidEmail(t *testing.T) {
	email := "invalid-email"
	name := "Test User"
	password := "password123"

	user, err := NewUser(email, name, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidEmail, err)
	assert.Nil(t, user)
}

func TestNewUser_NameTooShort(t *testing.T) {
	email := "test@example.com"
	name := "Te"
	password := "password123"

	user, err := NewUser(email, name, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidName, err)
	assert.Nil(t, user)
}

func TestNewUser_NameTooLong(t *testing.T) {
	email := "test@example.com"
	name := "This name is way too long and exceeds the maximum length allowed for a name in this system"
	password := "password123"

	user, err := NewUser(email, name, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidName, err)
	assert.Nil(t, user)
}

func TestNewUser_PasswordTooShort(t *testing.T) {
	email := "test@example.com"
	name := "Test User"
	password := "pass"

	user, err := NewUser(email, name, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
	assert.Nil(t, user)
}

func TestNewUserFromStorage(t *testing.T) {
	email := "test@example.com"
	id := "user123"
	name := "Test User"
	password := "password123"

	user := NewUserFromStorage(email, id, name, password)

	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, id, user.Id)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, password, user.Password)
	assert.True(t, user.IsInserted())
}

func TestIsInserted(t *testing.T) {
	user, _ := NewUser("test@example.com", "Test User", "password123")

	assert.False(t, user.IsInserted())
}
