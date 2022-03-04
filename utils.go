package twitter

import (
	"reflect"
	"strings"
)

func IsOneOfEnum(value string, enums ...interface{}) bool {
	valid := false
	for _, enum := range enums {
		v := reflect.ValueOf(enum)

		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).String() == value {
				valid = true
				break
			}
		}
	}

	return valid
}

func ParseURLParameters(parameters map[string][]string) (string, error) {
	var queryParams string

	if len(parameters) > 0 {
		queryParams = "?"
	}

	index := 0
	for key, value := range parameters {
		// TODO: validate if parameters are valid
		queryParams += key + "=" + strings.Join(value, ",")
		if index != (len(parameters) - 1) {
			queryParams += "&"
		}
		index++
	}
	return queryParams, nil
}
