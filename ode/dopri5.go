// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// DoPri5 implements the (explicit) Dormand-Prince 5(4) method
// Dormand-Prince 5(4), order=5, error_est_order=4, nstages=7
type DoPri5 struct {
	fcn Func
	dat *erkdata
}

// Init initialises structure
func (o *DoPri5) Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error) {
	o.fcn = fcn
	o.dat = &erkdata{
		Fprev: true,
		A: [][]float64{{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{1.0 / 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{3.0 / 40.0, 9.0 / 40.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{44.0 / 45.0, -56.0 / 15.0, 32.0 / 9.0, 0.0, 0.0, 0.0, 0.0},
			{19372.0 / 6561.0, -25360.0 / 2187.0, 64448.0 / 6561.0, -212.0 / 729.0, 0.0, 0.0, 0.0},
			{9017.0 / 3168.0, -355.0 / 33.0, 46732.0 / 5247.0, 49.0 / 176.0, -5103.0 / 18656.0, 0.0, 0.0},
			{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0}},
		B:  []float64{35.0 / 384.0, 0.0, 500.0 / 1113.0, 125.0 / 192.0, -2187.0 / 6784.0, 11.0 / 84.0, 0.0},
		Be: []float64{5179.0 / 57600.0, 0.0, 7571.0 / 16695.0, 393.0 / 640.0, -92097.0 / 339200.0, 187.0 / 2100.0, 1.0 / 40.0},
		C:  []float64{0.0, 1.0 / 5.0, 3.0 / 10.0, 4.0 / 5.0, 8.0 / 9.0, 1.0, 1.0},
		w:  la.NewVector(ndim),
	}
	return nil
}

// Info returns information about this method
func (o *DoPri5) Info() (fixedOnly, implicit bool, nstages int) {
	return false, false, 7
}

// Accept accepts update
func (o *DoPri5) Accept(y la.Vector, work *rkwork) {
	y.Apply(1, o.dat.w)
}

// Step steps update
func (o *DoPri5) Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) {
	return o.dat.step(h, x0, y0, stat, work, o.fcn)
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[DoPri5kind] = func() rkmethod { return new(DoPri5) }
}
