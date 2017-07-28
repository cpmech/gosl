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
