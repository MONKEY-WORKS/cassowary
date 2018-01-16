// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

type Variable struct {
	Value float64

	Name string

	owner *Param
}

func NewVariable(val float64) *Variable {
	return &Variable{
		Value: val,
	}
}

func (v *Variable) Param() *Param {
	return v.owner
}

func (v *Variable) applyUpdate(updated float64) bool {
	res := updated != v.Value
	v.Value = updated
	return res
}
