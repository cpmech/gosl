// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// QuadDiscreteTrapzXY approximates the area below the discrete curve defined by x and y points.
// Computations are carried out with the trapezoidal rule.
func QuadDiscreteTrapzXY(x, y []float64) (A float64) {
	if len(x) != len(y) {
		chk.Panic("length of x and y must be the same. %d != %d", len(x), len(y))
	}
	for i := 1; i < len(x); i++ {
		A += (x[i] - x[i-1]) * (y[i] + y[i-1]) / 2.0
	}
	return
}

// QuadDiscreteTrapzXF approximates the area below the discrete curve defined by x points and y
// function. Computations are carried out with the (very simple) trapezoidal rule.
func QuadDiscreteTrapzXF(x []float64, y fun.Ss) (A float64, err error) {
	var ya, yb float64
	for i := 1; i < len(x); i++ {
		ya, err = y(x[i-1])
		if err != nil {
			return
		}
		yb, err = y(x[i])
		if err != nil {
			return
		}
		A += (x[i] - x[i-1]) * (yb + ya) / 2.0
	}
	return
}

// QuadDiscreteTrapzRF approximates the area below the discrete curve defined by [xa,xy] range and y
// function. Computations are carried out with the (very simple) trapezoidal rule from xa to xb,
// with npts points
func QuadDiscreteTrapzRF(xa, xb float64, npts int, y fun.Ss) (A float64, err error) {
	if npts < 2 {
		chk.Panic("number of points must be at least 2", npts)
	}
	dx := (xb - xa) / float64(npts-1)
	var x0, x1, y0, y1 float64
	for i := 1; i < npts; i++ {
		x0 = xa + dx*float64(i-1)
		x1 = xa + dx*float64(i)
		y0, err = y(x0)
		if err != nil {
			return
		}
		y1, err = y(x1)
		if err != nil {
			return
		}
		A += (x1 - x0) * (y1 + y0) / 2.0
	}
	return
}

// QuadDiscreteTrapz2d approximates a double integral over the x-y plane with the elevation given by
// data points f[npts][npts]. Thus, the result is an estimate of the volume below the f[][] opints
// and the plane ortogonal to z @ x=0. The very simple trapezoidal method is used here.
//  Lx -- total length of plane along x
//  Ly -- total length of plane along y
//  f  -- elevations f(x,y)
func QuadDiscreteTrapz2d(Lx, Ly float64, f [][]float64) (V float64) {

	// check
	if len(f) < 2 {
		chk.Panic("len(f)=%d is incorrect; it must be at least 2", len(f))
	}
	m, n := len(f), len(f[0])

	// corners
	V = f[0][0] + f[m-1][0] + f[0][n-1] + f[m-1][n-1]

	// top/bottom: 2
	for j := 1; j < n-1; j++ {
		V += 2.0 * (f[0][j] + f[m-1][j])
	}

	// left/right: 2
	for i := 1; i < m-1; i++ {
		V += 2.0 * (f[i][0] + f[i][n-1])

		// centre: 4
		for j := 1; j < n-1; j++ {
			V += 4.0 * f[i][j]
		}
	}

	// final result
	V *= Lx * Ly / 4.0
	return
}
