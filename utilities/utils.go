package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
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

func CleanUpValue(target interface{}) interface{} {
	if val, ok := target.([]interface{}); ok {
		c := []primitive.E{}
		for _, v := range val {
			e := primitive.E{
				Key:   v.(map[string]interface{})["Key"].(string),
				Value: v.(map[string]interface{})["Value"],
			}
			c = append(c, e)
		}
		return primitive.D(c).Map()
	}
	if v, ok := target.(primitive.D); ok {
		return v.Map()
	}
	return target
}

func DeleteNilKeys(sch map[string]interface{}) map[string]interface{} {
	for key, value := range sch {
		if value == nil {
			delete(sch, key)
		}
	}
	return sch
}

func DeleteKeys(target *map[string]interface{}, keys []string) {
	for _, key := range keys {
		delete(*target, key)
	}
}

//FORMAT BODY
// These keys are hashed anyway and would make no sense to the user
// Sadly we cam never retrieve the application's api_key. I love it.
// if the user forgets his/her/{any fitting pronoun} api key, he/she/{any fitting pronoun} would have
// to change their secret key. So we can recompute all over again
func Transform(r map[string]interface{}) map[string]interface{} {
	keys := []string{"id", "app_key", "address", "api_key"}
	DeleteKeys(&r, keys)
	r["schema"] = r["schema"].(primitive.D).Map()
	for key, value := range r["schema"].(primitive.M) {
		r["schema"].(primitive.M)[key] = CleanUpValue(value)
	}
	return r
}
