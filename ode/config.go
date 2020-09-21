// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"gosl/chk"
	"gosl/la"
	"gosl/mpi"
	"gosl/utl"
)

// Config holds configuration parameters for the ODE solver
type Config struct {

	// parameters
	Hmin       float64 // minimum H allowed
	IniH       float64 // initial H
	NmaxIt     int     // max num iterations (allowed)
	NmaxSS     int     // max num substeps
	Mmin       float64 // min step multiplier
	Mmax       float64 // max step multiplier
	Mfac       float64 // step multiplier factor
	MfirstRej  float64 // coefficient to multiply stepsize if first step is rejected [0 ⇒ use dxnew]
	PredCtrl   bool    // use Gustafsson's predictive controller
	Eps        float64 // smallest number satisfying 1.0 + ϵ > 1.0
	ThetaMax   float64 // max theta to decide whether the Jacobian should be recomputed or not
	C1h        float64 // c1 of HW-VII p124 => min ratio to retain previous h
	C2h        float64 // c2 of HW-VII p124 => max ratio to retain previous h
	LerrStrat  int     // strategy to select local error computation method
	GoChan     bool    // allow use of go channels (threaded); e.g. to solve R and C systems concurrently
	CteTg      bool    // use constant tangent (Jacobian) in BwEuler
	UseRmsNorm bool    // use RMS norm instead of Euclidean in BwEuler
	Verbose    bool    // show messages, e.g. during iterations
	ZeroTrial  bool    // always start iterations with zero trial values (instead of collocation interpolation)
	StabBeta   float64 // Lund stabilisation coefficient β

	// stiffness detection
	StiffNstp  int     // number of steps to check stiff situation. 0 ⇒ no check. [default = 1]
	StiffRsMax float64 // maximum value of ρs [default = 0.5]
	StiffNyes  int     // number of "yes" stiff steps allowed [default = 15]
	StiffNnot  int     // number of "not" stiff steps to disregard stiffness [default = 6]

	// configurations for linear solver
	LinSolConfig *la.SparseConfig // configurations for sparse linear solver

	// output
	stepF     StepOutF  // function to process step output (of accepted steps) [may be nil]
	denseF    DenseOutF // function to process dense output [may be nil]
	denseDx   float64   // step size for dense output
	stepOut   bool      // perform output of (variable) steps
	denseOut  bool      // perform dense output is active
	denseNstp int       // number of dense steps

	// internal data
	method    string  // the ODE method
	stabBetaM float64 // factor to multiply stabilisation coefficient β

	// linear solver control
	lsKind string            // linear solver kind
	distr  bool              // MPI distributed execution
	comm   *mpi.Communicator // for MPI run (real linear solver)

	// tolerances
	atol  float64 // absolute tolerance
	rtol  float64 // relative tolerance
	fnewt float64 // Newton's iterations tolerance

	// coefficients
	rerrPrevMin float64 // min value of rerrPrev

	// fixed steps
	fixed       bool    // use fixed steps
	fixedH      float64 // value of fixed stepsize
	fixedNsteps int     // number of fixed steps
}

