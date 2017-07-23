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
	"github.com/cpmech/gosl/utl"
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
	chk.Float64(tst, "ΛN (Lebesgue constant)", 1e-15, o.EstimateLebesgue(), 3.106301040275436e+00)

	// check Kronecker property
	for i := 0; i < N+1; i++ {
		for j, x := range o.X {
			li := o.L(i, x)
			ana := 1.0
			if i != j {
				ana = 0
			}
			chk.AnaNum(tst, io.Sf("L^%d_%d(X[%d])", N, i, j), 1e-17, li, ana, false)
		}
	}

	// check Kronecker property (barycentic, notLog)
	o.CalcBary(false)
	for i := 0; i < N+1; i++ {
		for j, x := range o.X {
			li := o.L(i, x)
			ana := 1.0
			if i != j {
				ana = 0
			}
			chk.AnaNum(tst, io.Sf("L^%d_%d(X[%d])", N, i, j), 1e-17, li, ana, false)
		}
	}

	// compare formulae
	xx := utl.LinSpace(-1, 1, 11)
	for _, x := range xx {
		for i := 0; i < N+1; i++ {
			o.UseBary = true
			li1 := o.L(i, x)
			o.UseBary = false
			li2 := o.L(i, x)
			chk.AnaNum(tst, io.Sf("l%d", i), 1e-15, li1, li2, chk.Verbose)
		}
	}

	// check Kronecker property (barycentic, log)
	o.CalcBary(true)
	for i := 0; i < N+1; i++ {
		for j, x := range o.X {
			li := o.L(i, x)
			ana := 1.0
			if i != j {
				ana = 0
			}
			chk.AnaNum(tst, io.Sf("L^%d_%d(X[%d])", N, i, j), 1e-17, li, ana, false)
		}
	}

	// compare formulae
	for _, x := range xx {
		for i := 0; i < N+1; i++ {
			o.UseBary = true
			li1 := o.L(i, x)
			o.UseBary = false
			li2 := o.L(i, x)
			chk.AnaNum(tst, io.Sf("l%d", i), 1e-15, li1, li2, chk.Verbose)
		}
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

	// allocate structure and calculate U
	N := 5
	kind := UniformGridKind
	o, err := NewLagrangeInterp(N, kind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

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

	// allocate structure and calculate U
	N := 8
	kind := UniformGridKind
	o, err := NewLagrangeInterp(N, kind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

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

	// allocate structure and calculate U
	N := 8
	kind := ChebyGaussGridKind
	o, err := NewLagrangeInterp(N, kind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// check Lebesgue constants and compute max error
	ΛN := []float64{1.988854381999833e+00, 2.361856787767076e+00, 3.011792612349363e+00}
	for i, n := range []int{4, 8, 24} {
		p, err := NewLagrangeInterp(n, kind)
		chk.EP(err)
		chk.Float64(tst, "ΛN (Lebesgue constant)", 1e-13, p.EstimateLebesgue(), ΛN[i])
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
			p, err := NewLagrangeInterp(int(n), kind)
			chk.EP(err)
			err = p.CalcU(f)
			chk.EP(err)
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

func TestLagInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp04. Chebyshev-Gauss-Lobatto. Runge problem")

	// Runge function
	f := func(x float64) (float64, error) {
		return 1.0 / (1.0 + 16.0*x*x), nil
	}

	// allocate structure and calculate U
	N := 8
	kind := ChebyGaussLobGridKind
	o, err := NewLagrangeInterp(N, kind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// check Lebesgue constants and compute max error
	ΛN := []float64{1.798761778849085e+00, 2.274730699116020e+00, 2.984443326362511e+00}
	for i, n := range []int{4, 8, 24} {
		p, err := NewLagrangeInterp(n, kind)
		chk.EP(err)
		chk.Float64(tst, "ΛN (Lebesgue constant)", 1e-14, p.EstimateLebesgue(), ΛN[i])
	}

	if chk.Verbose {

		// plot nodal polynomial
		plt.Reset(true, nil)
		PlotLagInterpW(8, kind)
		plt.AxisYrange(-0.02, 0.02)
		plt.Save("/tmp/gosl/fun", "laginterp04a")

		// plot interpolation
		plt.Reset(true, nil)
		PlotLagInterpI([]int{4, 6, 8, 12, 16, 24}, kind, f)
		plt.AxisYrange(-1, 1)
		plt.Save("/tmp/gosl/fun", "laginterp04b")

		// plot error
		plt.Reset(true, nil)
		Nvalues := []float64{1, 4, 8, 16, 24, 40, 80, 100, 120, 140, 200}
		E := make([]float64, len(Nvalues))
		for i, n := range Nvalues {
			p, err := NewLagrangeInterp(int(n), kind)
			chk.EP(err)
			err = p.CalcU(f)
			chk.EP(err)
			E[i], _ = p.EstimateMaxErr(0, f)
		}
		plt.Plot(Nvalues, E, &plt.A{C: "red", M: ".", NoClip: true})
		plt.Grid(nil)
		plt.Gll("$N$", "$\\max[|f(x) - I^X_N\\{f\\}(x)|]$", nil)
		plt.HideTRborders()
		plt.SetYlog()
		plt.Save("/tmp/gosl/fun", "laginterp04c")
	}
}

func checkLam(tst *testing.T, o *LagrangeInterp, tol float64) {
	λ := make([]float64, o.N+1)
	for i := 0; i < o.N+1; i++ {
		λ[i] = 1
		for j := 0; j < o.N+1; j++ {
			if i != j {
				λ[i] *= (o.X[i] - o.X[j])
			}
		}
	}
	for i := 0; i < o.N+1; i++ {
		λ[i] = 1.0 / λ[i]
	}
	//if o.N == 4 {
	//λ = []float64{1, -2, 2, -2, 1} // N=4
	//}
	//if o.N == 6 {
	//d := 16.0 / 3.0
	//λ = []float64{8.0 / 3.0, -d, d, -d, d, -d, 8.0 / 3.0}
	//}
	chk.Array(tst, "λ", tol, o.Lam, λ)
}

func checkIandLam(tst *testing.T, N int, tolLam float64, f Ss) {

	// allocate structure and calculate U
	o, err := NewLagrangeInterp(N, ChebyGaussLobGridKind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

	// check interpolation
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// -------------- barycentric: std method ---------------------------

	// check λ (std method)
	io.Pl()
	o.CalcBary(false)
	io.Pforan("λ[std] = %v\n", o.Lam)
	checkLam(tst, o, 1e-17)

	// check interpolation
	io.Pl()
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// compare formulae
	io.Pl()
	xx := utl.LinSpace(-1, 1, 14)
	for _, x := range xx {
		for i := 0; i < o.N+1; i++ {
			o.UseBary = true
			i1, err := o.I(x, f)
			chk.EP(err)
			o.UseBary = false
			i2, err := o.I(x, f)
			chk.EP(err)
			chk.AnaNum(tst, io.Sf("I%d", i), 1e-15, i1, i2, false)
		}
	}

	// -------------- barycentric: log method ---------------------------

	// check λ (log method)
	io.Pl()
	o.CalcBary(true)
	io.Pforan("λ[log] = %v\n", o.Lam)
	checkLam(tst, o, tolLam)

	// check interpolation
	io.Pl()
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// compare formulae
	io.Pl()
	for _, x := range xx {

		for i := 0; i < o.N+1; i++ {
			o.UseBary = true
			i1, err := o.I(x, f)
			chk.EP(err)
			o.UseBary = false
			i2, err := o.I(x, f)
			chk.EP(err)
			chk.AnaNum(tst, io.Sf("I%d", i), 1e-15, i1, i2, false)
		}
	}
}

func TestLagInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp05. Barycentric formulae")

	// Runge function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// test
	Nvals := []int{3, 4, 5, 6, 7, 8}
	tols := []float64{1e-15, 1e-15, 1e-15, 1e-15, 1e-14, 1e-14}
	for k, N := range Nvals {
		io.Pf("\n\n-------------------------------- N = %d -----------------------------------------------\n\n", N)
		checkIandLam(tst, N, tols[k], f)
	}
}
