package application

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"personal-secretary-user-ap/internal/entity/jwtrefreshtoken"
	"personal-secretary-user-ap/internal/entity/user"
	appUser "personal-secretary-user-ap/internal/service/user"
	"personal-secretary-user-ap/pkg/logger"
	"sync"
)

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

		// Entity services
		jwtrefreshtoken.InitJwtRefreshTokenService(appResources.mongo.Database)
		user.InitUserService(appResources.mongo.Database)

		// App services
		appUser.InitUserService(appConfig.JWT.ExpirationMin, appConfig.JWT.Secret)
	})

	return applicationInstance, getInstanceError
}
