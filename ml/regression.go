// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import "github.com/cpmech/gosl/la"

// Regression defines an interface for regression models
type Regression interface {
	BackupParams()
	SetParam(i int, value float64)
	GetParam(i int) (value float64)
	RestoreParams()
	Predict(x la.Vector) (y float64)
	Cost() (c float64)
}
