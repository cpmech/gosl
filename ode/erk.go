// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/la"
)

// erkdata holds data for an explicit Runge-Kutta method
type erkdata struct {
	A  [][]float64 // a coefficients
	B  []float64   // b coefficients
	Be []float64   // be coefficients
	C  []float64   // c coefficients
}

// erkstep performs the step update of an explicit Runge-Kutta method
func erkstep(o *erkdata, nStages int, useFprev bool, sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	// update
	for i := 0; i < nStages; i++ {
		sol.u[i] = x0 + sol.h*o.C[i]
		sol.v[i].Apply(1, y0) // v[i] := y
		for j := 0; j < i; j++ {
			la.VecAdd(sol.v[i], 1, sol.v[i], sol.h*o.A[i][j], sol.f[j]) // v[i] += h*a[i][j]*f[j]
		}
		if i == 0 && useFprev && !sol.first {
			sol.f[i].Apply(1, sol.f[nStages-1]) // f[i] := f[nstg-1]
		} else {
			sol.Nfeval++
			err = sol.fcn(sol.f[i], sol.h, sol.u[i], sol.v[i])
			if err != nil {
				return
			}
		}
	}

	// error estimation
	var lerrm float64 // m component of local error estimate
	for m := 0; m < sol.ndim; m++ {
		lerrm = 0.0
		for i := 0; i < nStages; i++ {
			sol.w[0][m] += o.B[i] * sol.f[i][m] * sol.h
			lerrm += (o.Be[i] - o.B[i]) * sol.f[i][m] * sol.h
		}
		rerr += math.Pow(lerrm/sol.scal[m], 2.0)
	}
	rerr = max(math.Sqrt(rerr/float64(sol.ndim)), 1.0e-10)
	return
}

// expRKdat holds data to solve an ODE using the explicit Runge-Kutta method
type expRKdat struct {
	usefp bool        // method can use f from previous step
	a     [][]float64 // a coefficients
	b     []float64   // b coefficients
	be    []float64   // be coefficients
	c     []float64   // c coefficients
}

// erkAccept accepts update
func erkAccept(sol *Solver, y la.Vector) {
	y.Apply(1, sol.w[0]) // y := w (update y)
}

// erkStep performs one step update using the (explicit) Runge-Kutta method
func erkStep(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	for i := 0; i < sol.nstg; i++ {
		sol.u[i] = x0 + sol.h*sol.erkdat.c[i]
		sol.v[i].Apply(1, y0) // v[i] := y
		for j := 0; j < i; j++ {
			la.VecAdd(sol.v[i], 1, sol.v[i], sol.h*sol.erkdat.a[i][j], sol.f[j]) // v[i] += h*a[i][j]*f[j]
		}
		if i == 0 && sol.erkdat.usefp && !sol.first {
			sol.f[i].Apply(1, sol.f[sol.nstg-1]) // f[i] := f[nstg-1]
		} else {
			sol.Nfeval++
			err = sol.fcn(sol.f[i], sol.h, sol.u[i], sol.v[i])
			if err != nil {
				return
			}
		}
	}

	var lerrm float64 // m component of local error estimate
	for m := 0; m < sol.ndim; m++ {
		lerrm = 0.0
		for i := 0; i < sol.nstg; i++ {
			sol.w[0][m] += sol.erkdat.b[i] * sol.f[i][m] * sol.h
			lerrm += (sol.erkdat.be[i] - sol.erkdat.b[i]) * sol.f[i][m] * sol.h
		}
		rerr += math.Pow(lerrm/sol.scal[m], 2.0)
	}
	rerr = max(math.Sqrt(rerr/float64(sol.ndim)), 1.0e-10)

	return
}

// constants
var (
	// Modified-Euler 2(1), order=2, error_est_order=2, nstages=2
	ME2_a = [][]float64{{0.0, 0.0},
		{1.0, 0.0}}
	ME2_b  = []float64{1.0, 0.0}
	ME2_be = []float64{0.5, 0.5}
	ME2_c  = []float64{0.0, 1.0}

	// Dormand-Prince 5(4), order=5, error_est_order=4, nstages=7
	DP5_a = [][]float64{{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{1.0 / 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{3.0 / 40.0, 9.0 / 40.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{44.0 / 45.0, -56.0 / 15.0, 32.0 / 9.0, 0.0, 0.0, 0.0, 0.0},
		{19372.0 / 6561.0, -25360.0 / 2187.0, 64448.0 / 6561.0, -212.0 / 729.0, 0.0, 0.0, 0.0},
		{9017.0 / 3168.0, -355.0 / 33.0, 46732.0 / 5247.0, 49.0 / 176.0, -5103.0 / 18656.0, 0.0, 0.0},
		{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0}}
	DP5_b  = []float64{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0}
	DP5_be = []float64{5179.0 / 57600.0, 0.0, 7571.0 / 16695.0, 393.0 / 640.0, -92097.0 / 339200.0, 187.0 / 2100.0, 1.0 / 40.0}
	DP5_c  = []float64{0.0, 1.0 / 5.0, 3.0 / 10.0, 4.0 / 5.0, 8.0 / 9.0, 1.0, 1.0}
)
