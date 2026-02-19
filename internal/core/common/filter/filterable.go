package filter

import "github.com/Sanaruca/condominio/internal/core"

type Filterable interface {
	// ~struct{}
	FilterSpec() Spec
}

type Spec map[string]ValueType

func (s Spec) IsEmpty() bool {
	return len(s) == 0
}

func (s Spec) HasKey(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Spec) Keys() core.Set[string] {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return core.NewSet(keys...)
}
