package jwtrefreshtoken

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
	"time"
)

type jwtRefreshTokenRaw struct {
	CreatedAt time.Time `bson:"created_at"`
	Id        string    `bson:"_id"`
	UserId    string    `bson:"user_id"`
	Token     string    `bson:"token"`
}

type jwtRefreshTokenRepository struct {
	collection    *mongo.Collection
	loggerService *logger.Logger
}

func (repository *jwtRefreshTokenRepository) insert(entity *JwtRefreshToken) (*JwtRefreshToken, error) {
	raw := &jwtRefreshTokenRaw{
		CreatedAt: entity.GetCreatedAt(),
		Id:        uuid.NewSHA1(uuid.NameSpaceDNS, []byte(entity.GetToken())).String(),
		UserId:    entity.GetUserId(),
		Token:     entity.GetToken(),
	}

	_, err := repository.collection.InsertOne(context.TODO(), raw)
	if nil != err {
		msg := fmt.Sprintf("failed to insert refresh token: %v", err)
		repository.loggerService.Fatal(msg)
		return nil, errors.New(msg)
	}

	entity.id = raw.Id

	return entity, nil
}
