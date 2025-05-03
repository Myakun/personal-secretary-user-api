package user

import (
	"fmt"
	userEntityPackage "personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
)

var userServiceInstance *userService
var initUserServiceOnce sync.Once

type userService struct {
	logger *logger.Logger
}

type RegisterUserRequest struct {
	Email    string
	Name     string
	Password string
}

func (service *userService) Register(request RegisterUserRequest) (*userEntityPackage.User, error) {
	userEntity := userEntityPackage.NewUser(request.Email, "", request.Name, request.Password)
	userEntity, err := userEntityPackage.GetUserService().CreateUser(userEntity)
	if nil != err {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userEntity, nil
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetUserService() *userService {
	if nil == userServiceInstance {
		panic("user service is not initialized. Use InitUserService() to initialize.")
	}

	return userServiceInstance
}

func InitUserService() {
	initUserServiceOnce.Do(func() {
		userServiceInstance = &userService{
			logger: logger.GetLoggerService(),
		}
	})
}
