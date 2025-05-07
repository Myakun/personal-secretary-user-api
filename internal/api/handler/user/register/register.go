package register

import (
	"github.com/gin-gonic/gin"
	"net/http"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
	"personal-secretary-user-ap/internal/service/user"
)

const (
	RegisterHandlerLogTag = "API_REGISTER_HANDLER"
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
		context.Status(http.StatusBadRequest)
		return
	}

	registeredUser, err := user.GetUserService().Register(user.RegisterUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		loggerService.DebugWithLogTag("Failed to register user: "+err.Error(), RegisterHandlerLogTag)

		response, err := user.ConvertRegisterUserResponseErrorToDto(err)
		if nil != err {
			loggerService.FatalWithLogTag("Failed to convert register user error to DTO: "+err.Error(), RegisterHandlerLogTag)
			context.Status(http.StatusInternalServerError)
			return
		}

		context.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	dto := userEntityPackage.ConvertUserToDTo(registeredUser)
	context.JSON(http.StatusOK, dto)
}
