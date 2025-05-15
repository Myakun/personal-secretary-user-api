package jwtrefreshtoken

import (
	"sync"

	"github.com/Myakun/personal-secretary-user-api/pkg/logger"
)

var jwtRefreshTokenValidatorInstance *jwtRefreshTokenValidator
var initJwtRefreshTokenValidatorOnce sync.Once

type jwtRefreshTokenValidator struct {
	loggerService *logger.Logger
}

func (validator *jwtRefreshTokenValidator) Validate(jwtRefreshToken *JwtRefreshToken) error {
	// Empty implementation as per requirements
	return nil
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetJwtRefreshTokenValidator() *jwtRefreshTokenValidator {
	if nil == jwtRefreshTokenValidatorInstance {
		panic("jwt refresh token validator is not initialized. Use InitJwtRefreshTokenValidator() to initialize.")
	}

	return jwtRefreshTokenValidatorInstance
}

func InitJwtRefreshTokenValidator(loggerService *logger.Logger) {
	initJwtRefreshTokenValidatorOnce.Do(func() {
		jwtRefreshTokenValidatorInstance = &jwtRefreshTokenValidator{
			loggerService: loggerService,
		}
	})
}
