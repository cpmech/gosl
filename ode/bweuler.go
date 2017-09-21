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
	ndim  int             // problem dimension
	conf  *Config         // configurations
	work  *rkwork         // workspace
	stat  *Stat           // statistics
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
func (o *BwEuler) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet) {
	if M != nil {
		chk.Panic("Backward-Euler solver cannot handle M matrix yet\n")
	}
	o.ndim = ndim
	o.conf = conf
	o.work = work
	o.stat = stat
	o.fcn = fcn
	o.jac = jac
	o.dfdy = new(la.Triplet)
	o.drdy = new(la.Triplet)
	o.imat = new(la.Triplet)
	la.SpTriSetDiag(o.imat, ndim, 1)
	o.r = la.NewVector(ndim)
	o.dr = la.NewVector(ndim)
	o.ls = la.NewSparseSolver(o.conf.lsKind)
}

// Accept accepts update
func (o *BwEuler) Accept(y0 la.Vector, x0 float64) (dxnew float64) {
	return
}

// Reject processes step rejection
func (o *BwEuler) Reject() (dxnew float64) {
	return
}

// DenseOut produces dense output (after Accept)
func (o *BwEuler) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64) {
	chk.Panic("TODO")
}

// Step steps update
func (o *BwEuler) Step(x0 float64, y0 la.Vector) {

	// auxiliary
	h := o.work.h
	k := o.work.f[0]
	yOld := o.work.v[0]

	// new x
	x0 += h

	// previous y
	yOld.Apply(1, y0) // v := y_old

	// iterations
	var rmsnr float64 // rms norm of residual
	var it int
	for it = 0; it < o.conf.NmaxIt; it++ {

		// statistics about iterations
		if it+1 > o.stat.Nitmax {
			o.stat.Nitmax = it + 1
		}

		// trial f @ update y
		o.stat.Nfeval++
		o.fcn(k, h, x0, y0)

		// calculate residual
		rmsnr = 0.0
		for i := 0; i < o.ndim; i++ {
			o.r[i] = y0[i] - yOld[i] - h*k[i] // residual
			if o.conf.UseRmsNorm {
				rmsnr += math.Pow(o.r[i]/o.work.scal[i], 2.0)
			} else {
				rmsnr += o.r[i] * o.r[i]
			}
		}
		if o.conf.UseRmsNorm {
			rmsnr = math.Sqrt(rmsnr / float64(o.ndim))
		} else {
			rmsnr = math.Sqrt(rmsnr)
		}
		if o.conf.Verbose {
			io.Pfgrey("    residual = %10.5e    (tol = %10.5e)\n", rmsnr, o.conf.fnewt)
		}

		// converged
		if math.IsNaN(rmsnr) || math.IsInf(rmsnr, 0) {
			chk.Panic("residual is NaN or Inf. rmsnr = %v\n", rmsnr)
		}
		if rmsnr < o.conf.fnewt {
			break
		}

		// Jacobian matrix
		if o.work.first || !o.conf.CteTg {

			// stat
			o.stat.Njeval++

			// numerical Jacobian
			if o.jac == nil { // numerical
				num.Jacobian(o.dfdy, func(fy, yy la.Vector) {
					o.fcn(fy, h, x0, yy)
				}, y0, o.work.f[0], o.dr) // dr works here as workspace variable

				// analytical Jacobian
			} else {
				o.jac(o.dfdy, h, x0, y0)
			}

			// initialise drdy matrix
			if !o.ready {
				o.drdy.Init(o.ndim, o.ndim, o.imat.Len()+o.dfdy.Len())
			}

			// calculate drdy matrix
			la.SpTriAdd(o.drdy, 1, o.imat, -h, o.dfdy) // drdy = I - h ⋅ dfdy

			// initialise linear solver
			if !o.ready {
				o.ls.Init(o.drdy, o.conf.Symmetric, o.conf.LsVerbose, o.conf.Ordering, o.conf.Scaling, o.conf.comm)
				o.ready = true
			}

			// perform factorisation
			o.stat.Ndecomp++
			o.ls.Fact()
		}

		// solve linear system
		o.stat.Nlinsol++
		o.ls.Solve(o.dr, o.r, false) // dr := inv(drdy) * residual

		// update y
		for i := 0; i < o.ndim; i++ {
			y0[i] -= o.dr[i]
		}
	}

	// did not converge
	if it == o.conf.NmaxIt-1 {
		chk.Panic("convergence failed with nit = %d", it+1)
	}
}
