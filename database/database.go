package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

func logError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var DbCxt context.Context
var TestCollection *mongo.Collection

func ConnectDB() {
	u := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	DbCxt = ctx
	o := *options.Client().ApplyURI(u)
	client, err := mongo.Connect(ctx, &o)
	if err != nil {
		defer client.Disconnect(ctx)
		log.Fatalln(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	todoDB := client.Database("todo")
	TestCollection = todoDB.Collection("test")
}
