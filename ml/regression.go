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
	Model(xVec, theta la.Vector) float64        // model equation where xVec[1+nFeatures] (augmented vector) and theta[1+nFeatures]
	Cost(data *DataMatrix) float64              // computes cost
	Deriv(dCdTheta la.Vector, data *DataMatrix) // computes dCdθ for given data len(dCdθ) = 1+nFeatures
}
