package user

import (
	"errors"
	"fmt"
	"personal-secretary-user-ap/internal/common/jwt"
	"personal-secretary-user-ap/internal/common/validator"
	"personal-secretary-user-ap/internal/entity/jwtrefreshtoken"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
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
var userServiceInstance *userService
var initUserServiceOnce sync.Once

type userService struct {
	jwtConfig *JWTConfig
	logger    *logger.Logger
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token        string
	RefreshToken string
}

func (service *userService) LoginUser(request LoginUserRequest) (*LoginResult, error) {
	if !validator.ValidateEmail(request.Email) {
		return nil, NewLoginError(LoginErrorInvalidEmail)
	}

	user, err := userEntityPackage.GetUserService().FindOneByEmail(request.Email)
	if nil != err {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if nil == user {
		return nil, NewLoginError(LoginErrorUserNotFound)
	}

	err = userEntityPackage.GetUserService().VerifyPassword(user.GetPassword(), request.Password)
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

type RegisterUserRequest struct {
	Email    string
	Name     string
	Password string
}

func (service *userService) RegisterUser(request RegisterUserRequest) (*userEntityPackage.User, error) {
	userEntity := userEntityPackage.NewUser(request.Email, "", request.Name, request.Password)
	userEntity, err := userEntityPackage.GetUserService().CreateUser(userEntity)
	if nil != err {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userEntity, nil
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetUserService() *userService {
	if nil == userServiceInstance {
		panic("user service is not initialized. Use InitUserService() to initialize.")
	}

	return userServiceInstance
}

// InitJWTConfig initializes the JWT configuration
func InitJWTConfig(secret string, expirationMin int) {
	jwtConfig = &JWTConfig{
		Secret:        secret,
		ExpirationMin: expirationMin,
	}
}

func InitUserService(jwtExpirationMin int, jwtSecret string) {
	initUserServiceOnce.Do(func() {
		userServiceInstance = &userService{
			jwtConfig: &JWTConfig{
				ExpirationMin: jwtExpirationMin,
				Secret:        jwtSecret,
			},
			logger: logger.GetLoggerService(),
		}
	})
}
