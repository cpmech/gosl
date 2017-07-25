// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestChebyInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp01")

	// test function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate polynomials
	N := 8
	che, err := NewChebyInterp(N, true) // Gauss-Chebyshev
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}
	lob, err := NewChebyInterp(N, false) // Gauss-Lobatto
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	// check P
	for _, x := range utl.LinSpace(-1, 1, 7) {
		chk.AnaNum(tst, io.Sf("T%d(%+.3f)", N, x), 1e-15, ChebyshevT(N, x), che.HierarchicalT(8, x), chk.Verbose)
	}

	// Gauss-Chebyshev: check points
	xref := []float64{
		9.848077530122080e-01,
		8.660254037844387e-01,
		6.427876096865394e-01,
		3.420201433256687e-01,
		0,
		-3.420201433256688e-01,
		-6.427876096865394e-01,
		-8.660254037844387e-01,
		-9.848077530122080e-01,
	}
	chk.Array(tst, "Gauss-Chebyshev: X", 1e-15, che.X, xref)

	// Gauss-Chebyshev: check coefficients of interpolant
	che.CalcCoefI(f)
	cref := []float64{5.005025576289825e-01, -4.734690106554930e-01, 3.343030345866715e-01, 5.329760324967350e-01, 2.005496385333029e-01, -1.552357980491117e-01, -2.768837833165416e-01, -2.160862487215637e-01, -1.033306390240169e-01}
	chk.Array(tst, "Gauss-Chebyshev: CoefI", 1e-14, che.CoefI, cref)

	// Gauss-Chebyshev: check coefficients of projection
	che.CalcCoefP(f)
	cref = []float64{5.003559557885667e-01, -4.738396675676836e-01, 3.337904287575258e-01, 5.326202849023425e-01, 2.014887911962803e-01, -1.505413304349933e-01, -2.650525046501985e-01, -1.959021686279372e-01, -8.320914768336027e-02}
	chk.Array(tst, "Gauss-Chebyshev: CoefP", 1e-15, che.CoefP, cref)

	// Gauss-Lobatto: check points
	xref = []float64{
		1.000000000000000e+00,
		9.238795325112867e-01,
		7.071067811865475e-01,
		3.826834323650897e-01,
		0,
		-3.826834323650898e-01,
		-7.071067811865476e-01,
		-9.238795325112867e-01,
		-1.000000000000000e+00,
	}
	chk.Array(tst, "Gauss-Lobatto: X", 1e-15, lob.X, xref)

	// Gauss-Lobatto: check coefficients of interpolant
	lob.CalcCoefI(f)
	cref = []float64{4.998505262591759e-01, -4.745223967909372e-01, 3.345788625609180e-01, 5.372649935983991e-01, 2.133118551317398e-01, -1.303539051589940e-01, -2.449269176317319e-01, -2.036386455562469e-01, -8.320813059359007e-02}
	chk.Array(tst, "Gauss-Lobatto: CoefI", 1e-15, lob.CoefI, cref)

	// Gauss-Lobatto: check coefficients of projection
	lob.CalcCoefP(f)
	chk.Array(tst, "Gauss-Lobatto: CoefP", 1e-15, lob.CoefP, che.CoefP)

	// Gauss-Chebyshev: estimate error
	Eproj, _ := che.EstimateMaxErr(f, true)
	Einte, _ := che.EstimateMaxErr(f, false)
	io.Pforan("Gauss-Chebyshev: E{proj} = %v\n", Eproj)
	io.Pforan("Gauss-Chebyshev: E{inte} = %v\n", Einte)

	// Gauss-Lobatto: estimate error
	Eproj, _ = lob.EstimateMaxErr(f, true)
	Einte, _ = lob.EstimateMaxErr(f, false)
	io.Pforan("Gauss-Lobatto: E{proj} = %v\n", Eproj)
	io.Pforan("Gauss-Lobatto: E{inte} = %v\n", Einte)

	if chk.Verbose {

		// plot projection and interpolation
		X := utl.LinSpace(-1, 1, 201)
		Y := make([]float64, len(X))
		Yinte := make([]float64, len(X))
		Yproj := make([]float64, len(X))
		for i, x := range X {
			Y[i], _ = f(x)
			Yinte[i] = lob.I(x)
			Yproj[i] = lob.P(x)
		}
		plt.Reset(true, nil)
		plt.Plot(X, Y, &plt.A{C: "r", L: "$f$", NoClip: true})
		plt.Plot(X, Yinte, &plt.A{C: "g", L: "$I_N^{GL}f$", NoClip: true})
		plt.Plot(X, Yproj, &plt.A{C: "b", L: "$\\Pi_N^{w}f$", NoClip: true})
		plt.Gll("$x$", "$f(x)$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "chebyinterp01a")

		// plot error
		Nvalues := []float64{1, 8, 16, 24, 36, 40, 48, 60, 80, 100, 120}
		Yerr := make([]float64, len(Nvalues))
		for i, nn := range Nvalues {
			o, _ := NewChebyInterp(int(nn), false)
			o.CalcCoefP(f)
			Yerr[i], _ = o.EstimateMaxErr(f, true)
		}
		plt.Reset(true, nil)
		plt.Plot(Nvalues, Yerr, &plt.A{C: "r", M: "o", Void: true, NoClip: true})
		plt.SetYlog()
		plt.Gll("$N$", "$||f-\\Pi_N\\{f\\}||$", nil)
		plt.Save("/tmp/gosl/fun", "chebyinterp01b")
	}
}

func TestChebyInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp02")

	// test function
	f := func(x float64) (float64, error) {
		//return 1.0 / (1.0 + 4.0*x*x), nil
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate polynomials
	N := 6
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)

	// compute data
	err = o.CalcCoefI(f)
	chk.EP(err)

	// check interpolation @ nodes
	for k, x := range o.X {
		fk, err := f(x)
		chk.EP(err)
		chk.Float64(tst, io.Sf("I(x_%d)", k), 1e-14, o.I(x), fk)
	}

	// check conversion
	o.CalcConvMats()
	np := len(o.X) // number of points
	u := la.NewVector(np)
	for i, x := range o.X {
		u[i], err = f(x)
		chk.EP(err)
	}
	ub := la.NewVector(np)
	la.MatVecMul(ub, 1, o.C, u)
	io.Pf("u  = %.6f\n", u)
	io.Pfyel("ub = %.6f\n", ub)
	io.Pfyel("cf = %.6f\n", o.CoefI)
	chk.Array(tst, "ub", 1e-15, ub, o.CoefI)

	// check inversion
	uu := la.NewVector(np)
	la.MatVecMul(uu, 1, o.Ci, ub)
	io.Pf("uu = %.6f\n", uu)
	chk.Array(tst, "uu", 1e-15, uu, u)

	// plot
	if chk.Verbose {
		y0 := make([]float64, len(o.X))
		for i := range o.X {
			y0[i] = 0.2
		}
		xx := utl.LinSpace(-1, 1, 201)
		y1 := make([]float64, len(xx))
		y2 := make([]float64, len(xx))
		for i, x := range xx {
			y1[i], _ = f(x)
			y2[i] = o.I(x)
		}
		plt.Reset(true, nil)
		plt.Plot(o.X, y0, &plt.A{C: "k", M: "o", Ls: "none", Void: true, NoClip: true})
		plt.Plot(xx, y1, &plt.A{C: plt.C(0, 1), L: "$f$", NoClip: true})
		plt.Plot(xx, y2, &plt.A{C: plt.C(1, 1), L: "$I$", NoClip: true})
		plt.Gll("$x$", "$f(x)$", nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "chebyinterp02")
	}
}

// checkIandIs checks I, Is and ℓl @ nodes
func checkIandIs(tst *testing.T, N int, f Ss, tol float64, verb bool) {

	// allocate
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// compute coefficients
	o.CalcCoefI(f)
	o.CalcCoefIs(f)
	chk.EP(err)

	// check interpolations
	xx := utl.LinSpace(-1, 1, 11)
	for _, x := range xx {
		i1 := o.I(x)
		i2 := o.Il(x)
		chk.AnaNum(tst, "I(x) == Is(x)", 1e-14, i1, i2, chk.Verbose)
	}
	if verb {
		io.Pl()
	}

	// check ℓ @ notes
	for k, x := range o.X {
		for l := 0; l < o.N+1; l++ {
			res := o.L(l, x)
			if k == l {
				chk.AnaNum(tst, io.Sf("ℓ_%d(x_%d)==1", l, k), tol, res, 1, verb)
			} else {
				chk.AnaNum(tst, io.Sf("ℓ_%d(x_%d)==0", l, k), tol, res, 0, verb)
			}
		}
		if verb {
			io.Pl()
		}
	}
}

func TestChebyInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp03. ℓ and I(x) versus Is(x)")

	// test function
	f := func(x float64) (float64, error) {
		return math.Cos(math.Exp(2.0 * x)), nil
	}

	// allocate polynomial
	Nvals := []int{5, 6}
	for _, N := range Nvals {
		checkIandIs(tst, N, f, 1e-15, chk.Verbose)
	}

	// plot
	if chk.Verbose {
		N := 5
		o, err := NewChebyInterp(N, false) // Gauss-Lobatto
		chk.EP(err)
		npts := 201
		xx := utl.LinSpace(-1, 1, npts)
		yy := make([]float64, npts)
		plt.Reset(true, nil)
		for l := 0; l < N+1; l++ {
			for i := 0; i < npts; i++ {
				yy[i] = o.L(l, xx[i])
			}
			plt.Plot(xx, yy, &plt.A{C: plt.C(l, 1), L: io.Sf("l=%d", l), NoClip: true})
		}
		plt.HideTRborders()
		plt.Gll("$x$", "$\\psi_l(x)$", &plt.A{LegOut: true, LegNcol: 7, LegHlen: 2})
		plt.Save("/tmp/gosl/fun", "chebyinterp03")
	}
}

// cmpD1che compares D1 matrices with numerical differentiation and with each other
func cmpD1che(tst *testing.T, msg string, o *ChebyInterp, D1, D1trig [][]float64, tolD, tolCmp float64, verb bool) {
	io.Pf("\n\n . . . . . . . . . . . . . %s . . . . . . . . . . . . . . . . . . . . . \n\n", msg)
	m := len(D1)
	for j := 0; j < m; j++ {
		xj := o.X[j]
		for l := 0; l < m; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, D1[j][l], xj, 1e-2, verb, func(t float64) (float64, error) {
				return o.L(l, t), nil
			})
		}
	}
	if verb {
		io.Pl()
	}
	chk.Deep2(tst, "D1 vs. D1trig "+msg, tolCmp, D1, D1trig)
}

func checkD1che(tst *testing.T, N int, tolD, tolCmp float64, verb bool) {

	// allocate
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// noFlip,noNst
	o.Trig, o.Flip, o.Nst = false, false, false
	err = o.CalcD1()
	chk.EP(err)
	D1 := o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, false, false
	err = o.CalcD1()
	chk.EP(err)
	D1trig := o.D1.GetDeep2()
	cmpD1che(tst, "[noFlip,noNst]", o, D1, D1trig, tolD, tolCmp, verb)

	// flip,noNst
	o.Trig, o.Flip, o.Nst = false, true, false
	err = o.CalcD1()
	chk.EP(err)
	D1 = o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, true, false
	err = o.CalcD1()
	chk.EP(err)
	D1trig = o.D1.GetDeep2()
	cmpD1che(tst, "[flip,noNst]", o, D1, D1trig, tolD, tolCmp, verb)

	// nst
	o.Trig, o.Flip, o.Nst = false, false, true
	err = o.CalcD1()
	chk.EP(err)
	D1 = o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, false, true
	err = o.CalcD1()
	chk.EP(err)
	D1trig = o.D1.GetDeep2()
	cmpD1che(tst, "[---, nst]", o, D1, D1trig, tolD, tolCmp, verb)
}

func TestChebyInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp04. D1 matrix: first derivative of ℓ @ nodes")

	// run test
	Nvals := []int{3, 4, 5}
	tolsD := []float64{1e-5, 1e-5, 1e-4}
	tolsC := []float64{1e-15, 1e-15, 1e-14}
	for i, N := range Nvals {
		checkD1che(tst, N, tolsD[i], tolsC[i], chk.Verbose)
	}
}

func checkD2(tst *testing.T, N int, h, tolD float64, verb bool) {

	// allocate
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// check D2 matrix
	hh := h * h
	err = o.CalcD2()
	chk.EP(err)
	for j := 0; j < o.N+1; j++ {
		xj := o.X[j]
		for l := 0; l < o.N+1; l++ {
			LlxjBefore := o.L(l, xj-h)
			LlxjCurrent := o.L(l, xj)
			LlxjAfter := o.L(l, xj+h)
			dLldxAtXj := (LlxjBefore - 2.0*LlxjCurrent + LlxjAfter) / hh
			chk.AnaNum(tst, io.Sf("D2[%d,%d](%+.3f)", j, l, xj), tolD, o.D2.Get(j, l), dLldxAtXj, verb)
		}
	}
	if verb {
		io.Pl()
	}
}

func TestChebyInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp05. Second derivative of ℓ")

	// check D2
	k := 0
	hs := []float64{1e-2, 1e-3, 1e-3, 1e-3}
	tols := []float64{1e-9, 1e-5, 1e-4, 1e-3}
	for N := 3; N <= 6; N++ {
		checkD2(tst, N, hs[k], tols[k], chk.Verbose)
		k++
	}
}

// cmpD1ana compares D1 matrices with numerical differentiation and with each other
func cmpD1ana(tst *testing.T, msg string, o *ChebyInterp, f, dfdxAna Ss, tol float64, verb bool) {
	err := o.CalcCoefIs(f)
	chk.EP(err)
	u := o.CoefIs //f @ nodes: u = f(x_i)
	v := la.NewVector(o.N + 1)
	la.MatVecMul(v, 1, o.D1, u)
	maxDiff := o.CalcErrorD1(dfdxAna)
	if maxDiff > tol {
		tst.Errorf(msg+": maxDiff = %23g ⇒ D1 failed\n", maxDiff)
	} else {
		io.Pf(msg+": maxDiff = %23g ⇒ D1 is OK\n", maxDiff)
	}
}

func checkD1ana(tst *testing.T, N int, f, dfdxAna Ss, tol, tolTrig, tolNst float64, verb bool) {

	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)

	o.Trig, o.Flip, o.Nst = false, false, false
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[noTrig, noFlip, noNst]", o, f, dfdxAna, tol, verb)

	o.Trig, o.Flip, o.Nst = true, false, false
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[  trig, noFlip, noNst]", o, f, dfdxAna, tolTrig, verb)

	o.Trig, o.Flip, o.Nst = false, true, false
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[noTrig,   flip, noNst]", o, f, dfdxAna, tol, verb)

	o.Trig, o.Flip, o.Nst = true, true, false
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[  trig,   flip, noNst]", o, f, dfdxAna, tolTrig, verb)

	// nst

	o.Trig, o.Flip, o.Nst = false, false, true
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[noTrig,  ---  ,   nst]", o, f, dfdxAna, tolNst, verb)

	o.Trig, o.Flip, o.Nst = true, false, true
	err = o.CalcD1()
	chk.EP(err)
	cmpD1ana(tst, "[  trig,  ---  ,   nst]", o, f, dfdxAna, tolNst, verb)
}

func TestChebyInterp06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp06. I(x) versus Is(x). D1 matrices")

	// test function
	f := func(x float64) (float64, error) {
		//return math.Cos(math.Exp(2.0 * x)), nil
		//return 1.0 / (1.0 + 4.0*x*x), nil
		return math.Pow(x, 8), nil
	}
	g := func(x float64) (float64, error) {
		//return -2.0 * math.Exp(2.0*x) * math.Sin(math.Exp(2.0*x)), nil
		//return -8.0 * x / math.Pow(1.0+4.0*x*x, 2.0), nil
		return 8.0 * math.Pow(x, 7), nil
	}

	// allocate polynomial
	N := 8
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)

	// check interpolations
	o.CalcCoefI(f)
	o.CalcCoefIs(f)
	chk.EP(err)
	xx := utl.LinSpace(-1, 1, 11)
	for _, x := range xx {
		i1 := o.I(x)
		i2 := o.Il(x)
		chk.AnaNum(tst, "I(x) == Is(x)", 1e-14, i1, i2, chk.Verbose)
	}

	// check D1 matrices
	io.Pl()
	checkD1ana(tst, N, f, g, 1e-14, 1e-13, 1e-14, chk.Verbose)

	// plot
	if chk.Verbose && true {
		npts := 201
		xx := utl.LinSpace(-1, 1, npts)
		y1 := make([]float64, len(xx))
		y2 := make([]float64, len(xx))
		y3 := make([]float64, len(xx))
		y4 := make([]float64, len(xx))
		for i, x := range xx {
			y1[i], _ = f(x)
			y2[i] = o.I(x)
			y3[i] = o.Il(x)
			y4[i], _ = g(x)
		}

		o.CalcD1()
		u := o.CoefIs //f @ nodes: u = f(x_i)
		v := la.NewVector(o.N + 1)
		la.MatVecMul(v, 1, o.D1, u)

		plt.Reset(true, &plt.A{Prop: 1.5})

		plt.Subplot(2, 1, 1)
		plt.Plot(o.X, u, &plt.A{L: "$f(x_i)$", C: "r", Ls: "none", M: "o", Void: true, NoClip: true})
		plt.Plot(xx, y1, &plt.A{C: plt.C(0, 1), L: "$f$", NoClip: true})
		plt.Plot(xx, y2, &plt.A{C: plt.C(1, 1), L: "$I$", Lw: 3, NoClip: true})
		plt.Plot(xx, y3, &plt.A{C: plt.C(2, 1), L: "$Is$", M: "+", Me: 20, NoClip: true})
		plt.Gll("$x$", "$f(x)$", nil)
		plt.HideAllBorders()

		plt.Subplot(2, 1, 2)
		plt.Plot(xx, y4, &plt.A{C: plt.C(0, 1), L: "df/dx", NoClip: true})
		plt.Plot(o.X, v, &plt.A{C: "r", Ls: "none", M: ".", L: "d(Iu)/dx @ xi", NoClip: true})
		plt.Gll("$x$", "$g(x)$", nil)
		plt.HideAllBorders()

		plt.Save("/tmp/gosl/fun", "chebyinterp06")
	}
}

