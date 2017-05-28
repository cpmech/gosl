// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// QuadDiscreteSimpsonRF approximates the area below the discrete curve defined by [xa,xy] range and
// y function. Computations are carried out with the (very simple) Simpson method from xa to xb,
// with npts points
func QuadDiscreteSimpsonRF(a, b float64, n int, f fun.Ss) (res float64, err error) {
	if n < 2 || n%2 > 0 {
		err = chk.Err("number of subintervas should be even (n=%d)", n)
		return
	}
	fa, err := f(a)
	if err != nil {
		return
	}
	fb, err := f(b)
	if err != nil {
		return
	}
	var fx float64
	x := a
	h := (b - a) / float64(n)
	sum := fa + fb
	for i := 1; i < n; i++ {
		x += h
		fx, err = f(x)
		if err != nil {
			return
		}
		if i%2 == 1 { // i is odd
			sum += 4 * fx
		} else { // i is even
			sum += 2 * fx
		}
	}
	res = sum * h / 3.0
	return
}

// QuadDiscreteSimps2d approximates a double integral over the x-y plane with the elevation given by
// data points f[npts][npts]. Thus, the result is an estimate of the volume below the f[][] opints
// and the plane ortogonal to z @ x=0. The very simple Simpson's method is used here.
//  Lx -- total length of plane along x
//  Ly -- total length of plane along y
//  f  -- elevations f(x,y)
func QuadDiscreteSimps2d(Lx, Ly float64, f [][]float64) (V float64) {

	// check
	if len(f) < 2 {
		chk.Panic("len(f)=%d is incorrect; it must be at least 2", len(f))
	}
	m, n := len(f), len(f[0])

	// corners
	V = f[0][0] + f[m-1][0] + f[0][n-1] + f[m-1][n-1]

	// top/bottom: 4
	for j := 1; j < n-1; j += 2 {
		V += 4.0 * (f[0][j] + f[m-1][j])
	}

	// top/bottom: 2
	for j := 2; j < n-1; j += 2 {
		V += 2.0 * (f[0][j] + f[m-1][j])
	}

	// left/right: 4
	for i := 1; i < m-1; i += 2 {
		V += 4.0 * (f[i][0] + f[i][n-1])

		// centre: 16
		for j := 1; j < n-1; j += 2 {
			V += 16.0 * f[i][j]
		}

		// centre: 8a
		for j := 2; j < n-1; j += 2 {
			V += 8.0 * f[i][j]
		}
	}

	// left/right: 2
	for i := 2; i < m-1; i += 2 {
		V += 2.0 * (f[i][0] + f[i][n-1])

		// centre: 4
		for j := 2; j < n-1; j += 2 {
			V += 4.0 * f[i][j]
		}

		// centre: 8b
		for j := 1; j < n-1; j += 2 {
			V += 8.0 * f[i][j]
		}
	}

	// final result
	V *= Lx * Ly / 9.0
	return
}
