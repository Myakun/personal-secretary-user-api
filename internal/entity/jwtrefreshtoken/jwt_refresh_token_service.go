package jwtrefreshtoken

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
	"sync"
	"time"
)

var jwtRefreshTokenServiceInstance *jwtRefreshTokenService
var jwtRefreshTokenServiceOnce sync.Once

type jwtRefreshTokenService struct {
	jwtRefreshTokenRepository *jwtRefreshTokenRepository
	jwtRefreshTokenValidator  *jwtRefreshTokenValidator
}

func (service *jwtRefreshTokenService) CreateJwtRefreshTokenToken(entity *JwtRefreshToken) (*JwtRefreshToken, error) {
	err := service.jwtRefreshTokenValidator.Validate(entity)
	if nil != err {
		return nil, err
	}

	entity, err = service.jwtRefreshTokenRepository.insert(entity)
	if nil != err {
		return nil, fmt.Errorf("failed to create jwt refresh token: %w", err)
	}

	entity.setIsInserted(true)

	return entity, nil
}

func GenerateJwtRefreshTokenForUserId(userId string) *JwtRefreshToken {
	now := time.Now()
	data := now.String() + userId
	token := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(data)).String()

	return NewJwtRefreshToken(now, token, userId)
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
		InitJwtRefreshTokenValidator(logger.GetLoggerService())

		jwtRefreshTokenServiceInstance = &jwtRefreshTokenService{
			jwtRefreshTokenRepository: &jwtRefreshTokenRepository{
				collection:    db.Collection(TableName),
				loggerService: logger.GetLoggerService(),
			},
			jwtRefreshTokenValidator: GetJwtRefreshTokenValidator(),
		}
	})
}
