// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/fun"
)

// DerivCen5 approximates the derivative of f w.r.t x using central differences with 5 points.
func DerivCen5(x, h float64, f fun.Ss) (res float64, err error) {

	// first estimate
	res, round, trunc, err := centralDeriv5(x, h, f)
	if err != nil {
		return
	}
	errFirst := round + trunc

	// second estimate
	if round < trunc && (round > 0 && trunc > 0) {

		// Compute an optimised stepsize to minimize the total error,
		// using the scaling of the truncation error (O(h^2)) and rounding error (O(1/h)).
		hOpt := h * math.Pow(round/(2.0*trunc), 1.0/3.0)
		rOpt, roundOpt, truncOpt, err2 := centralDeriv5(x, hOpt, f)
		if err2 != nil {
			return 0, err2
		}
		errorOpt := roundOpt + truncOpt

		// Check that the new error is smaller, and that the new derivative
		// is consistent with the error bounds of the original estimate.
		if errorOpt < errFirst && math.Abs(rOpt-res) < 4.0*errFirst {
			res = rOpt
		}
	}
	return
}

// DerivFwd4 approximates the derivative of f w.r.t x using forward differences with 4 points.
func DerivFwd4(x, h float64, f fun.Ss) (res float64, err error) {

	// first estimate
	res, round, trunc, err := forwardDeriv4(x, h, f)
	if err != nil {
		return
	}
	errFirst := round + trunc

	// second estimate
	if round < trunc && (round > 0 && trunc > 0) {

		// Compute an optimised stepsize to minimize the total error,
		// using the scaling of the estimated truncation error (O(h)) and rounding error (O(1/h)).
		hOpt := h * math.Pow(round/(trunc), 1.0/2.0)
		rOpt, roundOpt, truncOpt, err2 := forwardDeriv4(x, hOpt, f)
		if err2 != nil {
			return 0, err2
		}
		errorOpt := roundOpt + truncOpt

		// Check that the new error is smaller, and that the new derivative
		// is consistent with the error bounds of the original estimate.
		if errorOpt < errFirst && math.Abs(rOpt-res) < 4.0*errFirst {
			res = rOpt
		}
	}
	return
}

// DerivBwd4 approximates the derivative of f w.r.t x using backward differences with 4 points.
func DerivBwd4(x, h float64, f fun.Ss) (res float64, err error) {
	return DerivFwd4(x, -h, f)
}

// lower level functions //////////////////////////////////////////////////////////////////////////

// centralDeriv5 computes the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
func centralDeriv5(x float64, h float64, f fun.Ss) (res, absErrRound, absErrTrunc float64, err error) {

	// constants
	EPS := 1e-15 // cannot be machine epsilon

	// Compute the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
	// Note that the central point is not used.
	// Compute the error using the difference between the 5-point and the 3-point rule (x-h,x,x+h).
	// Again the central point is not used.
	fm1, err := f(x - h)
	if err != nil {
		return
	}
	fp1, err := f(x + h)
	if err != nil {
		return
	}
	fmh, err := f(x - h/2)
	if err != nil {
		return
	}
	fph, err := f(x + h/2)
	if err != nil {
		return
	}
	r3 := 0.5 * (fp1 - fm1)
	r5 := (4.0/3.0)*(fph-fmh) - (1.0/3.0)*r3
	e3 := (math.Abs(fp1) + math.Abs(fm1)) * EPS
	e5 := 2.0*(math.Abs(fph)+math.Abs(fmh))*EPS + e3

	// The next term is due to finite precision in x+h = O (eps * x)
	dy := max(math.Abs(r3/h), math.Abs(r5/h)) * (math.Abs(x) / h) * EPS

	// The truncation error in the r5 approximation itself is O(h^4).
	// However, for safety, we estimate the error from r5-r3, which is O(h^2).
	// By scaling h we will minimise this estimated error, not the actual truncation error in r5.
	res = r5 / h
	absErrTrunc = math.Abs((r5 - r3) / h) // Estimated truncation error O(h^2)
	absErrRound = math.Abs(e5/h) + dy     // Rounding error (cancellations)
	return
}

// forwardDeriv4 compute the derivative using the 4-point rule (x+h/4, x+h/2, x+3h/4, x+h).
func forwardDeriv4(x, h float64, f fun.Ss) (res, absErrRound, absErrTrunc float64, err error) {

	// constants
	EPS := 1e-15 // cannot be machine epsilon

	// Compute the derivative using the 4-point rule (x+h/4, x+h/2, x+3h/4, x+h).
	// Compute the error using the difference between the 4-point and the 2-point rule (x+h/2,x+h).
	f1, err := f(x + h/4.0)
	if err != nil {
		return
	}
	f2, err := f(x + h/2.0)
	if err != nil {
		return
	}
	f3, err := f(x + (3.0/4.0)*h)
	if err != nil {
		return
	}
	f4, err := f(x + h)
	if err != nil {
		return
	}
	r2 := 2.0 * (f4 - f2)
	r4 := (22.0/3.0)*(f4-f3) - (62.0/3.0)*(f3-f2) + (52.0/3.0)*(f2-f1)

	// Estimate the rounding error for r4
	e4 := 2 * 20.67 * (math.Abs(f4) + math.Abs(f3) + math.Abs(f2) + math.Abs(f1)) * EPS

	// The next term is due to finite precision in x+h = O (eps * x)
	dy := max(math.Abs(r2/h), math.Abs(r4/h)) * math.Abs(x/h) * EPS

	// The truncation error in the r4 approximation itself is O(h^3).
	// However, for safety, we estimate the error from r4-r2, which is O(h).
	// By scaling h we will minimise this estimated error, not the actual truncation error in r4.
	res = r4 / h
	absErrTrunc = math.Abs((r4 - r2) / h) // Estimated truncation error O(h)
	absErrRound = math.Abs(e4/h) + dy
	return
}
