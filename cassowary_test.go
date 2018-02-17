// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

import (
	"fmt"
	"testing"
)

func TestParam(t *testing.T) {
	p := NewParam(22.0)

	if p.Value() != 22.0 {
		t.Error("Value did not match expected value", 22.0, ", was", p.Value(), "instead")
	}
}

func TestParamAdd(t *testing.T) {
	p := NewParam(22.0)
	res := p.Add(CM(22.0))

	if res.Value() != 44.0 {
		t.Error("Value did not match expected value", 44.0, ", was", res.Value(), "instead")
	}
}

func TestParamSub(t *testing.T) {
	p := NewParam(22.0)
	res := p.Sub(CM(20.0))

	if res.Value() != 2.0 {
		t.Error("Value did not match expected value", 2.0, ", was", res.Value(), "instead")
	}
}

func TestTerm(t *testing.T) {
	p := NewTerm(NewVariable(22.0), 2.0)
	expectedValue := 44.0

	if p.Value() != expectedValue {
		t.Error("Value did not match expected value", expectedValue, ", was", p.Value(), "instead")
	}
}

func TestExpression(t *testing.T) {
	terms := []*Term{
		NewTerm(NewVariable(22.0), 2.0),
		NewTerm(NewVariable(1.0), 1.0),
	}
	expression := NewExpression(terms, 40.0)
	expectedValue := 85.0

	if expression.Value() != expectedValue {
		t.Error("Value did not match expected value", expectedValue, ", was", expression.Value(), "instead")
	}
}

func TestExpression_Add(t *testing.T) {
	param1 := NewParam(10.0)
	param2 := NewParam(10.0)
	param3 := NewParam(22.0)

	expression := param1.Add(param2).Add(param3)
	expectedValue := 42.0

	if expression.Value() != expectedValue {
		t.Error("Value did not match expected value", expectedValue, ", was", expression.Value(), "instead")
	}
}

func expect(t *testing.T, member EquationMember, expectedValue float64) {
	if member.Value() != expectedValue {
		t.Error("Value did not match expected value", expectedValue, ", was", member.Value(), "instead")
	}
}

func TestExpression_Sums(t *testing.T) {
	expression := NewParam(10.0).Add(CM(5.0))
	expect(t, expression, 15.0)

	// constant
	expect(t, expression.Add(CM(2.0)), 17.0)
	expect(t, expression.Sub(CM(2.0)), 13.0)

	expect(t, expression, 15.0)

	// params
	param := NewParam(2.0)
	expect(t, expression.Add(param), 17.0)
	expect(t, expression.Sub(param), 13.0)

	expect(t, expression, 15.0)

	// terms
	term := NewTerm(param.Variable, 2.0)
	expect(t, expression.Add(term), 19.0)
	expect(t, expression.Sub(term), 11.0)

	expect(t, expression, 15.0)

	//expression
	expression2 := NewParam(7.0).Add(NewParam(3.0))
	expect(t, expression.Add(expression2), 25.0)
	expect(t, expression.Sub(expression2), 5.0)

	expect(t, expression, 15.0)
}

func TestTerm_Sums(t *testing.T) {
	term := NewTerm(NewVariable(12), 1.0)
	expect(t, term, 12.0)

	// constant
	expect(t, term.Add(CM(2.0)), 14.0)
	expect(t, term.Sub(CM(2.0)), 10.0)

	expect(t, term, 12.0)

	// params
	param := NewParam(2.0)
	expect(t, term.Add(param), 14.0)
	expect(t, term.Sub(param), 10.0)

	expect(t, term, 12.0)

	// terms
	term2 := NewTerm(NewVariable(1.0), 2.0)
	expect(t, term.Add(term2), 14.0)
	expect(t, term.Sub(term2), 10.0)

	expect(t, term, 12.0)

	//term
	expression2 := NewParam(1.0).Add(NewParam(1.0))
	expect(t, term.Add(expression2), 14.0)
	expect(t, term.Sub(expression2), 10.0)

	expect(t, term, 12.0)
}

func TestParam_Sums(t *testing.T) {
	param := NewParam(3.0)
	expect(t, param, 3.0)

	// constant
	expect(t, param.Add(CM(2.0)), 5.0)
	expect(t, param.Sub(CM(2.0)), 1.0)

	expect(t, param, 3.0)

	// params
	param2 := NewParam(2.0)
	expect(t, param.Add(param2), 5.0)
	expect(t, param.Sub(param2), 1.0)

	expect(t, param, 3.0)

	// terms
	term := NewTerm(NewVariable(1.0), 2.0)
	expect(t, param.Add(term), 5.0)
	expect(t, param.Sub(term), 1.0)

	expect(t, param, 3.0)

	//param
	expression2 := NewParam(1.0).Add(NewParam(1.0))
	expect(t, param.Add(expression2), 5.0)
	expect(t, param.Sub(expression2), 1.0)

	expect(t, param, 3.0)
}

