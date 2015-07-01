// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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

// Trapz2D computes a double integral over the x-y plane; thus resulting on
// the volume defined between the function f(x,y) and the plane ortogonal to z
func Trapz2D(dx, dy float64, f [][]float64) (V float64) {

	// check
	if len(f) < 2 {
		chk.Panic("len(f)=%d is incorrect; it must be at least 2", len(f))
	}
	m, n := len(f), len(f[0])

	// corners
	V = f[0][0] + f[m-1][0] + f[0][n-1] + f[m-1][n-1]
	//M := utl.IntsAlloc(m, n)
	//M[0][0] = 1
	//M[m-1][0] = 1
	//M[0][n-1] = 1
	//M[m-1][n-1] = 1

	// top/bottom: 2
	for j := 1; j < n-1; j++ {
		V += 2.0 * (f[0][j] + f[m-1][j])
		//M[0][j] = 2
		//M[m-1][j] = 2
	}

	// left/right: 2
	for i := 1; i < m-1; i++ {
		V += 2.0 * (f[i][0] + f[i][n-1])
		//M[i][0] = 2
		//M[i][n-1] = 2
	}

	// centre: 4
	for i := 1; i < m-1; i++ {
		for j := 1; j < n-1; j++ {
			V += 4.0 * f[i][j]
			//M[i][j] = 4
		}
	}

	// final result
	V *= dx * dy / 4.0

	// debug
	//for i := 0; i < m; i++ {
	//for j := 0; j < n; j++ {
	//io.Pf("%4d", M[i][j])
	//}
	//io.Pf("\n")
	//}
	return
}

// Simps2D computes a double integral over the x-y plane (Simpson's rule); thus resulting on
// the volume defined between the function f(x,y) and the plane ortogonal to z
func Simps2D(dx, dy float64, f [][]float64) (V float64) {

	// check
	if len(f) < 2 {
		chk.Panic("len(f)=%d is incorrect; it must be at least 2", len(f))
	}
	m, n := len(f), len(f[0])

	// corners
	V = f[0][0] + f[m-1][0] + f[0][n-1] + f[m-1][n-1]
	//M := utl.IntsAlloc(m, n)
	//M[0][0] = 1
	//M[m-1][0] = 1
	//M[0][n-1] = 1
	//M[m-1][n-1] = 1

	// top/bottom: 4
	for j := 1; j < n-1; j += 2 {
		V += 4.0 * (f[0][j] + f[m-1][j])
		//M[0][j] = 4
		//M[m-1][j] = 4
	}

	// top/bottom: 2
	for j := 2; j < n-1; j += 2 {
		V += 2.0 * (f[0][j] + f[m-1][j])
		//M[0][j] = 2
		//M[m-1][j] = 2
	}

	// left/right: 4
	for i := 1; i < m-1; i += 2 {
		V += 4.0 * (f[i][0] + f[i][n-1])
		//M[i][0] = 4
		//M[i][n-1] = 4
	}

	// left/right: 2
	for i := 2; i < m-1; i += 2 {
		V += 2.0 * (f[i][0] + f[i][n-1])
		//M[i][0] = 2
		//M[i][n-1] = 2
	}

	// centre: 4
	for i := 2; i < m-1; i += 2 {
		for j := 2; j < n-1; j += 2 {
			V += 4.0 * f[i][j]
			//M[i][j] = 4
		}
	}

	// centre: 8a
	for i := 1; i < m-1; i += 2 {
		for j := 2; j < n-1; j += 2 {
			V += 8.0 * f[i][j]
			//M[i][j] = 8
		}
	}

	// centre: 8b
	for i := 2; i < m-1; i += 2 {
		for j := 1; j < n-1; j += 2 {
			V += 8.0 * f[i][j]
			//M[i][j] = 8
		}
	}

	// centre: 16
	for i := 1; i < m-1; i += 2 {
		for j := 1; j < n-1; j += 2 {
			V += 16.0 * f[i][j]
			//M[i][j] = 16
		}
	}

	// final result
	V *= dx * dy / 9.0

	// debug
	//for i := 0; i < m; i++ {
	//for j := 0; j < n; j++ {
	//io.Pf("%4d", M[i][j])
	//}
	//io.Pf("\n")
	//}
	return
}

// error messages
var (
	_trapz_err1 = "trapz.go: Trapz: length of x and y must be the same. %d != %d"
	_trapz_err2 = "trapz.go: TrapzRange: number of points must be at least 2"
)
