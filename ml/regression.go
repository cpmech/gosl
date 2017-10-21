// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/la"
)

// RegData holds data for regression computations
type RegData struct {

	// input
	m int        // ndata: number of data points
	n int        // nparams: number of parameters = len(θ) = number of features + 1
	x *la.Matrix // x-data [ndata][nparams]
	y la.Vector  // y-data [ndata]
	θ la.Vector  // parameters θ [nparams]

	// workspace
	p  la.Vector // 1 + exp(-θ[k]⋅x[k])
	yb la.Vector // 1 - y
}

// NewRegData returns a new structure to hold Regression Data
func NewRegData(nData, nFeatures int) (o *RegData) {
	o = new(RegData)
	o.m = nData
	o.n = nFeatures + 1
	o.x = la.NewMatrix(o.m, o.n)
	o.y = la.NewVector(o.n)
	o.θ = la.NewVector(o.n)
	return
}

// SetX sets X value
func (o *RegData) SetX(iData, jFeature int, value float64) {
	o.x.Set(iData, 1+jFeature)
}
