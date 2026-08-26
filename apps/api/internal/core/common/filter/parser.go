package filter

import (
	"encoding/json"
	"fmt"

	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

// Parse transforma un input dinámico en un árbol de nodos (usa límite por defecto)
func Parse(input any) (Clause, error) {
	return parseWithDepth(input, 0, DEFAULT_MAX_DEPTH)
}

// Función interna que controla la recursividad
func parseWithDepth(input any, currentDepth, maxDepth int) (Clause, error) {
	// Validar capa de seguridad
	if currentDepth > maxDepth {
		return nil, fmt.Errorf(
			"límite de profundidad máxima excedido (%d): posible ataque o filtro demasiado complejo",
			maxDepth,
		)
	}

	logger.Debug("filter.Parse: Parseando %v a profundidad %d", input, currentDepth)

	// 1. Normalizar entrada a map[string]any
	rawMap, err := normalizeInput(input)
	if err != nil {
		return nil, err
	}

	return parseMap(rawMap, currentDepth, maxDepth)
}

func parseMap(m map[string]any, currentDepth, maxDepth int) (Clause, error) {
	var clauses []Clause

	for key, val := range m {
		switch key {
		case string(OPERATOR_AND), string(OPERATOR_OR):
			// Es un operador lógico que contiene un array
			childrenRaw, ok := val.([]any)
			if !ok {
				return nil, fmt.Errorf("el operador '%s' requiere un array", key)
			}

			logicalClause := &LogicalClause{Operator: LogicalOperator(key)}
			for _, child := range childrenRaw {
				// Pasamos el nivel de profundidad incrementado
				childClause, err := parseWithDepth(child, currentDepth+1, maxDepth)
				if err != nil {
					return nil, err
				}
				// Ignoramos subcláusulas vacías para no generar operadores lógicos vacíos
				if childClause != nil {
					logicalClause.Children = append(logicalClause.Children, childClause)
				}
			}
			// Un operador lógico sin hijos válidos se omite (equivale a "sin filtro")
			if len(logicalClause.Children) > 0 {
				clauses = append(clauses, logicalClause)
			}

		case string(OPERATOR_NOT):
			// Pasamos el nivel de profundidad incrementado
			childClause, err := parseWithDepth(val, currentDepth+1, maxDepth)
			if err != nil {
				return nil, err
			}
			// NOT sin operando válido se omite
			if childClause != nil {
				clauses = append(clauses, &LogicalClause{
					Operator: OPERATOR_NOT,
					Children: []Clause{childClause},
				})
			}

		default:
			// Es un campo convencional (ej: "age", "name")
			// Puede ser "age": 10 (Eq implícito) o "age": {"gt": 10}
			predicates, err := parseField(key, val)
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, predicates...)
		}
	}

	// Sin cláusulas válidas: no hay filtro que aplicar
	if len(clauses) == 0 {
		return nil, nil
	}

	// Si hay más de una clausula en este nivel del mapa, se asume AND implícito
	if len(clauses) == 1 {
		return clauses[0], nil
	}
	return &LogicalClause{Operator: OPERATOR_AND, Children: clauses}, nil
}

func parseField(field string, value any) ([]Clause, error) {
	// Caso Compuesto: "price": {"gt": 100, "lt": 200} -> Esto genera dos clausulas AND
	if valueMap, ok := value.(map[string]any); ok {
		var clauses []Clause
		for condStr, v := range valueMap {
			clauses = append(clauses, &PredicateClause{
				Field:     field,
				Condition: Condition(condStr),
				Value:     v,
			})
		}
		return clauses, nil
	}

	// Caso Simple (Igualdad implícita): "status": "active"
	return []Clause{&PredicateClause{
		Field:     field,
		Condition: CONDITON_EQ,
		Value:     value,
	}}, nil
}

func normalizeInput(v any) (map[string]any, error) {
	switch t := v.(type) {
	case map[string]any:
		return t, nil
	case string: // JSON String
		var res map[string]any
		err := json.Unmarshal([]byte(t), &res)
		return res, err
	default:
		// Fallback costoso pero seguro para tipos personalizados
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var res map[string]any
		err = json.Unmarshal(b, &res)
		return res, err
	}
}
