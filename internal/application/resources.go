package application

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ctx = context.Background()

type mongoResource struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type resources struct {
	mongo *mongoResource
}

func (resources *resources) close() []error {
	errorList := make([]error, 0)

	if nil == resources.mongo {
		errorList = append(errorList, errors.New("failed to close mongo: connection is nil"))
	} else {
		if err := resources.mongo.Client.Disconnect(ctx); err != nil {
			errorList = append(errorList, fmt.Errorf("failed to close mongo: %v", err))
		}
	}

	return errorList
}

//goland:noinspection GoExportedFuncWithUnexportedType
func resourcesFromConfig(config *config) (*resources, error) {
	// Mongo
	mongoConfig := config.Mongo
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port, mongoConfig.Database)
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(5 * time.Second)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to open mongo: %w", err)
	}

	ctxMongo, cancelMongo := context.WithTimeout(ctx, 5*time.Second)
	defer cancelMongo()
	err = mongoClient.Ping(ctxMongo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	mongoRes := &mongoResource{
		Client:   mongoClient,
		Database: mongoClient.Database(mongoConfig.Database),
	}

	resourcesObj := &resources{
		mongo: mongoRes,
	}

	return resourcesObj, nil
}
