package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func CheckForEmptyFieldsInMap(body map[string]interface{}, errors *[]string) {
	for key, value := range body {
		if _, ok := value.(float64); ok {
			continue
		}
		if reflect.ValueOf(value).Len() == 0 {
			*errors = append(*errors, key+" is missing ")
		}
	}
}

func MatchRequestToLengthInSchema(schema, request map[string]interface{}, errors *[]string) {
	//if there are some spill over errors from the previous validation
	//then don't bother validating anything anymore
	if len(*errors) > 0 {
		return
	}
	for key, value := range schema {
		val, ok := CleanUpValue(value).(primitive.M)
		if !ok {
			continue
		}
		//TODO check that min && max match are actual floats
		min := val["minLength"]
		max := val["maxLength"]
		if min == nil && max == nil {
			continue
		}
		// this should only work for strings right??? YES.
		str, ok := request[key].(string)
		if !ok {
			continue
		}
		if min != nil && float64(len(str)) < min.(float64) {
			*errors = append(*errors, fmt.Sprintf("%v should be at least %s characters long", key, strconv.Itoa(int(min.(float64)))))
		}
		if max != nil && float64(len(str)) > max.(float64) {
			*errors = append(*errors, fmt.Sprintf("%v should be not be more than %s characters", key, strconv.Itoa(int(max.(float64)))))
		}

	}
}
