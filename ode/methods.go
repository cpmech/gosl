// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// Runge-Kutta methods
var (
	// FwEulerKind specifies the Forward Euler method (explicit)
	FwEulerKind = io.NewEnum("FwEuler", "ode", "FE", "Forward Euler (explicit)")

	// BwEulerKind specifies the Backward Euler method (explicit)
	BwEulerKind = io.NewEnum("BwEuler", "ode", "BE", "Backward Euler (implicit)")

	// MoEulerKind specifies the Modified Euler method (explicit)
	MoEulerKind = io.NewEnum("MoEuler", "ode", "ME", "Modified Euler (explicit)")

	// DoPri5kind specifies the Dormand-Prince5 method (explicit)
	DoPri5kind = io.NewEnum("DoPri5", "ode", "DP", "Dormand-Prince5 (explicit)")

	// Radau5kind specifies the Radau5 method (implicit)
	Radau5kind = io.NewEnum("Radau5", "ode", "R", "Radau5 (implicit)")
)

// rkmethod defines the required functions of Runge-Kutta method
type rkmethod interface {
	Info() (fixedOnly, implicit bool, nstages int)                                        // information
	Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error)           // initialise
	Accept(y la.Vector, work *rkwork)                                                     // accept update
	Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) // step update
}

// rkmMaker defines a function that makes rkmethods
type rkmMaker func() rkmethod

// rkmDB implements a database of rkmethod makers
var rkmDB = make(map[io.Enum]rkmMaker)

// newRKmethod finds a rkmethod in database or panic
func newRKmethod(kind io.Enum) (rkmethod, error) {
	if maker, ok := rkmDB[kind]; ok {
		return maker(), nil
	}
	return nil, chk.Err("cannot find rkmethod named %q in database", kind)
}
