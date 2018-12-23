package services

import (
	"encoding/json"

	"github.com/authenticate/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const applicationCollection = "application"

type ApplicationService struct {
	Model *models.Model
}

func FindOne(name string, client *mongo.Client) (interface{}, error) {
	data, err := models.FindOne(map[string]string{"_id": name}, client, applicationCollection)
	if err != nil {
		return map[string]string{}, err
	}
	return data, err
}

func RegisterApplication(decoder *json.Decoder, client *mongo.Client) (interface{}, error) {
	var m models.ApplicationModel
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}
	_, err := models.Save(m, client, applicationCollection)
	return &m, err
}
