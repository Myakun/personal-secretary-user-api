package login

import (
	"errors"
	"fmt"
	"personal-secretary-user-ap/internal/service/user"
	"personal-secretary-user-ap/pkg/logger"
)

const (
	loginUserLogTag                  = "SERVICE_REQUEST_USER_LOGIN"
	loginErrorCodeInvalidEmail       = "invalid_email"
	loginErrorCodeInvalidCredentials = "invalid_credentials"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type loginUserResult struct {
	ErrorResponse   *ErrorResponse
	Success         bool
	SuccessResponse *SuccessResponse
}

func NewLoginUserResultWithErrorResponse(error string) *loginUserResult {
	return &loginUserResult{
		ErrorResponse: &ErrorResponse{Error: error},
		Success:       false,
	}
}

type SuccessResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginUser(request user.LoginUserRequest) (*loginUserResult, error) {
	loggerService := logger.GetLoggerService()

	loginResult, err := user.GetUserService().LoginUser(request)

	if nil != err {
		loggerService.DebugWithLogTag("Failed to login user: "+err.Error(), loginUserLogTag)

		var loginError *user.LoginError
		if errors.As(err, &loginError) {
			switch {
			case errors.Is(loginError, user.LoginErrorInvalidEmail):
				return NewLoginUserResultWithErrorResponse(loginErrorCodeInvalidEmail), nil
			case errors.Is(loginError, user.LoginErrorUserNotFound):
				return NewLoginUserResultWithErrorResponse(loginErrorCodeInvalidCredentials), nil
			case errors.Is(loginError, user.LoginErrorInvalidCredentials):
				return NewLoginUserResultWithErrorResponse(loginErrorCodeInvalidCredentials), nil
			default:
				msg := fmt.Sprintf("unknown login error: %v", loginError)
				loggerService.CriticalWithLogTag(msg, loginUserLogTag)
				return nil, errors.New(msg)
			}
		}

		return nil, err
	}

	// Convert the login result to a DTO presentation
	return &loginUserResult{
		Success: true,
		SuccessResponse: &SuccessResponse{
			Token:        loginResult.Token,
			RefreshToken: loginResult.RefreshToken,
		},
	}, nil
}
