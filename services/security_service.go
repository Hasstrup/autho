package services

import (
	"errors"
	"log"

	"github.com/authenticate/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func YieldAppFromApiKey(key, secret string, client *mongo.Client) (*models.WorkableApplication, error) {
	k, err := ParseFromJwToken(key)
	if err != nil {
		return nil, err
	}
	payload := extractPayload(k["payload"])
	_, err = Decrypt(payload, *Pass)
	hash, _ := Hash256(string(payload))
	query := map[string]string{"api_key": hash}
	result, err := models.FindWorkableApplication(query, client, "application")
	if err != nil {
		log.Println(err)
		return result, errors.New("Sorry we had problems finding the app with that key")
	}
	if result.Name == "" {
		return result, errors.New("There is no application with that api key provided")
	}
	if !CompareWithBcrypt(result.AppKey, secret) {
		return nil, errors.New("Hey you do not have permissions to do that")
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
