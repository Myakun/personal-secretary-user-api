package user

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
)

var userServiceInstance *userService
var initUserServiceOnce sync.Once

type userService struct {
	userRepository *userRepository
	userValidator  *userValidator
}

func ConvertUserToDTo(user *User) *UserDTO {
	return &UserDTO{
		Email: user.GetEmail(),
		Id:    user.GetId(),
		Name:  user.GetName(),
	}
}

func (service *userService) CreateUser(user *User) (*User, error) {
	err := service.userValidator.Validate(user)
	if nil != err {
		return nil, err
	}

	passwordHash, err := service.HashPassword(user.GetPassword())
	if nil != err {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.SetPassword(passwordHash)

	user, err = service.userRepository.insert(user)
	if nil != err {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.setIsInserted(true)

	return user, nil
}

func (service *userService) FindOneByEmail(email string) (*User, error) {
	entity, err := service.userRepository.FindOneByEmail(email)
	if nil != err {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return entity, nil
}

func (service *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if nil != err {
		return "", err
	}

	return string(bytes), err
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetUserService() *userService {
	if nil == userServiceInstance {
		panic("user service is not initialized. Use InitUserService() to initialize.")
	}

	return userServiceInstance
}

func InitUserService(db *mongo.Database) {
	initUserServiceOnce.Do(func() {
		InitUserValidator(logger.GetLoggerService())

		userServiceInstance = &userService{
			userRepository: &userRepository{
				collection:    db.Collection(TableName),
				loggerService: logger.GetLoggerService(),
			},
			userValidator: GetUserValidator(),
		}
	})
}
