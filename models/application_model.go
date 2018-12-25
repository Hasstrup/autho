package models

import utils "github.com/authenticate/utilities"

type ApplicationModel struct {
	ID       string                 `json: "id" bson:"_id"`
	Name     string                 `json:"name" bson:"name"`
	Schema   map[string]interface{} `json:"app_schema" bson:"app_schema"`
	Key      string                 `json:"app_key" bson:"app_key"`
	Database string                 `json:"database" bson:"database"`
	Address  string                 `json:"address" bson:"database"`
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
	}
	return results
}

func NewApplicationModel(fields []byte) *ApplicationModel {
	var m ApplicationModel
	utils.DecodeJSON(fields, m)
	return &m
}
