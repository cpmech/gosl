// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/la"
)

type ERKdat struct {
	usefp bool        // method can use f from previous step
	a     [][]float64 // a coefficients
	b     []float64   // b coefficients
	be    []float64   // be coefficients
	c     []float64   // c coefficients
}

func erk_accept(o *Solver, y []float64) {
	la.VecCopy(y, 1, o.w[0]) // update y
}

// explicit Runge-Kutta step function
func erk_step(o *Solver, y []float64, x float64, args ...interface{}) (rerr float64, err error) {

	for i := 0; i < o.nstg; i++ {
		o.u[i] = x + o.h*o.erkdat.c[i]
		la.VecCopy(o.v[i], 1, y)
		for j := 0; j < i; j++ {
			la.VecAdd(o.v[i], o.h*o.erkdat.a[i][j], o.f[j])
		}
		if i == 0 && o.erkdat.usefp && !o.first {
			la.VecCopy(o.f[i], 1, o.f[o.nstg-1])
		} else {
			o.Nfeval += 1
			err = o.fcn(o.f[i], o.h, o.u[i], o.v[i], args...)
			if err != nil {
				return
			}
		}
	}

	var lerrm float64 // m component of local error estimate
	for m := 0; m < o.ndim; m++ {
		lerrm = 0.0
		for i := 0; i < o.nstg; i++ {
			o.w[0][m] += o.erkdat.b[i] * o.f[i][m] * o.h
			lerrm += (o.erkdat.be[i] - o.erkdat.b[i]) * o.f[i][m] * o.h
		}
		rerr += math.Pow(lerrm/o.scal[m], 2.0)
	}
	rerr = max(math.Sqrt(rerr/float64(o.ndim)), 1.0e-10)

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
