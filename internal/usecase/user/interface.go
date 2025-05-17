package user

import (
	"context"

	userDomain "github.com/Myakun/personal-secretary-user-api/internal/domain/user"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (*userDomain.User, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
