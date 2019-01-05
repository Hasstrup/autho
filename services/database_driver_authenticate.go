package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	utils "github.com/authenticate/utilities"
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
		if val, ok := utils.CleanUpValue(val).(primitive.M); ok {
			if key != "password" && val["authenticable"] != nil && val["authenticable"].(bool) {
				target = key
			}
		}
	}
	if target == "" {
		return d, errors.New("Sorry we could not find any authenticable field here sadly")
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
	return fmt.Sprintf("SELECT password FROM %s WHERE %d=%v", DefaultCollectionAndTable, target, d.Payload[target].(string))
}

func (d *DatabaseDriver) AuthenticateMongo(target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, d.Address)
	if err != nil {
		return err
	}
	collection := client.Database(d.Database).Collection(DefaultCollectionAndTable)
	query := map[string]string{}
	query[target] = d.Payload[target].(string)
	result := map[string]interface{}
	r := collection.FindOne(ctx, query)
     if e := r.Err(); e != nil {
		if e == mongo.ErrNoDocuments { 
			return errors.New("Sorry we could not find a user matching the credentials sent")
		}
		return err
	 }
	 err = r.Decode(result)
	 // hmm really
	 if err != nil {
		 return err
	 }
	 if CompareWithBcrypt(r["password"].(string), d.Payload["password"].(string)) {
		 return nil
	 }
	 return errors.New("Invalid user/password combination")

}


func (d *DatabaseDriver) YieldToken(fields []string) string {
	target := map[string]interface{}
	for _, val := range fields {
		target := utils.CleanUpValue(d.Schema[val])
		if t, ok := target.(string); ok {
			target[val] = coerce(t, d.Payload[val])
		} else {
			target[val] = coerce(target.(primitive.M)["type"].(string), d.Payload[val])
		}
	}
	claims := jwt.MapClaims{"payload": target }
	return EncodeWithJwt(claims)
}
