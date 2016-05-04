// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"fmt"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// Solver implements an ODE solver
//  Note: Distr is automatically set ON by Init if mpi is on and there are more then one processor.
//        However, it can be set OFF after calling Init.
type Solver struct {

	// method
	method string  // method name
	step   stpfcn  // step function
	accept acptfcn // accept update function
	nstg   int     // number of stages

	// primary variables
	ndim       int          // size of y
	fcn        Cb_fcn       // dydx := f(x,y)
	jac        Cb_jac       // Jacobian: dfdy
	out        Cb_out       // output function
	hasM       bool         // has M matrix
	mTri       *la.Triplet  // M matrix in Triplet form
	mMat       *la.CCMatrix // M matrix
	silent     bool         // silent mode
	ZeroTrial  bool         // always start iterations with zero trial values (instead of collocation interpolation)
	Atol       float64      // absolute tolerance
	Rtol       float64      // relative tolerance
	IniH       float64      // initial H
	NmaxIt     int          // max num iterations (allowed)
	NmaxSS     int          // max num substeps
	Mmin       float64      // min step multiplier
	Mmax       float64      // max step multiplier
	Mfac       float64      // step multiplier factor
	PredCtrl   bool         // use Gustafsson's predictive controller
	ϵ          float64      // smallest number satisfying 1.0 + ϵ > 1.0
	θmax       float64      // max theta to decide whether the Jacobian should be recomputed or not
	C1h        float64      // c1 of HW-VII p124 => min ratio to retain previous h
	C2h        float64      // c2 of HW-VII p124 => max ratio to retain previous h
	LerrStrat  int          // strategy to select local error computation method
	Pll        bool         // parallel (threaded) execution
	CteTg      bool         // use constant tangent (Jacobian) in BwEuler
	UseRmsNorm bool         // use RMS norm instead of Euclidian in BwEuler
	Verbose    bool         // be more verbose, e.g. during iterations

	// derived variables
	Distr bool    // MPI distributed execution. automatically set ON in Init if mpi is on and there are more then one processor.
	root  bool    // if distributed, tells if this is the root processor
	fnewt float64 // Newton's iterations tolerance

	// stat variables
	nfeval    int // number of calls to fcn
	njeval    int // number of Jacobian matrix evaluations
	nsteps    int // total number of substeps
	naccepted int // number of accepted substeps
	nrejected int // number of rejected substeps
	ndecomp   int // number of matrix decompositions
	nlinsol   int // number of calls to linsolver
	nitmax    int // number max of iterations

	// control variables
	doinit    bool    // flag indicating 'do initialisation' within step function
	first     bool    // first substep
	last      bool    // last substep
	reject    bool    // reject step
	diverg    bool    // flag diverging step
	dvfac     float64 // dv factor
	η         float64 // eta tolerance
	jacIsOK   bool    // Jacobian is OK
	reuseJdec bool    // reuse current Jacobian and current decomposition
	reuseJ    bool    // reuse last Jacobian (only)
	nit       int     // current number of iterations
	hopt      float64 // optimal h after successful substepping
	θ         float64 // theta variable

	// step variables
	h, hprev float64   // step-size and previous step-size
	f0       []float64 // f(x,y) before step
	scal     []float64 // scal = Atol + Rtol*abs(y)

	// rk variables
	u     []float64   // u[stg]      = x + h*c[stg]
	v     [][]float64 // v[stg][dim] = ya[dim] + h*sum(a[stg][j]*f[j][dim], j, nstg)
	w, δw [][]float64 // workspace
	f     [][]float64 // f[stg][dim] = f(u[stg], v[stg][dim])

	// explicit rk variables
	erkdat ERKdat // explicit RK data

	// radau5 variables
	z             [][]float64 // Radau5
	ez, lerr, rhs []float64   // Radau5
	dfdyT         la.Triplet  // Jacobian (triplet)

	// interpolation (radau5)
	ycol [][]float64 // colocation values

	// for distributed solver
	rctriR       *la.Triplet
	rctriC       *la.TripletC
	lsolR, lsolC la.LinSol
}

