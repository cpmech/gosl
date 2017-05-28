// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "math"

type Cb_fx func(x float64, args ...interface{}) float64

// DerivFwd employs a 1st order forward difference to approximate the derivative of f(x) w.r.t x @ x
func DerivFwd(f Cb_fx, x float64, args ...interface{}) float64 {
	delta := math.Sqrt(EPS * max(CTE1, math.Abs(x)))
	return (f(x+delta, args...) - f(x, args...)) / delta
}

// DerivBwd employs a 1st order backward difference to approximate the derivative of f(x) w.r.t x @ x
func DerivBwd(f Cb_fx, x float64, args ...interface{}) float64 {
	delta := math.Sqrt(EPS * max(CTE1, math.Abs(x)))
	return (f(x, args...) - f(x-delta, args...)) / delta
}

// DerivCen uses a central difference method
func DerivCen(f Cb_fx, x float64, args ...interface{}) float64 {
	res, _ := DerivCentral(f, x, 1e-3, args...)
	return res
}

// DerivRange chooses fwd, bwd or cen according to whether x is near the boundary or not
func DerivRange(f Cb_fx, x, xmin, xmax float64, args ...interface{}) float64 {
	h := 1e-3
	if x < xmin+h {
		return DerivFwd(f, x, args...)
	}
	if x > xmax-h {
		return DerivBwd(f, x, args...)
	}
	return DerivCen(f, x, args...)
}

func central_deriv(f Cb_fx, x float64, h float64, args ...interface{}) (result, abserr_round, abserr_trunc float64) {

	/* Compute the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
	   Note that the central point is not used.
	   Compute the error using the difference between the 5-point and the 3-point rule (x-h,x,x+h).
	   Again the central point is not used.
	*/

	fm1 := f(x-h, args...)
	fp1 := f(x+h, args...)
	fmh := f(x-h/2, args...)
	fph := f(x+h/2, args...)
	r3 := 0.5 * (fp1 - fm1)
	r5 := (4.0/3.0)*(fph-fmh) - (1.0/3.0)*r3
	e3 := (math.Abs(fp1) + math.Abs(fm1)) * DBL_EPSILON
	e5 := 2.0*(math.Abs(fph)+math.Abs(fmh))*DBL_EPSILON + e3

	/* The next term is due to finite precision in x+h = O (eps * x) */

	dy := max(math.Abs(r3/h), math.Abs(r5/h)) * (math.Abs(x) / h) * DBL_EPSILON

	/* The truncation error in the r5 approximation itself is O(h^4).
	   However, for safety, we estimate the error from r5-r3, which is O(h^2).
	   By scaling h we will minimise this estimated error, not the actual truncation error in r5. */

	result = r5 / h
	abserr_trunc = math.Abs((r5 - r3) / h) // Estimated truncation error O(h^2)
	abserr_round = math.Abs(e5/h) + dy     // Rounding error (cancellations)

	return
}

func DerivCentral(f Cb_fx, x, h float64, args ...interface{}) (result, abserr float64) {

	r_0, round, trunc := central_deriv(f, x, h, args...)
	err := round + trunc

	if round < trunc && (round > 0 && trunc > 0) {

		/* Compute an optimised stepsize to minimize the total error, using the scaling of the truncation
		   error (O(h^2)) and rounding error (O(1/h)). */

		h_opt := h * math.Pow(round/(2.0*trunc), 1.0/3.0)
		r_opt, round_opt, trunc_opt := central_deriv(f, x, h_opt, args...)
		error_opt := round_opt + trunc_opt

		/* Check that the new error is smaller, and that the new derivative
		   is consistent with the error bounds of the original estimate. */

		if error_opt < err && math.Abs(r_opt-r_0) < 4.0*err {
			r_0 = r_opt
			err = error_opt
		}
	}

	result = r_0
	abserr = err
	return
}

func forward_deriv(f Cb_fx, x, h float64, args ...interface{}) (result, abserr_round, abserr_trunc float64) {

	/* Compute the derivative using the 4-point rule (x+h/4, x+h/2, x+3h/4, x+h).
	   Compute the error using the difference between the 4-point and the 2-point rule (x+h/2,x+h). */

	f1 := f(x+h/4.0, args...)
	f2 := f(x+h/2.0, args...)
	f3 := f(x+(3.0/4.0)*h, args...)
	f4 := f(x+h, args...)
	r2 := 2.0 * (f4 - f2)
	r4 := (22.0/3.0)*(f4-f3) - (62.0/3.0)*(f3-f2) + (52.0/3.0)*(f2-f1)

	/* Estimate the rounding error for r4 */

	e4 := 2 * 20.67 * (math.Abs(f4) + math.Abs(f3) + math.Abs(f2) + math.Abs(f1)) * DBL_EPSILON

	/* The next term is due to finite precision in x+h = O (eps * x) */

	dy := max(math.Abs(r2/h), math.Abs(r4/h)) * math.Abs(x/h) * DBL_EPSILON

	/* The truncation error in the r4 approximation itself is O(h^3). However, for safety, we estimate the
	   error from r4-r2, which is O(h). By scaling h we will minimise this estimated error, not
	   the actual truncation error in r4. */

	result = r4 / h
	abserr_trunc = math.Abs((r4 - r2) / h) // Estimated truncation error O(h)
	abserr_round = math.Abs(e4/h) + dy
	return
}

func DerivForward(f Cb_fx, x, h float64, args ...interface{}) (result, abserr float64) {

	r_0, round, trunc := forward_deriv(f, x, h, args...)
	err := round + trunc

	if round < trunc && (round > 0 && trunc > 0) {

		/* Compute an optimised stepsize to minimize the total error, using the scaling of the estimated
		   truncation error (O(h)) and rounding error (O(1/h)). */

		h_opt := h * math.Pow(round/(trunc), 1.0/2.0)
		r_opt, round_opt, trunc_opt := forward_deriv(f, x, h_opt, args...)
		error_opt := round_opt + trunc_opt

		/* Check that the new error is smaller, and that the new derivative
		   is consistent with the error bounds of the original estimate. */

		if error_opt < err && math.Abs(r_opt-r_0) < 4.0*err {
			r_0 = r_opt
			err = error_opt
		}
	}

	result = r_0
	abserr = err
	return
}

func DerivBackward(f Cb_fx, x, h float64, args ...interface{}) (result, abserr float64) {
	return DerivForward(f, x, -h, args...)
}
