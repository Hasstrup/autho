package models

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const database = "autho"

type PersistableInterface interface {
	Collection() string
	Fields() map[string]interface{}
}

type Model struct {
	data PersistableInterface
	db   *mongo.Client // you are going to change this moving forward
}

func NewModel(body PersistableInterface) *Model {
	return &Model{data: body}
}
func Save(body PersistableInterface, client *mongo.Client, c string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	collection := client.Database(database).Collection(c)
	defer cancel()
	res, err := collection.InsertOne(ctx, body.Fields())
	return res.InsertedID, err
}

func FindOne(query map[string]string, client *mongo.Client, c string) (interface{}, error) {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	var result struct{}
	err := collection.FindOne(ctx, query).Decode(&result)
	return &result, err
}

func FindAll(query interface{}, client *mongo.Client, c string) ([]interface{}, error) {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	var results []interface{}
	cur, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	return results, err
}

func UpdateOne(name string, changes map[string]interface{}, client *mongo.Client, c string) *mongo.SingleResult {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	// check if the application exists then do
	res := collection.FindOneAndUpdate(ctx, map[string]string{"name": name}, changes)
	return res
}

func DeleteOne(name string, client *mongo.Client, c string) *mongo.SingleResult {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	res := collection.FindOneAndDelete(ctx, map[string]string{"name": name})
	return res
}
func yieldCollection(timeout time.Duration, client *mongo.Client, c string) (context.CancelFunc, context.Context, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	collection := client.Database(database).Collection(c)
	return cancel, ctx, collection
}
