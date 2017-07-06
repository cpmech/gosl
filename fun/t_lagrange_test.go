// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func TestLagCardinal01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCardinal01. Lagrange cardinal polynomials")

	// allocate structure
	N := 5
	kind := UniformGridKind
	o, err := NewLagrangeInterp(N, kind)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	chk.Scalar(tst, "ΛN (Lebesgue constant)", 1e-15, o.EstimateLebesgue(), 3.106301040275436e+00)

	// check Kronecker property
	for i := 0; i < N+1; i++ {
		for j, x := range o.X {
			li := o.L(i, x)
			ana := 1.0
			if i != j {
				ana = 0
			}
			chk.AnaNum(tst, io.Sf("L^%d_%d(X[%d])", N, i, j), 1e-17, li, ana, chk.Verbose)
		}
		io.Pl()
	}

	// plot basis
	if chk.Verbose {
		plt.Reset(true, nil)
		PlotLagInterpL(N, kind)
		plt.Save("/tmp/gosl/fun", "lagcardinal01")
	}
}

func TestLagInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp01. Lagrange interpolation")

	// cos-exp function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate structure
	N := 5
	kind := UniformGridKind
	o, err := NewLagrangeInterp(N, kind)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		if err != nil {
			tst.Errorf("%v\n", err)
			return
		}
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// plot interpolation
	if chk.Verbose {
		plt.Reset(true, nil)
		PlotLagInterpI([]int{4, 6, 8, 12, 16, 24}, kind, f)
		plt.Save("/tmp/gosl/fun", "laginterp01")
	}
}

func TestLagInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp02. Lagrange interp. Runge problem")

	// Runge function
	f := func(x float64) (float64, error) {
		return 1.0 / (1.0 + 16.0*x*x), nil
	}

	// allocate structure
	N := 8
	kind := UniformGridKind
	o, err := NewLagrangeInterp(N, kind)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		if err != nil {
			tst.Errorf("%v\n", err)
			return
		}
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	if chk.Verbose {

		// plot nodal polynomial
		plt.Reset(true, nil)
		PlotLagInterpW(8, kind)
		plt.AxisYrange(-0.02, 0.02)
		plt.Save("/tmp/gosl/fun", "laginterp02a")

		// plot interpolation
		plt.Reset(true, nil)
		PlotLagInterpI([]int{4, 6, 8, 12, 16, 24}, kind, f)
		plt.AxisYrange(-1, 1)
		plt.Save("/tmp/gosl/fun", "laginterp02b")
	}
}

func TestLagInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp03. Chebyshev-Gauss. Runge problem")

	// Runge function
	f := func(x float64) (float64, error) {
		return 1.0 / (1.0 + 16.0*x*x), nil
	}

	// allocate structure
	N := 8
	kind := ChebyGaussGridKind
	o, err := NewLagrangeInterp(N, kind)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		if err != nil {
			tst.Errorf("%v\n", err)
			return
		}
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// check Lebesgue constants and compute max error
	ΛN := []float64{1.988854381999833e+00, 2.361856787767076e+00, 3.011792612349363e+00}
	for i, n := range []int{4, 8, 24} {
		p, _ := NewLagrangeInterp(n, kind)
		chk.Scalar(tst, "ΛN (Lebesgue constant)", 1e-13, p.EstimateLebesgue(), ΛN[i])
	}

	if chk.Verbose {

		// plot nodal polynomial
		plt.Reset(true, nil)
		PlotLagInterpW(8, kind)
		plt.AxisYrange(-0.02, 0.02)
		plt.Save("/tmp/gosl/fun", "laginterp03a")

		// plot interpolation
		plt.Reset(true, nil)
		PlotLagInterpI([]int{4, 6, 8, 12, 16, 24}, kind, f)
		plt.AxisYrange(-1, 1)
		plt.Save("/tmp/gosl/fun", "laginterp03b")

		// plot error
		plt.Reset(true, nil)
		Nvalues := []float64{1, 4, 8, 16, 24, 40, 80, 100, 120, 140, 200}
		E := make([]float64, len(Nvalues))
		for i, n := range Nvalues {
			p, _ := NewLagrangeInterp(int(n), kind)
			E[i], _ = p.EstimateMaxErr(0, f)
		}
		plt.Plot(Nvalues, E, &plt.A{C: "red", M: ".", NoClip: true})
		plt.Grid(nil)
		plt.Gll("$N$", "$\\max[|f(x) - I^X_N\\{f\\}(x)|]$", nil)
		plt.HideTRborders()
		plt.SetYlog()
		plt.Save("/tmp/gosl/fun", "laginterp03c")
	}
}
