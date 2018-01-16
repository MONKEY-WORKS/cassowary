package cassowary

type Relation int

const (
	EqualTo Relation = iota
	LessThanOrEqualTo
	GreaterThanOrEqualTo
)

type Constraint struct {
	relation   Relation
	expression *Expression

	priority Priority
}

func NewConstraint(exp *Expression, rel Relation) *Constraint {
	return &Constraint{
		relation:   rel,
		expression: exp,
		priority:   PriorityRequired,
	}
}
