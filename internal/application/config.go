package application

import (
	"github.com/kelseyhightower/envconfig"
	"strings"
)

type config struct {
	Api struct {
		Port int `envconfig:"APP_API_PORT" required:"true"`
	}

	Env string `envconfig:"APP_ENV" required:"true"`

	Mongo struct {
		Database string `envconfig:"APP_MONGO_DATABASE" required:"true"`
		Host     string `envconfig:"APP_MONGO_HOST" required:"true"`
		Password string `envconfig:"APP_MONGO_PASSWORD" required:"true"`
		Port     int    `envconfig:"APP_MONGO_PORT_APP" required:"true"`
		User     string `envconfig:"APP_MONGO_USER" required:"true"`
	}
}

//goland:noinspection GoExportedFuncWithUnexportedType
func configFromEnv() (*config, error) {
	var config config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	config.Env = strings.ToLower(config.Env)
	//config.Log.Level = strings.ToLower(config.Log.Level)

	return &config, nil
}
