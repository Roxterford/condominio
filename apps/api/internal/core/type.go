package core

type Validable interface {
	Validate() Error
}

func New[T Validable](t T) (T, Error) {
	err := t.Validate()
	if err != nil {
		return t, err
	}
	return t, nil
}
