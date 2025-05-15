package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setTestEnv() {
	os.Setenv("APP_API_PORT", "8080")
	os.Setenv("APP_ENV", "TEST")
	os.Setenv("APP_JWT_SECRET", "test-secret")
	os.Setenv("APP_JWT_EXPIRATION_MIN", "15")
	os.Setenv("APP_MONGO_DATABASE", "test-db")
	os.Setenv("APP_MONGO_HOST", "localhost")
	os.Setenv("APP_MONGO_PASSWORD", "password")
	os.Setenv("APP_MONGO_PORT_APP", "27017")
	os.Setenv("APP_MONGO_USER", "user")
}

func cleanupTestEnv(origEnv []string) {
	os.Clearenv()
	for _, e := range origEnv {
		parts := strings.SplitN(e, "=", 2)
		os.Setenv(parts[0], parts[1])
	}
}

func TestLoadFromEnv_Success(t *testing.T) {
	origEnv := os.Environ()
	setTestEnv()
	defer cleanupTestEnv(origEnv)

	cfg, err := LoadFromEnv()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Api.Port)
	assert.Equal(t, "test", cfg.Env) // Should be lowercase
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Equal(t, 15, cfg.JWT.ExpirationMin)
	assert.Equal(t, "test-db", cfg.Mongo.Database)
	assert.Equal(t, "localhost", cfg.Mongo.Host)
	assert.Equal(t, "password", cfg.Mongo.Password)
	assert.Equal(t, 27017, cfg.Mongo.Port)
	assert.Equal(t, "user", cfg.Mongo.User)
}

func TestLoadFromEnv_Error(t *testing.T) {
	origEnv := os.Environ()
	os.Clearenv()
	defer cleanupTestEnv(origEnv)

	cfg, err := LoadFromEnv()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadFromEnv_EnvToLowercase(t *testing.T) {
	origEnv := os.Environ()
	setTestEnv()
	defer cleanupTestEnv(origEnv)

	cfg, err := LoadFromEnv()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.Env)
}
