package vmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHost struct {
	Host string
	Port int
}

type MongoConfig struct {
	Hosts    []MongoHost
	UserName string
	Password string
}

type MongoClient struct {
	client *mongo.Client
}

// Collection is a handle to a MongoDB collection. It is safe for concurrent use by multiple goroutines.

func (m *MongoClient) Connect(ctx context.Context) (*mongo.Client, error) {

	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	defer mongoCli.Disconnect(ctx)

	// collection := mongoCli.Database("test").Collection("test")

	return mongoCli, nil
}

func (m *MongoClient) ConnectWithAuth(ctx context.Context) (*mongo.Client, error) {

	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	defer mongoCli.Disconnect(ctx)

	// collection := mongoCli.Database("test").Collection("test")

	// transaction
	//session, err := mongoCli.StartSession()
	//if err != nil {
	//	return nil, err
	//}
	//session.

	collection := mongoCli.Database("test").Collection("test")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	cur.All(ctx, &[]interface{}{})

	return mongoCli, nil
}
