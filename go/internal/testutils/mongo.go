package testutils

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	mongodbTestContainer "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SpawnMongoDBContainer(ctx context.Context, mongoVersion string) (*mongodbTestContainer.MongoDBContainer, string, error) {
	log.Println("[TEST] Starting MongoDB Container")
	var err error
	mongoContainer, err := mongodbTestContainer.RunContainer(ctx,
		testcontainers.WithImage("mongo:"+mongoVersion),
	)
	if err != nil {
		return nil, "", err
	}
	log.Println("[TEST] Waiting for MongoDB Container to be ready")
	for !mongoContainer.IsRunning() {
	}
	log.Println("[TEST] RabbitMQ Container is ready")
	mongoDBContainerURI, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		return nil, "", err
	}
	return mongoContainer, mongoDBContainerURI, nil
}

func TerminateMongoDBContainer(ctx context.Context, mongoContainer *mongodbTestContainer.MongoDBContainer) error {
	log.Println("[TEST] Terminating MongoDB Container")
	return mongoContainer.Terminate(ctx)
}

func NewMongoDBClient(mongoURI string) (*mongo.Client, error) {
	log.Println("[TEST] Connecting to MongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
}

func InsertOne(ctx context.Context, coll *mongo.Collection, document interface{}) error {
	_, err := coll.InsertOne(ctx, document)
	return err
}

func EmptyCollection(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.DeleteMany(ctx, primitive.M{})
	return err
}

func FindOneByID(ctx context.Context, coll *mongo.Collection, id string) (map[string]interface{}, error) {
	var doc map[string]interface{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	sgR, err := coll.FindOne(ctx, primitive.M{"_id": objectID}), nil
	if err != nil {
		return nil, err
	}
	err = sgR.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
