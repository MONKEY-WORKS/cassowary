// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

var _ EquationMember = &Param{}

type Param struct {
	Variable *Variable

	Context interface{}
}

func NewParam(val float64) *Param {
	param := &Param{
		Variable: &Variable{
			Value: val,
		},
	}
	param.Variable.owner = param

	return param
}

func NewParamWithContext(val float64, context interface{}) *Param {
	param := &Param{
		Variable: &Variable{
			Value: val,
		},
		Context: context,
	}
	param.Variable.owner = param

	return param
}

func (p *Param) IsConstant() bool {
	return false
}

func (p *Param) Value() float64 {
	return p.Variable.Value
}

func (p *Param) Add(member EquationMember) *Expression {
	return p.asExpression().Add(member)
}

func (p *Param) Sub(member EquationMember) *Expression {
	return p.asExpression().Sub(member)
}

func (p *Param) Div(member EquationMember) *Expression {
	return p.asExpression().Div(member)
}

func (p *Param) Mult(member EquationMember) *Expression {
	return p.asExpression().Mult(member)
}

func (p *Param) asExpression() *Expression {
	return NewExpression([]*Term{NewTerm(p.Variable, 1.0)}, 0.0)
}

func (p *Param) GreaterThanOrEqualTo(m EquationMember) *Constraint {
	return p.asExpression().GreaterThanOrEqualTo(m)
}

func (p *Param) LessThanOrEqualTo(m EquationMember) *Constraint {
	return p.asExpression().LessThanOrEqualTo(m)
}

func (p *Param) Equals(member EquationMember) *Constraint {
	return p.asExpression().createConstraint(member, EqualTo)
}
