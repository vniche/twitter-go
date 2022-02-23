package twitter

import "reflect"

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
