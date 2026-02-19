package filter

import (
	"github.com/Sanaruca/condominio/internal/core"
)

type Condition string

// ConditionSpec is a map that defines the conditions and the value types they
// support. The key of the map is the value type, and the value is another map
// that defines the conditions and their expected value types.
//
// Example:
//
//	{
//	  ValueTypeString: {
//	    Equals: ExpectedConditionValueTypes(core.NewSet(ValueTypeString)),
//	    Like:   ExpectedConditionValueTypes(core.NewSet(ValueTypeString)),
//	  },
//	  ValueTypeInt: {
//	    Equals: ExpectedConditionValueTypes(core.NewSet(ValueTypeInt)),
//	  },
//	}
type ConditionSpec map[ValueType]ConditionDefinition

// ConditionDefinition is a map that defines the conditions and their expected
// value types. The key of the map is the condition, and the value is the
// expected value types.
//
// Example:
//
//	{
//	  Equals: ExpectedConditionValueTypes(core.NewSet(ValueTypeString)),
//	  Like:   ExpectedConditionValueTypes(core.NewSet(ValueTypeString)),
//	}
type ConditionDefinition map[Condition]ExpectedConditionValueTypes

func (cd ConditionDefinition) Has(c Condition) bool {
	_, ok := cd[c]
	return ok
}

// ExpectedConditionValueTypes is a set that defines the expected value types for
// a condition. It is implemented as a wrapper around the core.Set type.
//
// Example:
//
//	ExpectedConditionValueTypes(core.NewSet(ValueTypeString, ValueTypeInt))
type ExpectedConditionValueTypes core.Set[ValueType]

var condition_spec = ConditionSpec{
	String:   string_conditions,
	Int:      numeric_conditions,
	Float:    numeric_conditions,
	TypeBool: boolean_conditions,
}

var string_conditions = ConditionDefinition{
	Equals:        ExpectedConditionValueTypes(core.NewSet(String)),
	Like:          ExpectedConditionValueTypes(core.NewSet(String)),
	Regex:         ExpectedConditionValueTypes(core.NewSet(String)),
	In:            ExpectedConditionValueTypes(core.NewSet(StringList)),
	LengthEquals:  ExpectedConditionValueTypes(core.NewSet(String)),
	LengthMin:     ExpectedConditionValueTypes(core.NewSet(String)),
	LengthMax:     ExpectedConditionValueTypes(core.NewSet(String)),
	LengthGreater: ExpectedConditionValueTypes(core.NewSet(String)),
	LengthLess:    ExpectedConditionValueTypes(core.NewSet(String)),
}

var numeric_conditions = ConditionDefinition{
	Greater:        ExpectedConditionValueTypes(core.NewSet(Int, Float)),
	Between:        ExpectedConditionValueTypes(core.NewSet(IntList, FloatList)),
	Less:           ExpectedConditionValueTypes(core.NewSet(Int, Float)),
	GreaterOrEqual: ExpectedConditionValueTypes(core.NewSet(Int, Float)),
	LessOrEqual:    ExpectedConditionValueTypes(core.NewSet(Int, Float)),
}

// TODO:
// var date_condition_spec = ConditionSpec{
// 	Greater:        ExpectedConditionTypeValues(core.NewSet(ValueTypeDate)),
// 	Less:           ExpectedConditionTypeValues(core.NewSet(ValueTypeDate)),
// 	GreaterOrEqual: ExpectedConditionTypeValues(core.NewSet(ValueTypeDate)),
// 	LessOrEqual:    ExpectedConditionTypeValues(core.NewSet(ValueTypeDate)),
// }

var boolean_conditions = ConditionDefinition{
	Equals: ExpectedConditionValueTypes(core.NewSet(TypeBool)),
}

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

//	func (c Condition) IsDateCondition() bool {
//		return IsDateCondition(c)
//	}
func (c Condition) IsBoolCondition() bool {
	return IsBoolCondition(c)
}

func IsStringCondition(c Condition) bool {
	return string_conditions.Has(c)
}
func IsNumericCondition(c Condition) bool {
	return numeric_conditions.Has(c)
}

//	func IsDateCondition(c Condition) bool {
//		return date_condition_spec.Has(c)
//	}
func IsBoolCondition(c Condition) bool {
	return boolean_conditions.Has(c)
}
