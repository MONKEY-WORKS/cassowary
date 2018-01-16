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

func (term *Term) Value() float64 {
	return term.coefficient * term.variable.Value
}

func (term *Term) IsConstant() bool {
	return false
}

func (term *Term) asExpression() *Expression {
	return NewExpression([]*Term{NewTerm(term.variable, term.coefficient)}, 0.0)
}

func (term *Term) Add(member EquationMember) *Expression {
	return term.asExpression().Add(member)
}

func (term *Term) Sub(member EquationMember) *Expression {
	return term.asExpression().Sub(member)
}

func (term *Term) Div(member EquationMember) *Expression {
	return term.asExpression().Div(member)
}

func (term *Term) Mult(member EquationMember) *Expression {
	return term.asExpression().Mult(member)
}

func (term *Term) GreaterThanOrEqualTo(m EquationMember) *Constraint {
	return term.asExpression().GreaterThanOrEqualTo(m)
}

func (term *Term) LessThanOrEqualTo(m EquationMember) *Constraint {
	return term.asExpression().LessThanOrEqualTo(m)
}

func (term *Term) Equals(member EquationMember) *Constraint {
	return term.asExpression().createConstraint(member, EqualTo)
}
