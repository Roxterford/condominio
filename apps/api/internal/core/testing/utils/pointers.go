package utils

import "time"

func StringPtr(s string) *string {
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func Ptr[T any](t T) *T {
	return &t
}
