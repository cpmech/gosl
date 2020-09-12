// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "gosl/io"

// Stat holds statistics and output data
type Stat struct {
	Nfeval    int     // number of calls to fcn
	Njeval    int     // number of Jacobian matrix evaluations
	Nsteps    int     // total number of substeps
	Naccepted int     // number of accepted substeps
	Nrejected int     // number of rejected substeps
	Ndecomp   int     // number of matrix decompositions
	Nlinsol   int     // number of calls to linsolver
	Nitmax    int     // number max of iterations
	Hopt      float64 // optimal step size at the end
	LsKind    string  // kind of linear solver used
	Implicit  bool    // method is implicit
}

// NewStat returns a new structure
func NewStat(lskind string, implicit bool) (o *Stat) {
	o = new(Stat)
	o.LsKind = lskind
	o.Implicit = implicit
	return
}

// Reset initialises Stat
func (o *Stat) Reset() {
	o.Nfeval = 0
	o.Njeval = 0
	o.Nsteps = 0
	o.Naccepted = 0
	o.Nrejected = 0
	o.Ndecomp = 0
	o.Nlinsol = 0
	o.Nitmax = 0
}

// Print prints information about the solution process
func (o *Stat) Print(extra bool) {
	io.Pf("number of F evaluations   =%6d\n", o.Nfeval)
	if o.Implicit {
		io.Pf("number of J evaluations   =%6d\n", o.Njeval)
	}
	io.Pf("total number of steps     =%6d\n", o.Nsteps)
	io.Pf("number of accepted steps  =%6d\n", o.Naccepted)
	io.Pf("number of rejected steps  =%6d\n", o.Nrejected)
	if o.Implicit {
		io.Pf("number of decompositions  =%6d\n", o.Ndecomp)
		io.Pf("number of lin solutions   =%6d\n", o.Nlinsol)
		io.Pf("max number of iterations  =%6d\n", o.Nitmax)
	}
	if extra {
		io.Pf("optimal step size Hopt    = %g\n", o.Hopt)
		io.Pf("kind of linear solver     = %q\n", o.LsKind)
	}
}
