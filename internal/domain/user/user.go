package user

import (
	"errors"

	"github.com/Myakun/personal-secretary-user-api/pkg/validator"
)

type User struct {
	Email      string
	Id         string
	isInserted bool
	Name       string
	Password   string
}

var (
	ErrInvalidEmail    = errors.New("invalid_email")
	ErrInvalidName     = errors.New("invalid_name")
	ErrInvalidPassword = errors.New("invalid_password")
)

const (
	NameMinLength     = 3
	NameMaxLength     = 50
	PasswordMinLength = 6
)

func NewUser(email string, name string, password string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	if !validator.ValidateEmail(email) {
		return nil, ErrInvalidEmail
	}

	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return nil, ErrInvalidName
	}

	if len(password) < PasswordMinLength {
		return nil, ErrInvalidPassword
	}

	return &User{
		Email:      email,
		isInserted: false,
		Name:       name,
		Password:   password,
	}, nil
}

func NewUserFromStorage(email string, id string, name string, password string) *User {
	return &User{
		Email:      email,
		Id:         id,
		isInserted: true,
		Name:       name,
		Password:   password,
	}
}

func (u *User) IsInserted() bool {
	return u.isInserted
}
