// Copyright 2016 The Chromium Authors, 2018 Elco Industrie Automation GmbH. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cassowary

type Priority int64

const (
	PriorityRequired Priority = 1000000000
	PriorityStrong   Priority = 1000000
	PriorityMedium   Priority = 1000
	PriorityWeak     Priority = 1
)
