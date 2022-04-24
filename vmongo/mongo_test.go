package vmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"testing"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestNewMongo(t *testing.T) {
	ctx := context.Background()
	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user1:pwd1@localhost:40000/test"))
	if err != nil {
		return
	}
	rand.Seed(time.Now().UnixNano())
	defer mongoCli.Disconnect(ctx)

	// collection := mongoCli.Database("test").Collection("test")

	// transaction
	//session, err := mongoCli.StartSession()
	//if err != nil {
	//	return nil, err
	//}
	//session.

	collection := mongoCli.Database("test").Collection("col")

	bsons := []interface{}{}
	for i := 0; i < 10000; i++ {
		bsons = append(bsons, bson.M{"name": randStr(3000)})
	}

	_, err = collection.InsertMany(ctx, bsons)
	if err != nil {
		fmt.Println(err)
	}

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return
	}

	var res []interface{}
	cur.All(ctx, &res)

	return
}
