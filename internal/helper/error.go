package helper

import "reflect"

func IsErrOfType(err, customErrorType error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(customErrorType)
}
