// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ml implements functions to develop Machine Learning algorithms
package ml

import (
	"github.com/cpmech/gosl/la"
)

// Regression defines the functions required to perform regression computations
type Regression interface {
	GetParams() (θ la.Vector, b float64)    // return a copy of parameters
	SetParams(θ la.Vector, b float64)       // set parameters
	SetTheta(iFeature int, value float64)   // set θ parameter
	SetBias(value float64)                  // set b parameter
	Model(x la.Vector) float64              // model equation. return y(x;θ,b)
	Cost(data *DataMatrix) float64          // computes cost
	Deriv(dCdθ la.Vector, data *DataMatrix) // computes dCdθ
}
