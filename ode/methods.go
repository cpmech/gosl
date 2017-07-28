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

// RKmethod defines the required functions of Runge-Kutta method
type RKmethod interface {
	Init(distr bool) (err error)                                        // initialise
	Accept(o *Solver, y la.Vector)                                      // accept update
	Step(o *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) // step update
}

// rkmMaker defines a function that makes RKmethods
type rkmMaker func() RKmethod

// rkmDB implements a database of RKmethod makers
var rkmDB = make(map[io.Enum]rkmMaker)

// NewRKmethod finds a RKmethod in database or panic
func NewRKmethod(kind io.Enum) RKmethod {
	if maker, ok := rkmDB[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find RKmethod named %q in database", kind)
	return nil
}

/*
func getMethod(kind io.Enum, distr bool) (name string, step stpfcn, accept acptfcn, nstg int, erkdat expRKdat) {
	switch kind {
	case FwEuler:
		step = fweulerStep
		accept = fweulerAccept
		nstg = 1
	case BwEuler:
		step = bweulerStep
		accept = bweulerAccept
		nstg = 1
	case MoEuler:
		step = erkStep
		accept = erkAccept
		nstg = 2
		erkdat = expRKdat{true, ME2_a, ME2_b, ME2_be, ME2_c}
	case DoPri5:
		step = erkStep
		accept = erkAccept
		nstg = 7
		erkdat = expRKdat{true, DP5_a, DP5_b, DP5_be, DP5_c}
	case Radau5:
		if distr {
			step = radau5_step_mpi
		} else {
			step = radau5_step
		}
		accept = radau5_accept
		nstg = 3
	default:
		chk.Panic("method %s is not available", kind)
	}
	return
}
*/
