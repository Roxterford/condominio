package filter

type Spec map[string]ValueType

func (s Spec) IsEmpty() bool {
	return len(s) == 0
}

func (s Spec) HasKey(key string) bool {
	_, ok := s[key]
	return ok
}
