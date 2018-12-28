package services

import (
	"errors"
	"log"

	"github.com/authenticate/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func YieldAppFromApiKey(key string, client *mongo.Client) (*map[string]interface{}, error) {
	k, err := ParseFromJwToken(key)
	if err != nil {
		return nil, err
	}
	payload := extractPayload(k["payload"])
	_, err = Decrypt(payload, *Pass)
	if err != nil {
		log.Printf("error %v", err.Error())
	}

	hash, _ := HashWithBcrypt(string(payload))
	query := map[string]string{
		"api_key": hash,
	}
	result, err := models.FindOne(query, client, "application")
	if err != nil {
		return result, errors.New("Sorry we had problems finding the app with that key")
	}
	item := *result
	if _, ok := item["_id"]; !ok {
		return result, errors.New("There is no application with that api key provided")
	}
	return result, nil
}

func extractPayload(g interface{}) []byte {
	n := []byte{}
	for _, val := range g.([]interface{}) {
		n = append(n, uint8(val.(float64)))
	}
	return n
}
