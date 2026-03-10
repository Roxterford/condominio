package utils

import "reflect"

func GetFieldNamesByStruct(t any, tag ...string) []string {
	var fieldNames []string
	val := reflect.ValueOf(t)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if len(tag) > 0 {
			// Si se pasa un tag, verificamos si el campo tiene ese tag
			if field.Tag.Get(tag[0]) != "" {
				fieldNames = append(fieldNames, field.Name)
			}
		} else {
			// Si no se pasa tag, simplemente agregamos el nombre del campo
			fieldNames = append(fieldNames, field.Name)
		}
	}
	return fieldNames
}
