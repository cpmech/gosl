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
}

// Init initialises structure
func (o *BwEuler) Init(distr bool) (err error) {
	return nil
}

// Accept accepts update
func (o *BwEuler) Accept(sol *Solver, y la.Vector) {
}

// Step steps update
func (o *BwEuler) Step(sol *Solver, y0 la.Vector, x0 float64) (rerr float64, err error) {

	// new x
	x0 += sol.h

	// previous y
	sol.v[0].Apply(1, y0) // v := y_old

	// iterations
	var rmsnr float64 // rms norm of residual
	var it int
	for it = 0; it < sol.NmaxIt; it++ {

		// max iterations ?
		sol.nit = it + 1
		if sol.nit > sol.Nitmax {
			sol.Nitmax = sol.nit
		}

		// calculate f @ update y
		sol.Nfeval++
		err = sol.fcn(sol.f[0], sol.h, x0, y0)
		if err != nil {
			return
		}

		// calculate residual
		rmsnr = 0.0
		for i := 0; i < sol.ndim; i++ {
			sol.w[0][i] = y0[i] - sol.v[0][i] - sol.h*sol.f[0][i] // w := residual
			if sol.UseRmsNorm {
				rmsnr += math.Pow(sol.w[0][i]/sol.scal[i], 2.0)
			} else {
				rmsnr += sol.w[0][i] * sol.w[0][i]
			}
		}
		if sol.UseRmsNorm {
			rmsnr = math.Sqrt(rmsnr / float64(sol.ndim))
		} else {
			rmsnr = math.Sqrt(rmsnr)
		}
		if sol.Verbose {
			io.Pfgrey("    residual = %10.5e    (tol = %10.5e)\n", rmsnr, sol.fnewt)
		}

		// converged
		if rmsnr < sol.fnewt {
			break
		}

		// Jacobian matrix
		if sol.doinit || !sol.CteTg {
			sol.Njeval++

			// calculate Jacobian
			if sol.jac == nil { // numerical
				err = num.Jacobian(&sol.dfdyT, func(fy, yy la.Vector) (e error) {
					e = sol.fcn(fy, sol.h, x0, yy)
					return
				}, y0, sol.f[0], sol.dw[0]) // δw works here as workspace variable
			} else { // analytical
				err = sol.jac(&sol.dfdyT, sol.h, x0, y0)
			}
			if err != nil {
				return
			}
			if sol.doinit {
				sol.rctriR = new(la.Triplet)
				sol.rctriR.Init(sol.ndim, sol.ndim, sol.mTri.Len()+sol.dfdyT.Len())
			}

			// calculate drdy matrix
			la.SpTriAdd(sol.rctriR, 1, sol.mTri, -sol.h, &sol.dfdyT) // rctriR := I - h * dfdy

			// initialise linear solver
			if sol.doinit {
				err = sol.lsolR.Init(sol.rctriR, sol.symmetric, sol.lsverbose, sol.ordering, sol.scaling, sol.comm)
				if err != nil {
					return
				}
			}

			// perform factorisation
			sol.Ndecomp++
			sol.lsolR.Fact()
		}

		// solve linear system
		sol.Nlinsol++
		sol.lsolR.Solve(sol.dw[0], sol.w[0], false) // δw := inv(rcmat) * residual

		// update y
		for i := 0; i < sol.ndim; i++ {
			y0[i] -= sol.dw[0][i]
		}
	}

	// did not converge
	if it == sol.NmaxIt-1 {
		chk.Panic("convergence failed with it = %d", it)
	}

	return 1e+20, err // must not be used with automatic substepping
}

// add method to database //////////////////////////////////////////////////////////////////////////

func init() {
	rkmDB[BwEulerKind] = func() RKmethod { return new(BwEuler) }
}
