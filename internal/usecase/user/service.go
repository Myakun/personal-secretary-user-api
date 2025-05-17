package user

import (
	"context"
	"errors"
	"fmt"

	userDomain "github.com/Myakun/personal-secretary-user-api/internal/domain/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/Myakun/personal-secretary-user-api/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	LoginErrorInvalidCredentials = errors.New("invalid_credentials")
	LoginErrorInvalidEmail       = errors.New("invalid_email")
	LoginErrorUserNotFound       = errors.New("user_not_found")
)

type LoginError struct {
	Err error
}

func NewLoginError(err error) *LoginError {
	return &LoginError{
		Err: err,
	}
}

func (e *LoginError) Error() string {
	return e.Err.Error()
}

func (e *LoginError) Unwrap() error {
	return e.Err
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	ExpirationMin int
}

var jwtConfig *JWTConfig

type userUseCase struct {
	jwtConfig *JWTConfig
	logger    logger.Logger
	userRepo  userDomain.UserRepository
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token        string
	RefreshToken string
}

/*
	func (service *UserUseCase) LoginUser(request LoginUserRequest) (*LoginResult, error) {
		if !validator.ValidateEmail(request.Email) {
			return nil, NewLoginError(LoginErrorInvalidEmail)
		}

		user, err := userEntity.GetUserService().FindOneByEmail(request.Email)
		if nil != err {
			return nil, fmt.Errorf("failed to find user by email: %w", err)
		}

		if nil == user {
			return nil, NewLoginError(LoginErrorUserNotFound)
		}

		err = userEntity.GetUserService().VerifyPassword(user.GetPassword(), request.Password)
		if nil != err {
			return nil, NewLoginError(LoginErrorInvalidCredentials)
		}

		if jwtConfig == nil {
			return nil, fmt.Errorf("JWT configuration is not initialized")
		}

		// Generate JWT token without user information
		jwtToken, err := jwt.GenerateToken(
			user.GetId(),
			"", // Remove email from token
			jwtConfig.Secret,
			jwtConfig.ExpirationMin,
		)
		if nil != err {
			return nil, fmt.Errorf("failed to generate JWT token: %w", err)
		}

		// Generate refresh token
		refreshToken := jwtrefreshtoken.GenerateJwtRefreshTokenForUserId(user.GetId())
		refreshToken, err = jwtrefreshtoken.GetJwtRefreshTokenService().CreateJwtRefreshTokenToken(refreshToken)
		if nil != err {
			return nil, fmt.Errorf("failed to create jwt refresh token: %w", err)
		}

		return &LoginResult{
			Token:        jwtToken,
			RefreshToken: refreshToken.GetToken(),
		}, nil
	}
*/
type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}

func NewUserUseCase(logger logger.Logger, userRepo userDomain.UserRepository) UserUseCase {
	return &userUseCase{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (uc *userUseCase) CreateUser(ctx context.Context, request CreateUserRequest) (*userDomain.User, error) {
	entity, err := userDomain.NewUser(request.Email, request.Name, request.Password)
	if err != nil {
		return nil, validator.NewValidationError(err)
	}

	err = uc.ValidateUser(ctx, entity)
	if err != nil {
		return nil, err
	}

	passwordHash, err := uc.HashPassword(entity.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	entity.Password = passwordHash

	entity, err = uc.userRepo.Save(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return entity, nil
}

func (uc *userUseCase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (uc *userUseCase) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// InitJWTConfig initializes the JWT configuration
func InitJWTConfig(secret string, expirationMin int) {
	jwtConfig = &JWTConfig{
		Secret:        secret,
		ExpirationMin: expirationMin,
	}
}
