// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"gosl/fun"
	"gosl/utl"
)

// DerivCen5 approximates the derivative df/dx using central differences with 5 points.
func DerivCen5(x, h float64, f fun.Ss) (res float64) {

	// first estimate
	res, round, trunc := centralDeriv5(x, h, f)
	errFirst := round + trunc

	// second estimate
	if round < trunc && (round > 0 && trunc > 0) {

		// Compute an optimised stepsize to minimize the total error,
		// using the scaling of the truncation error (O(h^2)) and rounding error (O(1/h)).
		hOpt := h * math.Pow(round/(2.0*trunc), 1.0/3.0)
		rOpt, roundOpt, truncOpt := centralDeriv5(x, hOpt, f)
		errorOpt := roundOpt + truncOpt

		// Check that the new error is smaller, and that the new derivative
		// is consistent with the error bounds of the original estimate.
		if errorOpt < errFirst && math.Abs(rOpt-res) < 4.0*errFirst {
			res = rOpt
		}
	}
	return
}

// DerivFwd4 approximates the derivative df/dx using forward differences with 4 points.
func DerivFwd4(x, h float64, f fun.Ss) (res float64) {

	// first estimate
	res, round, trunc := forwardDeriv4(x, h, f)
	errFirst := round + trunc

	// second estimate
	if round < trunc && (round > 0 && trunc > 0) {

		// Compute an optimised stepsize to minimize the total error,
		// using the scaling of the estimated truncation error (O(h)) and rounding error (O(1/h)).
		hOpt := h * math.Pow(round/(trunc), 1.0/2.0)
		rOpt, roundOpt, truncOpt := forwardDeriv4(x, hOpt, f)
		errorOpt := roundOpt + truncOpt

		// Check that the new error is smaller, and that the new derivative
		// is consistent with the error bounds of the original estimate.
		if errorOpt < errFirst && math.Abs(rOpt-res) < 4.0*errFirst {
			res = rOpt
		}
	}
	return
}

// DerivBwd4 approximates the derivative df/dx using backward differences with 4 points.
func DerivBwd4(x, h float64, f fun.Ss) (res float64) {
	return DerivFwd4(x, -h, f)
}

// SecondDerivCen3 approximates the second derivative d²f/dx² using central differences with 3 points
func SecondDerivCen3(x, h float64, f fun.Ss) float64 {
	return (f(x-h) - 2.0*f(x) + f(x+h)) / (h * h)
}

// SecondDerivCen5 approximates the second derivative d²f/dx² using central differences with 5 points
func SecondDerivCen5(x, h float64, f fun.Ss) float64 {
	return (-f(x-2.0*h) + 16.0*f(x-h) - 30.0*f(x) + 16.0*f(x+h) - f(x+2.0*h)) / (12.0 * h * h)
}

// SecondDerivMixedO2 approximates ∂²f/∂x∂y @ x={x,y} using O(h²) formula
func SecondDerivMixedO2(x, y, h float64, f fun.Sss) float64 {
	return (f(x-h, y-h) + f(x+h, y+h) - f(x+h, y-h) - f(x-h, y+h)) / (4.0 * h * h)
}

// SecondDerivMixedO4v1 approximates ∂²f/∂x∂y @ x={x,y} using O(h⁴) formula
// from http://www.holoborodko.com/pavel/numerical-methods/numerical-derivative/central-differences
func SecondDerivMixedO4v1(x, y, h float64, f fun.Sss) float64 {
	a1 := x + h
	a2 := x + 2.0*h
	α1 := x - h
	α2 := x - 2.0*h
	b1 := y + h
	b2 := y + 2.0*h
	β1 := y - h
	β2 := y - 2.0*h
	k := f(a1, β2) + f(a2, β1) + f(α2, b1) + f(α1, b2)
	l := f(α1, β2) + f(α2, β1) + f(a1, b2) + f(a2, b1)
	m := f(a2, β2) + f(α2, b2) - f(α2, β2) - f(a2, b2)
	n := f(α1, β1) + f(a1, b1) - f(a1, β1) - f(α1, b1)
	return (-63.0*k + 63.0*l + 44.0*m + 74.0*n) / (600.0 * h * h)
}

