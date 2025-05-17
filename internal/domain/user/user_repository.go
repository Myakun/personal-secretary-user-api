package user

import "context"

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) (*User, error)
}
