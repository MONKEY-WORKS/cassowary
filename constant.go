package cassowary

type ConstantMember struct {
	value float64
}

var _ EquationMember = &ConstantMember{}

func CM(value float64) *ConstantMember {
	return NewConstantMember(value)
}

func NewConstantMember(value float64) *ConstantMember {
	return &ConstantMember{
		value: value,
	}
}

func (c *ConstantMember) Value() float64 {
	return c.value
}

func (c *ConstantMember) IsConstant() bool {
	return true
}

func (p *ConstantMember) Add(member EquationMember) *Expression {
	return p.asExpression().Add(member)
}

func (p *ConstantMember) Sub(member EquationMember) *Expression {
	return p.asExpression().Sub(member)
}

func (p *ConstantMember) Div(member EquationMember) *Expression {
	return p.asExpression().Div(member)
}

func (p *ConstantMember) Mult(member EquationMember) *Expression {
	return p.asExpression().Mult(member)
}

func (c *ConstantMember) asExpression() *Expression {
	return NewExpression([]*Term{}, c.value)
}

func (c *ConstantMember) GreaterThanOrEqualTo(m EquationMember) *Constraint {
	return c.asExpression().GreaterThanOrEqualTo(m)
}

func (c *ConstantMember) LessThanOrEqualTo(m EquationMember) *Constraint {
	return c.asExpression().LessThanOrEqualTo(m)
}

func (c *ConstantMember) Equals(member EquationMember) *Constraint {
	return c.asExpression().createConstraint(member, EqualTo)
}