// Init initialises ODE structure with default values and allocate slices
func (o *Solver) Init(method string, ndim int, fcn Cb_fcn, jac Cb_jac, M *la.Triplet, out Cb_out, silent bool) {

	// primary variables
	o.method = method
	o.ndim = ndim
	o.fcn = fcn
	o.jac = jac
	o.out = out
	o.silent = silent
	o.ZeroTrial = false
	o.Atol = 1.0e-4
	o.Rtol = 1.0e-4
	o.IniH = 1.0e-4
	o.NmaxIt = 7
	o.NmaxSS = 1000
	o.Mmin = 0.125
	o.Mmax = 5.0
	o.Mfac = 0.9
	o.PredCtrl = true
	o.ϵ = 1.0e-16
	o.θmax = 1.0e-3
	o.C1h = 1.0
	o.C2h = 1.2
	o.LerrStrat = 3
	o.Pll = true
	o.UseRmsNorm = true
	o.SetTol(o.Atol, o.Rtol)

	// derived variables
	o.root = true
	o.Distr = false
	o.init_mpi()

	// M matrix
	if M != nil {
		o.mTri = M
		o.mMat = o.mTri.ToMatrix(nil)
		o.hasM = true
	} else {
		if o.method == "BwEuler" {
			M = new(la.Triplet)
			la.SpTriSetDiag(M, o.ndim, 1)
			o.mTri = M
			o.mMat = o.mTri.ToMatrix(nil)
			o.hasM = true
		}
	}

	// method
	switch method {
	case "FwEuler":
		o.step = fweuler_step
		o.accept = fweuler_accept
		o.nstg = 1
	case "BwEuler":
		o.step = bweuler_step
		o.accept = bweuler_accept
		o.nstg = 1
	case "MoEuler":
		o.step = erk_step
		o.accept = erk_accept
		o.nstg = 2
		o.erkdat = ERKdat{true, ME2_a, ME2_b, ME2_be, ME2_c}
	case "Dopri5":
		o.step = erk_step
		o.accept = erk_accept
		o.nstg = 7
		o.erkdat = ERKdat{true, DP5_a, DP5_b, DP5_be, DP5_c}
	case "Radau5":
		if o.Distr {
			o.step = radau5_step_mpi
		} else {
			o.step = radau5_step
		}
		o.accept = radau5_accept
		o.nstg = 3
	default:
		chk.Panic(_ode_err1, method)
	}

	// allocate step variables
	o.f0 = make([]float64, o.ndim)
	o.scal = make([]float64, o.ndim)

	// allocate rk variables
	o.u = make([]float64, o.nstg)
	o.v = make([][]float64, o.nstg)
	o.w = make([][]float64, o.nstg)
	o.δw = make([][]float64, o.nstg)
	o.f = make([][]float64, o.nstg)
	if method == "Radau5" {
		o.z = make([][]float64, o.nstg)
		o.ycol = make([][]float64, o.nstg)
		o.ez = make([]float64, o.ndim)
		o.lerr = make([]float64, o.ndim)
		o.rhs = make([]float64, o.ndim)
	}
	for i := 0; i < o.nstg; i++ {
		o.v[i] = make([]float64, o.ndim)
		o.w[i] = make([]float64, o.ndim)
		o.δw[i] = make([]float64, o.ndim)
		o.f[i] = make([]float64, o.ndim)
		if method == "Radau5" {
			o.z[i] = make([]float64, o.ndim)
			o.ycol[i] = make([]float64, o.ndim)
		}
	}
}

// SetTol sets tolerances according to Hairer and Wanner suggestions. This routine also
// checks for consistent values and only considers the case of scalars Atol and Rtol.
func (o *Solver) SetTol(atol, rtol float64) {
	o.Atol, o.Rtol = atol, rtol
	// check and change the tolerances
	β := 2.0 / 3.0
	if o.Atol <= 0.0 || o.Rtol <= 10.0*o.ϵ {
		chk.Panic(_ode_err4, o.Atol, o.Rtol)
	} else {
		quot := o.Atol / o.Rtol
		o.Rtol = 0.1 * math.Pow(o.Rtol, β)
		o.Atol = o.Rtol * quot
	}
}

