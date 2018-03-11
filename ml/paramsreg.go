// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import "github.com/cpmech/gosl/la"

// ParamsReg holds the θ and b parameters for regression computations
type ParamsReg struct {
	Theta     la.Vector // θ parameter [nFeatures]
	Bias      float64   // bias parameter
	thetaCopy la.Vector // copy of θ
	biasCopy  float64   // copy of b
}

// NewParamsReg returns a new object to hold regression parameters
func NewParamsReg(nFeatures int) (o *ParamsReg) {
	o = new(ParamsReg)
	o.Theta = la.NewVector(nFeatures)
	o.thetaCopy = la.NewVector(nFeatures)
	return o
}

// Backup creates an internal copy of parameters
func (o *ParamsReg) Backup() {
	copy(o.thetaCopy, o.Theta)
	o.biasCopy = o.Bias
}

// Restore restores an internal copy of parameters
func (o *ParamsReg) Restore() {
	copy(o.Theta, o.thetaCopy)
	o.Bias = o.biasCopy
}
