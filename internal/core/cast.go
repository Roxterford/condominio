package core

func Cast[T any](value any) (T, Error) {

	v, ok := value.(T)
	if !ok {
		return v, NewInvalidArgumentError("Valor '%s' no corresponde a ningun tipo valido", value)
	}

	return v, nil
}
