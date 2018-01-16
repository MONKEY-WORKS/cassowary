// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

var _ EquationMember = &Term{}

type Term struct {
	variable *Variable

	coefficient float64
}

func NewTerm(v *Variable, c float64) *Term {
	return &Term{
		v,
		c,
	}
}

func (c *Term) Value() float64 {
	return c.coefficient * c.variable.Value
}

func (c *Term) IsConstant() bool {
	return false
}

func (c *Term) asExpression() *Expression {
	return NewExpression([]*Term{NewTerm(c.variable, c.coefficient)}, 0.0)
}

func (p *Term) Add(member EquationMember) *Expression {
	return p.asExpression().Add(member)
}

func (p *Term) Sub(member EquationMember) *Expression {
	return p.asExpression().Sub(member)
}

func (p *Term) Div(member EquationMember) *Expression {
	return p.asExpression().Div(member)
}

func (p *Term) Mult(member EquationMember) *Expression {
	return p.asExpression().Mult(member)
}

func (c *Term) GreaterThanOrEqualTo(m EquationMember) *Constraint {
	return c.asExpression().GreaterThanOrEqualTo(m)
}

func (c *Term) LessThanOrEqualTo(m EquationMember) *Constraint {
	return c.asExpression().LessThanOrEqualTo(m)
}

func (c *Term) Equals(member EquationMember) *Constraint {
	return c.asExpression().createConstraint(member, EqualTo)
}
