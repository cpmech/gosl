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

	// check Kronecker property (barycentic)
	o.Bary = true
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

	// check Kronecker property
	o.Bary = false
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
			o.Bary = true
			li1 := o.L(i, x)
			o.Bary = false
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
	m := math.Pow(2, float64(o.N)-1) / float64(o.N)
	for i := 0; i < o.N+1; i++ {
		d := 1.0
		for j := 0; j < o.N+1; j++ {
			if i != j {
				d *= (o.X[i] - o.X[j])
			}
		}
		chk.AnaNum(tst, io.Sf("λ%d", i), tol, o.Lam(i), 1.0/d/m, chk.Verbose)
	}
}

func checkIandLam(tst *testing.T, N int, tolLam float64, f Ss) {

	// allocate structure and calculate U
	o, err := NewLagrangeInterp(N, ChebyGaussLobGridKind)
	chk.EP(err)
	err = o.CalcU(f)
	chk.EP(err)

	// check interpolation (std)
	o.Bary = false
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// check λ
	checkLam(tst, o, tolLam)

	// check interpolation (barycentric)
	o.Bary = true
	for i, x := range o.X {
		ynum, err := o.I(x, f)
		chk.EP(err)
		yana, _ := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// compare std and barycentric
	xx := utl.LinSpace(-1, 1, 14)
	for _, x := range xx {
		for i := 0; i < o.N+1; i++ {
			o.Bary = false
			i1, err := o.I(x, f)
			chk.EP(err)
			o.Bary = true
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
	tolsL := []float64{1e-15, 1e-15, 1e-15, 1e-15, 1e-14, 1e-14}
	for k, N := range Nvals {
		io.Pf("\n\n-------------------------------- N = %d -----------------------------------------------\n\n", N)
		checkIandLam(tst, N, tolsL[k], f)
	}
}

func checkD1lag(tst *testing.T, N int, tol float64) {

	// allocate structure
	o, err := NewLagrangeInterp(N, ChebyGaussLobGridKind)
	chk.EP(err)

	// calc and check D1
	err = o.CalcD1()
	chk.EP(err)
	for j := 0; j < N+1; j++ {
		xj := o.X[j]
		for l := 0; l < N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tol, o.D1.Get(j, l), xj, 1e-2, chk.Verbose, func(t float64) (float64, error) {
				return o.L(l, t), nil
			})
		}
	}
}

func TestLagInterp06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp06. D1")

	Nvals := []int{3, 4, 5, 6, 7, 8}
	tols := []float64{1e-10, 1e-9, 1e-9, 1e-9, 1e-9, 1e-8}
	for k, N := range Nvals {
		io.Pf("\n\n-------------------------------- N = %d -----------------------------------------------\n\n", N)
		checkD1lag(tst, N, tols[k])
	}
}
