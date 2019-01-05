package utils

import "github.com/mongodb/mongo-go-driver/bson/primitive"

func ExtractAuthenticableFields(schema map[string]interface{}) (map[string]interface{}, []string) {
	authFields := map[string]interface{}{}
	tokenFields := []string{}

	for key, value := range schema {
		val, ok := CleanUpValue(value).(primitive.M)
		if !ok {
			continue
		}
		if val["authenticable"] != nil {
			if val["authenticable"].(bool) {
				authFields[key] = value
			}
		}
		if val["tokenizable"] != nil {
			if val["tokenizable"].(bool) {
				tokenFields = append(tokenFields, key)
			}
		}
	}
	return authFields, tokenFields
}
