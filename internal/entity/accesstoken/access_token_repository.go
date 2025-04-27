package accesstoken

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
)

type accessTokenRaw struct {
	DeviceId string `bson:"device_id"`
	TeamId   int    `bson:"team_id"`
	Token    string `bson:"token"`
	UserId   int    `bson:"user_id"`
}

type accessTokenRepository struct {
	collection *mongo.Collection
	logger     *logger.Logger
}

func (repository *accessTokenRepository) FindOneByToken(token string) (*AccessToken, error) {
	filter := bson.M{"token": token}
	var raw accessTokenRaw
	err := repository.collection.FindOne(context.TODO(), filter).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		msg := fmt.Sprintf("failed to find token %s: %v", token, err)
		repository.logger.Fatal(msg)
		return nil, errors.New(msg)
	}

	return NewAccessToken(raw.DeviceId, raw.TeamId, raw.Token, raw.UserId), nil
}
