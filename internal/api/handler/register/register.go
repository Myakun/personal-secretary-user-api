package register

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal-secretary-user-ap/internal/presentation/user/register"
	"personal-secretary-user-ap/internal/service/user"
	"personal-secretary-user-ap/pkg/logger"
)

const (
	logTag = "API_USER_REGISTER"
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
		loggerService.DebugWithLogTag("Failed to bind JSON: "+err.Error(), logTag)
		context.Status(http.StatusBadRequest)
		return
	}

	result, err := register.RegisterUser(user.RegisterUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		loggerService.FatalWithLogTag("Failed to register user: "+err.Error(), logTag)
		context.Status(http.StatusInternalServerError)
		return
	}

	if result.Success {
		context.JSON(http.StatusOK, result.SuccessResponse)
		return
	}

	context.JSON(http.StatusUnprocessableEntity, result.ErrorResponse)
}
