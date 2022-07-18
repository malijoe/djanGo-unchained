package utils

import "reflect"

func IsZero(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}
