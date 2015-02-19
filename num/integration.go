// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "github.com/cpmech/gosl/chk"

// Trapz returns the area below the discrete curve defined by x and y.
// Computations are carried out with the trapezoidal rule.
func Trapz(x, y []float64) (A float64) {
	if len(x) != len(y) {
		chk.Panic(_trapz_err1, len(x), len(y))
	}
	for i := 1; i < len(x); i++ {
		A += (x[i] - x[i-1]) * (y[i] + y[i-1]) / 2.0
	}
	return
}

// TrapzF (function callback version) returns the area below the discrete curve defined by x and y.
// Computations are carried out with the trapezoidal rule.
func TrapzF(x []float64, y Cb_yx) (A float64) {
	for i := 1; i < len(x); i++ {
		A += (x[i] - x[i-1]) * (y(x[i]) + y(x[i-1])) / 2.0
	}
	return A
}

// TrapzRange (x-range and function callback version) returns the area below the discrete curve defined by x and y.
// Computations are carried out with the trapezoidal rule from xa to xb, with npts points
func TrapzRange(xa, xb float64, npts int, y Cb_yx) (A float64) {
	if npts < 2 {
		chk.Panic(_trapz_err2, npts)
	}
	dx := (xb - xa) / float64(npts-1)
	var x0, x1 float64
	for i := 1; i < npts; i++ {
		x0 = xa + dx*float64(i-1)
		x1 = xa + dx*float64(i)
		A += (x1 - x0) * (y(x1) + y(x0)) / 2.0
	}
	return A
}

// error messages
var (
	_trapz_err1 = "trapz.go: Trapz: length of x and y must be the same. %d != %d"
	_trapz_err2 = "trapz.go: TrapzRange: number of points must be at least 2"
)
