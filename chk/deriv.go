// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"fmt"
	"math"
	"testing"
)

const (
	DBL_EPSILON = 1.0e-15 // smallest number satisfying 1.0 + EPS > 1.0
)

// DerivVecScaCen checks the derivative of vector w.r.t scalar by comparing with numerical solution
// obtained with central differences (5-point rule)
//   tst     -- testing.T structure
//   msg     -- message about this test
//   tol     -- tolerance to compare dfdxAna with dfdxNum
//   dfdxAna -- [vector] analytical (or other kind) derivative dfdx
//   xAt     -- [scalar] position to compute dfdx
//   dx      -- stepsize; e.g. 1e-3
//   fcn     -- [vector] function f(x). x is scalar
func DerivVecScaCen(tst *testing.T, msg string, tol float64, dfdxAna []float64, xAt, dx float64, verbose bool, fcn func(f []float64, x float64)) {
	res := make([]float64, len(dfdxAna))
	for i, ana := range dfdxAna {
		fi := func(x float64) float64 {
			fcn(res, x)
			return res[i]
		}
		DerivScaScaCen(tst, fmt.Sprintf("%s%d", msg, i), tol, ana, xAt, dx, verbose, fi)
	}
}

// DerivScaVecCen checks the derivative of scalar w.r.t vector by comparing with numerical solution
// obtained with central differences (5-point rule)
//   tst     -- testing.T structure
//   msg     -- message about this test
//   tol     -- tolerance to compare dfdxAna with dfdxNum
//   dfdxAna -- [vector] analytical (or other kind) derivative dfdx
//   xAt     -- [vector] position to compute dfdx
//   dx      -- stepsize; e.g. 1e-3
//   fcn     -- [scalar] function f(x). x is vector
func DerivScaVecCen(tst *testing.T, msg string, tol float64, dfdxAna, xAt []float64, dx float64, verbose bool, fcn func(x []float64) float64) {
	ndim := len(xAt)
	if len(dfdxAna) != ndim {
		tst.Errorf("len(dfdxAna) != len(xAt)\n")
		return
	}
	xTmp := make([]float64, ndim)
	for i := 0; i < ndim; i++ {
		copy(xTmp, xAt)
		fi := func(x float64) float64 {
			xTmp[i] = x
			return fcn(xTmp)
		}
		DerivScaScaCen(tst, fmt.Sprintf("%s%d", msg, i), tol, dfdxAna[i], xAt[i], dx, verbose, fi)
	}
}

// DerivScaScaCen checks the derivative of scalar w.r.t scalar by comparing with numerical solution
// obtained with central differences (5-point rule)
//   tst     -- testing.T structure
//   msg     -- message about this test
//   tol     -- tolerance to compare dfdxAna with dfdxNum
//   dfdxAna -- [scalar] analytical (or other kind) derivative dfdx
//   xAt     -- [scalar] position to compute dfdx
//   dx      -- stepsize; e.g. 1e-3
//   fcn     -- [scalar] function f(x). x is scalar
func DerivScaScaCen(tst *testing.T, msg string, tol, dfdxAna, xAt, dx float64, verbose bool, fcn func(x float64) float64) {

	// call centralDeriv first
	r_0, round, trunc := centralDeriv(fcn, xAt, dx)
	err := round + trunc

	// check rounding error
	if round < trunc && (round > 0 && trunc > 0) {

		// compute an optimised stepsize to minimize the total error, using the scaling of the
		// truncation error (O(h^2)) and rounding error (O(1/h)).
		h_opt := dx * math.Pow(round/(2.0*trunc), 1.0/3.0)
		r_opt, round_opt, trunc_opt := centralDeriv(fcn, xAt, h_opt)
		error_opt := round_opt + trunc_opt

		// check that the new error is smaller, and that the new derivative is consistent with the
		// error bounds of the original estimate.
		if error_opt < err && math.Abs(r_opt-r_0) < 4.0*err {
			r_0 = r_opt
			err = error_opt
		}
	}

	// compare
	AnaNum(tst, msg, tol, dfdxAna, r_0, verbose)
}

// centralDeriv Computes the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
func centralDeriv(f func(x float64) float64, x float64, h float64) (result, abserr_round, abserr_trunc float64) {

	// compute the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).  Note that the
	// central point is not used.  Compute the error using the difference between the 5-point and the
	// 3-point rule (x-h,x,x+h).  Again the central point is not used.
	fm1 := f(x - h)
	fp1 := f(x + h)
	fmh := f(x - h/2)
	fph := f(x + h/2)
	r3 := 0.5 * (fp1 - fm1)
	r5 := (4.0/3.0)*(fph-fmh) - (1.0/3.0)*r3
	e3 := (math.Abs(fp1) + math.Abs(fm1)) * DBL_EPSILON
	e5 := 2.0*(math.Abs(fph)+math.Abs(fmh))*DBL_EPSILON + e3

	// the next term is due to finite precision in x+h = O (eps * x)
	dy := max(math.Abs(r3/h), math.Abs(r5/h)) * (math.Abs(x) / h) * DBL_EPSILON

	// the truncation error in the r5 approximation itself is O(h^4).  However, for safety, we estimate
	// the error from r5-r3, which is O(h^2).  By scaling h we will minimise this estimated error, not
	// the actual truncation error in r5.
	result = r5 / h
	abserr_trunc = math.Abs((r5 - r3) / h) // Estimated truncation error O(h^2)
	abserr_round = math.Abs(e5/h) + dy     // Rounding error (cancellations)
	return
}
