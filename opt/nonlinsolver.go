// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "github.com/cpmech/gosl/la"

// NonLinSolver solves (unconstrained) nonlinear optimization problems
type NonLinSolver interface {
	Min(x la.Vector, args map[string]float64) (fmin float64) // computes minimum and updates x @ min
}

// NewNonLinSolver returns new object
//  method --   "conjgrad"
//              "powell"
//              "graddesc"
func NewNonLinSolver(method string, prob *Problem) (o *NonLinSolver) {
	o = new(NonLinSolver)
	//switch method{
	//case "conjgrad":
	//return NewConjGrad(ndim, prob.Ffcn, prob
	//}
	return
}
