package filter

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

func NewExpression(field string, condition Condition, value any) Expression {
	return Expression{
		field:     field,
		condition: condition,
		value:     walue{value},
	}
}

func NewExpressionFromQuery(key string, value any) (Expression, core.Error) {
	logger.Debug("[%s]: %v", key, value)
	switch v := value.(type) {
	case map[string]any:
		// Caso compuesto: <campo>: { <condicion>: <valor> }

		// Solo tomaremos la primera condición
		for conditionKey, val := range v {
			cond := Condition(conditionKey)
			expr := NewExpression(key, cond, val)
			// Si la expresión es válida, la retornamos de inmediato
			return expr, expr.Validate()
		}

	default:
		// Caso simple: <campo>: <valor> -> Se asume Equals
		return Expression{
			field:     key,
			condition: Equals,
			value:     walue{value},
		}, nil
	}

	// Si llegamos aquí es porque el mapa estaba vacío
	// o ninguna condición dentro del mapa fue válida.
	return Expression{}, core.NewInvalidArgumentError("INVALID_EXPRESSION: No se pudo crear una expresión válida")
}

type Expression struct {
	field     string
	condition Condition
	value     Value
}

func (ex Expression) Field() string {
	return ex.field
}

func (ex Expression) Condition() Condition {
	return ex.condition
}

func (ex Expression) Value() any {
	return ex.value
}

func (ex Expression) Validate() core.Error {
	return ex.validate(nil)
}

func (ex Expression) ValidateWithSpec(spec Spec) core.Error {
	return ex.validate(spec)
}

func (ex Expression) validate(spec Spec) core.Error {

	switch Operator(ex.field) {
	case NOT, AND, OR:
		return errors.New(errors.INVALID_ARGUMENT, "INVALID_EXPRESSION: Campo '%s' es un operador condicional reservado", ex.field)
	}

	specJson, _ := json.MarshalIndent(spec, "", "    ")
	logger.Debug("Expression.validate: Spec: " + string(specJson))
	logger.Debug("Expression.validate: Spec is empty: %v", spec.IsEmpty())
	logger.Debug("Expression.validate: Spec has key '%s': %v", ex.field, spec.HasKey(ex.field))

	if !spec.IsEmpty() {
		if !spec.HasKey(ex.field) {
			return core.NewInvalidArgumentError("INVALID_EXPRESSION: Campo '%s' no permitido", ex.field)
		}
		value_type := spec[ex.field]
		if !condition_spec[value_type].Has(ex.condition) {
			return core.NewInvalidArgumentError("INVALID_EXPRESSION: Campo '%s' de tipo '%s' no espera una condición '%s'", ex.field, value_type, ex.condition)
		}
	} else {
		// Validaremos que la condición sea válida basandonos en el tipo de valor recibido
		value_type := ex.value.Type()
		logger.Warning("Expression.validate: No se proporcionó especificaciones de filtro para el campo '%s', se evaluará la expresión '{ %s: { %s: %v } }' basado en el tipo de valor recibido como '%s'.", ex.field, ex.field, ex.condition, ex.value.Raw(), value_type)
		definition := condition_spec[value_type]

		if definition == nil {
			return core.NewInvalidArgumentError("INVALID_EXPRESSION: El tipo de valor '%s' no está definido en las especificacion de las condicionales", value_type)
		}

		if !definition.Has(ex.condition) {
			return core.NewInvalidArgumentError("INVALID_EXPRESSION: El tipo de valor '%s' no espera una condición '%s'", value_type, ex.condition)
		}
	}

	return nil

}
