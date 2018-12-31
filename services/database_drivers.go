package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	utils "github.com/authenticate/utilities"
	_ "github.com/lib/pq"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const DefaultCollectionAndTable = "users"

type DatabaseDriver struct {
	Name     string
	Database string
	Schema   map[string]interface{}
	Address  string
	Payload  map[string]interface{}
}

func (d *DatabaseDriver) Write() (*DatabaseDriver, error) {
	if d.Database != "postgres" && d.Database != "mongodb" {
		return d, errors.New("That's an invalid db sent")
	}
	if d.Database == "postgres" {
		if id, err := d.WriteToPostgres(); err != nil {
			log.Println(err)
			return d, err
		} else {
			d.Payload["id"] = id
			return d, nil
		}
	}
	InsertedID, err := d.WriteToMongo()
	if err != nil {
		return d, err
	}
	d.Payload["id"] = InsertedID
	return d, nil
}

func (d *DatabaseDriver) WriteToPostgres() (int64, error) {
	db, err := sql.Open("postgres", d.Address)
	if err != nil {
		log.Println(err) //Perhaps a flag that tells whether to run in dev mode
		return 0, err
	}
	keys, values := BuildQuery(d.Schema, d.Payload)
	query := `INSERT INTO ` + DefaultCollectionAndTable + " " + keys + ` VALUES` + values
	log.Println(query) //dev mode
	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	// TODO: not sure why but I dont think postgres supports the LastInsertedId thing
	// sucks but need to find a workaround
	id, err := result.LastInsertId()
	return id, nil
}

func (d *DatabaseDriver) WriteToMongo() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, d.Address)
	if err != nil {
		return nil, err
	}
	collection := client.Database(d.Name).Collection(DefaultCollectionAndTable)
	res, err := collection.InsertOne(ctx, d.yieldMongoFields())
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (d *DatabaseDriver) yieldMongoFields() map[string]interface{} {
	payload := map[string]interface{}{}
	for key := range d.Schema {
		if key == "password" {
			hash, _ := HashWithBcrypt(d.Payload[key].(string))
			d.Payload[key] = hash
		}
		payload[key] = d.Payload[key]
	}
	return payload
}

func NewDatabaseDriver(app, body map[string]interface{}) *DatabaseDriver {
	return &DatabaseDriver{
		Name:     app["name"].(string),
		Database: app["database"].(string),
		Schema:   app["Schema"].(map[string]interface{}),
		Address:  app["address"].(string),
		Payload:  body,
	}
}

func BuildQuery(sch, payload map[string]interface{}) (string, string) {
	str := "("
	values := " ("
	for key, val := range sch {
		str += (" " + key + ", ")
		if key == "password" {
			hash, _ := HashWithBcrypt(payload[key].(string))
			payload[key] = hash
		}
		if dataType, ok := val.(string); ok {
			payload[key] = coerce(dataType, payload[key])
		} else {
			dt := utils.CleanUpValue(val).(primitive.M)["type"].(string)
			payload[key] = coerce(dt, payload[key])
		}
		values += (" " + payload[key].(string) + ",")
	}
	str = strings.TrimSuffix(str, ", ")
	str += ")"
	return str, strings.TrimSuffix(values, ",") + ")"
}

func coerce(t string, field interface{}) string {
	var str string
	if t == "number" {
		str = strconv.Itoa(field.(int))
	} else {
		str = "'" + field.(string) + "'"
	}
	return str
}
