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

func FindOne(query map[string]string, client *mongo.Client) (interface{}, error) {
	data, err := models.FindOne(query, client, applicationCollection)
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

func FindAllApplications(client *mongo.Client) []interface{} {
	results, _ := models.FindAll(map[string]interface{}{}, client, applicationCollection)
	return results
}

func FindOneApplication(query map[string]string, client *mongo.Client) map[string]interface{} {
	result, _ := models.FindOne(query, client, applicationCollection)
	return *result
}
