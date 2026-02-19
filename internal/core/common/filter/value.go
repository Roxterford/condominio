package filter

import (
	"reflect"

	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

type ValueType string

const (
	String      ValueType = "string"
	StringList  ValueType = "[]string"
	Int         ValueType = "int"
	IntList     ValueType = "[]int"
	Float       ValueType = "float"
	FloatList   ValueType = "[]float"
	TypeBool    ValueType = "bool"
	Unknown     ValueType = "unknown"
	UnknownList ValueType = "[]unknown"
)

type Value interface {
	String() (string, bool)
	StringList() ([]string, bool)
	Int() (int, bool)
	IntList() ([]int, bool)
	Float() (float64, bool)
	FloatList() ([]float64, bool)
	Bool() (bool, bool)
	Type() ValueType
	Raw() any
}

// Value Wrapper
type walue struct {
	value any
}

func (wrapper walue) Raw() any {
	return wrapper.value
}

func (wrapper walue) String() (string, bool) {
	v, ok := wrapper.value.(string)
	return v, ok
}

func (wrapper walue) StringList() ([]string, bool) {
	v, ok := wrapper.value.([]string)
	return v, ok
}

func (wrapper walue) Int() (int, bool) {
	v, ok := wrapper.value.(int)
	return v, ok
}

func (wrapper walue) IntList() ([]int, bool) {
	v, ok := wrapper.value.([]int)
	return v, ok
}

func (wrapper walue) Float() (float64, bool) {
	v, ok := wrapper.value.(float64)
	return v, ok
}

func (wrapper walue) FloatList() ([]float64, bool) {
	v, ok := wrapper.value.([]float64)
	return v, ok
}

func (wrapper walue) Bool() (bool, bool) {
	v, ok := wrapper.value.(bool)
	return v, ok
}

func (wrapper walue) Type() ValueType {
	logger.Debug("Value `%v` type: %s", wrapper.value, reflect.TypeOf(wrapper.value))
	switch wrapper.value.(type) {
	case string:
		return String
	case []string:
		return StringList
	case int:
		return Int
	case []int:
		return IntList
	case float64:
		return Float
	case []float64:
		return FloatList
	case bool:
		return TypeBool
	case []any:
		if len(wrapper.value.([]any)) == 0 {
			return UnknownList
		}
		switch wrapper.value.([]any)[0].(type) {
		case string:
			return StringList
		case int:
			return IntList
		case float64:
			return FloatList
		default:
			return UnknownList
		}
	}

	return Unknown
}
