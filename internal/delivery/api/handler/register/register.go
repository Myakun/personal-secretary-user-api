package register

import (
	"net/http"

	"github.com/Myakun/personal-secretary-user-api/internal/presentation/user/registration"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	logger           logger.Logger
	userRegistration registration.UserRegistration
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const (
	logTag = "API_USER_REGISTER"
)

func NewRegisterHandler(logger logger.Logger, userRegistration registration.UserRegistration) *RegisterHandler {
	return &RegisterHandler{
		logger:           logger,
		userRegistration: userRegistration,
	}
}

func (h *RegisterHandler) Register(ctx *gin.Context) {
	var request registerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.DebugWithTagW("Failed to bind JSON", logTag, "error", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	result, err := h.userRegistration.RegisterUser(ctx.Request.Context(), registration.RegisterUserRequest{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	})

	if nil != err {
		h.logger.DebugWithTagW("Failed to register user", logTag, "error", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if result.Success {
		ctx.JSON(http.StatusOK, result.SuccessResponse)
		return
	}

	ctx.JSON(http.StatusUnprocessableEntity, result.ErrorResponse)
}
