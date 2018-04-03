// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
)

// GradDesc implements a simple gradient-descent optimizer
//  NOTE: Check Convergence to see how to set convergence parameters,
//        max iteration number, or to enable and access history of iterations
type GradDesc struct {

	// merge properties
	Convergence // auxiliary object to check convergence

	// configuration
	Alpha float64 // rate to take descents

	// internal
	dfdx la.Vector // gradient vector
}

// NewGradDesc returns a new multidimensional optimizer using GradDesc's method (no derivatives required)
//   Ndim -- length(x)
//   Ffcn -- objective function: y = f({x})
//   Gfcn -- gradient function: g = dy/d{x} = deriv(f({x}), {x})
func NewGradDesc(ndim int, Ffcn fun.Sv, Gfcn fun.Vv) (o *GradDesc) {
	o = new(GradDesc)
	o.InitConvergence(Ffcn, Gfcn)
	o.Alpha = 1e-3
	o.dfdx = la.NewVector(ndim)
	return
}

// Min solves minimization problem
//
//  Input:
//    x -- [size] initial starting point (will be modified)
//
//  Output:
//    fmin -- f(x@min) minimum f({x}) found
//    x -- [modify input] position of minimum f({x})
//
func (o *GradDesc) Min(x la.Vector) (fmin float64) {

	// initializations
	o.NumFeval, o.NumGeval = 0, 0
	fmin = o.Ffcn(x)
	fprev := fmin

	// history
	if o.UseHist {
		o.InitHist(x)
	}

	// iterations
	for o.NumIter = 0; o.NumIter < o.MaxIt; o.NumIter++ {

		// compute and check gradient
		o.Gfcn(o.dfdx, x)
		if o.Gconvergence(fprev, x, o.dfdx) {
			return
		}

		// perform descent
		la.VecAdd(x, 1, x, -o.Alpha, o.dfdx) // x := x - α⋅dfdx

		// compute objective function
		fprev = fmin
		fmin = o.Ffcn(x)

		// history
		if o.UseHist {
			o.uhist.Apply(-o.Alpha, o.dfdx)
			o.Hist.Append(fmin, x, o.uhist)
		}

		// compute and check objective function
		if o.Fconvergence(fprev, fmin) {
			return
		}
	}

	// did not converge
	chk.Panic("fail to converge after %d iterations\n", o.NumIter)
	return
}
