package accesstoken

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
)

var accessTokenServiceInstance *accessTokenService
var initAccessTokenServiceOnce sync.Once

type accessTokenService struct {
	accessTokenRepository accessTokenRepository
}

func (service *accessTokenService) FindOneByToken(token string) (*AccessToken, error) {
	entity, err := service.accessTokenRepository.FindOneByToken(token)
	if nil != err {
		return nil, fmt.Errorf("failed to find access token by token: %w", err)
	}

	return entity, nil
}

func InitAccessTokenService(db *mongo.Database) {
	initAccessTokenServiceOnce.Do(func() {
		accessTokenServiceInstance = &accessTokenService{
			accessTokenRepository: accessTokenRepository{
				collection: db.Collection("mobile_app_access_tokens"),
				logger:     logger.GetLoggerService(),
			},
		}
	})
}

//goland:noinspection GoExportedFuncWithUnexportedType
func GetAccessTokenService() *accessTokenService {
	if accessTokenServiceInstance == nil {
		panic("AccessToken service is not initialized. Use InitAccessTokenService() to initialize.")
	}

	return accessTokenServiceInstance
}
