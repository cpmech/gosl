// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// FwEuler implements the (explicit) Forward Euler method
type FwEuler struct {
	fcn Func
}

// Info returns information about this method
func (o *FwEuler) Info() (fixedOnly, implicit bool, nstages int) {
	return true, false, 1
}

// Init initialises structure
func (o *FwEuler) Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error) {
	if M != nil {
		err = chk.Err("Forward-Euler solver cannot handle M matrix yet\n")
		return
	}
	o.fcn = fcn
	return
}

// Accept accepts update
func (o *FwEuler) Accept(y la.Vector, work *rkwork) {
}

// Step steps update
func (o *FwEuler) Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) {
	stat.Nfeval++
	err = o.fcn(work.f[0], h, x0, y0)
	if err != nil {
		return
	}
	for i := 0; i < work.ndim; i++ {
		y0[i] += h * work.f[0][i]
	}
	return 1e+20, err // must not be used with automatic substepping
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[FwEulerKind] = func() rkmethod { return new(FwEuler) }
}