func TestConstant_Sums(t *testing.T) {
	constant := CM(3.0)
	expect(t, constant, 3.0)

	// constant
	expect(t, constant.Add(CM(2.0)), 5.0)
	expect(t, constant.Sub(CM(2.0)), 1.0)

	expect(t, constant, 3.0)

	// params
	param := NewParam(2.0)
	expect(t, constant.Add(param), 5.0)
	expect(t, constant.Sub(param), 1.0)

	expect(t, constant, 3.0)

	// terms
	term := NewTerm(NewVariable(1.0), 2.0)
	expect(t, constant.Add(term), 5.0)
	expect(t, constant.Sub(term), 1.0)

	expect(t, constant, 3.0)

	//constant
	expression2 := NewParam(1.0).Add(NewParam(1.0))
	expect(t, constant.Add(expression2), 5.0)
	expect(t, constant.Sub(expression2), 1.0)

	expect(t, constant, 3.0)
}

func TestSimpleMultiplication(t *testing.T) {
	c := CM(20.0)
	res := c.Mult(CM(2.0))
	expect(t, res, 40.0)

	p := NewParam(20.0)
	res = p.Mult(CM(2.0))
	expect(t, res, 40.0)

	term := NewTerm(p.Variable, 1.0)
	res = term.Mult(CM(2.0))
	expect(t, res, 40.0)

	e := NewExpression([]*Term{term}, 0.0)
	res = e.Mult(CM(2.0))
	expect(t, res, 40.0)
}

func TestSimpleDivision(t *testing.T) {
	c := CM(20.0)
	res := c.Div(CM(2.0))
	expect(t, res, 10.0)

	p := NewParam(20.0)
	res = p.Div(CM(2.0))
	expect(t, res, 10.0)

	term := NewTerm(p.Variable, 1.0)
	res = term.Div(CM(2.0))
	expect(t, res, 10.0)

	e := NewExpression([]*Term{term}, 0.0)
	res = e.Div(CM(2.0))
	expect(t, res, 10.0)
}

func TestConstraint(t *testing.T) {
	left := NewParam(2)
	right := NewParam(10)

	c := right.Sub(left).GreaterThanOrEqualTo(CM(20))
	if c.relation != GreaterThanOrEqualTo {
		t.Error("Relation does not match expected one")
	}
	if c.expression.constant != -20.0 {
		t.Error("Constraint expression constant does not match expected value")
	}

	c2 := right.Sub(left).Equals(CM(30))
	if c2.relation != EqualTo {
		t.Error("Relation does not match expected one")
	}
	if c2.expression.constant != -30.0 {
		t.Error("Constraint expression constant does not match expected value")
	}

	c3 := right.Sub(left).LessThanOrEqualTo(CM(30))
	if c3.relation != LessThanOrEqualTo {
		t.Error("Relation does not match expected one")
	}
	if c3.expression.constant != -30.0 {
		t.Error("Constraint expression constant does not match expected value")
	}
}

func TestConstraintComplexSetup(t *testing.T) {
	expression := NewParam(200).Sub(NewParam(100))

	expression2 := NewExpression([]*Term{NewTerm(NewVariable(2.0), 1.0)}, 20.0)

	c := expression.GreaterThanOrEqualTo(expression2)
	if c.relation != GreaterThanOrEqualTo {
		t.Error("Relation does not match expected one")
	}
	if len(c.expression.terms) != 3 {
		t.Error("Constraint expression term count does not match expected one", 3, ", was", len(c.expression.terms))
	}
	if c.expression.constant != -20.0 {
		t.Error("Constraint expression constant does not match expected value")
	}
}

func TestSolver(t *testing.T) {
	s := NewSolver()

	left := NewParam(2)
	right := NewParam(100)

	c1 := right.Sub(left).GreaterThanOrEqualTo(CM(200))

	err := s.AddConstraint(c1)

	if err != nil {
		t.Error(err)
	}
}

func TestSingleVariable(t *testing.T) {
	left := NewParam(-20)

	s := NewSolver()
	s.AddConstraint(left.GreaterThanOrEqualTo(CM(5)))
	s.FlushUpdates()

	fmt.Println(left.Value())
}

func TestMoreComplexSetup(t *testing.T) {
	containerWidth := CM(400)
	childCompWidth := NewParam(100)

	div := containerWidth.Div(CM(2))
	constraint := childCompWidth.Equals(div)

	s := NewSolver()
	err := s.AddConstraint(constraint)

	if err != nil {
		t.Error(err)
	}

	s.FlushUpdates()

	fmt.Println(childCompWidth.Value())
}

func TestMidpoints(t *testing.T) {
	left := NewParam(0)
	left.Variable.Name = "Left"
	right := NewParam(0)
	right.Variable.Name = "Right"
	mid := NewParam(0)
	mid.Variable.Name = "Mid"

	s := NewSolver()

	s.AddConstraint(right.Add(left).Equals(mid.Mult(CM(2.0))))
	s.AddConstraint(right.Sub(left).GreaterThanOrEqualTo(CM(100)))
	s.AddConstraint(left.GreaterThanOrEqualTo(CM(0)))

	s.FlushUpdates()

	if left.Value() != 0 {
		t.Error("Left valze is wrong")
	}
	if mid.Value() != 50 {
		t.Error("Mid value is wrong")
	}
	if right.Value() != 100 {
		t.Error("Right value is wrong")
	}
}

