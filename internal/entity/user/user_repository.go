package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"personal-secretary-user-ap/internal/service/logger"
)

type userRaw struct {
	Email    string `bson:"email"`
	Id       string `bson:"_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

func (raw *userRaw) ToUser() *User {
	user := NewUser(raw.Email, raw.Id, raw.Name, raw.Password)
	user.setIsInserted(true)

	return user
}

type userRepository struct {
	collection    *mongo.Collection
	loggerService *logger.Logger
}

func (repository *userRepository) insert(user *User) (*User, error) {
	raw := &userRaw{
		Email:    user.GetEmail(),
		Id:       uuid.NewSHA1(uuid.NameSpaceDNS, []byte(user.GetEmail())).String(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
	}

	_, err := repository.collection.InsertOne(context.TODO(), raw)
	if nil != err {
		msg := fmt.Sprintf("failed to insert user: %v", err)
		repository.loggerService.Fatal(msg)
		return nil, errors.New(msg)
	}

	user.id = raw.Id

	return user, nil
}

func (repository *userRepository) FindOneByEmail(email string) (*User, error) {
	filter := bson.M{"email": email}
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