// SecondDerivMixedO4v2 approximates ∂²f/∂x∂y @ x={x,y} using O(h⁴) formula
// from http://www.holoborodko.com/pavel/numerical-methods/numerical-derivative/central-differences
func SecondDerivMixedO4v2(x, y, h float64, f fun.Sss) float64 {
	a1 := x + h
	a2 := x + 2.0*h
	α1 := x - h
	α2 := x - 2.0*h
	b1 := y + h
	b2 := y + 2.0*h
	β1 := y - h
	β2 := y - 2.0*h
	k := f(a1, β2) + f(a2, β1) + f(α2, b1) + f(α1, b2)
	l := f(α1, β2) + f(α2, β1) + f(a1, b2) + f(a2, b1)
	m := f(a2, β2) + f(α2, b2) - f(α2, β2) - f(a2, b2)
	n := f(α1, β1) + f(a1, b1) - f(a1, β1) - f(α1, b1)
	return (8.0*k - 8.0*l - m + 64.0*n) / (144.0 * h * h)
}

// lower level functions //////////////////////////////////////////////////////////////////////////

// centralDeriv5 computes the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
func centralDeriv5(x float64, h float64, f fun.Ss) (res, absErrRound, absErrTrunc float64) {

	// constants
	EPS := 1e-15 // cannot be machine epsilon

	// Compute the derivative using the 5-point rule (x-h, x-h/2, x, x+h/2, x+h).
	// Note that the central point is not used.
	// Compute the error using the difference between the 5-point and the 3-point rule (x-h,x,x+h).
	// Again the central point is not used.
	fm1 := f(x - h)
	fp1 := f(x + h)
	fmh := f(x - h/2)
	fph := f(x + h/2)
	r3 := 0.5 * (fp1 - fm1)
	r5 := (4.0/3.0)*(fph-fmh) - (1.0/3.0)*r3
	e3 := (math.Abs(fp1) + math.Abs(fm1)) * EPS
	e5 := 2.0*(math.Abs(fph)+math.Abs(fmh))*EPS + e3

	// The next term is due to finite precision in x+h = O (eps * x)
	dy := utl.Max(math.Abs(r3/h), math.Abs(r5/h)) * (math.Abs(x) / h) * EPS

	// The truncation error in the r5 approximation itself is O(h^4).
	// However, for safety, we estimate the error from r5-r3, which is O(h^2).
	// By scaling h we will minimise this estimated error, not the actual truncation error in r5.
	res = r5 / h
	absErrTrunc = math.Abs((r5 - r3) / h) // Estimated truncation error O(h^2)
	absErrRound = math.Abs(e5/h) + dy     // Rounding error (cancellations)
	return
}

// forwardDeriv4 compute the derivative using the 4-point rule (x+h/4, x+h/2, x+3h/4, x+h).
func forwardDeriv4(x, h float64, f fun.Ss) (res, absErrRound, absErrTrunc float64) {

	// constants
	EPS := 1e-15 // cannot be machine epsilon

	// Compute the derivative using the 4-point rule (x+h/4, x+h/2, x+3h/4, x+h).
	// Compute the error using the difference between the 4-point and the 2-point rule (x+h/2,x+h).
	f1 := f(x + h/4.0)
	f2 := f(x + h/2.0)
	f3 := f(x + (3.0/4.0)*h)
	f4 := f(x + h)
	r2 := 2.0 * (f4 - f2)
	r4 := (22.0/3.0)*(f4-f3) - (62.0/3.0)*(f3-f2) + (52.0/3.0)*(f2-f1)

	// Estimate the rounding error for r4
	e4 := 2 * 20.67 * (math.Abs(f4) + math.Abs(f3) + math.Abs(f2) + math.Abs(f1)) * EPS

	// The next term is due to finite precision in x+h = O (eps * x)
	dy := utl.Max(math.Abs(r2/h), math.Abs(r4/h)) * math.Abs(x/h) * EPS

	// The truncation error in the r4 approximation itself is O(h^3).
	// However, for safety, we estimate the error from r4-r2, which is O(h).
	// By scaling h we will minimise this estimated error, not the actual truncation error in r4.
	res = r4 / h
	absErrTrunc = math.Abs((r4 - r2) / h) // Estimated truncation error O(h)
	absErrRound = math.Abs(e4/h) + dy
	return
}
