package filter

import (
	"github.com/Sanaruca/condominio/internal/core"
)

type Filter[T Filterable] struct {
	Query Query
}

func (f Filter[T]) Validate() core.Error {
	return f.Query.Validate()
}

type Filterable interface {
	// ~struct{}
	FilterableKeys() []string
}

// type FilterValue interface {
// 	~string | ~int | ~bool | time.Time
// }

// type Filter struct {
// 	key       string
// 	condition Condition
// 	value     any
// }

// type FilterSpec[T Filterable] struct {
// 	obj     T
// 	filters []Filter
// }

// func (f FilterSpec[T]) Validate() core.Error {
// 	valid_keys := f.obj.FilterableKeys()
// 	for _, filter := range f.filters {
// 		if !slices.Contains(valid_keys, filter.key) {
// 			return core.NewInvalidArgumentError("Campo '%s' no corresponde a ningun campo valido", filter.key)
// 		}

// 		if err := filter.Validate(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (f *FilterSpec[T]) Sanitize() {
// 	f.filters = slices.CompactFunc(f.filters, func(a, b Filter) bool {
// 		return a == b
// 	})
// }

// func (f Filter) Validate() core.Error {

// 	if f.key == "" {
// 		return core.NewInvalidArgumentError("Campo no puede estar vacio")
// 	}

// 	switch f.value.(type) {
// 	case string:
// 		if !f.condition.IsStringCondition() {
// 			return core.NewInvalidArgumentError("Condicion '%s' no es valida para valores de tipo string", f.condition)
// 		}
// 	case int:
// 		if !f.condition.IsNumericCondition() {
// 			return core.NewInvalidArgumentError("Condicion '%s' no es valida para valores de tipo int", f.condition)
// 		}
// 	case bool:
// 		if !f.condition.IsNumericCondition() {
// 			return core.NewInvalidArgumentError("Condicion '%s' no es valida para valores de tipo bool", f.condition)
// 		}
// 	case time.Time:
// 		if !f.condition.IsDateCondition() {
// 			return core.NewInvalidArgumentError("Condicion '%s' no es valida para valores de tipo Date", f.condition)
// 		}
// 	}

// 	return nil
// }

// func (f Filter) StringValue() (string, core.Error) {
// 	return core.Cast[string](f.value)
// }

// func (f Filter) IntValue() (int, core.Error) {
// 	return core.Cast[int](f.value)
// }

// func (f Filter) BoolValue() (bool, core.Error) {
// 	return core.Cast[bool](f.value)
// }
