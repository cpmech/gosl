// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	// initialise generator
	rnd.Init(1234)

	// Halton points
	ndim := 2
	npts := 100
	Xhps := rnd.HaltonPoints(ndim, npts)

	// Latin Hypercube
	dupfactor := 5
	lhs := rnd.LatinIHS(ndim, npts, dupfactor)
	xmin := []float64{0, 0}
	xmax := []float64{1, 1}
	Xlhs := rnd.HypercubeCoords(lhs, xmin, xmax)

	// plot
	plt.Reset(true, &plt.A{WidthPt: 300})
	plt.Plot(Xhps[0], Xhps[1], &plt.A{C: "b", M: ".", Ls: "none", L: "Halton"})
	plt.Plot(Xlhs[0], Xlhs[1], &plt.A{C: "r", M: "o", Ls: "none", L: "LHS", Void: true})
	plt.Equal()
	plt.Gll("$x_0$", "$x_1$", &plt.A{LegOut: true, LegNcol: 2})
	plt.Save("/tmp/gosl", "rnd_haltonAndLatin01")
}
