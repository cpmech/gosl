// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// define function and derivative function
	y_fcn := func(x float64) float64 { return math.Sin(x) }
	dydx_fcn := func(x float64) float64 { return math.Cos(x) }
	d2ydx2_fcn := func(x float64) float64 { return -math.Sin(x) }

	// run test for 11 points
	X := utl.LinSpace(0, 2*math.Pi, 11)
	io.Pf("          %8s %23s %23s %23s\n", "x", "analytical", "numerical", "error")
	for _, x := range X {

		// analytical derivatives
		dydx_ana := dydx_fcn(x)
		d2ydx2_ana := d2ydx2_fcn(x)

		// numerical derivative: dydx
		dydx_num, _ := num.DerivCentral(func(t float64, args ...interface{}) float64 {
			return y_fcn(t)
		}, x, 1e-3)

		// numerical derivative d2ydx2
		d2ydx2_num, _ := num.DerivCentral(func(t float64, args ...interface{}) float64 {
			return dydx_fcn(t)
		}, x, 1e-3)

		// check
		chk.PrintAnaNum(io.Sf("dy/dx   @ %.6f", x), 1e-10, dydx_ana, dydx_num, true)
		chk.PrintAnaNum(io.Sf("d²y/dx² @ %.6f", x), 1e-10, d2ydx2_ana, d2ydx2_num, true)
	}

	// generate 101 points for plotting
	X = utl.LinSpace(0, 2*math.Pi, 101)
	Y := make([]float64, len(X))
	for i, x := range X {
		Y[i] = y_fcn(x)
	}

	// plot
	plt.Reset(false, nil)
	plt.Plot(X, Y, &plt.A{C: "b", M: ".", Me: 10, L: "y(x)=sin(x)"})
	plt.Gll("x", "y", nil)
	plt.Save("/tmp/gosl", "num_deriv01")
}
