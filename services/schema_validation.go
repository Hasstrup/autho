package services

import (
	"reflect"

	utils "github.com/authenticate/utilities"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

var typeMap = map[string]reflect.Kind{
	"string": reflect.String,
	"number": reflect.Int,
}

func ValidateRequestAgainstSchema(schema, body map[string]interface{}, key string) []string {
	errors := []string{}
	for key, data := range schema {
		value := utils.CleanUpValue(data)
		fieldMightBeRequired := reflect.TypeOf(value).Kind() == reflect.Map
		if fieldMightBeRequired {
			v := value.(primitive.M)
			if required, present := v[key]; present && required.(bool) {
				val, ok := body[key]
				if !ok {
					errors = append(errors, key+" is a required field for this application")
					continue
				} else {
					if !isValidType(v["type"].(string), val) {
						errors = append(errors, key+" has to be a "+v["type"].(string))
						continue
					}
				}
			}
			if body[key] != nil {
				if isValidType(v["type"].(string), body[key]) {
					continue
				} else {
					errors = append(errors, key+" has to be a "+v["type"].(string))
				}
			}
			continue
		}
		if val, ok := body[key]; ok {
			if isValidType(value.(string), val) {
				continue
			} else {
				errors = append(errors, key+" has to be a "+value.(string))
			}
		}

	}
	return errors
}

func isValidType(key string, value interface{}) bool {
	return reflect.TypeOf(value).Kind() == typeMap[key]
}
