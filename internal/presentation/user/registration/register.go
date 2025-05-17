package registration

import (
	"context"
	"errors"
	"fmt"

	userDomain "github.com/Myakun/personal-secretary-user-api/internal/domain/user"
	//	"github.com/Myakun/personal-secretary-user-api/internal/presentation/user/login"
	userUseCase "github.com/Myakun/personal-secretary-user-api/internal/usecase/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/Myakun/personal-secretary-user-api/pkg/validator"
)

const (
	errCodeEmailExists     = "email_exists"
	errCodeInvalidEmail    = "invalid_email"
	errCodeInvalidName     = "invalid_name"
	errCodeInvalidPassword = "invalid_password"
	logTag                 = "PRESENTATION_USER_REGISTER"
)

type userRegistration struct {
	logger      logger.Logger
	userUseCase userUseCase.UserUseCase
}

func NewUserRegistration(logger logger.Logger, userUseCase userUseCase.UserUseCase) UserRegistration {
	return &userRegistration{
		logger:      logger,
		userUseCase: userUseCase,
	}
}

func (ur *userRegistration) RegisterUser(ctx context.Context, request RegisterUserRequest) (*RegisterUserResult, error) {
	registeredUser, err := ur.userUseCase.CreateUser(ctx, userUseCase.CreateUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		ur.logger.DebugWithTagW("Failed to register user", logTag, "error", err)

		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			var response *ErrorResponse
			switch {
			case errors.Is(validationErr, userDomain.ErrInvalidEmail):
				response = newErrorResponse(errCodeInvalidEmail)
			case errors.Is(validationErr, userUseCase.ValidationErrEmailAlreadyExists):
				response = newErrorResponse(errCodeEmailExists)
			case errors.Is(validationErr, userDomain.ErrInvalidName):
				response = newErrorResponse(errCodeInvalidName)
			case errors.Is(validationErr, userDomain.ErrInvalidPassword):
				response = newErrorResponse(errCodeInvalidPassword)
			default:
				msg := fmt.Sprintf("unknown validation error: %v", validationErr)
				ur.logger.ErrorWithTag(msg, logTag)
				return nil, errors.New(msg)
			}

			return &RegisterUserResult{
				ErrorResponse: response,
				Success:       false,
			}, nil
		}

		return nil, err
	}

	/*loginResult, err := login.LoginUser(user.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})*/

	if nil != err {
		// TODO: handle error
	}

	fmt.Println(registeredUser)

	return &RegisterUserResult{
		Success: true,
		SuccessResponse: &SuccessResponse{
			Token: "13",
			//	User:         userEntityPackage.ConvertUserToDTo(registeredUser),
			//	Token:        loginResult.SuccessResponse.Token,
			//	RefreshToken: loginResult.SuccessResponse.RefreshToken,
		},
	}, nil
}
