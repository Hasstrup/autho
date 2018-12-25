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
			validateSchema(body["app_schema"].(map[string]interface{}), ch)
		}
	}
	*counter++
	ch <- true
}

func validateSchema(schema map[string]interface{}, ch chan interface{}) {
	for key, value := range schema {
		_, present := value.(string)
		field, ok := value.(map[string]interface{})
		if !present && !ok {
			ch <- "The datatype provided for schema '" + key + "' is invalid"
			continue
		}
		if present && !utils.Contains(value, stringTypes) {
			ch <- "Currently we think it makes sense to support just strings and numbers for app_schem fields"
			continue
		}
		if ok {
			ty, ok := field["type"]
			if !ok {
				ch <- "Please supply the type for the schema (" + key + ") field"
				continue
			}
			if ok && !utils.Contains(ty, stringTypes) {
				ch <- "The type field for schema field (" + key + ") should be 'string' or 'number'"
				continue
			}
			if required, present := field["required"]; present {
				if _, isBool := required.(bool); !isBool {
					ch <- "The required field for schema (" + key + ") should be a boolean :)"
				}
			}

		}
	}
}
