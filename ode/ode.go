// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ode implements solvers for ordinary differential equations, including explicit and
// implicit Runge-Kutta methods; e.g. the fantastic Radau5 method by
// Hairer, Norsett & Wanner [1, 2].
//   References:
//     [1] Hairer E, Nørsett SP, Wanner G (1993). Solving Ordinary Differential Equations I:
//         Nonstiff Problems. Springer Series in Computational Mathematics, Vol. 8, Berlin,
//         Germany, 523 p.
//     [2] Hairer E, Wanner G (1996). Solving Ordinary Differential Equations II: Stiff and
//         Differential-Algebraic Problems. Springer Series in Computational Mathematics,
//         Vol. 14, Berlin, Germany, 614 p.
package ode

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/mpi"
	"github.com/cpmech/gosl/utl"
)

// Solver implements an ODE solver
//  Note: Distr is automatically set ON by Init if mpi is on and there are more then one processor.
//        However, it can be set OFF after calling Init.
type Solver struct {

	// method data
	method io.Enum  // method kind
	rkm    RKmethod // Runge-Kutta method

	// primary variables
	ndim int          // size of y
	fcn  Func         // dydx := f(x,y)
	jac  JacF         // Jacobian: dfdy
	out  OutF         // output function
	hasM bool         // has M matrix
	mTri *la.Triplet  // M matrix in Triplet form
	mMat *la.CCMatrix // M matrix

	// flags
	ZeroTrial  bool    // always start iterations with zero trial values (instead of collocation interpolation)
	Atol       float64 // absolute tolerance
	Rtol       float64 // relative tolerance
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
	Pll        bool    // parallel (threaded) execution
	CteTg      bool    // use constant tangent (Jacobian) in BwEuler
	UseRmsNorm bool    // use RMS norm instead of Euclidian in BwEuler
	Verbose    bool    // be more verbose, e.g. during iterations
	SaveXY     bool    // save X values in an array (e.g. for plotting)

	// output
	IdxSave int         // current index in Xvalues and Yvalues == last output
	Hvalues []float64   // h values if SaveXY is true [IdxSave]
	Xvalues []float64   // X values if SaveXY is true [IdxSave]
	Yvalues [][]float64 // Y values if SaveXY is true [ndim][IdxSave]

	// derived variables
	Distr bool    // MPI distributed execution. automatically set ON in Init if mpi is on and there are more then one processor.
	root  bool    // if distributed, tells if this is the root processor
	fnewt float64 // Newton's iterations tolerance

	// stat variables
	Nfeval    int // number of calls to fcn
	Njeval    int // number of Jacobian matrix evaluations
	Nsteps    int // total number of substeps
	Naccepted int // number of accepted substeps
	Nrejected int // number of rejected substeps
	Ndecomp   int // number of matrix decompositions
	Nlinsol   int // number of calls to linsolver
	Nitmax    int // number max of iterations

	// control variables
	doinit    bool    // flag indicating 'do initialisation' within step function
	first     bool    // first substep
	last      bool    // last substep
	reject    bool    // reject step
	diverg    bool    // flag diverging step
	dvfac     float64 // dv factor
	eta       float64 // eta tolerance
	jacIsOK   bool    // Jacobian is OK
	reuseJdec bool    // reuse current Jacobian and current decomposition
	reuseJ    bool    // reuse last Jacobian (only)
	nit       int     // current number of iterations
	hopt      float64 // optimal h after successful substepping
	theta     float64 // theta variable

	// step variables
	h, hprev float64   // step-size and previous step-size
	f0       la.Vector // f(x,y) before step
	scal     la.Vector // scal = Atol + Rtol*abs(y)

	// rk variables
	u  la.Vector   // u[stg]      = x + h*c[stg]
	v  []la.Vector // v[stg][dim] = ya[dim] + h*sum(a[stg][j]*f[j][dim], j, nstg)
	w  []la.Vector // w[stg][dim] workspace
	dw []la.Vector // dw[stg][dim] workspace
	f  []la.Vector // f[stg][dim] = f(u[stg], v[stg][dim])

	// complex variables
	v12  la.VectorC // join 1 and 2: complex(v[1],v[2])
	dw12 la.VectorC // join 1 and 2: complex(dw[1],d2[2])

	// radau5 variables
	z     []la.Vector // Radau5
	ez    la.Vector   // Radau5
	lerr  la.Vector   // Radau5
	rhs   la.Vector   // Radau5
	dfdyT la.Triplet  // Jacobian (triplet)

	// interpolation (radau5)
	ycol []la.Vector // colocation values

	// linear systems solver
	symmetric bool              // symmetric
	lsverbose bool              // verbose
	ordering  string            // ordering
	scaling   string            // scaling
	comm      *mpi.Communicator // communicator
	rctriR    *la.Triplet       // matrix for the real part
	rctriC    *la.TripletC      // matrix for the complex part
	lsolR     la.SparseSolver   // solver for the real part
	lsolC     la.SparseSolverC  // solver for the complex part
}