// Solve solves from (xa,ya) to (xb,yb) => find yb (stored in y)
func (o *Solver) Solve(y []float64, x, xb, Δx float64, fixstp bool, args ...interface{}) (err error) {

	// check
	if xb < x {
		err = chk.Err(_ode_err3, xb, x)
		return
	}

	// derived variables
	o.fnewt = max(10.0*o.ϵ/o.Rtol, min(0.03, math.Sqrt(o.Rtol)))

	// initial step size
	Δx = min(Δx, xb-x)
	if fixstp {
		o.h = Δx
	} else {
		o.h = min(Δx, o.IniH)
	}
	o.hprev = o.h

	// output initial state
	if o.out != nil {
		o.out(true, o.h, x, y, args...)
	}

	// stat variables
	o.nfeval = 0
	o.njeval = 0
	o.nsteps = 0
	o.naccepted = 0
	o.nrejected = 0
	o.ndecomp = 0
	o.nlinsol = 0
	o.nitmax = 0

	// control variables
	o.doinit = true
	o.first = true
	o.last = false
	o.reject = false
	o.diverg = false
	o.dvfac = 0
	o.η = 1.0
	o.jacIsOK = false
	o.reuseJdec = false
	o.reuseJ = false
	o.nit = 0
	o.hopt = o.h
	o.θ = o.θmax

	// local error indicator
	var rerr float64

	// linear solver
	lsname := "umfpack"
	if o.Distr {
		lsname = "mumps"
	}
	o.lsolR = la.GetSolver(lsname)
	o.lsolC = la.GetSolver(lsname)

	// clean up and show stat before leaving
	defer func() {
		o.lsolR.Clean()
		o.lsolC.Clean()
		if !o.silent {
			o.Stat()
		}
	}()

	// first scaling variable
	la.VecScaleAbs(o.scal, o.Atol, o.Rtol, y) // o.scal := o.Atol + o.Rtol * abs(y)

	// fixed steps
	if fixstp {
		la.VecCopy(o.w[0], 1, y) // copy initial values to worksapce
		if o.Verbose {
			io.Pfgreen("x = %v\n", x)
		}
		for x < xb {
			//if x + o.h > xb { o.h = xb - x }
			if o.jac == nil { // numerical Jacobian
				if o.method == "Radau5" {
					o.nfeval += 1
					o.fcn(o.f0, o.h, x, y, args...)
				}
			}
			o.reuseJdec = false
			o.reuseJ = false
			o.jacIsOK = false
			o.step(o, y, x, args...)
			o.nsteps += 1
			o.doinit = false
			o.first = false
			o.hprev = o.h
			x += o.h
			o.accept(o, y)
			if o.out != nil {
				o.out(false, o.h, x, y, args...)
			}
			if o.Verbose {
				io.Pfgreen("x = %v\n", x)
			}
		}
		return
	}

	// first function evaluation
	o.nfeval += 1
	o.fcn(o.f0, o.h, x, y, args...) // o.f0 := f(x,y)

	// time loop
	var dxmax, xstep, fac, div, dxnew, facgus, old_h, old_rerr float64
	var dxratio float64
	var failed bool
	for x < xb {
		dxmax, xstep = Δx, x+Δx
		failed = false
		for iss := 0; iss < o.NmaxSS+1; iss++ {

			// total number of substeps
			o.nsteps += 1

			// error: did not converge
			if iss == o.NmaxSS {
				failed = true
				break
			}

			// converged?
			if x-xstep >= 0.0 {
				break
			}

			// step update
			rerr, err = o.step(o, y, x, args...)

			// initialise only once
			o.doinit = false

			// iterations diverging ?
			if o.diverg {
				o.diverg = false
				o.reject = true
				o.last = false
				o.h = o.dvfac * o.h
				continue
			}

			// step size change
			fac = min(o.Mfac, o.Mfac*float64(1+2*o.NmaxIt)/float64(o.nit+2*o.NmaxIt))
			div = max(o.Mmin, min(o.Mmax, math.Pow(rerr, 0.25)/fac))
			dxnew = o.h / div

			// accepted
			if rerr < 1.0 {

				// set flags
				o.naccepted += 1
				o.first = false
				o.jacIsOK = false

				// update x and y
				o.hprev = o.h
				x += o.h
				o.accept(o, y)

				// output
				if o.out != nil {
					o.out(false, o.h, x, y, args...)
				}

				// converged ?
				if o.last {
					o.hopt = o.h // optimal h
					break
				}

				// predictive controller of Gustafsson
				if o.PredCtrl {
					if o.naccepted > 1 {
						facgus = (old_h / o.h) * math.Pow(math.Pow(rerr, 2.0)/old_rerr, 0.25) / o.Mfac
						facgus = max(o.Mmin, min(o.Mmax, facgus))
						div = max(div, facgus)
						dxnew = o.h / div
					}
					old_h = o.h
					old_rerr = max(1.0e-2, rerr)
				}

				// calc new scal and f0
				la.VecScaleAbs(o.scal, o.Atol, o.Rtol, y) // o.scal := o.Atol + o.Rtol * abs(y)
				o.nfeval += 1
				o.fcn(o.f0, o.h, x, y, args...) // o.f0 := f(x,y)

				// new step size
				dxnew = min(dxnew, dxmax)
				if o.reject { // do not alow o.h to grow if previous was a reject
					dxnew = min(o.h, dxnew)
				}
				o.reject = false

				// do not reuse current Jacobian and decomposition by default
				o.reuseJdec = false

				// last step ?
				if x+dxnew-xstep >= 0.0 {
					o.last = true
					o.h = xstep - x
				} else {
					dxratio = dxnew / o.h
					o.reuseJdec = (o.θ <= o.θmax && dxratio >= o.C1h && dxratio <= o.C2h)
					if !o.reuseJdec {
						o.h = dxnew
					}
				}

				// check θ to decide if at least the Jacobian can be reused
				if !o.reuseJdec {
					o.reuseJ = (o.θ <= o.θmax)
				}

				// rejected
			} else {
				// set flags
				if o.naccepted > 0 {
					o.nrejected += 1
				}
				o.reject = true
				o.last = false

				// new step size
				if o.first {
					o.h = 0.1 * o.h
				} else {
					o.h = dxnew
				}

				// last step
				if x+o.h > xstep {
					o.h = xstep - x
				}
			}
		}

		// sub-stepping failed
		if failed {
			err = chk.Err(_ode_err2, o.NmaxSS)
			break
		}
	}
	return
}

