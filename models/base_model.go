package models

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const database = "autho"

type abstractInterface interface {
	Collection() string
	Fields() interface{}
}

type Model struct {
	data abstractInterface
	db   *mongo.Client // you are going to change this moving forward
}

func (m *Model) Save(body interface{}) (interface{}, error) {
	cancel, ctx, collection := m.yieldCollection(3)
	defer cancel()
	res, err := collection.InsertOne(ctx, m.data.Fields())
	return res.InsertedID, err
}

func (m *Model) FindOne(query map[string]interface{}) (interface{}, error) {
	cancel, ctx, collection := m.yieldCollection(30)
	defer cancel()
	var result struct{}
	err := collection.FindOne(ctx, query).Decode(&result)
	return &result, err
}

func (m *Model) FindAll(query interface{}) ([]interface{}, error) {
	cancel, ctx, collection := m.yieldCollection(30)
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

func (m *Model) UpdateOne(name string, changes map[string]interface{}) *mongo.SingleResult {
	cancel, ctx, collection := m.yieldCollection(30)
	defer cancel()
	// check if the application exists then do
	res := collection.FindOneAndUpdate(ctx, map[string]string{"name": name}, changes)
	return res
}

func (m *Model) DeleteOne(name string) *mongo.SingleResult {
	cancel, ctx, collection := m.yieldCollection(30)
	defer cancel()
	res := collection.FindOneAndDelete(ctx, map[string]string{"name": name})
	return res
}
func (m *Model) yieldCollection(timeout time.Duration) (context.CancelFunc, context.Context, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	collection := m.db.Database(database).Collection(m.data.Collection())
	return cancel, ctx, collection
}
