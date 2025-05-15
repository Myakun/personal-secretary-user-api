package login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	loginService "personal-secretary-user-ap/internal/presentation/user/login"
	"personal-secretary-user-ap/internal/service/user"
	"personal-secretary-user-ap/pkg/logger"
)

const (
	logTag = "API_USER_LOGIN"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(context *gin.Context) {
	loggerService := logger.GetLoggerService()

	var request loginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		loggerService.DebugWithLogTag("Failed to bind JSON: "+err.Error(), logTag)
		context.Status(http.StatusBadRequest)
		return
	}

	result, err := loginService.LoginUser(user.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})

	if nil != err {
		loggerService.FatalWithLogTag("Failed to login user: "+err.Error(), logTag)
		context.Status(http.StatusInternalServerError)
		return
	}

	if result.Success {
		context.JSON(http.StatusOK, result.SuccessResponse)
		return
	}

	context.JSON(http.StatusUnprocessableEntity, result.ErrorResponse)
}