func (o *Solver) Stat() {
	if !o.root {
		return
	}
	io.Pf("number of F evaluations   =%6d\n", o.nfeval)
	io.Pf("number of J evaluations   =%6d\n", o.njeval)
	io.Pf("total number of steps     =%6d\n", o.nsteps)
	io.Pf("number of accepted steps  =%6d\n", o.naccepted)
	io.Pf("number of rejected steps  =%6d\n", o.nrejected)
	io.Pf("number of decompositions  =%6d\n", o.ndecomp)
	io.Pf("number of lin solutions   =%6d\n", o.nlinsol)
	io.Pf("max number of iterations  =%6d\n", o.nitmax)
}

func (o *Solver) GetStat() (s string) {
	s = fmt.Sprintf("number of F evaluations   =%6d\n", o.nfeval)
	s += fmt.Sprintf("number of J evaluations   =%6d\n", o.njeval)
	s += fmt.Sprintf("total number of steps     =%6d\n", o.nsteps)
	s += fmt.Sprintf("number of accepted steps  =%6d\n", o.naccepted)
	s += fmt.Sprintf("number of rejected steps  =%6d\n", o.nrejected)
	s += fmt.Sprintf("number of decompositions  =%6d\n", o.ndecomp)
	s += fmt.Sprintf("number of lin solutions   =%6d\n", o.nlinsol)
	s += fmt.Sprintf("max number of iterations  =%6d\n", o.nitmax)
	return
}

// auxiliary functions
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// error messages
var (
	_ode_err1 = "ode.go: ODE.Init: method %s is not available"
	_ode_err2 = "ode.go: ODE.Solve: substepping did not converge after %d steps\n"
	_ode_err3 = "ode.go: ODE.Solve: xb == %v must be greater than x == %v\n"
	_ode_err4 = "ode.go: ODE.SetHWtol: tolerances are too small: Atol=%v, Rtol=%v"
)
