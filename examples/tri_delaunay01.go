// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm/tri"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	// fix seed
	rnd.Init(1358)

	// generate cloud of points
	nx, ny := 6, 6
	dx := 1.0 / float64(nx-1)
	dy := 1.0 / float64(ny-1)
	X := make([]float64, nx*ny)
	Y := make([]float64, nx*ny)
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			n := i + j*nx
			X[n] = float64(i) * (dx * rnd.Float64(0.5, 1.0))
			Y[n] = float64(j) * (dy * rnd.Float64(0.5, 1.0))
		}
	}

	// generate
	V, C := tri.Delaunay(X, Y, false)

	// plot
	plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
	tri.DrawVC(V, C, &plt.A{C: "orange", Ls: "-", NoClip: true})
	plt.Plot(X, Y, &plt.A{C: "k", Ls: "none", M: ".", NoClip: true})
	plt.Gll("x", "y", nil)
	plt.Equal()
	plt.HideAllBorders()
	plt.Save("/tmp/gosl", "tri_delaunay01")
}
