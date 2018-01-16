// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

import (
	"fmt"
)

var _ EquationMember = &Expression{}

type Expression struct {
	terms []*Term

	constant float64
}

func NewExpression(terms []*Term, constant float64) *Expression {
	return &Expression{
		terms:    terms,
		constant: constant,
	}
}

func FromExpression(exp *Expression) *Expression {
	terms := make([]*Term, len(exp.terms))
	copy(terms, exp.terms)
	return &Expression{
		terms:    terms,
		constant: exp.constant,
	}
}

func (c *Expression) Value() float64 {
	value := c.constant
	for _, t := range c.terms {
		value += t.Value()
	}
	return value
}

func (c *Expression) IsConstant() bool {
	return len(c.terms) == 0
}

func (c *Expression) asExpression() *Expression {
	return c
}

func (c *Expression) Add(member EquationMember) *Expression {
	if cm, ok := member.(*ConstantMember); ok {
		return NewExpression(c.terms, c.constant+cm.value)
	}

	if param, ok := member.(*Param); ok {
		return NewExpression(append(c.terms, NewTerm(param.Variable, 1.0)), c.constant)
	}

	if term, ok := member.(*Term); ok {
		return NewExpression(append(c.terms, term), c.constant)
	}

	if exp, ok := member.(*Expression); ok {
		newArray := make([]*Term, len(c.terms)+len(exp.terms))
		copy(newArray[0:len(c.terms)], c.terms)
		copy(newArray[len(c.terms):], exp.terms)

		return NewExpression(newArray, c.constant+exp.constant)
	}

	fmt.Println("Unknown EquationMember " + fmt.Sprintf("%T", member))
	panic("Unknown EquationMember " + fmt.Sprintf("%T", member))
}

func (c *Expression) Sub(member EquationMember) *Expression {
	if cm, ok := member.(*ConstantMember); ok {
		return NewExpression(c.terms, c.constant-cm.value)
	}

	if param, ok := member.(*Param); ok {
		return NewExpression(append(c.terms, NewTerm(param.Variable, -1.0)), c.constant)
	}

	if term, ok := member.(*Term); ok {
		return NewExpression(append(c.terms, NewTerm(term.variable, -term.coefficient)), c.constant)
	}

	if exp, ok := member.(*Expression); ok {
		offset := len(c.terms)
		newArray := make([]*Term, offset+len(exp.terms))
		copy(newArray[0:offset], c.terms)
		for i, t := range exp.terms {
			newArray[offset+i] = NewTerm(t.variable, -t.coefficient)
		}

		return NewExpression(newArray, c.constant-exp.constant)
	}

	panic("Unknown EquationMember " + fmt.Sprintf("%T", member))
}

type multiplication struct {
	multiplier   *Expression
	multiplicand float64
}

func (c *Expression) Mult(member EquationMember) *Expression {
	args := c.findMultiplierAndMultiplicand(member)
	if args == nil {
		return nil //, errors.New("Could not find constant multiplicand or multiplier")
	}

	return args.multiplier.applyMultiplicand(args.multiplicand) //, nil
}

func (c *Expression) Div(member EquationMember) *Expression {
	if !member.IsConstant() {
		return nil //, errors.New("The divisor was not a constant expression")
	}

	return c.applyMultiplicand(1.0 / member.Value()) //, nil
}

func (c *Expression) findMultiplierAndMultiplicand(member EquationMember) *multiplication {
	if !c.IsConstant() && !member.IsConstant() {
		return nil
	}

	if c.IsConstant() {
		return &multiplication{
			member.asExpression(),
			c.Value(),
		}
	} else { // so member is constant
		return &multiplication{
			c, // TODO(Johannes): maybe copy is better
			member.Value(),
		}
	}
}

func (c *Expression) applyMultiplicand(multiplicand float64) *Expression {

	terms := make([]*Term, len(c.terms))
	for i, srcTerm := range c.terms {
		terms[i] = NewTerm(srcTerm.variable, srcTerm.coefficient*multiplicand)
	}

	return NewExpression(terms, c.constant*multiplicand)
}

func (c *Expression) GreaterThanOrEqualTo(member EquationMember) *Constraint {
	return c.createConstraint(member, GreaterThanOrEqualTo)
}

func (c *Expression) LessThanOrEqualTo(member EquationMember) *Constraint {
	return c.createConstraint(member, LessThanOrEqualTo)
}

func (c *Expression) Equals(member EquationMember) *Constraint {
	return c.createConstraint(member, EqualTo)
}

func (c *Expression) createConstraint(member EquationMember, rel Relation) *Constraint {
	newTerms := make([]*Term, len(c.terms))
	copy(newTerms, c.terms)

	if cm, ok := member.(*ConstantMember); ok {
		return NewConstraint(NewExpression(newTerms, c.constant-cm.value), rel)
	}

	if param, ok := member.(*Param); ok {
		return NewConstraint(NewExpression(append(newTerms, NewTerm(param.Variable, -1.0)), c.constant), rel)
	}

	if term, ok := member.(*Term); ok {
		return NewConstraint(NewExpression(append(c.terms, NewTerm(term.variable, -term.coefficient)), c.constant), rel)
	}

	if exp, ok := member.(*Expression); ok {
		for _, t := range exp.terms {
			newTerms = append(newTerms, NewTerm(t.variable, -t.coefficient))
		}

		return NewConstraint(NewExpression(newTerms, c.constant-exp.constant), rel)
	}

	panic("You should never reach this point ;)")
}
