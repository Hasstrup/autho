package services

import (
	"errors"
	"flag"

	"github.com/authenticate/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const applicationCollection = "application"

// TODO: shift this to an os.LookUp instead thanks
var Pass = flag.String("passcode", "Thisshouldnevereverbeused", "The ultimate key to encrypting everything")

type ApplicationService struct {
	Model *models.Model
}

func FindOne(query map[string]string, client *mongo.Client) (interface{}, error) {
	data, err := models.FindOne(query, client, applicationCollection)
	if err != nil {
		return map[string]string{}, err
	}
	return data, err
}

func RegisterApplication(m *models.ApplicationModel, client *mongo.Client) (interface{}, error) {
	if itExists(m.Name, client) {
		return nil, errors.New("Sorry the name is already taken :( ")
	}
	m.Address, _ = HashWithBcrypt(m.Address)
	m.Key, _ = HashWithBcrypt(m.Key)
	encryptedKey, _ := Encrypt([]byte(m.Name+"--"+m.Address), *Pass)
	// Hash the api key right before saving
	m.ApiKey, _ = Hash256(string(encryptedKey))
	_, err := models.Save(m, client, applicationCollection)
	claims := jwt.MapClaims{"payload": CustomSlice(encryptedKey)}
	m.ApiKey = EncodeWithJwt(claims)
	return &m, err
}

func FindAllApplications(client *mongo.Client) []interface{} {
	results, _ := models.FindAll(map[string]interface{}{}, client, applicationCollection)
	return results
}

func FindOneApplication(query map[string]string, client *mongo.Client) map[string]interface{} {
	result, _ := models.FindOne(query, client, applicationCollection)
	return *result
}

func itExists(name string, client *mongo.Client) bool {
	item := FindOneApplication(map[string]string{"name": name}, client)
	_, ok := item["_id"]
	return ok
}
