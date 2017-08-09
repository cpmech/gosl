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
	Free()                                                                                // free memory
	Info() (fixedOnly, implicit bool, nstages int)                                        // information
	Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error)           // initialise
	Accept(y la.Vector, work *rkwork)                                                     // accept update
	ContOut(yOut, y la.Vector, xOut, x, h float64)                                        // continuous output (after Accept)
	Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) // step update
}

// rkmMaker defines a function that makes rkmethods
type rkmMaker func() rkmethod

// rkmDB implements a database of rkmethod makers
var rkmDB = make(map[string]rkmMaker)

// newRKmethod finds a rkmethod in database or panic
func newRKmethod(kind string) (rkmethod, error) {
	if maker, ok := rkmDB[kind]; ok {
		return maker(), nil
	}
	return nil, chk.Err("cannot find rkmethod named %q in database", kind)
}
