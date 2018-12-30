package models

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const database = "autho"

type PersistableInterface interface {
	Collection() string
	Fields() map[string]interface{}
}

type Model struct {
	data PersistableInterface
	db   *mongo.Client
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

func FindOne(query map[string]string, client *mongo.Client, c string) (*map[string]interface{}, error) {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	var result map[string]interface{}
	err := collection.FindOne(ctx, query).Decode(&result)
	return &result, err
}

func FindWorkableApplication(query map[string]string, client *mongo.Client, c string) (*WorkableApplication, error) {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	var r WorkableApplication
	err := collection.FindOne(ctx, query).Decode(&r)
	return &r, err
}

// TODO: This has to be role based in the future - only the admin running this locally should be able to get
//this
func FindAll(query interface{}, client *mongo.Client, c string) ([]interface{}, error) {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
	var results []interface{}
	cur, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var result map[string]interface{}
		err := cur.Decode(&result)
		if err != nil {
			// lol can't fatal anything here o. Remember to quietly fail
			log.Fatal(err)
		}
		results = append(results, result)
	}
	return results, err
}

func UpdateOne(name string, changes map[string]interface{}, client *mongo.Client, c string) *mongo.SingleResult {
	cancel, ctx, collection := yieldCollection(30, client, c)
	defer cancel()
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
