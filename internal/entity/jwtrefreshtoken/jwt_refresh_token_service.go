package jwtrefreshtoken

import (
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
)

var jwtRefreshTokenServiceInstance *jwtRefreshTokenService
var jwtRefreshTokenServiceOnce sync.Once

type jwtRefreshTokenService struct {
	jwrRefreshTokenRepository *jwtRefreshTokenRepository
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetJwtRefreshTokenService() *jwtRefreshTokenService {
	if nil == jwtRefreshTokenServiceInstance {
		panic("jwt refresh token service is not initialized. Use InitJwtRefreshTokenService() to initialize.")
	}

	return jwtRefreshTokenServiceInstance
}

func InitJwtRefreshTokenService(db *mongo.Database) {
	jwtRefreshTokenServiceOnce.Do(func() {
		jwtRefreshTokenServiceInstance = &jwtRefreshTokenService{
			jwrRefreshTokenRepository: &jwtRefreshTokenRepository{
				collection:    db.Collection(TableName),
				loggerService: logger.GetLoggerService(),
			},
		}
	})
}
