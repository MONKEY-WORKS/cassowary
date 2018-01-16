// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

type Relation int

const (
	EqualTo Relation = iota
	LessThanOrEqualTo
	GreaterThanOrEqualTo
)

type Constraint struct {
	relation Relation

	expression *Expression

	priority Priority
}

func NewConstraint(exp *Expression, rel Relation) *Constraint {
	return &Constraint{
		relation:   rel,
		expression: exp,
		priority:   PriorityRequired,
	}
}
