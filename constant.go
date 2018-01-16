// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

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

func (c *ConstantMember) Add(member EquationMember) *Expression {
	return c.asExpression().Add(member)
}

func (c *ConstantMember) Sub(member EquationMember) *Expression {
	return c.asExpression().Sub(member)
}

func (c *ConstantMember) Div(member EquationMember) *Expression {
	return c.asExpression().Div(member)
}

func (c *ConstantMember) Mult(member EquationMember) *Expression {
	return c.asExpression().Mult(member)
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
