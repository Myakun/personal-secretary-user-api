package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Myakun/personal-secretary-user-api/internal/domain/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewUserRepository(collection *mongo.Collection, logger logger.Logger) user.UserRepository {
	return &mongoUserRepository{
		collection: collection,
		logger:     logger,
	}
}

func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var raw rawUser

	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		msg := "GetByEmail query failed"
		r.logger.ErrorW(msg, "error", err, "email", email)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return raw.ToUser(), nil
}

func (r *mongoUserRepository) GetById(ctx context.Context, id string) (*user.User, error) {
	var raw rawUser

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		msg := "GetById query failed"
		r.logger.ErrorW(msg, "error", err, "id", id)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return raw.ToUser(), nil
}

func (r *mongoUserRepository) Save(ctx context.Context, entity *user.User) (*user.User, error) {
	if entity == nil {
		msg := "user is nil"
		r.logger.Error(msg)
		return nil, errors.New(msg)
	}

	raw := &rawUser{
		Email:    entity.Email,
		Name:     entity.Name,
		Password: entity.Password,
	}

	if !entity.IsInserted() {
		raw.Id = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(entity.Email)).String()
	}

	_, err := r.collection.InsertOne(ctx, raw)
	if nil != err {
		msg := fmt.Sprintf("failed to insert user: %v", err)
		r.logger.Error(msg)
		return nil, errors.New(msg)
	}

	return entity, nil
}
