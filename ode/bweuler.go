// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// BwEuler implements the (implicit) Backward Euler method
type BwEuler struct {
	conf  *Config         // configurations
	fcn   Func            // dy/dx := f(x,y)
	jac   JacF            // Jacobian function: df/dy(x,y)
	dfdy  *la.Triplet     // df/dy matrix
	drdy  *la.Triplet     // linear system matrix: drdy = I - h ⋅ dfdy
	imat  *la.Triplet     // I matrix in triplet format
	r     la.Vector       // residual
	dr    la.Vector       // increment of residual
	ls    la.SparseSolver // linear solver
	ready bool            // matrices and solver are ready
}

// add method to database
func init() {
	rkmDB["bweuler"] = func() rkmethod { return new(BwEuler) }
}

// Free releases memory
func (o *BwEuler) Free() {
	if o.ls != nil {
		o.ls.Free()
	}
}

// Info returns information about this method
func (o *BwEuler) Info() (fixedOnly, implicit bool, nstages int) {
	return true, true, 1
}

// Init initialises structure
func (o *BwEuler) Init(conf *Config, ndim int, fcn Func, jac JacF, M *la.Triplet) (err error) {
	if M != nil {
		err = chk.Err("Backward-Euler solver cannot handle M matrix yet\n")
		return
	}
	o.conf = conf
	o.fcn = fcn
	o.jac = jac
	o.dfdy = new(la.Triplet)
	o.drdy = new(la.Triplet)
	o.imat = new(la.Triplet)
	la.SpTriSetDiag(o.imat, ndim, 1)
	o.r = la.NewVector(ndim)
	o.dr = la.NewVector(ndim)
	o.ls = la.NewSparseSolver(o.conf.lsKind)
	return
}

// Accept accepts update
func (o *BwEuler) Accept(y la.Vector, work *rkwork) {
}

// ContOut produces continuous output (after Accept)
func (o *BwEuler) ContOut(yOut, y la.Vector, xOut, x, h float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *BwEuler) Step(h, x0 float64, y0 la.Vector, stat *Stat, work *rkwork) (rerr float64, err error) {

	// new x
	x0 += h

	// previous y
	work.v[0].Apply(1, y0) // v := y_old

	// iterations
	var rmsnr float64 // rms norm of residual
	var it int
	for it = 0; it < o.conf.NmaxIt; it++ {

		// statistics about iterations
		if it+1 > stat.Nitmax {
			stat.Nitmax = it + 1
		}

		// trial f @ update y
		stat.Nfeval++
		err = o.fcn(work.f[0], h, x0, y0)
		if err != nil {
			return
		}

		// calculate residual
		rmsnr = 0.0
		for i := 0; i < work.ndim; i++ {
			o.r[i] = y0[i] - work.v[0][i] - h*work.f[0][i] // residual
			if o.conf.UseRmsNorm {
				rmsnr += math.Pow(o.r[i]/work.scal[i], 2.0)
			} else {
				rmsnr += o.r[i] * o.r[i]
			}
		}
		if o.conf.UseRmsNorm {
			rmsnr = math.Sqrt(rmsnr / float64(work.ndim))
		} else {
			rmsnr = math.Sqrt(rmsnr)
		}
		if o.conf.Verbose {
			io.Pfgrey("    residual = %10.5e    (tol = %10.5e)\n", rmsnr, o.conf.fnewt)
		}

		// converged
		if math.IsNaN(rmsnr) || math.IsInf(rmsnr, 0) {
			err = chk.Err("residual is NaN or Inf. rmsnr = %v\n", rmsnr)
			return
		}
		if rmsnr < o.conf.fnewt {
			break
		}

		// Jacobian matrix
		if work.first || !o.conf.CteTg {

			// stat
			stat.Njeval++

			// numerical Jacobian
			if o.jac == nil { // numerical
				err = num.Jacobian(o.dfdy, func(fy, yy la.Vector) (e error) {
					e = o.fcn(fy, h, x0, yy)
					return
				}, y0, work.f[0], o.dr) // dr works here as workspace variable

				// analytical Jacobian
			} else {
				err = o.jac(o.dfdy, h, x0, y0)
			}

			// check
			if err != nil {
				return
			}

			// initialise drdy matrix
			if !o.ready {
				o.drdy.Init(work.ndim, work.ndim, o.imat.Len()+o.dfdy.Len())
			}

			// calculate drdy matrix
			la.SpTriAdd(o.drdy, 1, o.imat, -h, o.dfdy) // drdy = I - h ⋅ dfdy

			// initialise linear solver
			if !o.ready {
				err = o.ls.Init(o.drdy, o.conf.Symmetric, o.conf.LsVerbose, o.conf.Ordering, o.conf.Scaling, o.conf.comm)
				if err != nil {
					return
				}
				o.ready = true
			}

			// perform factorisation
			stat.Ndecomp++
			err = o.ls.Fact()
			if err != nil {
				return
			}
		}

		// solve linear system
		stat.Nlinsol++
		o.ls.Solve(o.dr, o.r, false) // dr := inv(drdy) * residual

		// update y
		for i := 0; i < work.ndim; i++ {
			y0[i] -= o.dr[i]
		}
	}

	// did not converge
	if it == o.conf.NmaxIt-1 {
		err = chk.Err("convergence failed with nit = %d", it+1)
		return
	}
	return 1e+20, err // must not be used with automatic substepping
}
