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

func checkD1(tst *testing.T, N int, tolPsi, tolD, tolCmp float64, verb bool) {

	// allocate
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// check ψl @ nodes
	for k, x := range o.X {
		for l := 0; l < o.N+1; l++ {
			res := o.PsiLobDirect(l, x)
			if k == l {
				chk.AnaNum(tst, io.Sf("ψ_%d(x_%d)==1", l, k), tolPsi, res, 1, verb)
			} else {
				chk.AnaNum(tst, io.Sf("ψ_%d(x_%d)==0", l, k), tolPsi, res, 0, verb)
			}
		}
		if verb {
			io.Pl()
		}
	}

	// ------------------------------- no flip

	// check D1 matrix (noFlip)
	err = o.CalcD1(false, false)
	chk.EP(err)
	D1 := o.D1.GetDeep2()
	for j := 0; j < o.N+1; j++ {
		xj := o.X[j]
		for l := 0; l < o.N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, o.D1.Get(j, l), xj, 1e-2, verb, func(t float64) (float64, error) {
				return o.PsiLobDirect(l, t), nil
			})
		}
	}
	if verb {
		io.Pl()
	}

	// check D1 matrix (trigo, noFlip)
	err = o.CalcD1(true, false)
	chk.EP(err)
	D1trigo := o.D1.GetDeep2()
	for j := 0; j < o.N+1; j++ {
		xj := o.X[j]
		for l := 0; l < o.N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, o.D1.Get(j, l), xj, 1e-2, verb, func(t float64) (float64, error) {
				return o.PsiLobDirect(l, t), nil
			})
		}
	}
	if verb {
		io.Pl()
	}

	// compare D1
	chk.Deep2(tst, "D1 [noflip]", tolCmp, D1, D1trigo)

	// ------------------------------- flip
	io.Pl()
	io.Pl()

	// check D1 matrix (noFlip)
	err = o.CalcD1(false, true)
	chk.EP(err)
	D1 = o.D1.GetDeep2()
	for j := 0; j < o.N+1; j++ {
		xj := o.X[j]
		for l := 0; l < o.N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, o.D1.Get(j, l), xj, 1e-2, verb, func(t float64) (float64, error) {
				return o.PsiLobDirect(l, t), nil
			})
		}
	}
	if verb {
		io.Pl()
	}

	// check D1 matrix (trigo, noFlip)
	err = o.CalcD1(true, true)
	chk.EP(err)
	D1trigo = o.D1.GetDeep2()
	for j := 0; j < o.N+1; j++ {
		xj := o.X[j]
		for l := 0; l < o.N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, o.D1.Get(j, l), xj, 1e-2, verb, func(t float64) (float64, error) {
				return o.PsiLobDirect(l, t), nil
			})
		}
	}
	if verb {
		io.Pl()
	}

	// compare D1
	chk.Deep2(tst, "D1 [flip]", tolCmp, D1, D1trigo)
}

func TestChebyInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp03. ψ and first derivative of ψ")

	// run test
	k := 0
	tolPsi := 1e-15
	tolsD := []float64{1e-5, 1e-5, 1e-4, 1e-4}
	tolsC := []float64{1e-15, 1e-15, 1e-14, 1e-14}
	for N := 3; N <= 6; N++ {
		checkD1(tst, N, tolPsi, tolsD[k], tolsC[k], chk.Verbose)
		k++
	}

	// plot
	if chk.Verbose {
		N := 4
		o, err := NewChebyInterp(N, false) // Gauss-Lobatto
		chk.EP(err)
		npts := 201
		xx := utl.LinSpace(-1, 1, npts)
		yy := make([]float64, npts)
		plt.Reset(true, nil)
		for l := 0; l < N+1; l++ {
			for i := 0; i < npts; i++ {
				yy[i] = o.PsiLobDirect(l, xx[i])
			}
			plt.Plot(xx, yy, &plt.A{C: plt.C(l, 1), L: io.Sf("l=%d", l), NoClip: true})
		}
		plt.HideTRborders()
		plt.Gll("$x$", "$\\psi_l(x)$", &plt.A{LegOut: true, LegNcol: 7, LegHlen: 2})
		plt.Save("/tmp/gosl/fun", "chebyinterp03")
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
			ψlxjBefore := o.PsiLobDirect(l, xj-h)
			ψlxjCurrent := o.PsiLobDirect(l, xj)
			ψlxjAfter := o.PsiLobDirect(l, xj+h)
			dψldxAtXj := (ψlxjBefore - 2.0*ψlxjCurrent + ψlxjAfter) / hh
			chk.AnaNum(tst, io.Sf("D2[%d,%d](%+.3f)", j, l, xj), tolD, o.D2.Get(j, l), dψldxAtXj, verb)
		}
	}
	if verb {
		io.Pl()
	}
}

func TestChebyInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp04. Second derivative of ψ")

	// check D2
	k := 0
	hs := []float64{1e-2, 1e-3, 1e-3, 1e-3}
	tols := []float64{1e-9, 1e-5, 1e-4, 1e-3}
	for N := 3; N <= 6; N++ {
		checkD2(tst, N, hs[k], tols[k], chk.Verbose)
		k++
	}
}

func TestChebyInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp05. I(x) versus Is(x). Flip vs NoFlip")

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

	// compute coefficients
	o.CalcCoefI(f)
	o.CalcCoefIs(f)
	chk.EP(err)
	u := o.CoefIs //f @ nodes: u = f(x_i)

	// check interpolations
	xx := utl.LinSpace(-1, 1, 11)
	for _, x := range xx {
		i1 := o.I(x)
		i2 := o.Is(x)
		chk.AnaNum(tst, "I(x) == Is(x)", 1e-14, i1, i2, chk.Verbose)
	}

	// derivative of interpolation @ x_i
	trigo, flip := false, false
	err = o.CalcD1(trigo, flip)
	chk.EP(err)
	v := la.NewVector(o.N + 1)
	la.MatVecMul(v, 1, o.D1, u)

	// check error on derivative
	maxDiff := o.CalcErrorD1(f, g)
	io.Pforan("maxDiff = %v\n", maxDiff)
	if maxDiff > 1e-14 {
		tst.Errorf("maxDiff = %g ⇒ D1 failed\n", maxDiff)
	}

	// compare flip vs non-flip
	D1flip := o.D1.GetDeep2()
	flip = true
	err = o.CalcD1(trigo, flip)
	chk.EP(err)
	chk.Deep2(tst, "D1 flip vs noFlip", 1e-13, D1flip, o.D1.GetDeep2())

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
			y3[i] = o.Is(x)
			y4[i], _ = g(x)
		}
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

		plt.Save("/tmp/gosl/fun", "chebyinterp05")
	}
}

func calcD1error(N int, f, dfdxAna Ss, trigo, flip bool) (maxDiff float64) {

	// allocate polynomial
	o, err := NewChebyInterp(N, false) // Gauss-Lobatto
	chk.EP(err)

	// compute coefficients
	err = o.CalcCoefIs(f)
	chk.EP(err)

	// compute D1 matrix
	err = o.CalcD1(trigo, flip)
	chk.EP(err)

	// compute error
	maxDiff = o.CalcErrorD1(f, dfdxAna)
	return
}

func TestChebyInterp06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp06. round-off errors")

	// test function
	f := func(x float64) (float64, error) {
		return math.Pow(x, 8), nil
	}
	g := func(x float64) (float64, error) {
		return 8.0 * math.Pow(x, 7), nil
	}

	if chk.Verbose {

		// check
		Nvals := []int{16, 32, 50, 64, 100, 128, 250, 256, 500, 512, 1000, 1024, 2000, 2048}
		nn := make([]float64, len(Nvals))
		eeA := make([]float64, len(Nvals))
		eeB := make([]float64, len(Nvals))
		eeC := make([]float64, len(Nvals))
		eeD := make([]float64, len(Nvals))
		for i, N := range Nvals {
			nn[i] = float64(N)
			eeA[i] = calcD1error(N, f, g, false, false)
			eeB[i] = calcD1error(N, f, g, true, false)
			eeC[i] = calcD1error(N, f, g, false, true)
			eeD[i] = calcD1error(N, f, g, true, true)
			io.Pf("%4d: maxdiff = %24v  %24v\n", N, eeA[i], eeB[i])
		}

		// plot
		plt.Reset(true, nil)
		plt.Plot(nn, eeA, &plt.A{C: "b", L: "std,noFlip", M: "s", Me: 1, NoClip: true})
		plt.Plot(nn, eeB, &plt.A{C: "r", L: "tri,noFlip", M: "+", Me: 1, NoClip: true})
		plt.Plot(nn, eeC, &plt.A{C: "c", L: "std,flip", M: ".", Me: 1, NoClip: true})
		plt.Plot(nn, eeD, &plt.A{C: "m", L: "tri,flip", M: "*", Me: 1, NoClip: true})
		plt.Gll("$N$", "$||Df-df/dx||_\\infty$", nil)
		plt.SetYlog()
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "chebyinterp06")
	}
}
