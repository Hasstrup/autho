package services

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/authenticate/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func YieldAppFromApiKey(key, secret string, client *mongo.Client) (*models.WorkableApplication, error) {
	k, err := ParseFromJwToken(key)
	if err != nil {
		return nil, err
	}
	payload := extractPayload(k["payload"])
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
	p, err := Decrypt(payload, *Pass)
	result.Address = strings.Split(string(p), "--")[1]
	return result, nil
}

func extractPayload(g interface{}) []byte {
	n := []byte{}
	for _, val := range g.([]interface{}) {
		n = append(n, uint8(val.(float64)))
	}
	return n
}

func RootUserOnly(key string) error {
	pass, ok := os.LookupEnv("AUTHOAPPKEY")
	if !ok {
		return errors.New("Hey you need to set the AUTHOAPPKEY as an env variable")
	}
	if pass != key {
		return errors.New("The key provided doesn't match the Autho App Key. Try again.")
	}
	return nil
}

func ComputeApiKey(name, addy string) (string, string) {
	final := name + "--" + addy
	raw, _ := Encrypt([]byte(final), *Pass)
	final := EncodeWithJwt(jwt.MapClaims{"payload": CustomSlice(raw)})
	return final, Hash256(string(raw))
}
