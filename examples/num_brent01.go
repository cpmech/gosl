// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// y(x) function
	yx := func(x float64) float64 {
		return math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
	}

	// range: be sure to enclose root
	xa, xb := 0.0, 0.11

	// initialise solver
	solver := num.NewBrent(yx, nil)

	// solve
	xo := solver.Root(xa, xb)

	// output
	yo := yx(xo)
	io.Pf("\n")
	io.Pf("x      = %v\n", xo)
	io.Pf("f(x)   = %v\n", yo)
	io.Pf("nfeval = %v\n", solver.NumFeval)
	io.Pf("niter. = %v\n", solver.NumIter)

	// plotting
	npts := 101
	X := make([]float64, npts)
	Y := make([]float64, npts)
	for i := 0; i < npts; i++ {
		X[i] = xa + float64(i)*(xb-xa)/float64(npts-1)
		Y[i] = yx(X[i])
	}
	plt.Reset(false, nil)
	plt.AxHline(0, nil)
	plt.Plot(X, Y, &plt.A{C: "g"})
	plt.PlotOne(xo, yo, &plt.A{C: "r", M: "."})
	plt.Gll("x", "y(x)", nil)
	plt.Save("/tmp/gosl", "num_brent01")
}
