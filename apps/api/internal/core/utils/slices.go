package utils

func SliceFirst[T any](fields []T) T {

	var defailt T

	if len(fields) == 0 {
		return defailt
	}

	return fields[0]
}
