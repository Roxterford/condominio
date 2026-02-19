package filter

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/core"
)

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

func NewQuery() Query {
	return Query{}
}

func NewQueryFrom(v any) (Query, bool) {
	switch t := v.(type) {
	case Query:
		return t, true
	case map[string]any:
		return Query(t), true
	case string:
		var res map[string]any
		err := json.Unmarshal([]byte(t), &res)
		return Query(res), err == nil
	case []byte:
		var res map[string]any
		err := json.Unmarshal(t, &res)
		return Query(res), err == nil
	default:
		data, _ := json.Marshal(v)

		var res map[string]any
		err := json.Unmarshal(data, &res)
		if err != nil {
			return nil, false
		}
		return Query(res), true
	}
}

func (q Query) Get(key string) any {
	return q[key]
}

func (q *Query) Set(key string, value any) *Query {
	(*q)[key] = value
	return q
}

func (q Query) Keys() []string {

	keys := make([]string, 0, len(q))
	for key := range q {
		if !IsOperator(key) {
			keys = append(keys, key)
		}
	}

	return keys
}

func (q Query) Operators() []string {

	operators := make([]string, 0, len(q))
	for key := range q {
		if IsOperator(key) {
			operators = append(operators, key)
		}
	}
	return operators
}

func (q Query) Validate() core.Error {
	return q.validate(nil)
}

func (q Query) ValidateWithSpec(spec Spec) core.Error {
	return q.validate(spec)
}

func (q Query) validate(spec Spec) core.Error {

	for key, value := range q {
		switch key {
		case string(AND), string(OR):
			// Esperamos un slice de sub-consultas: <operador>: [{...}, {...}]
			children, ok := value.([]any)
			if !ok {
				return core.NewInvalidArgumentError("Operador '%s' debe ser un array", key)
			}
			for index, child := range children {
				sub_query, ok := NewQueryFrom(child)
				if !ok {
					return core.NewInvalidArgumentError("Condicion '%s[%d]' debe ser un objeto", key, index)
				}
				if err := sub_query.validate(spec); err != nil {
					return err
				}
			}
		case string(NOT):
			// "not" suele recibir un solo objeto: "not": { "status": "active" }
			sub_query, ok := NewQueryFrom(value)
			if !ok {
				return core.NewInvalidArgumentError("Operador 'not' debe ser un objeto")
			}
			return Query(sub_query).validate(spec)

		default:
			// Es un formato simple o compuesto: <campo>: { <condicion>: <valor> } o <campo>: <valor>
			if err := processField(key, value, spec); err != nil {
				return err
			}
		}
	}
	return nil
}

func processField(field string, value any, spec Spec) core.Error {

	// Verificamos si el campo está en la lista blanca
	if !spec.IsEmpty() {
		if !spec.HasKey(field) {
			return core.NewInvalidArgumentError("Campo '%s' no permitido", field)
		}
	}

	switch v := value.(type) {
	case map[string]any:
		// Caso compuesto: <campo>: { <condicion>: <valor> }
		for condition_key, val := range v {
			cond := Condition(condition_key)
			if err := NewExpression(field, cond, val).Validate(); err != nil {
				return err
			}
		}
	default:
		// Caso simple: <campo>: <valor>
		if err := NewExpression(field, Equals, value).Validate(); err != nil {
			return err
		}
	}
	return nil
}
