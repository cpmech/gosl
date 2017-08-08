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
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Solver implements an ODE solver
type Solver struct {

	// main
	Conf *Config // configuration parameters
	Stat *Stat   // statistics
	Out  *Output // output

	// input
	ndim int  // size of y
	fcn  Func // dy/dx := f(x,y)
	jac  JacF // Jacobian: df/dy

	// method, info and workspace
	rkm       rkmethod // Runge-Kutta method
	fixedOnly bool     // method can only be used with fixed steps
	implicit  bool     // method is implicit
	work      *rkwork  // Runge-Kutta workspace
}

// NewSolver returns a new ODE structure with default values and allocated slices
//  NOTE: remember to call Free() to release allocated resources (e.g. from the linear solvers)
func NewSolver(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet, ofcn OutF) (o *Solver, err error) {

	// main
	o = new(Solver)
	o.Conf = conf
	o.Stat = NewStat()
	o.Out = NewOutput(ofcn)
	if conf.SaveXY {
		o.Out.Resize(conf.NmaxSS + 1)
	}

	// input
	o.ndim = ndim
	o.fcn = fcn
	o.jac = jac

	// method
	o.rkm, err = newRKmethod(o.Conf.Method)
	if err != nil {
		return
	}
	err = o.rkm.Init(o.Conf, ndim, fcn, jac, M)
	if err != nil {
		return
	}

	// information
	var nstg int
	o.fixedOnly, o.implicit, nstg = o.rkm.Info()

	// workspace
	o.work = newRKwork(nstg, o.ndim)
	return
}

// Free releases allocated memory (e.g. by the linear solvers)
func (o *Solver) Free() {
	if o.rkm != nil {
		o.rkm.Free()
	}
}

// Solve solves dy/dx = f(x,y) from x to xf with initial y given in y
func (o *Solver) Solve(y la.Vector, x, xf float64) (err error) {

	// check
	if xf < x {
		err = chk.Err("xf=%v must be greater than x=%v\n", xf, x)
		return
	}

	// initial step size
	h := xf - x
	fixed := false
	if o.Conf.FixedStp > 0 || o.fixedOnly {
		if o.Conf.FixedStp < o.Conf.Hmin {
			o.Conf.FixedStp = o.Conf.IniH
		}
		h = utl.Min(h, o.Conf.FixedStp)
		fixed = true
	} else {
		h = utl.Min(h, o.Conf.IniH)
	}

	// stat and output
	o.Stat.Reset()
	o.Stat.Hopt = h
	o.Out.Execute(h, x, y)

	// set control flags
	o.work.first = true

	// first scaling variable
	la.VecScaleAbs(o.work.scal, o.Conf.atol, o.Conf.rtol, y) // scal = atol + rtol * abs(y)

	// fixed steps //////////////////////////////
	if fixed {
		if o.Conf.Verbose {
			io.Pfgreen("x = %v\n", x)
			io.Pf("y = %v\n", y)
		}
		for x < xf {
			if o.implicit && o.jac == nil { // f0 for numerical Jacobian
				o.Stat.Nfeval++
				o.fcn(o.work.f0, h, x, y)
			}
			_, err = o.rkm.Step(h, x, y, o.Stat, o.work)
			if err != nil {
				return
			}
			o.Stat.Nsteps++
			o.work.first = false
			x += h
			o.rkm.Accept(y, o.work)
			o.Out.Execute(h, x, y)
			if o.Conf.Verbose {
				io.Pfgreen("x = %v\n", x)
				io.Pf("y = %v\n", y)
			}
		}
		return
	}

	// variable steps //////////////////////////////

	// control variables
	o.work.reuseJdec = false
	o.work.reuseJ = false
	o.work.jacIsOK = false
	o.work.hprev = h
	o.work.nit = 0
	o.work.eta = 1.0
	o.work.theta = o.Conf.ThetaMax
	o.work.dvfac = 0.0
	o.work.diverg = false
	o.work.reject = false

	// first function evaluation
	o.Stat.Nfeval++
	o.fcn(o.work.f0, h, x, y) // o.f0 := f(x,y)

	// time loop
	Δx := xf - x
	var dxmax, xstep, div, dxnew, oldH, oldRerr, dxratio, rerr float64
	var last, failed bool
	for x < xf {
		dxmax, xstep = Δx, x+Δx
		failed = false
		for iss := 0; iss < o.Conf.NmaxSS+1; iss++ {

			// total number of substeps
			o.Stat.Nsteps++

			// error: did not converge
			if iss == o.Conf.NmaxSS {
				failed = true
				break
			}

			// converged?
			if x-xstep >= 0.0 {
				break
			}

			// step update
			rerr, err = o.rkm.Step(h, x, y, o.Stat, o.work)

			// iterations diverging ?
			if o.work.diverg {
				o.work.diverg = false
				o.work.reject = true
				last = false
				h *= o.work.dvfac
				continue
			}

			// step size change
			dxnew, div = o.Conf.dxnew(h, rerr, o.work.nit)

			// accepted
			if rerr < 1.0 {

				// set flags
				o.Stat.Naccepted++
				o.work.first = false
				o.work.jacIsOK = false

				// update x and y
				o.work.hprev = h
				x += h
				o.rkm.Accept(y, o.work)

				// output
				o.Out.Execute(h, x, y)

				// converged ?
				if last {
					o.Stat.Hopt = h // optimal h
					break
				}

				// predictive controller of Gustafsson
				if o.Conf.PredCtrl {
					if o.Stat.Naccepted > 1 {
						dxnew = o.Conf.dxnewGus(div, oldH, h, oldRerr, rerr)
					}
					oldH = h
					oldRerr = utl.Max(1.0e-2, rerr)
				}

				// calc new scal and f0
				la.VecScaleAbs(o.work.scal, o.Conf.atol, o.Conf.rtol, y)
				o.Stat.Nfeval++
				o.fcn(o.work.f0, h, x, y) // o.f0 := f(x,y)

				// new step size
				dxnew = utl.Min(dxnew, dxmax)
				if o.work.reject { // do not alow h to grow if previous was a reject
					dxnew = utl.Min(h, dxnew)
				}
				o.work.reject = false

				// do not reuse current Jacobian and decomposition by default
				o.work.reuseJdec = false

				// last step ?
				if x+dxnew-xstep >= 0.0 {
					last = true
					h = xstep - x
				} else {
					dxratio = dxnew / h
					o.work.reuseJdec = o.work.theta <= o.Conf.ThetaMax && dxratio >= o.Conf.C1h && dxratio <= o.Conf.C2h
					if !o.work.reuseJdec {
						h = dxnew
					}
				}

				// check θ to decide if at least the Jacobian can be reused
				if !o.work.reuseJdec {
					o.work.reuseJ = o.work.theta <= o.Conf.ThetaMax
				}

				// rejected
			} else {

				// set flags
				if o.Stat.Naccepted > 0 {
					o.Stat.Nrejected++
				}
				o.work.reject = true
				last = false

				// new step size
				if o.work.first {
					h = 0.1 * h
				} else {
					h = dxnew
				}

				// last step
				if x+h > xstep {
					h = xstep - x
				}
			}
		}

		// sub-stepping failed
		if failed {
			err = chk.Err("substepping did not converge after %d steps\n", o.Conf.NmaxSS)
			break
		}
	}
	return
}
