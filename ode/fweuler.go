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
	ndim int     // problem dimension
	conf *Config // configurations
	work *rkwork // workspace
	stat *Stat   // statistics
	fcn  Func    // dy/dx := f(x,y)
}

// add method to database
func init() {
	rkmDB["fweuler"] = func() rkmethod { return new(FwEuler) }
}

// Free releases memory
func (o *FwEuler) Free() {}

// Info returns information about this method
func (o *FwEuler) Info() (fixedOnly, implicit bool, nstages int) {
	return true, false, 1
}

// Init initializes structure
func (o *FwEuler) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) {
	if M != nil {
		chk.Panic("Forward-Euler solver cannot handle M matrix yet\n")
	}
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn
}

// Accept accepts update
func (o *FwEuler) Accept(y0 la.Vector, x0 float64) (dxnew float64) {
	return
}

// Reject processes step rejection
func (o *FwEuler) Reject() (dxnew float64) {
	return
}

// DenseOut produces dense output (after Accept)
func (o *FwEuler) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *FwEuler) Step(x0 float64, y0 la.Vector) {
	o.stat.Nfeval++
	o.fcn(o.work.f[0], o.work.h, x0, y0)
	for i := 0; i < o.ndim; i++ {
		y0[i] += o.work.h * o.work.f[0][i]
	}
}
