package models

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

func RegisterDatabase(str string) *mongo.Client {
	flag.Parse()
	client, err := mongo.NewClient(str)
	if err != nil {
		log.Fatalf("Failed to connect to the mongo database: Error %v", err.Error())
		panic("failed to connect")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	PopulateCollectionIndexes(client)
	return client
}

func PopulateCollectionIndexes(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("autho").Collection("application")
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "name", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	}
	createOptions := options.CreateIndexes().SetMaxTime(5 * time.Second)
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{index}, createOptions)
	if err != nil {
		//HMMM should we let this slide? Perhaps not.
		panic(err)
	}
}
