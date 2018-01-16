// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

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

	Mult(m EquationMember) *Expression

	Div(m EquationMember) *Expression
}
