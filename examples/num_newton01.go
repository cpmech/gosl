// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// Function: y(x) = fx[0] with x = xvec[0]
	fcn := func(fx, xvec la.Vector) {
		x := xvec[0]
		fx[0] = math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
	}

	// Jacobian: dfdx(x) function
	Jfcn := func(dfdx *la.Matrix, x la.Vector) {
		dfdx.Set(0, 0, 3.0*x[0]*x[0]-2.0*0.165*x[0])
		return
	}

	// trial solution
	xguess := 0.03

	// initialise solver
	neq := 1      // number of equations
	useDn := true // use dense Jacobian
	numJ := false // numerical Jacobian
	var o num.NlSolver
	o.Init(neq, fcn, nil, Jfcn, useDn, numJ, nil)

	// solve
	xvec := []float64{xguess}
	o.Solve(xvec, false)

	// output
	fx := []float64{123}
	fcn(fx, xvec)
	xo, yo := xvec[0], fx[0]
	io.Pf("\n")
	io.Pf("x      = %v\n", xo)
	io.Pf("f(x)   = %v\n", yo)
	io.Pf("nfeval = %v\n", o.NFeval)
	io.Pf("niter. = %v\n", o.It)

	// plotting
	xa, xb := 0.0, 0.11
	npts := 101
	X := make([]float64, npts)
	Y := make([]float64, npts)
	for i := 0; i < npts; i++ {
		xvec[0] = xa + float64(i)*(xb-xa)/float64(npts-1)
		fcn(fx, xvec)
		X[i] = xvec[0]
		Y[i] = fx[0]
	}
	plt.Reset(true, nil)
	plt.AxHline(0, &plt.A{C: "k"})
	plt.Plot(X, Y, nil)
	plt.PlotOne(xo, yo, &plt.A{C: "r", M: "."})
	plt.HideTRborders()
	plt.Gll("x", "y(x)", nil)
	plt.Save("/tmp/gosl", "num_newton01")
}
