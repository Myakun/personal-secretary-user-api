package registration

import (
	"context"
)

type UserRegistration interface {
	RegisterUser(ctx context.Context, request RegisterUserRequest) (*RegisterUserResult, error)
}

type RegisterUserRequest struct {
	Email    string
	Name     string
	Password string
}
