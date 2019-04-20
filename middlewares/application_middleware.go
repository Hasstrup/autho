package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/authenticate/services"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/authenticate/drivers"

	utils "github.com/authenticate/utilities"
)

var PermittedDatabaseTypes = []interface{}{"postgres", "mongodb"}

func SanitizeApplicationRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]interface{}
		var errors []string
		ch := make(chan interface{})
		numberOfChecksDone, expectedNumberOfChecks := 0, 4
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			utils.RespondWithJSON(w, 400, map[string]string{"error": "Unable to parse the input sent"})
			return
		}

		go CheckForRequiredFieldsInRequestBody(request, ch, &numberOfChecksDone)
		go checkForEmptyValuesInBody(request, ch, &numberOfChecksDone)
		go validateDatabaseStructure(request["database"], ch, &numberOfChecksDone)
		go PingDatabaseAddress(request["database"], request["address"], ch, &numberOfChecksDone)

		for msg := range ch {
			switch msg.(type) {
			case string:
				errors = append(errors, msg.(string))
			case bool:
				if numberOfChecksDone == expectedNumberOfChecks {
					close(ch)
				}
			}
		}
		if len(errors) > 0 {
			utils.RespondWithJSON(w, 400, map[string][]string{"errors": errors})
			return
		}
		//reset the content of r.Body
		b, _ := json.Marshal(request)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		next.ServeHTTP(w, r)
	})
}

func checkForEmptyValuesInBody(body map[string]interface{}, ch chan interface{}, counter *int) {
	for key, value := range body {
		if !checkForLength(reflect.TypeOf(value), value) {
			ch <- key + " is missing or invalid"
		}
	}
	*counter++
	ch <- true
}

func checkForLength(t reflect.Type, entity interface{}) bool {
	switch t.Kind() {
	case reflect.String, reflect.Slice, reflect.Map:
		return reflect.ValueOf(entity).Len() > 0
	default:
		return false
	}
}

func validateDatabaseStructure(e interface{}, ch chan interface{}, counter *int) {
	if _, ok := e.(string); !ok {
		ch <- "Database must be a string"
	} else {
		if !utils.Contains(e, PermittedDatabaseTypes) {
			ch <- "Currently we only support postgresql and mongodb"
		}
	}
	*counter++
	ch <- true
}

func isValidDataType(database, value string) (bool, error) {
	var permittedTypes = map[string][]interface{}{
		"postgres": []interface{}{"string", "number"},
		"mongodb":  []interface{}{"string", "number"},
	}
	if _, ok := permittedTypes[database]; !ok {
		return false, errors.New("The database provided isn't supported yet :)")
	}
	if !utils.Contains(value, permittedTypes[database]) {
		return false, nil
	}
	return true, nil
}

func ValidationPipeline(key string) func(errors *[]string, values ...interface{}) {
	validationMap := map[string]func(errors *[]string, values ...interface{}){
		"database": func(errors *[]string, values ...interface{}) {
			str, ok := values[0].(string)
			if !ok {
				*errors = append(*errors, "Database needs to be a string")
				return
			}
			if str != "mongodb" && str != "postgres" {
				*errors = append(*errors, "Hey we only support postgres and mongodb")
			}
		},
		"address": func(errors *[]string, values ...interface{}) {
			if len(values) < 2 {
				*errors = append(*errors, "Please provide the database (mongodb and postgresql) and the address")
				return
			}
			for _, val := range values {
				if _, ok := val.(string); !ok {
					*errors = append(*errors, "The database and address has to be a string (please provide both of them)")
					return
				}
			}
			err := drivers.YieldDrivers(values[0].(string))(values[1].(string))
			if err != nil {
				*errors = append(*errors, err.Error())
			}
		},
		"app_schema": func(errors *[]string, values ...interface{}) {
			ch := make(chan interface{})
			go ValidateSchema(values[0].(map[string]interface{}), ch, true)
			for msg := range ch {
				*errors = append(*errors, msg.(string))
			}
		},
		"name": func(errors *[]string, values ...interface{}) {
			if _, ok := values[0].(string); !ok {
				*errors = append(*errors, "Name has to be a string please")
			}
			result := services.FindOneApplication(map[string]string{"name": values[0].(string)}, values[1].(*mongo.Client))
			if result["_id"] != nil {
				*errors = append(*errors, "The name is already taken sadly :(")
			}
		},
		"app_key": func(errors *[]string, values ...interface{}) {
			if _, ok := values[0].(string); !ok {
				*errors = append(*errors, "App key has to be a string please")
			}
		},
	}
	if f, ok := validationMap[key]; ok {
		return f
	}
	return func(errors *[]string, values ...interface{}) {}
}
