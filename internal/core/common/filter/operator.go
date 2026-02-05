package filter

type Operator string

const (
	AND Operator = "and"
	OR  Operator = "or"
	NOT Operator = "not"
)
