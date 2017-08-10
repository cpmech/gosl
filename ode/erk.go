// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// ExplicitRK implements explicit Runge-Kutta methods
//   The methods available are:
//     moeuler -- Modified-Euler 2(1), order=2, nstages=2, error_est_order=2
//     dopri5  -- Dormand-Prince 5(4), order=5, nstages=7, error_est_order=4
type ExplicitRK struct {
	Fprev bool        // can use previous f
	A     [][]float64 // a coefficients
	B     []float64   // b coefficients
	Be    []float64   // be coefficients
	C     []float64   // c coefficients
	Nstg  int         // number of stages
	w     la.Vector   // step update (normalised variable starting from zero)
	fcn   Func        // dy/dx = f(x,y) function
}

// add methods to database
func init() {
	rkmDB["moeuler"] = func() rkmethod { return newERK("moeuler") }
	rkmDB["dopri5"] = func() rkmethod { return newERK("dopri5") }
}

// Free releases memory
func (o *ExplicitRK) Free() {}

// Info returns information about this method
func (o *ExplicitRK) Info() (fixedOnly, implicit bool, nstages int) {
	return false, false, o.Nstg
}

// Init initialises structure
func (o *ExplicitRK) Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error) {
	o.fcn = fcn
	o.w = la.NewVector(ndim)
	return nil
}

// Accept accepts update
func (o *ExplicitRK) Accept(y la.Vector, work *rkwork) {
	y.Apply(1, o.w)
}

// ContOut produces continuous output (after Accept)
func (o *ExplicitRK) ContOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *ExplicitRK) Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) {

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
			err = o.fcn(work.f[i], h, work.u[i], work.v[i])
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

func newERK(kind string) rkmethod {
	o := new(ExplicitRK)
	o.Fprev = true
	switch kind {
	case "moeuler":
		o.A = [][]float64{
			{0.0, 0.0},
			{1.0, 0.0},
		}
		o.B = []float64{1.0, 0.0}
		o.Be = []float64{0.5, 0.5}
		o.C = []float64{0.0, 1.0}
		o.Nstg = 2
	case "dopri5":
		o.A = [][]float64{
			{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{3.0 / 40.0, 9.0 / 40.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{44.0 / 45.0, -56.0 / 15.0, 32.0 / 9.0, 0.0, 0.0, 0.0, 0.0},
			{19372.0 / 6561.0, -25360.0 / 2187.0, 64448.0 / 6561.0, -212.0 / 729.0, 0.0, 0.0, 0.0},
			{9017.0 / 3168.0, -355.0 / 33.0, 46732.0 / 5247.0, 49.0 / 176.0, -5103.0 / 18656.0, 0.0, 0.0},
			{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0},
		}
		o.B = []float64{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0}
		o.Be = []float64{5179.0 / 57600.0, 0.0, 7571.0 / 16695.0, 393.0 / 640.0, -92097.0 / 339200.0, 187.0 / 2100.0, 1.0 / 40.0}
		o.C = []float64{0.0, 1.0 / 5.0, 3.0 / 10.0, 4.0 / 5.0, 8.0 / 9.0, 1.0, 1.0}
		o.Nstg = 7
	default:
		return nil
	}
	return o
}
