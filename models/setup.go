package models

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
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
	return client
}
