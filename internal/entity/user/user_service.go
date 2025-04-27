package user

import (
	"errors"
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

func (service *userService) CreateUser(user *User) (*User, error) {
	isValid, err := service.userValidator.Validate(user)
	if !isValid {
		return nil, err
	}

	//passwordHash, err := service.HashPassword(user.GetPassword())

	return nil, nil
}

func (service *userService) FindOneByEmail(email string) (*User, error) {
	return nil, errors.New("not implemented")

	entity, err := service.userRepository.FindOneByEmail(email)
	if nil != err {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return entity, nil
}

func (service *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
				collection:    db.Collection(tableName),
				loggerService: logger.GetLoggerService(),
			},
			userValidator: GetUserValidator(),
		}
	})
}
