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

	// constants
	Fprev bool        // can use previous f
	A     [][]float64 // a coefficients
	B     []float64   // b coefficients
	Be    []float64   // be coefficients
	C     []float64   // c coefficients
	Nstg  int         // number of stages
	q     int         // order of error estimator; e.g. DoPri5(4) ⇒ q = 4 (=min(order(y1),order(y1bar))

	// data
	ndim int       // problem dimension
	conf *Config   // configuration
	work *rkwork   // workspace
	stat *Stat     // statistics
	fcn  Func      // dy/dx = f(x,y) function
	w    la.Vector // local workspace
	beta float64   // factor to stabilize step
	n    float64   // exponent n = 1/(q+1) (or 1/(q+1)-0.75⋅β) of rerrⁿ
	dmin float64   // dmin = 1/Mmin
	dmax float64   // dmax = 1/Mmax
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
func (o *ExplicitRK) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) (err error) {
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn
	o.w = la.NewVector(o.ndim)
	if conf.Method == "dopri5" {
		o.beta = conf.DP5beta
	}
	o.n = 1.0/float64(o.q+1) - o.beta*0.75
	o.dmin = 1.0 / o.conf.Mmin
	o.dmax = 1.0 / o.conf.Mmax
	return nil
}

// Accept accepts update and computes next stepsize
func (o *ExplicitRK) Accept(y la.Vector) (dxnew float64) {

	// update y and k0
	y.Apply(1, o.w)
	o.work.f[0].Apply(1, o.work.f[o.Nstg-1]) // k0 := ks for next step

	// estimate new stepsize
	d := math.Pow(o.work.rerr, o.n)
	if o.beta > 0 { // lund-stabilization
		d = d / math.Pow(o.work.rerrPrev, o.beta)
	}
	d = utl.Max(o.dmax, utl.Min(o.dmin, d/o.conf.Mfac)) // we require  fac1 <= hnew/h <= fac2
	dxnew = o.work.h / d
	return
}

// Reject processes step rejection and computes next stepsize
func (o *ExplicitRK) Reject() (dxnew float64) {

	// estimate new stepsize
	d := math.Pow(o.work.rerr, o.n) / o.conf.Mfac
	dxnew = o.work.h / utl.Min(o.dmin, d)
	return
}

// ContOut produces continuous output (after Accept)
func (o *ExplicitRK) ContOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *ExplicitRK) Step(xa float64, ya la.Vector) (err error) {

	// auxiliary
	h := o.work.h
	k := o.work.f
	v := o.work.v

	// compute k0 (otherwise, use k0 saved in Accept)
	if o.work.first || !o.Fprev { // do it also if cannot reuse previous ks
		u0 := xa + h*o.C[0]
		o.stat.Nfeval++
		err = o.fcn(k[0], h, u0, ya) // k0 := f(ui,vi)
		if err != nil {
			return
		}
	}

	// compute ki
	var ui float64
	for i := 1; i < o.work.nstg; i++ {
		ui = xa + h*o.C[i]
		v[i].Apply(1, ya)        // vi := ya
		for j := 0; j < i; j++ { // lower diagonal ⇒ explicit
			la.VecAdd(v[i], 1, v[i], h*o.A[i][j], k[j]) // vi += h⋅aij⋅kj
		}
		o.stat.Nfeval++
		err = o.fcn(k[i], h, ui, v[i]) // ki := f(ui,vi)
		if err != nil {
			return
		}
	}

	// error estimation
	var kh, sum, lerrm, ratio float64 // m component of local error estimate
	for m := 0; m < o.ndim; m++ {
		o.w[m] = ya[m]
		lerrm = 0.0 // must be zeroed for each m
		for i := 0; i < o.Nstg; i++ {
			kh = k[i][m] * h
			o.w[m] += o.B[i] * kh
			lerrm += (o.Be[i] - o.B[i]) * kh
		}
		sk := o.conf.atol + o.conf.rtol*utl.Max(math.Abs(ya[m]), math.Abs(o.w[m]))
		ratio = lerrm / sk
		sum += ratio * ratio
	}
	o.work.rerr = utl.Max(math.Sqrt(sum/float64(o.ndim)), 1.0e-10)
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
		o.q = 1
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
		o.q = 4
	default:
		return nil
	}
	return o
}
