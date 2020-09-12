// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"gosl/chk"
	"gosl/io"
	"gosl/la"
	"gosl/utl"
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

// add optimizer to database
func init() {
	nlsMakersDB["graddesc"] = func(prob *Problem) NonLinSolver { return NewGradDesc(prob) }
}

// NewGradDesc returns a new multidimensional optimizer using GradDesc's method (no derivatives required)
func NewGradDesc(prob *Problem) (o *GradDesc) {
	o = new(GradDesc)
	o.InitConvergence(prob.Ffcn, prob.Gfcn)
	o.Alpha = 1e-3
	o.dfdx = la.NewVector(prob.Ndim)
	return
}

// Min solves minimization problem
//
//  Input:
//    x -- [ndim] initial starting point (will be modified)
//    params -- [may be nil] optional parameters. e.g. "alpha", "maxit". Example:
//                 params := utl.NewParams(
//                     &utl.P{N: "alpha", V: 0.5},
//                     &utl.P{N: "maxit", V: 1000},
//                     &utl.P{N: "ftol", V: 1e-2},
//                     &utl.P{N: "gtol", V: 1e-2},
//                     &utl.P{N: "hist", V: 1},
//                     &utl.P{N: "verb", V: 1},
//                 )
//  Output:
//    fmin -- f(x@min) minimum f({x}) found
//    x -- [modify input] position of minimum f({x})
//
func (o *GradDesc) Min(x la.Vector, params utl.Params) (fmin float64) {

	// set parameters
	o.Convergence.SetParams(params)
	o.Alpha = params.GetValueOrDefault("alpha", o.Alpha)
	io.Pforan("α = %v\n", o.Alpha)
	io.Pforan("nit = %v\n", o.MaxIt)
	io.Pforan("ftol = %v\n", o.Convergence.Ftol)

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
