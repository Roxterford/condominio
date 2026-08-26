package utils

import (
	"reflect"
	"time"
)

func GetCurrentYear() int {
	year, _, _ := time.Now().Date()
	return year
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func GetPointer[T any](value T) *T {
	return &value
}
