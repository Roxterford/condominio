package filter

import "github.com/Sanaruca/condominio/internal/core"

// Query representa el mapa dinámico que envía el usuario
//
//   - Formato (base): { <campo | operador>: <query> | <valor> | { <condicion>: <valor> } }
//
//   - Formato (simple): { <campo>: <valor> }
//     Ej. { "status": "active" }
//
//   - Formato (compuesto): { <campo>: { <condicion>: <valor> } }
//     Ej. { "status": { "eq": "active" } }
//
//   - Formato (operadores): { <operador>: <query> }
//     Ej. { "and": [{ "status": "active" }, { "price": { "gt": 100 } }] }
//
//     Operadores:
//
//   - AND: { "and": [{...}, {...}] }
//
//   - OR: { "or": [{...}, {...}] }
//
//   - NOT: { "not": {...} }
//
//     Campos:
//
//   - <campo>: <valor> (igualdad implícita)
//
//   - <campo>: { <condicion>: <valor> }
//
//     Ejemplos:
//
//   - { "status": "active" }
//
//   - { "price": { "gt": 100 } }
//
//   - { "and": [{ "status": "active" }, { "price": { "gt": 100 } }] }
type Query map[string]any

func (q Query) Validate() core.Error {

	for key, value := range q {
		switch key {
		case string(AND), string(OR):
			// Esperamos un slice de sub-consultas: <operador>: [{...}, {...}]
			children, ok := value.([]any)
			if !ok {
				return core.NewInvalidArgumentError("Operador '%s' debe ser un array", key)
			}
			for _, child := range children {
				sub_query, ok := child.(map[string]any)
				if !ok {
					return core.NewInvalidArgumentError("Campo '%s' debe ser un objeto", key)
				}
				if err := Query(sub_query).Validate(); err != nil {
					return err
				}
			}
		case string(NOT):
			// "not" suele recibir un solo objeto: "not": { "status": "active" }
			sub_query, ok := value.(map[string]any)
			if !ok {
				return core.NewInvalidArgumentError("Operador 'not' debe ser un objeto")
			}
			return Query(sub_query).Validate()

		default:
			// Es un formato simple o compuesto: <campo>: { <condicion>: <valor> } o <campo>: <valor>
			if err := processField(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func processField(field string, value any) core.Error {
	// 1. Validar si el campo es permitido (Whitelist)
	// if !isFilterable(field) { return ... }

	switch v := value.(type) {
	case map[string]any:
		// Caso compuesto: <campo>: { <condicion>: <valor> }
		for condition_key, val := range v {
			cond := Condition(condition_key)
			if err := validateTypeAndCondition(field, cond, val); err != nil {
				return err
			}
		}
	default:
		// Caso simple: <campo>: <valor>
		if err := validateTypeAndCondition(field, Equals, v); err != nil {
			return err
		}
	}
	return nil
}

func validateTypeAndCondition(field string, cond Condition, value any) core.Error {
	switch value.(type) {
	case string:
		if !cond.IsStringCondition() {
			return core.NewInvalidArgumentError("Condicion '%s' no valida para el campo string '%s'", cond, field)
		}
	case int, float64: // JSON deserealiza números como float64 por defecto
		if !cond.IsNumericCondition() {
			return core.NewInvalidArgumentError("Condicion '%s' no valida para el campo numerico '%s'", cond, field)
		}
	case bool:
		if !cond.IsBoolCondition() {
			return core.NewInvalidArgumentError("Condicion '%s' no valida para el campo bool '%s'", cond, field)
		}
	}
	return nil
}
