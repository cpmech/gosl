// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// Solve solves ODE problem using standard parameters
//
//  INPUT:
//   method    -- the method
//   fcn       -- function d{y}/dx := {f}(h=dx, x, {y})
//   jac       -- Jacobian function d{f}/d{y} := [J](h=dx, x, {y}) [may be nil]
//   y         -- current {y} @ x=0
//   xf        -- final x
//   dx        -- stepsize. [may be used for dense output]
//   atol      -- absolute tolerance; use 0 for default [default = 1e-4] (for fixedStp=false)
//   rtol      -- relative tolerance; use 0 for default [default = 1e-4] (for fixedStp=false)
//   numJac    -- use numerical Jacobian if if jac is non nil
//   fixedStep -- fixed steps
//   saveStep  -- save steps
//   saveDense -- save many steps (dense output) [using dx]
//
//  OUTPUT:
//   y    -- modified y vector with final {y}
//   stat -- statistics
//   out  -- output with all steps results with save==true
//
func Solve(method string, fcn Func, jac JacF, y la.Vector, xf, dx, atol, rtol float64,
	numJac, fixedStep, saveStep, saveDense bool) (stat *Stat, out *Output) {

	// configuration
	conf := NewConfig(method, "")
	if atol > 0 && rtol > 0 {
		conf.SetTols(atol, rtol)
	}
	if fixedStep {
		conf.SetFixedH(dx, xf)
	}
	if saveStep {
		conf.SetStepOut(true, nil)
	}
	if saveDense {
		conf.SetDenseOut(true, dx, xf, nil)
	}

	// allocate solver
	J := jac
	if numJac {
		J = nil
	}
	sol := NewSolver(len(y), conf, fcn, J, nil)
	defer sol.Free()

	// solve ODE
	sol.Solve(y, 0.0, xf)

	// results
	stat = sol.Stat
	out = sol.Out
	return
}

// Dopri5simple solves ODE using DoPri5 method without options for saving results and others
func Dopri5simple(fcn Func, y la.Vector, xf, tol float64) {
	Solve("dopri5", fcn, nil, y, xf, 0, tol, tol, false, false, false, false)
}

// Dopri8simple solves ODE using DoPri8 method without options for saving results and others
func Dopri8simple(fcn Func, y la.Vector, xf, tol float64) {
	Solve("dopri8", fcn, nil, y, xf, 0, tol, tol, false, false, false, false)
}

// Radau5simple solves ODE using Radau5 method without options for saving results and others
func Radau5simple(fcn Func, jac JacF, y la.Vector, xf, tol float64) {
	numJac := jac == nil
	Solve("radau5", fcn, jac, y, xf, 0, tol, tol, numJac, false, false, false)
}
