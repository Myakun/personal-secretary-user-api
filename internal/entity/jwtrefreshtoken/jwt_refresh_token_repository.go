package jwtrefreshtoken

import (
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
)

type jwtRefreshTokenRepository struct {
	collection    *mongo.Collection
	loggerService *logger.Logger
}
