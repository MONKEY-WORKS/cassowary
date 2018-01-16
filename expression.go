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

func (exp *Expression) Value() float64 {
	value := exp.constant
	for _, t := range exp.terms {
		value += t.Value()
	}
	return value
}

func (exp *Expression) IsConstant() bool {
	return len(exp.terms) == 0
}

func (exp *Expression) asExpression() *Expression {
	return exp
}

func (exp *Expression) Add(member EquationMember) *Expression {
	if cm, ok := member.(*ConstantMember); ok {
		return NewExpression(exp.terms, exp.constant+cm.value)
	}

	if param, ok := member.(*Param); ok {
		return NewExpression(append(exp.terms, NewTerm(param.Variable, 1.0)), exp.constant)
	}

	if term, ok := member.(*Term); ok {
		return NewExpression(append(exp.terms, term), exp.constant)
	}

	if exp, ok := member.(*Expression); ok {
		newArray := make([]*Term, len(exp.terms)+len(exp.terms))
		copy(newArray[0:len(exp.terms)], exp.terms)
		copy(newArray[len(exp.terms):], exp.terms)

		return NewExpression(newArray, exp.constant+exp.constant)
	}

	fmt.Println("Unknown EquationMember " + fmt.Sprintf("%T", member))
	panic("Unknown EquationMember " + fmt.Sprintf("%T", member))
}

func (exp *Expression) Sub(member EquationMember) *Expression {
	if cm, ok := member.(*ConstantMember); ok {
		return NewExpression(exp.terms, exp.constant-cm.value)
	}

	if param, ok := member.(*Param); ok {
		return NewExpression(append(exp.terms, NewTerm(param.Variable, -1.0)), exp.constant)
	}

	if term, ok := member.(*Term); ok {
		return NewExpression(append(exp.terms, NewTerm(term.variable, -term.coefficient)), exp.constant)
	}

	if exp, ok := member.(*Expression); ok {
		offset := len(exp.terms)
		newArray := make([]*Term, offset+len(exp.terms))
		copy(newArray[0:offset], exp.terms)
		for i, t := range exp.terms {
			newArray[offset+i] = NewTerm(t.variable, -t.coefficient)
		}

		return NewExpression(newArray, exp.constant-exp.constant)
	}

	panic("Unknown EquationMember " + fmt.Sprintf("%T", member))
}

type multiplication struct {
	multiplier   *Expression
	multiplicand float64
}

func (exp *Expression) Mult(member EquationMember) *Expression {
	args := exp.findMultiplierAndMultiplicand(member)
	if args == nil {
		return nil //, errors.New("Could not find constant multiplicand or multiplier")
	}

	return args.multiplier.applyMultiplicand(args.multiplicand) //, nil
}

func (exp *Expression) Div(member EquationMember) *Expression {
	if !member.IsConstant() {
		return nil //, errors.New("The divisor was not a constant expression")
	}

	return exp.applyMultiplicand(1.0 / member.Value()) //, nil
}

func (exp *Expression) findMultiplierAndMultiplicand(member EquationMember) *multiplication {
	if !exp.IsConstant() && !member.IsConstant() {
		return nil
	}

	if exp.IsConstant() {
		return &multiplication{
			member.asExpression(),
			exp.Value(),
		}
	} else { // so member is constant
		return &multiplication{
			exp, // TODO(Johannes): maybe copy is better
			member.Value(),
		}
	}
}

func (exp *Expression) applyMultiplicand(multiplicand float64) *Expression {

	terms := make([]*Term, len(exp.terms))
	for i, srcTerm := range exp.terms {
		terms[i] = NewTerm(srcTerm.variable, srcTerm.coefficient*multiplicand)
	}

	return NewExpression(terms, exp.constant*multiplicand)
}

func (exp *Expression) GreaterThanOrEqualTo(member EquationMember) *Constraint {
	return exp.createConstraint(member, GreaterThanOrEqualTo)
}

func (exp *Expression) LessThanOrEqualTo(member EquationMember) *Constraint {
	return exp.createConstraint(member, LessThanOrEqualTo)
}

func (exp *Expression) Equals(member EquationMember) *Constraint {
	return exp.createConstraint(member, EqualTo)
}

func (exp *Expression) createConstraint(member EquationMember, rel Relation) *Constraint {
	newTerms := make([]*Term, len(exp.terms))
	copy(newTerms, exp.terms)

	if cm, ok := member.(*ConstantMember); ok {
		return NewConstraint(NewExpression(newTerms, exp.constant-cm.value), rel)
	}

	if param, ok := member.(*Param); ok {
		return NewConstraint(NewExpression(append(newTerms, NewTerm(param.Variable, -1.0)), exp.constant), rel)
	}

	if term, ok := member.(*Term); ok {
		return NewConstraint(NewExpression(append(exp.terms, NewTerm(term.variable, -term.coefficient)), exp.constant), rel)
	}

	if exp, ok := member.(*Expression); ok {
		for _, t := range exp.terms {
			newTerms = append(newTerms, NewTerm(t.variable, -t.coefficient))
		}

		return NewConstraint(NewExpression(newTerms, exp.constant-exp.constant), rel)
	}

	panic("You should never reach this point ;)")
}
