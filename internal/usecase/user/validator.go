package user

import (
	"context"
	"errors"
	"fmt"

	userDomain "github.com/Myakun/personal-secretary-user-api/internal/domain/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/validator"
)

var (
	ValidationErrEmailAlreadyExists = errors.New("email_already_exists")
)

func (uc *userUseCase) ValidateUser(ctx context.Context, user *userDomain.User) error {
	err := uc.ValidateUserEmail(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) ValidateUserEmail(ctx context.Context, user *userDomain.User) error {
	existingUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to validate user email: %w", err)
	}

	if nil != existingUser && user.Id != existingUser.Id {
		return validator.NewValidationError(ValidationErrEmailAlreadyExists)
	}

	return nil
}
