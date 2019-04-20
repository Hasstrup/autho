package services

import (
	"context"
	"testing"
	"time"

	"github.com/authenticate/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var body = map[string]interface{}{
	"name": "test",
	"app_schema": map[string]interface{}{
		"username": "string",
	},
	"app_key":  "key",
	"database": "test_database",
	"api_key":  "test_api_key",
	"address":  "test address",
}

var model = models.NewApplicationModel(body)

func TestMain(t *testing.T) {
	client := SeedDatabase()
	// clear up db after test
	defer func() {
		if err := recover(); err != nil {
			dropDatabase(client)
		}
	}()
	t.Run("FindOneApplication", testFindOneApplication(client))
	dropDatabase(client)
}

func testFindOneApplication(client *mongo.Client) func(t *testing.T) {
	return func(t *testing.T) {
		query := map[string]string{
			"name": "Test Application",
		}
		record := FindOneApplication(query, client)
		expectedKey := "This is an application key"
		if key, ok := record["app_key"].(string); !ok || key != expectedKey {
			t.Errorf("expected app key to be This is an application key")
		}
	}
}

func SeedDatabase() *mongo.Client {
	data := map[string]interface{}{
		"name":     "Test Application",
		"app_key":  "This is an application key",
		"database": "here we go",
	}

	client, err := initializeDatabase()
	if err != nil {
		panic(err)
	}
	collection := client.Database("authenticate_test").Collection("application")
	if _, err := collection.InsertOne(context.Background(), data); err != nil {
		panic(err)
	}
	return client
}

func initializeDatabase() (*mongo.Client, error) {
	uri := "mongodb://localhost:27017"
	client, err := mongo.NewClient(uri)
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

func dropDatabase(client *mongo.Client) {
	client.Database("authenticate_test").Drop(context.Background())
}
