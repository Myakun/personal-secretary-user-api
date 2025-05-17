package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawUser_ToUser(t *testing.T) {
	email := "test@example.com"
	id := "user123"
	name := "Test User"
	password := "password123"

	raw := &rawUser{
		Email:    email,
		Id:       id,
		Name:     name,
		Password: password,
	}

	user := raw.ToUser()

	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, id, user.Id)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, password, user.Password)
	assert.True(t, user.IsInserted())
}

func TestRawUser_ToUser_EmptyFields(t *testing.T) {
	raw := &rawUser{}

	user := raw.ToUser()

	assert.NotNil(t, user)
	assert.Empty(t, user.Email)
	assert.Empty(t, user.Id)
	assert.Empty(t, user.Name)
	assert.Empty(t, user.Password)
	assert.True(t, user.IsInserted())
}
