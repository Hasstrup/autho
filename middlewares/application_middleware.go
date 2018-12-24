package middlewares

import (
	"encoding/json"
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
		go validateDatabaseStructure(request["database"], ch, &numberOfChecksDone)
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
