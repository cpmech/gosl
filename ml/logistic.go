// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/la"
)

// Logistic implements functions to perform the Logistic Regression
type Logistic struct {

	// input
	Ndata   int        // number of data points = m
	Nparams int        // number of parameters = n = len(θ)
	X       *la.Matrix // x-data [Ndata][Nparams]
	Y       la.Vector  // y-data [Ndata]
	Theta   la.Vector  // parameters θ [Nparams]

	// workspace
	logP la.Vector // log(p[k]) with p[k] = 1 + exp(-θ[k]⋅x[k]) [Ndata]
	ybar la.Vector // ybar = 1 - y [Ndata]
	tmp  la.Vector // temporary vector [Ndata]
}

// NewLogistic returns a new structure
func NewLogistic(ndata, nparams int) (o *Logistic) {
	o = new(Logistic)
	o.Ndata = ndata
	o.Nparams = nparams
	o.X = la.NewMatrix(o.Ndata, o.Nparams)
	o.Y = la.NewVector(o.Nparams)
	o.Theta = la.NewVector(o.Nparams)
	o.logP = la.NewVector(o.Nparams)
	o.ybar = la.NewVector(o.Nparams)
	o.tmp = la.NewVector(o.Ndata)
	return
}

// Start restarts internal auxiliary arrays

// Cost computes the total cost
func (o *Logistic) Cost() {
	la.MatVecMul(o.tmp, 1, o.X, o.Theta)
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// p computes g(θ[k]⋅x[k])
func (o *Logistic) h(x la.Vector) float64 {
	return 0
}

// g computes the (vectorial) logistic function
//   Input:
//     z -- vector z
//   Output:
//     l -- vector l such that l[k] = [1 + exp(-z[k])]⁻¹
func (o *Logistic) g(l, z la.Vector) {
	for k := 0; k < len(z); k++ {
		l[k] = 1.0 / (1.0 + math.Exp(-z[k]))
	}
}
