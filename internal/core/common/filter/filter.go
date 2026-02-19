package filter

import (
	"github.com/Sanaruca/condominio/internal/core"
)

// TODO: La query pasada como argumento deberia ser validada, sin enbargo, ya
// que esto es un a funcion recursiva, debemos asegurarnos que query.Validate()
// se ejecute una sola vez. El metodo Validate() de Query se encarga de validar
// toda su estrura.
func New[T Filterable](query Query) (Filter[T], core.Error) {
	var filterable T = *new(T)

	var filter Filter[T] = Filter[T]{
		filterable: filterable,
	}

	if query == nil {
		filter.query = NewQuery()
	} else {

		filter.query = query

	}

	for key, value := range filter.query {

		switch key {
		case string(AND), string(OR):

			expressions := value.([]any)

			for _, expression := range expressions {

				query, ok := NewQueryFrom(expression)

				if !ok {
					return Filter[T]{}, core.NewInvalidArgumentError("Invalid expression")
				}

				if key == string(AND) {
					filter_and, err := New[T](query)
					if err != nil {
						return Filter[T]{}, err
					}
					filter.And = append(filter.And, filter_and)
				} else {
					filter_or, err := New[T](query)
					if err != nil {
						return Filter[T]{}, err
					}
					filter.Or = append(filter.Or, filter_or)
				}

			}

		case string(NOT):

			query, ok := NewQueryFrom(value)
			if !ok {
				return Filter[T]{}, core.NewInvalidArgumentError("Invalid expression")
			}

			sub_filter, err := New[T](query)
			if err != nil {
				return Filter[T]{}, err
			}

			filter.Not = &sub_filter

		default:

			expression, err := NewExpressionFromQuery(key, value)
			if err != nil {
				return Filter[T]{}, err
			}
			if err := expression.ValidateWithSpec(filterable.FilterSpec()); err != nil {
				return Filter[T]{}, err
			}
			filter.Expression = &expression

			filter.And = append(filter.And, Filter[T]{filterable: filterable, Expression: &expression})
		}

	}

	return filter, nil

}

// A filter is a collection of expressions that can be used to filter a collection of items.
//
// A Valid Filter exprect a Filterable type that has at least one filterable key.
type Filter[T Filterable] struct {
	filterable T
	query      Query

	And []Filter[T]
	Or  []Filter[T]
	Not *Filter[T]

	Expression *Expression
}

func (f Filter[T]) Keys() []string {
	return f.query.Keys()
}

// A filter is valid if it has at least one filterable key
func (f Filter[T]) Validate() core.Error {

	if f.filterable.FilterSpec().IsEmpty() {
		return core.NewInvalidArgumentError("No se encontraron campos validos para filtrar")
	}
	return f.query.ValidateWithSpec(f.filterable.FilterSpec())
}
