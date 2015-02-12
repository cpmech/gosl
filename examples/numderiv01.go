// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/num"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// define function and derivative function
	y_fcn := func(x float64) float64 { return math.Sin(x) }
	dydx_fcn := func(x float64) float64 { return math.Cos(x) }

	// run test for 11 points
	X := utl.LinSpace(0, 2*math.Pi, 11)
	for _, x := range X {

		// analytical derivative
		dydx_ana := dydx_fcn(x)

		// numerical derivative
		dydx_num, _ := num.DerivCentral(func(t float64, args ...interface{}) float64 {
			return y_fcn(t)
		}, x, 1e-3)

		// check
		utl.AnaNum(utl.Sf("dydx @ %.6f", x), 1e-10, dydx_ana, dydx_num, true)
	}

	// generate 101 points
	X = utl.LinSpace(0, 2*math.Pi, 101)
	Y := make([]float64, len(X))
	for i, x := range X {
		Y[i] = y_fcn(x)
	}

	// plot
	plt.Plot(X, Y, "'b.-'")
	plt.Gll("x", "y", "")
	plt.Show()
}
