// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// MoEuler implements the (explicit) Modified Euler method
// Modified-Euler 2(1), order=2, error_est_order=2, nstages=2
type MoEuler struct {
	dat *erkdata
}

// Init initialises structure
func (o *MoEuler) Init(distr bool) (err error) {
	o.dat = &erkdata{
		A: [][]float64{
			{0.0, 0.0},
			{1.0, 0.0},
		},
		B:  []float64{1.0, 0.0},
		Be: []float64{0.5, 0.5},
		C:  []float64{0.0, 1.0},
	}
	return nil
}

// Accept accepts update
func (o *MoEuler) Accept(sol *Solver, y la.Vector) {
	y.Apply(1, sol.w[0]) // y := w (update y)
}

// Step steps update
func (o *MoEuler) Step(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {
	return erkstep(o.dat, 2, true, sol, y0, x0)
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[MoEulerKind] = func() RKmethod { return new(MoEuler) }
}
