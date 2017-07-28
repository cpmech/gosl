// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// FwEuler implements the (explicit) Forward Euler method
type FwEuler struct {
}

// Init initialises structure
func (o *FwEuler) Init(distr bool) (err error) {
	return nil
}

// Accept accepts update
func (o *FwEuler) Accept(sol *Solver, y la.Vector) {
}

// Nstages returns the number of stages
func (o *FwEuler) Nstages() int {
	return 1
}

// Step steps update
func (o *FwEuler) Step(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {
	sol.Nfeval++
	err = sol.fcn(sol.f[0], sol.h, x0, y0)
	if err != nil {
		return
	}
	for i := 0; i < sol.ndim; i++ {
		y0[i] += sol.h * sol.f[0][i]
	}
	return 1e+20, err // must not be used with automatic substepping
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[FwEulerKind] = func() RKmethod { return new(FwEuler) }
}
