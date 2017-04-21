// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/fdm"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// create grid
	var g fdm.Grid2d
	g.Init(2.0, 14.0, 2.0, 8.0, 21, 11)

	// callback function
	fxy := func(x, y float64) float64 { return x*x + y*y }

	// generate data
	X, Y, F := g.Generate(fxy, nil)

	// clear figure, apply default configuration values,
	// and set height/width proportion to 0.5
	plt.Reset(true, &plt.A{Prop: 0.5})

	// draw contour
	plt.ContourF(X, Y, F, nil)

	// setup axes and save figure as PNG file
	plt.Equal()
	plt.Gll("x", "y", nil)
	plt.Save("/tmp/gosl", "fdm_grid2d")
}
