package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	utils "github.com/authenticate/utilities"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func (d *DatabaseDriver) Authenticate() (*DatabaseDriver, error) {
	if d.Database != "postgres" && d.Database != "mongodb" {
		return d, errors.New("That's an invalid db sent")
	}
	var target string
	var err error
	// check for the authenticable field
	for key, value := range d.Schema {
		if val, ok := utils.CleanUpValue(value).(primitive.M); ok {
			if key != "password" && val["authenticable"] != nil && val["authenticable"].(bool) {
				target = key
			}
		}
	}
	if target == "" {
		return d, errors.New("Sorry we could not find any authenticable field here in the schema, so we can't process this request")
	}
	if d.Database == "postgres" {
		err = d.AuthenticatePostgres(target)
	} else {
		err = d.AuthenticateMongo(target)
	}

	return d, err

}

func (d *DatabaseDriver) AuthenticatePostgres(target string) error {
	db, err := sql.Open("postgres", d.Address)
	if err != nil {
		log.Println(err)
		return err
	}
	if d.Payload[target] == nil {
		return errors.New(fmt.Sprintf("Please provide a %v field for this request to be processed", target))
	}
	query := d.buildSelectionQuery(target)
	log.Println(query)
	var password string
	err = db.QueryRowContext(context.Background(), query).Scan(&password)
	switch {
	case err == sql.ErrNoRows:
		return errors.New(fmt.Sprintf("Sorry we couldn't find a user with that %v", target))
	case err == nil:
		if CompareWithBcrypt(password, d.Payload["password"].(string)) {
			return nil
		} else {
			return errors.New("Invalid password for user")
		}
	}
	return err
}

func (d *DatabaseDriver) buildSelectionQuery(target string) string {
	// original query was return fmt.Sprintf("SELECT password FROM %v where %s=%d", DefaultCollectionAndTable, target, d.Payload[target].(string))
	return "SELECT password FROM " + DefaultCollectionAndTable + " WHERE " + target + "=" + "'" + d.Payload[target].(string) + "'"
}

func (d *DatabaseDriver) AuthenticateMongo(target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, d.Address)
	if err != nil {
		return err
	}
	collection := client.Database(d.Name).Collection(DefaultCollectionAndTable)
	query := map[string]string{}
	str, ok := d.Payload[target].(string)
	if !ok {
		return errors.New(fmt.Sprintf("Please provide a %v field for this request to be processed", target))
	}
	query[target] = str
	result := map[string]interface{}{}
	r := collection.FindOne(ctx, query)
	err = r.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("Sorry we could not find a user matching the credentials sent")
		}
		return err
	}
	if CompareWithBcrypt(result["password"].(string), d.Payload["password"].(string)) {
		return nil
	}
	return errors.New("Invalid user/password combination")

}

func (d *DatabaseDriver) YieldToken(fields []string) string {
	target := map[string]interface{}{}
	for _, key := range fields {
		to := utils.CleanUpValue(d.Schema[key])
		if t, ok := to.(string); ok {
			target[key] = coerce(t, d.Payload[key])
		} else {
			target[key] = coerce(to.(primitive.M)["type"].(string), d.Payload[key])
		}
	}
	claims := jwt.MapClaims{"payload": target}
	return EncodeWithJwt(claims)
}
