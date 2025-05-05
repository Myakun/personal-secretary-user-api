package user

import (
	"github.com/stretchr/testify/assert"
	userEntity "personal-secretary-user-ap/internal/entity/user"
	"testing"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	id := "user123"
	name := "Test User"
	password := "password123"

	user := userEntity.NewUser(email, id, name, password)

	assert.Equal(t, email, user.GetEmail())
	assert.Equal(t, id, user.GetId())
	assert.Equal(t, name, user.GetName())
	assert.Equal(t, password, user.GetPassword())
	assert.False(t, user.IsInserted())
}

func TestGetEmail(t *testing.T) {
	email := "test@example.com"
	user := userEntity.NewUser(email, "id", "name", "password")

	result := user.GetEmail()

	assert.Equal(t, email, result)
}

func TestGetId(t *testing.T) {
	id := "user123"
	user := userEntity.NewUser("email", id, "name", "password")

	result := user.GetId()

	assert.Equal(t, id, result)
}

func TestGetName(t *testing.T) {
	name := "Test User"
	user := userEntity.NewUser("email", "id", name, "password")

	result := user.GetName()

	assert.Equal(t, name, result)
}

func TestGetPassword(t *testing.T) {
	password := "password123"
	user := userEntity.NewUser("email", "id", "name", password)

	result := user.GetPassword()

	assert.Equal(t, password, result)
}

func TestIsInserted(t *testing.T) {
	user := userEntity.NewUser("email", "id", "name", "password")

	assert.False(t, user.IsInserted())
}

func TestSetPassword(t *testing.T) {
	initialPassword := "initial"
	newPassword := "newPassword123"
	user := userEntity.NewUser("email", "id", "name", initialPassword)

	user.SetPassword(newPassword)

	assert.Equal(t, newPassword, user.GetPassword())
}
