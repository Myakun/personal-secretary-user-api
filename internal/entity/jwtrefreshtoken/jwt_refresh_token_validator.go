package jwtrefreshtoken

import (
	"personal-secretary-user-ap/pkg/logger"
	"sync"
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