func TestMoreComplexSetup2(t *testing.T) {
	containerWidth := NewVariable(1024)
	containerHeight := NewVariable(768)

	widthTerm := NewTerm(containerWidth, 1)
	heightTerm := NewTerm(containerHeight, 1)

	childX := NewParam(50)
	childY := NewParam(50)
	childCompWidth := NewParam(200)
	childCompHeight := NewParam(100)

	c1 := childX.Equals(widthTerm.Mult(CM(50.0 / 1024.0)))
	c2 := childY.Equals(heightTerm.Mult(CM(50.0 / 768.0)))
	c3 := childCompWidth.Equals(widthTerm.Mult(CM(200.0 / 1024.0)))
	c4 := childCompHeight.Equals(heightTerm.Mult(CM(100.0 / 1024.0)))

	s := NewSolver()
	err := s.AddConstraint(c1)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c2)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c3)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c4)
	if err != nil {
		t.Error(err)
	}

	s.AddEditVariable(containerWidth, float64(PriorityStrong))
	s.SuggestValueForVariable(containerWidth, 2048)

	s.FlushUpdates()

	fmt.Println(containerWidth.Value, childX.Value())
}

func TestMoreComplexSetup3(t *testing.T) {
	containerWidth := NewVariable(1024)
	containerHeight := NewVariable(768)

	widthTerm := NewTerm(containerWidth, 1)
	heightTerm := NewTerm(containerHeight, 1)

	childX := NewParam(50)
	childY := NewParam(50)
	childCompWidth := NewParam(200)
	childCompHeight := NewParam(100)

	c1 := childX.Equals(widthTerm.Mult(CM(50.0 / 1024.0)))
	c2 := childY.Equals(heightTerm.Mult(CM(50.0 / 768.0)))
	c3 := childCompWidth.Equals(widthTerm.Mult(CM(200.0 / 1024.0)))
	c3.Priority = PriorityStrong

	c4 := childCompHeight.Equals(heightTerm.Mult(CM(100.0 / 1024.0)))
	c5 := childCompWidth.GreaterThanOrEqualTo(CM(500))
	c5.Priority = PriorityWeak

	s := NewSolver()
	err := s.AddConstraint(c5)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c1)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c2)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c3)
	if err != nil {
		t.Error(err)
	}
	err = s.AddConstraint(c4)
	if err != nil {
		t.Error(err)
	}

	s.AddEditVariable(containerWidth, float64(PriorityStrong))
	s.SuggestValueForVariable(containerWidth, 2048)

	s.FlushUpdates()

	fmt.Println(containerWidth.Value, childCompWidth.Value())
}

func TestConstraintUpdate(t *testing.T) {
	left := NewParam(2.0)
	right := NewParam(100.0)

	c1 := right.Sub(left).GreaterThanOrEqualTo(CM(200))
	c2 := right.GreaterThanOrEqualTo(left)

	s := NewSolver()
	s.AddConstraint(c1)
	s.AddConstraint(c2)
	s.RemoveConstraint(c1)
}

func TestSolutionWithOptiomize(t *testing.T) {
	p1 := NewParam(0)
	p2 := NewParam(0)
	p3 := NewParam(0)

	container := NewParam(0)

	solver := NewSolver()
	solver.AddEditVariable(container.Variable, float64(PriorityStrong))
	solver.SuggestValueForVariable(container.Variable, 100.0)

	c1 := p1.GreaterThanOrEqualTo(CM(30.0))
	c1.Priority = PriorityStrong

	c2 := p1.Equals(p3)
	c2.Priority = PriorityMedium

	c3 := p2.Equals(CM(2.0).Mult(p1))

	c4 := container.Equals(p1.Add(p2).Add(p3))

	solver.AddConstraint(c1)
	fmt.Println(solver.rows)
	solver.AddConstraint(c2)
	fmt.Println(solver.rows)
	solver.AddConstraint(c3)
	fmt.Println(solver.rows)
	solver.AddConstraint(c4)
	//solver.AddConstraints(c1, c2, c3, c4)
	solver.FlushUpdates()

	if container.Value() != 100.0 {
		t.Error("Container value does not match expected one ", 100, ", was", container.Value())
	}

	if p1.Value() != 30.0 {
		t.Error("P1 value does not match expected one ", 30, ", was", p1.Value())
	}

	if p2.Value() != 60.0 {
		t.Error("P2 value does not match expected one ", 60, ", was", p2.Value())
	}

	if p3.Value() != 10.0 {
		t.Error("P3 value does not match expected one ", 10, ", was", p3.Value())
	}
}
