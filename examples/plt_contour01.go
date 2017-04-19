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

	// grid size
	xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

	// mesh grid, scalar and vector field
	X, Y, F, U, V := utl.MeshGrid2dFG(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (f, u, v float64) {

		// scalar field
		m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
		f = -math.Pow(m, 2.0)

		// gradient. u=dfdx, v=dfdy
		u = 4.0 * math.Cos(x) * math.Sin(x) * m
		v = 4.0 * math.Cos(y) * math.Sin(y) * m
		return
	})

	// plot
	plt.Reset(false, nil)
	plt.ContourF(X, Y, F, &plt.A{CmapIdx: 4, Nlevels: 15})
	plt.Quiver(X, Y, U, V, &plt.A{C: "r"})
	plt.Gll("x", "y", nil)
	plt.Equal()
	plt.Save("/tmp/gosl", "plt_contour01")
}
