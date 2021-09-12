// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"time"

	"github.com/cpmech/gosl/io"
)

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

	// benchmark
	NanosecondsStep       int64 // maximum time elapsed during steps [nanoseconds]
	NanosecondsJeval      int64 // maximum time elapsed during Jacobian evaluation [nanoseconds]
	NanosecondsIniSol     int64 // maximum time elapsed during initialization of the linear solver [nanoseconds]
	NanosecondsFact       int64 // maximum time elapsed during factorization (if any) [nanoseconds]
	NanosecondsLinSol     int64 // maximum time elapsed during solution of linear system (if any) [nanoseconds]
	NanosecondsErrorEstim int64 // maximum time elapsed during the error estimate [nanoseconds]
	NanosecondsTotal      int64 // total time elapsed during solution of ODE system [nanoseconds]
}

// NewStat returns a new structure
func NewStat(lskind string, implicit bool) (o *Stat) {
	o = new(Stat)
	o.LsKind = lskind
	o.Implicit = implicit
	return
}

// Reset initializes Stat
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
// options -- noExtraInfo, noElapsedTime [all default to true]
func (o *Stat) Print(options ...bool) {
	noExtraInfo := false
	noElapsedTime := false
	if len(options) > 0 {
		noExtraInfo = options[0]
	}
	if len(options) > 1 {
		noElapsedTime = options[1]
	}
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
	if !noExtraInfo {
		io.Pf("optimal step size Hopt    = %g\n", o.Hopt)
		io.Pf("kind of linear solver     = %q\n", o.LsKind)
	}
	if !noElapsedTime {
		durStep, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsStep))
		durTotal, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsTotal))
		io.Pf("elapsed time:        Step = %v\n", durStep)
		if o.Implicit {
			durJeval, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsJeval))
			durIniSo, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsIniSol))
			durFacto, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsFact))
			durLinSo, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsLinSol))
			durEstim, _ := time.ParseDuration(io.Sf("%dns", o.NanosecondsErrorEstim))
			io.Pf("elapsed time:       Jeval = %v\n", durJeval)
			io.Pf("elapsed time:      IniSol = %v\n", durIniSo)
			io.Pf("elapsed time:        Fact = %v\n", durFacto)
			io.Pf("elapsed time:      LinSol = %v\n", durLinSo)
			io.Pf("elapsed time:  ErrorEstim = %v\n", durEstim)
		}
		io.Pf("elapsed time:       Total = %v\n", durTotal)
	}
}

func (o *Stat) updateNanosecondsStep(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsStep {
		o.NanosecondsStep = delta
	}
}

func (o *Stat) updateNanosecondsJeval(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsJeval {
		o.NanosecondsJeval = delta
	}
}

func (o *Stat) updateNanosecondsIniSol(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsIniSol {
		o.NanosecondsIniSol = delta
	}
}

func (o *Stat) updateNanosecondsFact(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsFact {
		o.NanosecondsFact = delta
	}
}

func (o *Stat) updateNanosecondsLinSol(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsLinSol {
		o.NanosecondsLinSol = delta
	}
}

func (o *Stat) updateNanosecondsErrorEstim(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsErrorEstim {
		o.NanosecondsErrorEstim = delta
	}
}

func (o *Stat) updateNanosecondsTotal(start time.Time) {
	delta := time.Now().Sub(start).Nanoseconds()
	if delta > o.NanosecondsTotal {
		o.NanosecondsTotal = delta
	}
}
