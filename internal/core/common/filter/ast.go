package filter

// Clause es la interfaz que todo elemento del filtro debe cumplir
type Clause interface {
	// Accept permite aplicar el patrón Visitor (útil para generar SQL, validar, etc.)
	Accept(visitor Visitor) error
}

// LogicalClause representa AND, OR, NOT
type LogicalClause struct {
	Operator LogicalOperator
	Children []Clause // NOT usará solo el primer elemento
}

func (l *LogicalClause) Accept(v Visitor) error {
	return v.VisitLogical(l)
}

// PredicateClause representa una condición final: status == "active"
type PredicateClause struct {
	Field     string
	Condition Condition
	Value     any
}

func (p *PredicateClause) Accept(v Visitor) error {
	return v.VisitPredicate(p)
}

// Visitor define la interfaz para recorrer el árbol
type Visitor interface {
	VisitLogical(*LogicalClause) error
	VisitPredicate(*PredicateClause) error
}
