// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// erkdata holds data for an explicit Runge-Kutta method
type erkdata struct {
	Fprev bool        // can use previous f
	A     [][]float64 // a coefficients
	B     []float64   // b coefficients
	Be    []float64   // be coefficients
	C     []float64   // c coefficients
	w     la.Vector   // step update (normalised variable starting from zero)
}

// step performs the step update of an explicit Runge-Kutta method
func (o *erkdata) step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork, fcn Func) (rerr float64, err error) {

	// update
	for i := 0; i < work.nstg; i++ {
		work.u[i] = x0 + h*o.C[i]
		work.v[i].Apply(1, y0) // v[i] := y
		for j := 0; j < i; j++ {
			la.VecAdd(work.v[i], 1, work.v[i], h*o.A[i][j], work.f[j]) // v[i] += h*a[i][j]*f[j]
		}
		if i == 0 && o.Fprev && !work.first {
			work.f[i].Apply(1, work.f[work.nstg-1]) // f[i] := f[nstg-1]
		} else {
			stat.Nfeval++
			err = fcn(work.f[i], h, work.u[i], work.v[i])
			if err != nil {
				return
			}
		}
	}

	// error estimation
	var lerrm, ratio float64 // m component of local error estimate
	for m := 0; m < work.ndim; m++ {
		lerrm = 0.0
		for i := 0; i < work.nstg; i++ {
			o.w[m] += o.B[i] * work.f[i][m] * h
			lerrm += (o.Be[i] - o.B[i]) * work.f[i][m] * h
		}
		ratio = lerrm / work.scal[m]
		rerr += ratio * ratio
	}
	rerr = utl.Max(math.Sqrt(rerr/float64(work.ndim)), 1.0e-10)
	return
}