// NewSolver returns a new ODE structure with default values and allocated slices
func NewSolver(method io.Enum, ndim int, fcn Func, jac JacF, M *la.Triplet, out OutF) (o *Solver) {

	// new structure
	o = new(Solver)

	// primary variables
	o.ndim = ndim
	o.fcn = fcn
	o.jac = jac
	o.out = out
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
	o.Eps = 1.0e-16
	o.ThetaMax = 1.0e-3
	o.C1h = 1.0
	o.C2h = 1.2
	o.LerrStrat = 3
	o.Pll = true
	o.UseRmsNorm = true
	o.SetTol(o.Atol, o.Rtol)

	// derived variables
	o.root = true
	o.Distr = false
	if mpi.IsOn() {
		o.comm = mpi.NewCommunicator(nil)
		o.root = (o.comm.Rank() == 0)
		if o.comm.Size() > 1 {
			o.Distr = true
		}
	}

	// M matrix
	if M != nil {
		o.mTri = M
		o.mMat = o.mTri.ToMatrix(nil)
		o.hasM = true
	} else {
		if method == BwEulerKind {
			M = new(la.Triplet)
			la.SpTriSetDiag(M, o.ndim, 1)
			o.mTri = M
			o.mMat = o.mTri.ToMatrix(nil)
			o.hasM = true
		}
	}

	// method
	o.method = method
	o.rkm = NewRKmethod(method)
	o.rkm.Init(o.Distr)
	nstg := o.rkm.Nstages()

	// allocate step variables
	o.f0 = la.NewVector(o.ndim)
	o.scal = la.NewVector(o.ndim)

	// allocate rk variables
	o.u = la.NewVector(nstg)
	o.v = make([]la.Vector, nstg)
	o.w = make([]la.Vector, nstg)
	o.dw = make([]la.Vector, nstg)
	o.f = make([]la.Vector, nstg)
	if method == Radau5kind {
		o.z = make([]la.Vector, nstg)
		o.ycol = make([]la.Vector, nstg)
		o.ez = la.NewVector(o.ndim)
		o.lerr = la.NewVector(o.ndim)
		o.rhs = la.NewVector(o.ndim)
		o.v12 = la.NewVectorC(o.ndim)
		o.dw12 = la.NewVectorC(o.ndim)
	}
	for i := 0; i < nstg; i++ {
		o.v[i] = la.NewVector(o.ndim)
		o.w[i] = la.NewVector(o.ndim)
		o.dw[i] = la.NewVector(o.ndim)
		o.f[i] = la.NewVector(o.ndim)
		if method == Radau5kind {
			o.z[i] = la.NewVector(o.ndim)
			o.ycol[i] = la.NewVector(o.ndim)
		}
	}
	return
}

// SetTol sets tolerances according to Hairer and Wanner suggestions. This routine also
// checks for consistent values and only considers the case of scalars Atol and Rtol.
func (o *Solver) SetTol(atol, rtol float64) {
	o.Atol, o.Rtol = atol, rtol
	// check and change the tolerances
	β := 2.0 / 3.0
	if o.Atol <= 0.0 || o.Rtol <= 10.0*o.Eps {
		chk.Panic("tolerances are too small: Atol=%v, Rtol=%v", o.Atol, o.Rtol)
	} else {
		quot := o.Atol / o.Rtol
		o.Rtol = 0.1 * math.Pow(o.Rtol, β)
		o.Atol = o.Rtol * quot
	}
}

