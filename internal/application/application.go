package application

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"personal-secretary-user-ap/internal/common/entity"
	"personal-secretary-user-ap/internal/entity/accesstoken"
	"personal-secretary-user-ap/internal/entity/user"
	"personal-secretary-user-ap/internal/service/logger"
	appUser "personal-secretary-user-ap/internal/service/user"
	"sync"
)

const Version = "1.0.0"

type application struct {
	config    *config
	env       *env
	logger    *logger.Logger
	resources *resources
}

func (a *application) Close() {
	a.resources.close()
	//errors := 	a.resources.close()
	// TODO log if errors
}

func (a *application) GetConfig() *config {
	return a.config
}

func (a *application) GetLogger() *logger.Logger {
	return a.logger
}

func (a *application) GetEnv() *env {
	return a.env
}

func (a *application) GetResources() *resources {
	return a.resources
}

var applicationInstance *application
var getInstanceError error
var getInstanceOnce sync.Once

//goland:noinspection GoExportedFuncWithUnexportedType
func GetInstance(envFile *string) (*application, error) {
	getInstanceOnce.Do(func() {
		err := godotenv.Load(*envFile)
		if nil != err {
			log.Fatalf("Error loading %s file: %v", *envFile, err)
		}

		appConfig, err := configFromEnv()
		if nil != err {
			getInstanceError = fmt.Errorf("failed to load config: %w", err)
			return
		}

		appEnv, err := EnvFromString(appConfig.Env)
		if nil != err {
			getInstanceError = fmt.Errorf("failed to parse env: %w", err)
			return
		}

		loggerService, err := logger.InitLoggerService()
		if nil != err {
			getInstanceError = fmt.Errorf("failed to initialize logger service: %w", err)
			return
		}

		appResources, err := resourcesFromConfig(appConfig)
		if nil != err {
			getInstanceError = fmt.Errorf("failed to initialize resources: %w", err)
			return
		}

		applicationInstance = &application{
			config:    appConfig,
			logger:    loggerService,
			env:       appEnv,
			resources: appResources,
		}

		// Initialize entity services
		accesstoken.InitAccessTokenService(appResources.mongo.Database)
		user.InitUserService(appResources.mongo.Database)

		// Initialize app services
		appUser.InitUserService()

		userEntity := user.NewUser("ss@ss.ss", "", "ss", "ss")
		userEntity, err = user.GetUserService().CreateUser(userEntity)
		if nil != err {
			var validationErr *entity.ValidationError
			if errors.As(err, &validationErr) {
				fmt.Println("Validation error:", err)
			} else {
				fmt.Println("Other error:", err)
			}
			getInstanceError = fmt.Errorf("failed to create user: %w", err)
			return
		}
	})

	return applicationInstance, getInstanceError
}
