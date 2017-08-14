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
	"github.com/cpmech/gosl/utl"
)

// Solver implements an ODE solver
type Solver struct {

	// structures
	conf *Config // configuration parameters
	out  *Output // output handler
	Stat *Stat   // statistics

	// problem definition
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
//
//  INPUT:
//    ndim -- problem dimension
//    conf -- configuration parameters
//    out  -- output handler [may be nil]
//    fcn  -- f(x,y) = dy/dx function
//    jac  -- Jacobian: df/dy function [may be nil ⇒ use numerical Jacobian, if neccessary]
//    M    -- "mass" matrix, such that M ⋅ dy/dx = f(x,y) [may be nil]
//
//  NOTE: remember to call Free() to release allocated resources (e.g. from the linear solvers)
//
func NewSolver(ndim int, conf *Config, out *Output, fcn Func, jac JacF, M *la.Triplet) (o *Solver, err error) {

	// main
	o = new(Solver)
	o.conf = conf
	o.out = out

	// problem definition
	o.ndim = ndim
	o.fcn = fcn
	o.jac = jac

	// allocate method
	o.rkm, err = newRKmethod(o.conf.method)
	if err != nil {
		return
	}

	// information
	var nstg int
	o.fixedOnly, o.implicit, nstg = o.rkm.Info()

	// stat
	o.Stat = NewStat(o.conf.lsKind, o.implicit)

	// workspace
	o.work = newRKwork(nstg, o.ndim)

	// initialise method
	err = o.rkm.Init(ndim, o.conf, o.work, o.Stat, fcn, jac, M)
	if err != nil {
		return
	}

	// connect continuous output function
	if o.out != nil {
		o.out.cout = o.rkm.ContOut
	}
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
	if o.fixedOnly && !o.conf.fixed {
		err = chk.Err("method %q can only be used with fixed steps. make sure to call conf.SetFixedH > 0", o.conf.method)
		return
	}

	// initial step size
	o.work.h = xf - x
	if o.conf.fixed {
		o.work.h = o.conf.fixedH
	} else {
		o.work.h = utl.Min(o.work.h, o.conf.IniH)
	}

	// stat and output
	o.Stat.Reset()
	o.Stat.Hopt = o.work.h
	if o.out != nil {
		stop, e := o.out.Execute(0, false, o.work.h, x, y)
		if stop || e != nil {
			err = e
			return
		}
	}

	// set control flags
	o.work.first = true

	// first scaling variable
	la.VecScaleAbs(o.work.scal, o.conf.atol, o.conf.rtol, y) // scal = atol + rtol * abs(y)

	// make sure that final x is equal to xf in the end
	defer func() {
		if err == nil {
			if math.Abs(x-xf) > 1e-15 {
				err = chk.Err("internal error: x must be equal to xf in the end. x-xf=%v\n", x-xf)
			}
		}
	}()

	// fixed steps //////////////////////////////
	if o.conf.fixed {
		istep := 1
		if o.conf.Verbose {
			io.Pfgreen("x = %v\n", x)
			io.Pf("y = %v\n", y)
		}
		for n := 0; n < o.conf.fixedNsteps; n++ {
			if o.implicit && o.jac == nil { // f0 for numerical Jacobian
				o.Stat.Nfeval++
				o.fcn(o.work.f0, o.work.h, x, y)
			}
			err = o.rkm.Step(x, y)
			if err != nil {
				return
			}
			o.Stat.Nsteps++
			o.work.first = false
			x = float64(n+1) * o.work.h
			o.rkm.Accept(y)
			if o.out != nil {
				stop, e := o.out.Execute(istep, false, o.work.h, x, y)
				if stop || e != nil {
					err = e
					return
				}
			}
			if o.conf.Verbose {
				io.Pfgreen("x = %v\n", x)
				io.Pf("y = %v\n", y)
			}
			istep++
		}
		return
	}

	// variable steps //////////////////////////////

	// control variables
	o.work.reuseJdec = false
	o.work.reuseJ = false
	o.work.jacIsOK = false
	o.work.hPrev = o.work.h
	o.work.nit = 0
	o.work.eta = 1.0
	o.work.theta = o.conf.ThetaMax
	o.work.dvfac = 0.0
	o.work.diverg = false
	o.work.reject = false
	o.work.rerrPrev = 1e-4

	// first function evaluation
	o.Stat.Nfeval++
	o.fcn(o.work.f0, o.work.h, x, y) // o.f0 := f(x,y)

	// time loop
	Δx := xf - x
	var dxmax, xstep, dxnew, dxratio float64
	var last, failed bool
	for x < xf {
		dxmax, xstep = Δx, x+Δx
		failed = false
		for iss := 0; iss < o.conf.NmaxSS+1; iss++ {

			// total number of substeps
			o.Stat.Nsteps++

			// error: did not converge
			if iss == o.conf.NmaxSS {
				failed = true
				break
			}

			// converged?
			if x-xstep >= 0.0 {
				break
			}

			// step update
			err = o.rkm.Step(x, y)
			if err != nil {
				return
			}

			// iterations diverging ?
			if o.work.diverg {
				o.work.diverg = false
				o.work.reject = true
				last = false
				o.work.h *= o.work.dvfac
				continue
			}

			// accepted
			if o.work.rerr < 1.0 {

				// set flags
				o.Stat.Naccepted++
				o.work.first = false
				o.work.jacIsOK = false

				// update x and y
				x += o.work.h
				dxnew = o.rkm.Accept(y)

				// output
				if o.out != nil {
					stop, e := o.out.Execute(o.Stat.Naccepted, last, o.work.h, x, y)
					if stop || e != nil {
						err = e
						return
					}
				}

				// converged ?
				if last {
					o.Stat.Hopt = o.work.h // optimal stepsize
					break
				}

				// save previous stepsize and relative error
				o.work.hPrev = o.work.h
				o.work.rerrPrev = utl.Max(o.conf.rerrPrevMin, o.work.rerr)

				// calc new scal and f0
				if o.implicit {
					la.VecScaleAbs(o.work.scal, o.conf.atol, o.conf.rtol, y)
					o.Stat.Nfeval++
					o.fcn(o.work.f0, o.work.h, x, y) // o.f0 := f(x,y)
				}

				// check new step size
				dxnew = utl.Min(dxnew, dxmax)
				if o.work.reject { // do not alow h to grow if previous was a reject
					dxnew = utl.Min(o.work.h, dxnew)
				}
				o.work.reject = false

				// do not reuse current Jacobian and decomposition by default
				o.work.reuseJdec = false

				// last step ?
				if x+dxnew-xstep >= 0.0 {
					last = true
					o.work.h = xstep - x
				} else {
					if o.implicit {
						dxratio = dxnew / o.work.h
						o.work.reuseJdec = o.work.theta <= o.conf.ThetaMax && dxratio >= o.conf.C1h && dxratio <= o.conf.C2h
						if !o.work.reuseJdec {
							o.work.h = dxnew
						}
					} else {
						o.work.h = dxnew
					}
				}

				// check θ to decide if at least the Jacobian can be reused
				if o.implicit {
					if !o.work.reuseJdec {
						o.work.reuseJ = o.work.theta <= o.conf.ThetaMax
					}
				}

				// rejected
			} else {

				// set flags
				if o.Stat.Naccepted > 0 {
					o.Stat.Nrejected++
				}
				o.work.reject = true
				last = false

				// compute next stepsize
				dxnew = o.rkm.Reject()

				// new step size
				if o.work.first {
					o.work.h = 0.1 * o.work.h
				} else {
					o.work.h = dxnew
				}

				// last step
				if x+o.work.h > xstep {
					o.work.h = xstep - x
				}
			}
		}

		// sub-stepping failed
		if failed {
			err = chk.Err("substepping did not converge after %d steps\n", o.conf.NmaxSS)
			break
		}
	}
	return
}
