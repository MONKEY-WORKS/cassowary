package cassowary

type EquationMember interface {
	asExpression() *Expression

	IsConstant() bool

	Value() float64

	GreaterThanOrEqualTo(m EquationMember) *Constraint

	LessThanOrEqualTo(m EquationMember) *Constraint

	Equals(m EquationMember) *Constraint

	Add(m EquationMember) *Expression

	Sub(m EquationMember) *Expression

	Mult(m EquationMember) *Expression //, error)

	Div(m EquationMember) *Expression //, error)
}
