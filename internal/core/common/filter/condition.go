package filter

import (
	"github.com/Sanaruca/condominio/internal/core"
)

type Condition string

const (
	Equals         Condition = "eq"
	Like           Condition = "li"
	Regex          Condition = "rx"
	In             Condition = "in"
	Between        Condition = "bt"
	Greater        Condition = "gt"
	Less           Condition = "lt"
	GreaterOrEqual Condition = "ge"
	LessOrEqual    Condition = "le"
	// alias for GreaterOrEqual
	Min Condition = "ge"
	// alias for LessOrEqual
	Max Condition = "le"
	// Length conditions
	LengthEquals         Condition = "leq"
	LengthGreaterOrEqual Condition = "lge"
	LengthGreater        Condition = "lgt"
	LengthLessOrEqual    Condition = "lle"
	LengthLess           Condition = "llt"
	// alias for LengthGreaterOrEqual
	LengthMin Condition = "lge"
	// alias for LengthLessOrEqual
	LengthMax Condition = "lle"
)

func (c Condition) Validate() core.Error {
	switch c {
	case Equals,
		Like,
		Regex,
		In,
		Between,
		Greater,
		Less,
		GreaterOrEqual,
		LessOrEqual,
		LengthEquals,
		LengthMin,
		LengthMax,
		LengthGreater,
		LengthLess:
		return nil
	default:
		return core.NewInvalidArgumentError("Operador '%s' no soportado", string(c))
	}
}

func (c Condition) String() string {
	return string(c)
}

func (c Condition) IsStringCondition() bool {
	return IsStringCondition(c)
}
func (c Condition) IsNumericCondition() bool {
	return IsNumericCondition(c)
}
func (c Condition) IsDateCondition() bool {
	return IsDateCondition(c)
}
func (c Condition) IsBoolCondition() bool {
	return IsBoolCondition(c)
}

func IsStringCondition(c Condition) bool {
	switch c {
	case Equals, Like, Regex, In, LengthEquals, LengthMin, LengthMax, LengthGreater, LengthLess:
		return true
	default:
		return false
	}
}
func IsNumericCondition(c Condition) bool {
	switch c {
	case Equals, Greater, Less, GreaterOrEqual, LessOrEqual:
		return true
	default:
		return false
	}
}
func IsDateCondition(c Condition) bool {
	switch c {
	case Equals, Greater, Less, GreaterOrEqual, LessOrEqual:
		return true
	default:
		return false
	}
}
func IsBoolCondition(c Condition) bool {
	switch c {
	case Equals:
		return true
	default:
		return false
	}
}
