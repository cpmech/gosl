// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// Runge function
	f := func(x float64) (float64, error) {
		return 1.0 / (1.0 + 16.0*x*x), nil
	}

	// allocate interpolator structure
	N := 8
	kind := fun.ChebyGaussGridKind
	lag, err := fun.NewLagrangeInterp(N, kind)
	if err != nil {
		io.Pf("ERROR: %v\n", err)
		return
	}

	// interpolate one point
	res, err := lag.I(0.0, f)
	if err != nil {
		io.Pf("ERROR: %v\n", err)
		return
	}
	io.Pf("I{f}(x=0) = %v\n", res)

	// plot nodal polynomial
	plt.Reset(true, nil)
	fun.PlotLagInterpW(N, kind)
	plt.AxisYrange(-0.02, 0.02)
	plt.Save("/tmp/gosl", "fun_laginterp01a")

	// plot interpolation
	plt.Reset(true, nil)
	fun.PlotLagInterpI([]int{4, 6, 8, 12, 16, 24}, kind, f)
	plt.AxisYrange(-1, 1)
	plt.Save("/tmp/gosl", "fun_laginterp01b")

	// plot error
	plt.Reset(true, nil)
	Nvalues := []float64{1, 4, 8, 16, 24, 40, 80, 100, 120, 140, 200}
	E := make([]float64, len(Nvalues))
	for i, n := range Nvalues {
		p, _ := fun.NewLagrangeInterp(int(n), kind)
		E[i], _ = p.EstimateMaxErr(10000, f)
	}
	plt.Plot(Nvalues, E, &plt.A{C: "red", M: ".", NoClip: true})
	plt.Grid(nil)
	plt.Gll("$N$", "$\\max[|f(x) - I^X_N\\{f\\}(x)|]$", nil)
	plt.HideTRborders()
	plt.SetYlog()
	plt.Save("/tmp/gosl", "fun_laginterp01c")
}
