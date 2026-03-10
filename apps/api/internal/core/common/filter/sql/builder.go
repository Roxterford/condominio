package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

// SQLBuilder genera consultas SQL a partir del AST.
type SQLBuilder struct {
	sql          strings.Builder
	args         []any
	ctx          context.Context
	currentDepth int
	maxDepth     int
}

// SQLBuilderOption define un patrón de opciones para configurar el builder
type SQLBuilderOption func(*SQLBuilder)

// WithContext inyecta un contexto para cancelación temprana
func WithContext(ctx context.Context) SQLBuilderOption {
	return func(b *SQLBuilder) { b.ctx = ctx }
}

// WithMaxDepth define la profundidad máxima del AST para evitar DDoS
func WithMaxDepth(depth int) SQLBuilderOption {
	return func(b *SQLBuilder) { b.maxDepth = depth }
}

// NewSQLBuilder inicializa el builder con valores seguros por defecto.
func NewSQLBuilder(opts ...SQLBuilderOption) *SQLBuilder {
	b := &SQLBuilder{
		args:     make([]any, 0, 8),
		ctx:      context.Background(),
		maxDepth: filter.DEFAULT_MAX_DEPTH, // Límite por defecto seguro para la mayoría de casos de negocio
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *SQLBuilder) Build(root filter.Clause) (string, []any, error) {
	if root == nil {
		return "", nil, nil
	}
	if err := root.Accept(b); err != nil {
		return "", nil, err
	}
	return b.sql.String(), b.args, nil
}

func (b *SQLBuilder) VisitLogical(clause *filter.LogicalClause) error {
	// 1. Verificación de Contexto (Cancelación)
	if err := b.ctx.Err(); err != nil {
		return fmt.Errorf("construcción de SQL cancelada: %w", err)
	}

	// 2. Verificación de Profundidad (DDoS)
	b.currentDepth++
	if b.currentDepth > b.maxDepth {
		return fmt.Errorf("profundidad máxima del filtro excedida (%d)", b.maxDepth)
	}
	defer func() { b.currentDepth-- }() // Restablecemos al salir del scope

	if len(clause.Children) == 0 {
		return nil
	}

	// 3. Manejo de NOT y optimización con Leyes de De Morgan
	if clause.Operator == filter.OPERATOR_NOT {
		return b.handleDeMorgan(clause.Children[0])
	}

	// 4. Procesamiento normal de AND / OR
	var sql_operator string
	switch clause.Operator {
	case filter.OPERATOR_AND:
		sql_operator = " AND "
	case filter.OPERATOR_OR:
		sql_operator = " OR "
	default:
		return fmt.Errorf("operador lógico desconocido: '%s'", clause.Operator)
	}

	b.sql.WriteString("(")
	for i, sub_clause := range clause.Children {
		if i > 0 {
			b.sql.WriteString(sql_operator)
		}
		if err := sub_clause.Accept(b); err != nil {
			return err
		}
	}
	b.sql.WriteString(")")

	return nil
}

func (b *SQLBuilder) VisitPredicate(clause *filter.PredicateClause) error {
	return b.writePredicate(clause.Field, clause.Condition, clause.Value)
}

// writePredicate está separado para ser reutilizado por la optimización de De Morgan
func (b *SQLBuilder) writePredicate(field string, cond filter.Condition, value any) error {
	switch cond {
	case filter.CONDITON_EQ:
		b.sql.WriteString(fmt.Sprintf("%s = ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_NEQ:
		b.sql.WriteString(fmt.Sprintf("%s != ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_GT:
		b.sql.WriteString(fmt.Sprintf("%s > ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_GTE:
		b.sql.WriteString(fmt.Sprintf("%s >= ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_LT:
		b.sql.WriteString(fmt.Sprintf("%s < ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_LTE:
		b.sql.WriteString(fmt.Sprintf("%s <= ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_LIKE:
		b.sql.WriteString(fmt.Sprintf("%s LIKE ?", field))
		b.args = append(b.args, value)
	case filter.CONDITON_IN:
		return b.handleInCondition(field, value)
	default:
		return fmt.Errorf("condición no soportada: '%s'", cond)
	}
	return nil
}

// handleDeMorgan aplica las leyes lógicas para invertir predicados o propagar negaciones
func (b *SQLBuilder) handleDeMorgan(child filter.Clause) error {
	switch c := child.(type) {
	case *filter.PredicateClause:
		// Aplicamos De Morgan en el nodo hoja invirtiendo el operador
		inverseCond, err := invertCondition(c.Condition)
		if err != nil {
			return err
		}
		// Escribimos el predicado ya invertido (sin usar la palabra NOT en SQL)
		return b.writePredicate(c.Field, inverseCond, c.Value)

	case *filter.LogicalClause:
		// Para simplificar la inyección de SQL, si el hijo es otro nodo lógico,
		// delegamos la evaluación a la base de datos agrupando la negación.
		// Un paso extra (muy avanzado) sería reescribir el AST mutando el árbol
		// antes de procesarlo, pero hacerlo on-the-fly de forma limpia es así:
		b.sql.WriteString("NOT (")
		if err := c.Accept(b); err != nil {
			return err
		}
		b.sql.WriteString(")")
		return nil
	}
	return fmt.Errorf("tipo de cláusula no soportado en De Morgan")
}

func (b *SQLBuilder) handleInCondition(field string, value any) error {
	var placeholders []string

	// Evaluamos los tipos estáticos más comunes de forma explícita
	switch v := value.(type) {
	case []any:
		if len(v) == 0 {
			b.sql.WriteString("1=0")
			return nil
		}
		for _, val := range v {
			placeholders = append(placeholders, "?")
			b.args = append(b.args, val)
		}
	case []string: // Soporte por si se construye el AST manualmente sin JSON
		if len(v) == 0 {
			b.sql.WriteString("1=0")
			return nil
		}
		for _, val := range v {
			placeholders = append(placeholders, "?")
			b.args = append(b.args, val)
		}
	case []int:
		if len(v) == 0 {
			b.sql.WriteString("1=0")
			return nil
		}
		for _, val := range v {
			placeholders = append(placeholders, "?")
			b.args = append(b.args, val)
		}
	default:
		return fmt.Errorf("la condición IN requiere un slice ([]any, []string, []int), se recibió %T", value)
	}

	b.sql.WriteString(fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", ")))
	return nil
}

// invertCondition es el core algebraico para De Morgan
func invertCondition(cond filter.Condition) (filter.Condition, error) {
	switch cond {
	case filter.CONDITON_EQ:
		return filter.CONDITON_NEQ, nil
	case filter.CONDITON_NEQ:
		return filter.CONDITON_EQ, nil
	case filter.CONDITON_GT:
		return filter.CONDITON_LTE, nil
	case filter.CONDITON_GTE:
		return filter.CONDITON_LT, nil
	case filter.CONDITON_LT:
		return filter.CONDITON_GTE, nil
	case filter.CONDITON_LTE:
		return filter.CONDITON_GT, nil
	case filter.CONDITON_LIKE:
		return "", fmt.Errorf("la inversión de LIKE (NOT LIKE) requiere lógica dedicada")
	case filter.CONDITON_IN:
		return "", fmt.Errorf("la inversión de IN (NOT IN) requiere lógica dedicada")
	default:
		return "", fmt.Errorf("imposible invertir condición: '%s'", cond)
	}
}