func calcD1errorChe(N int, f, dfdxAna Ss, trig, flip, nst bool) (maxDiff float64) {

	// allocate polynomial
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)

	// compute coefficients
	err = o.CalcCoefIs(f)
	chk.EP(err)

	// compute D1 matrix
	o.Trig, o.Flip, o.Nst = trig, flip, nst
	err = o.CalcD1()
	chk.EP(err)

	// compute error
	maxDiff = o.CalcErrorD1(dfdxAna)
	return
}

func TestChebyInterp07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp07. round-off errors")

	// test function
	f := func(x float64) (float64, error) {
		return math.Pow(x, 8), nil
		//return math.Sin(8.0*x) / math.Pow(x+1.1, 1.5), nil
	}
	g := func(x float64) (float64, error) {
		//d := math.Pow(x+1.1, 1.5)
		//return (8*math.Cos(8*x))/d - (3*math.Sin(8*x))/(2*(1.1+x)*d), nil
		return 8.0 * math.Pow(x, 7), nil
	}

	if chk.Verbose {
		var dummy bool

		// check
		Nvals := []int{16, 32, 50, 64, 100, 128, 250, 256, 500, 512, 1000, 1024, 2000, 2048}
		//Nvals := utl.IntRange3(16, 2048, 10)
		nn := make([]float64, len(Nvals))
		eeA := make([]float64, len(Nvals))
		eeB := make([]float64, len(Nvals))
		eeC := make([]float64, len(Nvals))
		eeD := make([]float64, len(Nvals))
		eeE := make([]float64, len(Nvals))
		eeF := make([]float64, len(Nvals))
		for i, N := range Nvals {
			nn[i] = float64(N)

			eeA[i] = calcD1errorChe(N, f, g, false, false, false)
			eeB[i] = calcD1errorChe(N, f, g, true, false, false)
			eeC[i] = calcD1errorChe(N, f, g, false, true, false)
			eeD[i] = calcD1errorChe(N, f, g, true, true, false)
			eeE[i] = calcD1errorChe(N, f, g, false, dummy, true)
			eeF[i] = calcD1errorChe(N, f, g, true, dummy, true)

			io.Pf("%4d: %.2e  %.2e  %.2e  %.2e  %.2e  %.2e\n", N,
				eeA[i], eeB[i], eeC[i], eeD[i], eeE[i], eeF[i])

		}

		// plot
		plt.Reset(true, nil)

		plt.Plot(nn, eeA, &plt.A{C: "b", L: "std,nof", M: "s", Me: 1, NoClip: true})
		plt.Plot(nn, eeB, &plt.A{C: "r", L: "tri,nof", M: "+", Me: 1, NoClip: true})
		plt.Plot(nn, eeC, &plt.A{C: "c", L: "std,fli", M: ".", Me: 1, NoClip: true})
		plt.Plot(nn, eeD, &plt.A{C: "m", L: "tri,fli", M: "*", Me: 1, NoClip: true})
		plt.Plot(nn, eeE, &plt.A{C: "y", L: "std,nst", M: "s", Me: 1, NoClip: true})
		plt.Plot(nn, eeF, &plt.A{C: "k", L: "tri,nst", M: "+", Me: 1, NoClip: true})

		plt.Gll("$N$", "$||Df-df/dx||_\\infty$", &plt.A{LegOut: true, LegNcol: 3, LegHlen: 3})
		plt.SetYlog()
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "chebyinterp07")
	}
}
