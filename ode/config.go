// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/utl"
)

// Config holds configuration parameters for the ODE solver
type Config struct {

	// parameters
	Method     string  // the ODE method
	FixedStp   float64 // if >0, use fixed steps instead of automatic substepping
	ZeroTrial  bool    // always start iterations with zero trial values (instead of collocation interpolation)
	Hmin       float64 // minimum H allowed
	IniH       float64 // initial H
	NmaxIt     int     // max num iterations (allowed)
	NmaxSS     int     // max num substeps
	Mmin       float64 // min step multiplier
	Mmax       float64 // max step multiplier
	Mfac       float64 // step multiplier factor
	PredCtrl   bool    // use Gustafsson's predictive controller
	Eps        float64 // smallest number satisfying 1.0 + ϵ > 1.0
	ThetaMax   float64 // max theta to decide whether the Jacobian should be recomputed or not
	C1h        float64 // c1 of HW-VII p124 => min ratio to retain previous h
	C2h        float64 // c2 of HW-VII p124 => max ratio to retain previous h
	LerrStrat  int     // strategy to select local error computation method
	GoChan     bool    // allow use of go channels (threaded); e.g. to solve R and C systems concurrently
	CteTg      bool    // use constant tangent (Jacobian) in BwEuler
	UseRmsNorm bool    // use RMS norm instead of Euclidian in BwEuler
	Verbose    bool    // show messages, e.g. during iterations
	SaveXY     bool    // save X values in an array (e.g. for plotting)

	// linear solver
	Symmetric bool   // assume symmetric matrix
	LsVerbose bool   // show linear solver messages
	Ordering  string // ordering for linear solver
	Scaling   string // scaling for linear solver

	// linear solver control
	comm   *mpi.Communicator // for MPI run (real linear solver)
	lsKind string            // linear solver kind
	distr  bool              // MPI distributed execution

	// tolerances
	atol  float64 // absolute tolerance
	rtol  float64 // relative tolerance
	fnewt float64 // Newton's iterations tolerance
}

// NewConfig returns a new [default] set of configuration parameters
//   method -- the ODE method: e.g. fweuler, bweuler, radau5, moeuler, dopri5
//   comm   -- communicator for the linear solver [may be nil]
//   lsKind -- kind of linear solver: "umfpack" or "mumps" [may be empty]
//   NOTE: (1) if comm == nil, the linear solver will be "umfpack" by default
//         (2) if comm != nil and comm.Size() == 1, you can use either "umfpack" or "mumps"
//         (3) if comm != nil and comm.Size() > 1, the linear solver will be set to "mumps" automatically
func NewConfig(method string, lsKind string, comm *mpi.Communicator) (o *Config, err error) {

	// parameters
	o = new(Config)
	o.Method = method
	o.FixedStp = 0
	o.ZeroTrial = false
	o.Hmin = 1.0e-10
	o.IniH = 1.0e-4
	o.NmaxIt = 7
	o.NmaxSS = 1000
	o.Mmin = 0.125
	o.Mmax = 5.0
	o.Mfac = 0.9
	o.PredCtrl = true
	o.Eps = 1.0e-16
	o.ThetaMax = 1.0e-3
	o.C1h = 1.0
	o.C2h = 1.2
	o.LerrStrat = 3
	o.GoChan = true
	o.CteTg = false
	o.UseRmsNorm = true
	o.Verbose = false
	o.SaveXY = false

	// linear solver control
	if comm == nil || lsKind == "" {
		lsKind = "umfpack"
	}
	if comm != nil {
		if comm.Size() > 1 {
			lsKind = "mumps"
			o.distr = true
		}
	}
	o.lsKind = lsKind
	o.comm = comm

	// set tolerances
	err = o.SetTol(1e-4, 1e-4)
	return
}

// SetTol sets tolerances according to Hairer and Wanner' suggestions
//   atol   -- absolute tolerance; use 0 for default [default = 1e-4]
//   rtol   -- relative tolerance; use 0 for default [default = 1e-4]
func (o *Config) SetTol(atol, rtol float64) (err error) {

	// check
	if atol <= 0.0 || atol <= 10.0*o.Eps {
		return chk.Err("tolerances are too small: Atol=%v, Rtol=%v", atol, atol)
	}

	// set
	o.atol, o.rtol = atol, rtol

	// check and change the tolerances
	β := 2.0 / 3.0
	quot := o.atol / o.rtol
	o.rtol = 0.1 * math.Pow(o.rtol, β)
	o.atol = o.rtol * quot

	// tolerance for iterations
	o.fnewt = utl.Max(10.0*o.Eps/o.rtol, utl.Min(0.03, math.Sqrt(o.rtol)))
	return
}

// dxnew computes standard dx estimate
func (o *Config) dxnew(h, rerr float64, nit int) (dx, div float64) {
	fac := utl.Min(o.Mfac, o.Mfac*float64(1+2*o.NmaxIt)/float64(nit+2*o.NmaxIt))
	div = utl.Max(o.Mmin, utl.Min(o.Mmax, math.Pow(rerr, 0.25)/fac))
	dx = h / div
	return
}

// dxnewGus computes dx estimate using predictive controller of Gustafsson
func (o *Config) dxnewGus(div, oldH, h, oldRerr, rerr float64) float64 {
	r2 := rerr * rerr
	fac := (oldH / h) * math.Pow(r2/oldRerr, 0.25) / o.Mfac
	fac = utl.Max(o.Mmin, utl.Min(o.Mmax, fac))
	div = utl.Max(div, fac)
	return h / div
}
