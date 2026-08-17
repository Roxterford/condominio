package model

import (
	"reflect"
	"strings"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

func applyFilter[T filter.Filterable](input any) filter.Filter[T] {
	nill := *filter.NewFilter[T](nil)
	if input == nil {
		return nill
	}

	inputMap, err := structToMap(input)
	if err != nil {
		return nill
	}

	return *filter.NewFilter[T](inputMap)
}

// structToMap convierte un struct a map[string]any, manejando correctamente
// campos Omittable de gqlgen distinguiendo entre "no enviado" y "null explícito".
func structToMap(input any) (map[string]any, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	result := make(map[string]any)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Ignorar campos no exportados
		if !fieldType.IsExported() {
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		fieldName := extractJSONFieldName(jsonTag)
		if fieldName == "" || fieldName == "-" {
			continue
		}

		// Detectar campos Omittable por la presencia del método IsSet
		if isSetMethod := field.MethodByName("IsSet"); isSetMethod.IsValid() {
			isSet := isSetMethod.Call(nil)[0].Bool()
			if !isSet {
				continue // Campo no enviado — omitir del mapa
			}
			// Campo enviado (aunque sea null) — leer el valor
			valueMethod := field.MethodByName("Value")
			value := valueMethod.Call(nil)[0]
			if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
				if value.IsNil() {
					result[fieldName] = nil
				} else {
					result[fieldName] = value.Elem().Interface()
				}
			} else {
				result[fieldName] = value.Interface()
			}
			continue
		}

		// Desreferenciar punteros
		fieldValue := field
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue // Puntero nil — omitir del mapa
			}
			fieldValue = fieldValue.Elem()
		}

		// Recursión para structs anidados
		if fieldValue.Kind() == reflect.Struct {
			nested, err := structToMap(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			if len(nested) > 0 {
				result[fieldName] = nested
			}
			continue
		}

		// Slices: omitir nil/vacíos y convertir a []any (como haría json.Unmarshal)
		if fieldValue.Kind() == reflect.Slice {
			if fieldValue.IsNil() || fieldValue.Len() == 0 {
				continue
			}
			anySlice := make([]any, fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				anySlice[j] = reflectValueToAny(elem)
			}
			result[fieldName] = anySlice
			continue
		}

		// Tipos primitivos
		result[fieldName] = fieldValue.Interface()
	}

	return result, nil
}

func extractJSONFieldName(tag string) string {
	if tag == "" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

// reflectValueToAny convierte un valor reflect a any, manejando recursión para structs
// y punteros a structs (como los elementos de []*OperacionFilter en slices).
func reflectValueToAny(v reflect.Value) any {
	// Desreferenciar punteros e interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		nested, err := structToMap(v.Interface())
		if err != nil {
			return nil
		}
		return nested
	}

	return v.Interface()
}
