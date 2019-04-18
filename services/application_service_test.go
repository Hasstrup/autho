package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/authenticate/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

var body = map[string]interface{}{
	"name": "test",
	"app_schema": map[string]interface{}{
		"username": "string",
	},
	"app_key":  "key",
	"database": "test_database",
	"api_key":  "test_api_key",
}

var model = models.NewApplicationModel(body)
var service = ApplicationService{
	Model: model,
}

func TestMain(t *testing.T) {
	//setUp
	SeedDatabase()
	t.Run("find_one", testFindOne())
}

func testFindOne() func(t *testing.T) {
	return func(t *testing.T) {
		fmt.Println("Okay")
	}
}

func SeedDatabase() {
	data := map[string]interface{}{
		"name":     "Test Application",
		"app_key":  "This is an application key",
		"database": "here we go",
	}

	if client, err := initializeDatabase(); err != nil {
		panic(err)
	}
	collection := client.Database("authenticate_test").Collection("application")
	_, err := collection.insertMany(context.Background(), []interface{}{data})
	if err != nil {
		panic(err)
	}
}

func initializeDatabase() (*mongo.Client, error) {
	uri := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
