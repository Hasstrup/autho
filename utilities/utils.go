package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func EncodeJSON(data interface{}) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return d
}

func DecodeJSON(source []byte, target interface{}) {
	err := json.Unmarshal(source, &target)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func Contains(el interface{}, ar []interface{}) bool {
	for _, val := range ar {
		if reflect.TypeOf(val).Name() == reflect.TypeOf(el).Name() &&
			reflect.ValueOf(el).Interface() == reflect.ValueOf(val).Interface() {
			return true
		}
	}
	return false
}
