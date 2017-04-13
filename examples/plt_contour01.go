// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// scalar field
	fcn := func(x, y float64) float64 {
		return -math.Pow(math.Pow(math.Cos(x), 2.0)+math.Pow(math.Cos(y), 2.0), 2.0)
	}

	// gradient. u=dfdx, v=dfdy
	grad := func(x, y float64) (u, v float64) {
		m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
		u = 4.0 * math.Cos(x) * math.Sin(x) * m
		v = 4.0 * math.Cos(y) * math.Sin(y) * m
		return
	}

	// grid size
	xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

	// mesh grid
	X, Y := utl.MeshGrid2D(xmin, xmax, xmin, xmax, N, N)

	// compute f(x,y) and components of gradient
	F := utl.DblsAlloc(N, N)
	U := utl.DblsAlloc(N, N)
	V := utl.DblsAlloc(N, N)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			F[i][j] = fcn(X[i][j], Y[i][j])
			U[i][j], V[i][j] = grad(X[i][j], Y[i][j])
		}
	}

	// plot
	plt.SetForPng(0.75, 600, 150)
	plt.Contour(X, Y, F, "levels=20, cmapidx=4")
	plt.Quiver(X, Y, U, V, "color='red'")
	plt.Gll("x", "y", "")
	plt.Equal()
	plt.SaveD("/tmp/gosl", "plt_contour01.png")
}
