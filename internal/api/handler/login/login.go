package login

import (
	"net/http"

	loginService "github.com/Myakun/personal-secretary-user-api/internal/presentation/user/login"
	"github.com/Myakun/personal-secretary-user-api/internal/service/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"

	"github.com/gin-gonic/gin"
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
