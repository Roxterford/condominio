package model

import (
	"reflect"
	"strings"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

// conditionTypeMappers asocia el tipo de condición GraphQL con su
// ValueType de dominio. Se indexa por reflect.Type para búsquedas O(1).
var conditionTypeMappers = map[reflect.Type]filter.ValueType{
	reflect.TypeOf(StringCondition{}):  filter.TypeString,
	reflect.TypeOf(IntCondition{}):     filter.TypeInt,
	reflect.TypeOf(BooleanCondition{}): filter.TypeBool,
}

// GetFilterSpec construye un filter.Spec a partir de un struct de input
// GraphQL, mapeando la etiqueta JSON de cada campo a su tipo de valor
// de dominio según el tipo de la condición que lo contiene.
func GetFilterSpec[T any](typo T) filter.Spec {
	spec := filter.Spec{}

	t := reflect.TypeOf(typo)
	if t == nil {
		return spec
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return spec
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		name := jsonFieldName(field)
		if name == "-" {
			continue
		}

		fieldType := field.Type
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Solo los tipos de condición son válidos en el Spec plano.
		// Los demás campos (and/or/not) son recursiones manejadas por
		// el validador mediante el árbol de cláusulas, por lo que se omiten.
		valueType, ok := conditionTypeMappers[fieldType]
		if !ok {
			continue
		}

		spec[name] = valueType
	}

	return spec
}

// jsonFieldName extrae el nombre de campo desde la etiqueta JSON,
// fallando al nombre del campo Go cuando no hay etiqueta explícita.
func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	name := tag
	if idx := strings.Index(tag, ","); idx >= 0 {
		name = tag[:idx]
	}

	if name == "" {
		return field.Name
	}

	return name
}
