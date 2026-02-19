package filter

type Operator string

const (
	AND Operator = "and"
	OR  Operator = "or"
	NOT Operator = "not"
)

func IsOperator(s string) bool {
	return s == string(AND) || s == string(OR) || s == string(NOT)
}