// NewConfig returns a new [default] set of configuration parameters
//   method -- the ODE method: e.g. fweuler, bweuler, radau5, moeuler, dopri5
//   lsKind -- kind of linear solver: "umfpack" or "mumps" [may be empty]
//   comm   -- communicator for the linear solver [may be nil]
//   NOTE: (1) if comm == nil, the linear solver will be "umfpack" by default
//         (2) if comm != nil and comm.Size() == 1, you can use either "umfpack" or "mumps"
//         (3) if comm != nil and comm.Size() > 1, the linear solver will be set to "mumps" automatically
func NewConfig(method string, lsKind string, comm *mpi.Communicator) (o *Config) {

	// check kind of linear solver
	if lsKind != "" && lsKind != "umfpack" && lsKind != "mumps" {
		chk.Panic("lsKind must be empty or \"umfpack\" or \"mumps\"")
	}
	if lsKind == "" {
		lsKind = "umfpack"
	}
	if comm == nil {
		lsKind = "umfpack"
	} else {
		if comm.Size() > 1 {
			lsKind = "mumps"
		}
	}

	// parameters
	o = new(Config)
	o.ZeroTrial = false
	o.Hmin = 1.0e-10
	o.IniH = 1.0e-4
	o.NmaxIt = 7
	o.NmaxSS = 1000
	o.Mmin = 0.125
	o.Mmax = 5.0
	o.Mfac = 0.9
	o.MfirstRej = 0.1
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

	// stiffness detection
	o.StiffNstp = 0
	o.StiffRsMax = 0.5
	o.StiffNyes = 15
	o.StiffNnot = 6

	// configurations for linear solver
	o.LinSolConfig = la.NewSparseConfig(comm)

	// internal data
	o.method = method

	// linear solver control
	o.lsKind = lsKind
	o.distr = false
	o.comm = comm
	if lsKind == "mumps" {
		o.distr = true
	}

	// set tolerances
	o.SetTols(1e-4, 1e-4)

	// coefficients
	o.rerrPrevMin = 1e-4
	switch method {
	case "radau5":
		o.rerrPrevMin = 1e-2
	case "dopri5":
		o.StabBeta = 0.04
		o.stabBetaM = 0.75
	case "dopri8":
		o.stabBetaM = 0.2
	}
	return
}

// SetTols sets tolerances according to Hairer and Wanner' suggestions
//   atol   -- absolute tolerance; use 0 for default [default = 1e-4]
//   rtol   -- relative tolerance; use 0 for default [default = 1e-4]
func (o *Config) SetTols(atol, rtol float64) {

	// check
	if atol <= 0.0 || atol <= 10.0*o.Eps {
		chk.Panic("tolerances are too small: Atol=%v, Rtol=%v", atol, atol)
	}

	// set
	o.atol, o.rtol = atol, rtol

	// check and change the tolerances [radau5 only]
	if o.method == "radau5" {
		β := 2.0 / 3.0
		quot := o.atol / o.rtol
		o.rtol = 0.1 * math.Pow(o.rtol, β)
		o.atol = o.rtol * quot
	}

	// tolerance for iterations
	o.fnewt = utl.Max(10.0*o.Eps/o.rtol, utl.Min(0.03, math.Sqrt(o.rtol)))
}

// SetTol sets both tolerances: Atol and Rtol
func (o *Config) SetTol(atolAndRtol float64) {
	o.SetTols(atolAndRtol, atolAndRtol)
}

// SetFixedH calculates the number of steps, the exact stepsize h, and set to use fixed stepsize
func (o *Config) SetFixedH(dxApprox, xf float64) {
	o.fixed = true
	o.fixedNsteps = int(math.Ceil(xf / dxApprox))
	o.fixedH = xf / float64(o.fixedNsteps)
	xfinal := float64(o.fixedNsteps) * o.fixedH
	if math.Abs(xfinal-xf) > 1e-14 {
		chk.Panic("_internal_: xfinal should be equal to xf. xfinal-xf=%25.18e\n", xfinal-xf)
	}
}

// SetStepOut activates output of (variable) steps
//  save -- save all values
//  out  -- function to be during step output [may be nil]
func (o *Config) SetStepOut(save bool, out StepOutF) {
	o.stepOut = save
	o.stepF = out
}

// SetDenseOut activates dense output
//  save -- save all values
//  out  -- function to be during dense output [may be nil]
func (o *Config) SetDenseOut(save bool, dxOut, xf float64, out DenseOutF) {
	if dxOut > 0 {
		o.denseOut = save
		o.denseF = out
		o.denseNstp = int(math.Ceil(xf / dxOut))
		o.denseDx = dxOut
	}
}
