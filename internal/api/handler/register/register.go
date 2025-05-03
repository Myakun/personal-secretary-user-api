package register

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"personal-secretary-user-ap/internal/common/entity"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
	"personal-secretary-user-ap/internal/service/user"
)

const (
	ErrorCodeEmailExists     = "email_exists"
	ErrorCodeInvalidEmail    = "invalid_email"
	ErrorCodeInvalidInput    = "invalid_input"
	ErrorCodeInvalidPassword = "invalid_password"
	RegisterHandlerLogTag    = "REGISTER_HANDLER"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(context *gin.Context) {
	loggerService := logger.GetLoggerService()
	var request registerRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		loggerService.DebugWithLogTag("Failed to bind JSON: "+err.Error(), RegisterHandlerLogTag)
		errorResponse(context, ErrorCodeInvalidInput)
		return
	}

	registeredUser, err := user.GetUserService().Register(user.RegisterUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			switch {
			case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidEmail):
				errorResponse(context, ErrorCodeInvalidEmail)
				return
			case errors.Is(validationErr, userEntityPackage.ValidationErrorEmailAlreadyExists):
				errorResponse(context, ErrorCodeEmailExists)
				return
			case errors.Is(validationErr, userEntityPackage.ValidationErrorInvalidPassword):
				errorResponse(context, ErrorCodeInvalidPassword)
				return
			default:
				loggerService.CriticalWithLogTag("Unknown validation error: "+validationErr.Error(), RegisterHandlerLogTag)
				context.Status(http.StatusBadRequest)
				return
			}
		}

		loggerService.DebugWithLogTag("Failed to register user: "+err.Error(), RegisterHandlerLogTag)
		context.Status(http.StatusInternalServerError)
		return
	}

	dto := userEntityPackage.ConvertUserToDTo(registeredUser)
	context.JSON(http.StatusOK, dto)
}

func errorResponse(context *gin.Context, error string) {
	context.JSON(http.StatusBadRequest, gin.H{"error": error})
}
