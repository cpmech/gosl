// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// NlSolverConfig holds the configuration input for NlSolver
type NlSolverConfig struct {
	// input
	Verbose          bool // show messages
	ConstantJacobian bool // constant Jacobian (Modified Newton's method)
	LineSearch       bool // use line search
	LineSearchMaxIt  int  // line search maximum iterations
	MaxIterations    int  // Newton's method maximum iterations
	EnforceConvRate  bool // check and enforce convergence rate

	// function to be called during each output
	OutCallback func(x []float64) // output callback function

	// configurations for linear solver
	LinSolConfig *la.SparseConfig // configurations for sparse linear solver

	// internal
	useDenseSolver      bool // use dense solver instead of Umfpack (sparse)
	hasJacobianFunction bool // false => use numerical Jacobian (with sparse solver)

	// tolerances
	atol  float64 // absolute tolerance
	rtol  float64 // relative tolerance
	ftol  float64 // minimum value of fx
	fnewt float64 // [derived] Newton's method tolerance
}

// NewNlSolverConfig creates a new NlSolverConfig
// Default values:
//   CteJac      = false
//   LinSearch   = false
//   LinSchMaxIt = 20
//   MaxIt       = 20
//   ChkConv     = false
//   Atol        = 1e-8
//   Rtol        = 1e-8
//   Ftol        = 1e-9
func NewNlSolverConfig() (o *NlSolverConfig) {

	// input
	o = new(NlSolverConfig)
	o.Verbose = false
	o.ConstantJacobian = false
	o.LineSearch = false
	o.LineSearchMaxIt = 20
	o.MaxIterations = 20
	o.EnforceConvRate = false

	// configurations for linear solver
	o.LinSolConfig = la.NewSparseConfig(nil)

	// internal
	o.useDenseSolver = false
	o.hasJacobianFunction = false

	// tolerances
	o.SetTolerances(1e-8, 1e-8, 1e-9)
	return
}

// SetTolerances sets all tolerances
func (o *NlSolverConfig) SetTolerances(atol, rtol, ftol float64) {
	o.atol = atol
	o.rtol = rtol
	o.ftol = ftol
	o.fnewt = utl.Max(10.0*MACHEPS/rtol, utl.Min(0.03, math.Sqrt(rtol)))
}
