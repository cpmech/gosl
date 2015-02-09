// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"math"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/utl"
)

func bweuler_accept(o *ODE, y []float64) {
}

// backward-Euler
func bweuler_step(o *ODE, y []float64, x float64, args ...interface{}) (rerr float64, err error) {

	// new x
	x += o.h

	// previous y
	la.VecCopy(o.v[0], 1, y) // v := y_old

	// iterations
	var rmsnr float64 // rms norm of residual
	var it int
	for it = 0; it < o.NmaxIt; it++ {

		// max iterations ?
		o.nit = it + 1
		if o.nit > o.nitmax {
			o.nitmax = o.nit
		}

		// calculate f @ update y
		o.nfeval += 1
		err = o.fcn(o.f[0], x, y, args...)
		if err != nil {
			return
		}

		// calculate residual
		rmsnr = 0.0
		for i := 0; i < o.ndim; i++ {
			o.w[0][i] = y[i] - o.v[0][i] - o.h*o.f[0][i] // w := residual
			if o.UseRmsNorm {
				rmsnr += math.Pow(o.w[0][i]/o.scal[i], 2.0)
			} else {
				rmsnr += o.w[0][i] * o.w[0][i]
			}
		}
		if o.UseRmsNorm {
			rmsnr = math.Sqrt(rmsnr / float64(o.ndim))
		} else {
			rmsnr = math.Sqrt(rmsnr)
		}
		if o.Verbose {
			utl.Pfgrey("    residual = %10.5e    (tol = %10.5e)\n", rmsnr, o.fnewt)
		}

		// converged
		if rmsnr < o.fnewt {
			break
		}

		// Jacobian matrix
		if o.doinit || !o.CteTg {
			o.njeval += 1

			// calculate Jacobian
			if o.jac == nil { // numerical
				err = num.Jacobian(&o.dfdyT, func(fy, yy []float64) (e error) {
					e = o.fcn(fy, x, yy, args...)
					return
				}, y, o.f[0], o.δw[0], o.Distr) // δw works here as workspace variable
			} else { // analytical
				err = o.jac(&o.dfdyT, x, y, args...)
			}
			if err != nil {
				return
			}
			// debug
			//if true {
			//utl.Pfblue2("J = %v\n", o.dfdyT.ToMatrix(nil).ToDense()[0])
			//}
			if o.doinit {
				o.rctriR = new(la.Triplet)
				o.rctriR.Init(o.ndim, o.ndim, o.mTri.Len()+o.dfdyT.Len())
			}

			// calculate drdy matrix
			la.SpTriAdd(o.rctriR, 1, o.mTri, -o.h, &o.dfdyT) // rctriR := I - h * dfdy
			//la.PrintMat("rcmat", o.rctriR.ToMatrix(nil).ToDense(), "%8.3f", false)

			// initialise linear solver
			if o.doinit {
				err = o.lsolR.InitR(o.rctriR, false, false, false)
				if err != nil {
					return
				}
			}

			// perform factorisation
			o.ndecomp += 1
			o.lsolR.Fact()
		}

		// solve linear system
		o.nlinsol += 1
		o.lsolR.SolveR(o.δw[0], o.w[0], false) // δw := inv(rcmat) * residual

		// update y
		for i := 0; i < o.ndim; i++ {
			y[i] -= o.δw[0][i]
		}
	}

	// did not converge
	if it == o.NmaxIt-1 {
		utl.Panic("bweuler_step failed with it = %d", it)
	}

	return 1e+20, err // must not be used with automatic substepping
}
