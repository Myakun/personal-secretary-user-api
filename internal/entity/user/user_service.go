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
	loggerService  *logger.Logger
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

func (service *userService) CreateUser(entity *User) (*User, error) {
	err := service.userValidator.Validate(entity)
	if nil != err {
		return nil, err
	}

	passwordHash, err := service.HashPassword(entity.GetPassword())
	if nil != err {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	entity.SetPassword(passwordHash)

	entity, err = service.userRepository.insert(entity)
	if nil != err {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	entity.setIsInserted(true)

	return entity, nil
}

func (service *userService) FindOneByEmail(email string) (*User, error) {
	entity, err := service.userRepository.FindOneByEmail(email)
	if nil != err {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return entity, nil
}

func (service *userService) FindOneById(id string) (*User, error) {
	entity, err := service.userRepository.FindOneById(id)
	if nil != err {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return entity, nil
}

func (service *userService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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
		loggerService := logger.GetLoggerService()

		userServiceInstance = &userService{
			loggerService: loggerService,
			userRepository: &userRepository{
				collection:    db.Collection(TableName),
				loggerService: loggerService,
			},
			userValidator: GetUserValidator(),
		}
	})
}
