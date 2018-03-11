// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import "github.com/cpmech/gosl/la"

// ParamsReg holds the θ and b parameters for regression computations
type ParamsReg struct {
	Theta      la.Vector // θ parameter [nFeatures]
	Bias       float64   // bias parameter
	Lambda     float64   // regularization parameter
	Degree     int       // degree of polynomial
	thetaCopy  la.Vector // copy of θ
	biasCopy   float64   // copy of b
	lambdaCopy float64   // copy of λ
	degreeCopy int       // copy of degree
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
	o.lambdaCopy = o.Lambda
	o.degreeCopy = o.Degree
}

// Restore restores an internal copy of parameters
func (o *ParamsReg) Restore() {
	copy(o.Theta, o.thetaCopy)
	o.Bias = o.biasCopy
	o.Lambda = o.lambdaCopy
	o.Degree = o.degreeCopy
}
