// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// MoEuler implements the (explicit) Modified Euler method
// Modified-Euler 2(1), order=2, error_est_order=2, nstages=2
type MoEuler struct {
	fcn Func
	dat *erkdata
}

// Free releases memory
func (o *MoEuler) Free() {}

// Info returns information about this method
func (o *MoEuler) Info() (fixedOnly, implicit bool, nstages int) {
	return false, false, 2
}

// Init initialises structure
func (o *MoEuler) Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error) {
	o.fcn = fcn
	o.dat = &erkdata{
		Fprev: true,
		A: [][]float64{
			{0.0, 0.0},
			{1.0, 0.0},
		},
		B:  []float64{1.0, 0.0},
		Be: []float64{0.5, 0.5},
		C:  []float64{0.0, 1.0},
		w:  la.NewVector(ndim),
	}
	return nil
}

// Accept accepts update
func (o *MoEuler) Accept(y la.Vector, work *rkwork) {
	y.Apply(1, o.dat.w)
}

// Step steps update
func (o *MoEuler) Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) {
	return o.dat.step(h, x0, y0, stat, work, o.fcn)
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[MoEulerKind] = func() rkmethod { return new(MoEuler) }
}
