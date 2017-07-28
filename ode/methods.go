// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// ODE methods
var (
	// FwEuler specifies the Forward Euler method (explicit)
	FwEuler = io.NewEnum("FwEuler", "ode", "FE", "Forward Euler (explicit)")

	// BwEuler specifies the Backward Euler method (explicit)
	BwEuler = io.NewEnum("BwEuler", "ode", "BE", "Backward Euler (explicit)")

	// MoEuler specifies the Modified Euler method (explicit)
	MoEuler = io.NewEnum("MoEuler", "ode", "ME", "Modified Euler (explicit)")

	// DoPri5 specifies the Dormand-Prince5 method (explicit)
	DoPri5 = io.NewEnum("DoPri5", "ode", "DP", "Dormand-Prince5 (explicit)")

	// Radau5 specifies the Radau5 method (implicit)
	Radau5 = io.NewEnum("Radau5", "ode", "R", "Foward Euler (implicit)")
)

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
