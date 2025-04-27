package user

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
)

type userRaw struct {
	collection *mongo.Collection
	Email      string
	Id         string
	Name       string
	Password   string
}

func (raw *userRaw) ToUser() *User {
	return NewUser(raw.Email, raw.Id, raw.Name, raw.Password)
}

type userRepository struct {
	collection    *mongo.Collection
	loggerService *logger.Logger
}

func (repository *userRepository) FindOneByEmail(email string) (*User, error) {
	filter := bson.M{"esmail": email}
	var raw userRaw
	err := repository.collection.FindOne(context.TODO(), filter).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		msg := fmt.Sprintf("failed to run query: %v", err)
		repository.loggerService.Fatal(msg)
		return nil, errors.New(msg)
	}

	return raw.ToUser(), nil
}
