package register

import (
	"errors"
	"fmt"
	"personal-secretary-user-ap/internal/common/entity"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/presentation/user/login"
	"personal-secretary-user-ap/internal/service/user"
	"personal-secretary-user-ap/pkg/logger"
)

const (
	registerUserLogTag               = "SERVICE_REQUEST_USER_REGISTER"
	registerErrorCodeEmailExists     = "email_exists"
	registerErrorCodeInvalidEmail    = "invalid_email"
	registerErrorCodeInvalidPassword = "invalid_password"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type registerUserResult struct {
	ErrorResponse   *ErrorResponse
	Success         bool
	SuccessResponse *SuccessResponse
}

type SuccessResponse struct {
	Token        string                     `json:"token"`
	RefreshToken string                     `json:"refresh_token"`
	User         *userEntityPackage.UserDTO `json:"user"`
}

func RegisterUser(request user.RegisterUserRequest) (*registerUserResult, error) {
	loggerService := logger.GetLoggerService()

	registeredUser, err := user.GetUserService().RegisterUser(user.RegisterUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		loggerService.DebugWithLogTag("Failed to register user: "+err.Error(), registerUserLogTag)

		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			var response *ErrorResponse
			switch {
			case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidEmail):
				response = &ErrorResponse{Error: registerErrorCodeInvalidEmail}
			case errors.Is(validationErr, userEntityPackage.ValidationErrorEmailAlreadyExists):
				response = &ErrorResponse{Error: registerErrorCodeEmailExists}
			case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidPassword):
				response = &ErrorResponse{Error: registerErrorCodeInvalidPassword}
			default:
				msg := fmt.Sprintf("unknown validation error: %v", validationErr)
				loggerService.CriticalWithLogTag(msg, registerUserLogTag)
				return nil, errors.New(msg)
			}

			return &registerUserResult{
				ErrorResponse: response,
				Success:       false,
			}, nil
		}

		return nil, err
	}

	loginResult, err := login.LoginUser(user.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})

	if nil != err {
		// TODO: handle error
	}

	return &registerUserResult{
		Success: true,
		SuccessResponse: &SuccessResponse{
			User:         userEntityPackage.ConvertUserToDTo(registeredUser),
			Token:        loginResult.SuccessResponse.Token,
			RefreshToken: loginResult.SuccessResponse.RefreshToken,
		},
	}, nil
}
