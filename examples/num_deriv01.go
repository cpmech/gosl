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
	yFcn := func(x float64) float64 { return math.Sin(x) }
	dydxFcn := func(x float64) float64 { return math.Cos(x) }
	d2ydx2Fcn := func(x float64) float64 { return -math.Sin(x) }

	// run test for 11 points
	X := utl.LinSpace(0, 2*math.Pi, 11)
	io.Pf("          %8s %23s %23s %23s\n", "x", "analytical", "numerical", "error")
	for _, x := range X {

		// analytical derivatives
		dydxAna := dydxFcn(x)
		d2ydx2Ana := d2ydx2Fcn(x)

		// numerical derivative: dydx
		dydxNum, _ := num.DerivCen5(x, 1e-3, func(t float64) (float64, error) {
			return yFcn(t), nil
		})

		// numerical derivative d2ydx2
		d2ydx2Num, _ := num.DerivCen5(x, 1e-3, func(t float64) (float64, error) {
			return dydxFcn(t), nil
		})

		// check
		chk.PrintAnaNum(io.Sf("dy/dx   @ %.6f", x), 1e-10, dydxAna, dydxNum, true)
		chk.PrintAnaNum(io.Sf("d²y/dx² @ %.6f", x), 1e-10, d2ydx2Ana, d2ydx2Num, true)
	}

	// generate 101 points for plotting
	X = utl.LinSpace(0, 2*math.Pi, 101)
	Y := make([]float64, len(X))
	for i, x := range X {
		Y[i] = yFcn(x)
	}

	// plot
	plt.Reset(false, nil)
	plt.Plot(X, Y, &plt.A{C: "b", M: ".", Me: 10, L: "y(x)=sin(x)"})
	plt.Gll("x", "y", nil)
	plt.Save("/tmp/gosl", "num_deriv01")
}
