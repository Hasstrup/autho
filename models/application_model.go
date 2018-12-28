package models

type ApplicationModel struct {
	ID       string                 `json: "id" bson:"_id"`
	Name     string                 `json:"name" bson:"name"`
	Schema   map[string]interface{} `json:"app_schema" bson:"app_schema"`
	Key      string                 `json:"app_key" bson:"app_key"`
	Database string                 `json:"database" bson:"database"`
	Address  string                 `json:"address" bson:"database"`
	ApiKey   string                 `json"api_key" bson:"api_key"`
}

func (m ApplicationModel) Collection() string {
	return "application"
}

func (m ApplicationModel) Fields() map[string]interface{} {
	results := map[string]interface{}{
		"id":       m.ID,
		"name":     m.Name,
		"app_key":  m.Key,
		"database": m.Database,
		"schema":   m.Schema,
		"address":  m.Address,
		"api_key":  m.ApiKey,
	}
	return results
}

func NewApplicationModel(body map[string]interface{}) *ApplicationModel {
	return &ApplicationModel{
		Name:     body["name"].(string),
		Schema:   body["app_schema"].(map[string]interface{}),
		Key:      body["app_key"].(string),
		Database: body["database"].(string),
		Address:  body["address"].(string),
	}
}