// Solve solves from (xa,ya) to (xb,yb) => find yb (stored in y)
func (o *Solver) Solve(y la.Vector, x, xb, Δx float64, fixstp bool) (err error) {

	// check
	if xb < x {
		err = chk.Err("xb == %v must be greater than x == %v\n", xb, x)
		return
	}

	// derived variables
	o.fnewt = max(10.0*o.Eps/o.Rtol, min(0.03, math.Sqrt(o.Rtol)))

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
		o.out(true, o.h, x, y)
	}

	// save X
	o.IdxSave = 0
	if o.SaveXY {
		o.Hvalues = make([]float64, o.NmaxSS+1)
		o.Xvalues = make([]float64, o.NmaxSS+1)
		o.Yvalues = utl.Alloc(o.ndim, o.NmaxSS+1)
		o.Xvalues[o.IdxSave] = x
		for i := 0; i < o.ndim; i++ {
			o.Yvalues[i][o.IdxSave] = y[i]
		}
		o.IdxSave++
	}

	// stat variables
	o.Nfeval = 0
	o.Njeval = 0
	o.Nsteps = 0
	o.Naccepted = 0
	o.Nrejected = 0
	o.Ndecomp = 0
	o.Nlinsol = 0
	o.Nitmax = 0

	// control variables
	o.doinit = true
	o.first = true
	o.last = false
	o.reject = false
	o.diverg = false
	o.dvfac = 0
	o.eta = 1.0
	o.jacIsOK = false
	o.reuseJdec = false
	o.reuseJ = false
	o.nit = 0
	o.hopt = o.h
	o.theta = o.ThetaMax

	// local error indicator
	var rerr float64

	// linear solver
	lsname := "umfpack"
	if o.Distr {
		lsname = "mumps"
	}
	o.lsolR = la.NewSparseSolver(lsname)
	o.lsolC = la.NewSparseSolverC(lsname)

	// free memory and show stat before leaving
	defer func() {
		o.lsolR.Free()
		o.lsolC.Free()
	}()

	// first scaling variable
	la.VecScaleAbs(o.scal, o.Atol, o.Rtol, y) // o.scal := o.Atol + o.Rtol * abs(y)

	// fixed steps
	if fixstp {
		o.w[0].Apply(1, y) // w0 := y (copy initial values to worksapce)
		if o.Verbose {
			io.Pfgreen("x = %v\n", x)
		}
		for x < xb {
			//if x + o.h > xb { o.h = xb - x }
			if o.jac == nil { // numerical Jacobian
				if o.method == Radau5kind {
					o.Nfeval++
					o.fcn(o.f0, o.h, x, y)
				}
			}
			o.reuseJdec = false
			o.reuseJ = false
			o.jacIsOK = false
			o.rkm.Step(o, y, x)
			o.Nsteps++
			o.doinit = false
			o.first = false
			o.hprev = o.h
			x += o.h
			o.rkm.Accept(o, y)
			if o.out != nil {
				o.out(false, o.h, x, y)
			}
			if o.SaveXY {
				if o.IdxSave < o.NmaxSS {
					o.Hvalues[o.IdxSave] = o.h
					o.Xvalues[o.IdxSave] = x
					for i := 0; i < o.ndim; i++ {
						o.Yvalues[i][o.IdxSave] = y[i]
					}
					o.IdxSave++
				}
			}
			if o.Verbose {
				io.Pfgreen("x = %v\n", x)
			}
		}
		return
	}

	// first function evaluation
	o.Nfeval++
	o.fcn(o.f0, o.h, x, y) // o.f0 := f(x,y)

	// time loop
	var dxmax, xstep, fac, div, dxnew, facgus, oldH, oldRerr float64
	var dxratio float64
	var failed bool
	for x < xb {
		dxmax, xstep = Δx, x+Δx
		failed = false
		for iss := 0; iss < o.NmaxSS+1; iss++ {

			// total number of substeps
			o.Nsteps++

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
			rerr, err = o.rkm.Step(o, y, x)

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
				o.Naccepted++
				o.first = false
				o.jacIsOK = false

				// update x and y
				o.hprev = o.h
				x += o.h
				o.rkm.Accept(o, y)

				// output
				if o.out != nil {
					o.out(false, o.h, x, y)
				}

				// save X value
				if o.SaveXY {
					o.Hvalues[o.IdxSave] = o.h
					o.Xvalues[o.IdxSave] = x
					for i := 0; i < o.ndim; i++ {
						o.Yvalues[i][o.IdxSave] = y[i]
					}
					o.IdxSave++
				}

				// converged ?
				if o.last {
					o.hopt = o.h // optimal h
					break
				}

				// predictive controller of Gustafsson
				if o.PredCtrl {
					if o.Naccepted > 1 {
						facgus = (oldH / o.h) * math.Pow(math.Pow(rerr, 2.0)/oldRerr, 0.25) / o.Mfac
						facgus = max(o.Mmin, min(o.Mmax, facgus))
						div = max(div, facgus)
						dxnew = o.h / div
					}
					oldH = o.h
					oldRerr = max(1.0e-2, rerr)
				}

				// calc new scal and f0
				la.VecScaleAbs(o.scal, o.Atol, o.Rtol, y) // o.scal := o.Atol + o.Rtol * abs(y)
				o.Nfeval++
				o.fcn(o.f0, o.h, x, y) // o.f0 := f(x,y)

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
					o.reuseJdec = (o.theta <= o.ThetaMax && dxratio >= o.C1h && dxratio <= o.C2h)
					if !o.reuseJdec {
						o.h = dxnew
					}
				}

				// check θ to decide if at least the Jacobian can be reused
				if !o.reuseJdec {
					o.reuseJ = (o.theta <= o.ThetaMax)
				}

				// rejected
			} else {
				// set flags
				if o.Naccepted > 0 {
					o.Nrejected++
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
			err = chk.Err("substepping did not converge after %d steps\n", o.NmaxSS)
			break
		}
	}
	return
}

// Stat prints "statistical" information about the solution process
func (o *Solver) Stat() {
	io.Pf("number of F evaluations   =%6d\n", o.Nfeval)
	io.Pf("number of J evaluations   =%6d\n", o.Njeval)
	io.Pf("total number of steps     =%6d\n", o.Nsteps)
	io.Pf("number of accepted steps  =%6d\n", o.Naccepted)
	io.Pf("number of rejected steps  =%6d\n", o.Nrejected)
	io.Pf("number of decompositions  =%6d\n", o.Ndecomp)
	io.Pf("number of lin solutions   =%6d\n", o.Nlinsol)
	io.Pf("max number of iterations  =%6d\n", o.Nitmax)
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
