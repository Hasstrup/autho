package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	utils "github.com/authenticate/utilities"
)

var PermittedDatabaseTypes = []interface{}{"postgres", "mongodb"}

func SanitizeApplicationRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]interface{}
		var errors []string
		ch := make(chan interface{})
		numberOfChecksDone, expectedNumberOfChecks := 0, 2
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			utils.RespondWithJSON(w, 400, map[string]string{"error": "Unable to parse the input sent"})
			return
		}
		go checkForEmptyValuesInBody(request, ch, &numberOfChecksDone)
		go validateDatabaseStructure(request["database"], ch, &numberOfChecksDone, request["app_schema"])
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

func validateDatabaseStructure(e interface{}, ch chan interface{}, counter *int, schema interface{}) {
	if _, ok := e.(string); !ok {
		ch <- "Database must be a string"
	} else {
		if !utils.Contains(e, PermittedDatabaseTypes) {
			ch <- "Currently we only support postgresql and mongodb"
		} else {
			// check that the schema is a map
			if reflect.TypeOf(schema).Name() != "map" {
				ch <- "app_schema must be an object containing the fields that are either required"
			} else {
				for key, value := range schema.(map[string]interface{}) {
					// TODO: refactor this whole thing, so messy and long I could barely breathe
					// reading this thing
					if val, ok := value.(string); ok {
						valid, err := isValidDataType(e.(string), val)
						if err != nil {
							ch <- err.Error()
							break
						}
						if !valid {
							ch <- key + " is not a supported datatype for the database provided"
						}
						continue
					}
					if val, ok := value.(map[string]interface{}); ok {
						if _, present := val["type"]; !present {
							ch <- "Please specify the datatype for " + key
						} else {
							// HMMMMM potential bug here
							str, ok := val["type"].(string)
							if valid, _ := isValidDataType(e.(string), str); ok && valid {
								continue
							}
							ch <- "The type specified for " + key + " must either be 'string' or 'number'"
						}

					}
				}
			}
		}
	}
	*counter++
	ch <- true
}

func isValidDataType(database, value string) (bool, error) {
	var permittedTypes map[string][]interface{}
	permittedTypes["postgres"] = []interface{}{"string", "number"}
	permittedTypes["mongodb"] = []interface{}{"string", "number"}
	if _, ok := permittedTypes[database]; !ok {
		return false, errors.New("The database provided isn't supported yet :)")
	}
	if !utils.Contains(value, permittedTypes[database]) {
		return false, nil
	}
	return true, nil
}
