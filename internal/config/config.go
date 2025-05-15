package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Api struct {
		Port int `envconfig:"APP_API_PORT" required:"true"`
	}

	Env string `envconfig:"APP_ENV" required:"true"`

	JWT struct {
		Secret        string `envconfig:"APP_JWT_SECRET" required:"true"`
		ExpirationMin int    `envconfig:"APP_JWT_EXPIRATION_MIN" required:"true"`
	}

	Mongo struct {
		Database string `envconfig:"APP_MONGO_DATABASE" required:"true"`
		Host     string `envconfig:"APP_MONGO_HOST" required:"true"`
		Password string `envconfig:"APP_MONGO_PASSWORD" required:"true"`
		Port     int    `envconfig:"APP_MONGO_PORT_APP" required:"true"`
		User     string `envconfig:"APP_MONGO_USER" required:"true"`
	}
}

func LoadFromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	config.Env = strings.ToLower(config.Env)

	return &config, nil
}
