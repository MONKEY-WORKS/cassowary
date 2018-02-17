// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

// Ensure that ConstantMember implements EquationMember interface
var _ EquationMember = &ConstantMember{}

// ConstantMember describes a constant part of an Expression which won't change when added to a Solver
type ConstantMember struct {
	value float64
}

// CM just creates a new ConstantMember containing the given constant specified by value
func CM(value float64) *ConstantMember {
	return NewConstantMember(value)
}

// NewConstantMember creates a new ConstantMember containing the given constant specified by value
func NewConstantMember(value float64) *ConstantMember {
	return &ConstantMember{
		value: value,
	}
}

// Value returns the value of this constant
func (c *ConstantMember) Value() float64 {
	return c.value
}

// IsConstant returns true if this EquationMember is a constant (Which is always true for constants)
func (c *ConstantMember) IsConstant() bool {
	return true
}

// Add creates an expression which represents the sum of this ConstantMember and another EquationMember
func (c *ConstantMember) Add(member EquationMember) *Expression {
	return c.asExpression().Add(member)
}

// Sub creates an expression which represents the
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
