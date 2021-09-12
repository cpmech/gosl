// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// rkmethod defines the required functions of Runge-Kutta method
type rkmethod interface {
	Free()                                                                                    // free memory
	Info() (fixedOnly, implicit bool, nstages int)                                            // information
	Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) // initialize
	Accept(y0 la.Vector, x0 float64) (dxnew float64)                                          // accept update (must compute rerr)
	Reject() (dxnew float64)                                                                  // process step rejection (must compute rerr)
	DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64)                         // dense output (after Accept)
	Step(x0 float64, y0 la.Vector)                                                            // step update
}

// rkmMaker defines a function that makes rkmethods
type rkmMaker func() rkmethod

// rkmDB implements a database of rkmethod makers
var rkmDB = make(map[string]rkmMaker)

// newRKmethod finds a rkmethod in database or panic
func newRKmethod(kind string) rkmethod {
	if maker, ok := rkmDB[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find rkmethod named %q in database\n", kind)
	return nil
}
