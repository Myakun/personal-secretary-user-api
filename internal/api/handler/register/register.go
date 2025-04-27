package register

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal-secretary-user-ap/internal/service/logger"
	"personal-secretary-user-ap/internal/service/user"
)

const (
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
		switch err {
		case user.ErrorInvalidEmail:
			errorResponse(context, ErrorCodeInvalidEmail)
			return
		case user.ErrorInvalidPassword:
			errorResponse(context, ErrorCodeInvalidPassword)
			return
		default:
			loggerService.DebugWithLogTag("Failed to register user: "+err.Error(), RegisterHandlerLogTag)
			context.Status(http.StatusInternalServerError)
			return
		}
	}

	// ✅ Passed validation
	// Here you would typically create user in database, generate token etc.
	context.JSON(http.StatusOK, registeredUser)
}

func errorResponse(context *gin.Context, error string) {
	context.JSON(http.StatusBadRequest, gin.H{"error": error})
}
