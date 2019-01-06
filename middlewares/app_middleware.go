package middlewares

import (
	"reflect"

	"github.com/authenticate/drivers"
	utils "github.com/authenticate/utilities"
)

var requiredSchemaFields = map[string]reflect.Kind{
	"name":       reflect.String,
	"database":   reflect.String,
	"address":    reflect.String,
	"app_schema": reflect.Map,
	"app_key":    reflect.String,
}

var stringTypes = []interface{}{"string", "number"}

var SupportedOptions = map[string]map[string]interface{}{
	"type": map[string]interface{}{
		"type":     reflect.String,
		"message":  "Type has to be present and has to be a string",
		"required": true,
	},
	"maxLength": map[string]interface{}{
		"type":     reflect.Float64,
		"message":  "The maximum length for this field has to be a number",
		"required": false,
	},
	"minLength": map[string]interface{}{
		"type":     reflect.Float64,
		"message":  "The minimum length for this field has to be a number",
		"required": false,
	},
	"required": map[string]interface{}{
		"type":     reflect.Bool,
		"message":  "Required has to be a boolean",
		"required": false,
	},
	"authenticable": map[string]interface{}{
		"type":     reflect.Bool,
		"message":  "Authenticable (If to be used in logging in users) has to a boolean",
		"required": false,
	},
	"tokenizable": map[string]interface{}{
		"type":     reflect.Bool,
		"message":  "Tokenizable (If to be used in forming a token) has to be a boolean",
		"required": false,
	},
}

func PingDatabaseAddress(database, address interface{}, ch chan interface{}, counter *int) {
	db, valid := database.(string)
	addy, ok := address.(string)
	if !ok || !valid {
		*counter++
		ch <- true
		return
	}
	if !utils.Contains(database, PermittedDatabaseTypes) {
		ch <- "Sorry the database has to be 'postgres' or 'mongodb' :("
	} else {
		if err := drivers.YieldDrivers(db)(addy); err != nil {
			ch <- err.Error()
		}
	}
	*counter++
	ch <- true
}

func CheckForRequiredFieldsInRequestBody(body map[string]interface{}, ch chan interface{}, counter *int) {
	for key, value := range requiredSchemaFields {
		data, ok := body[key]
		if !ok {
			ch <- "Hey you need to send in the '" + key + "' field in the request body"
			continue
		}
		if reflect.TypeOf(data).Kind() != value {
			ch <- "Please send in the correct datatype for the '" + key + "' field "
			continue
		}
		if key == "app_schema" {
			ValidateSchema(body["app_schema"].(map[string]interface{}), ch)
		}
	}
	*counter++
	ch <- true
}

func ValidateSchema(schema map[string]interface{}, ch chan interface{}, opts ...interface{}) {

	for key, value := range schema {
		_, present := value.(string)
		field, ok := value.(map[string]interface{})
		//If it is anything but a string or a map
		if !present && !ok {
			ch <- "The datatype provided for schema '" + key + "' is invalid"
			continue
		}
		if present && !utils.Contains(value, stringTypes) {
			ch <- "Currently we think it makes sense to support just strings and numbers for app_schem fields"
			continue
		}
		if ok {
			for k, v := range SupportedOptions {
				if !isValidType(field[k], v["type"].(reflect.Kind), v["required"].(bool)) {
					ch <- "( " + key + " ) " + v["message"].(string)
				}
			}
		}
	}
	if len(opts) > 0 && opts[0].(bool) {
		close(ch)
	}

}

func isValidType(value interface{}, desired reflect.Kind, required bool) bool {
	if value == nil && !required {
		return true
	} else if value == nil && required {
		return false
	}
	return reflect.TypeOf(value).Kind() == desired
}
