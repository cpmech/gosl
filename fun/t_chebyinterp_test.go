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
	"github.com/cpmech/gosl/utl"
)

func TestChebyInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp01")

	// test function
	f := func(x float64) float64 {
		return math.Cos(math.Exp(2.0 * x))
	}

	// allocate polynomials
	N := 8
	che := NewChebyInterp(N, true)  // Gauss-Chebyshev
	lob := NewChebyInterp(N, false) // Gauss-Lobatto

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
}

func TestChebyInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp02")

	// test function
	f := func(x float64) float64 {
		//return 1.0 / (1.0 + 4.0*x*x), nil
		return math.Cos(math.Exp(2.0 * x))
	}

	// allocate polynomials
	N := 6
	o := NewChebyInterp(N, false) // Gauss-Lobatto

	// compute data
	o.CalcCoefI(f)

	// check interpolation @ nodes
	for k, x := range o.X {
		fk := f(x)
		chk.Float64(tst, io.Sf("I(x_%d)", k), 1e-14, o.I(x), fk)
	}

	// check conversion
	o.CalcConvMats()
	np := len(o.X) // number of points
	u := la.NewVector(np)
	for i, x := range o.X {
		u[i] = f(x)
	}
	ub := la.NewVector(np)
	la.MatVecMul(ub, 1, o.C, u)
	io.Pf("u  = %.6f\n", u)
	io.Pfyel("ub = %.6f\n", ub)
	io.Pfyel("cf = %.6f\n", o.CoefI)
	chk.Array(tst, "ub", 1.05e-15, ub, o.CoefI)

	// check inversion
	uu := la.NewVector(np)
	la.MatVecMul(uu, 1, o.Ci, ub)
	io.Pf("uu = %.6f\n", uu)
	chk.Array(tst, "uu", 1e-15, uu, u)
}

// checkIandIs checks I, Is and ℓl @ nodes
func checkIandIs(tst *testing.T, N int, f Ss, tol float64, verb bool) {

	// allocate
	o := NewChebyInterp(N, false) // Gauss-Lobatto
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// compute coefficients
	o.CalcCoefI(f)
	o.CalcCoefIs(f)

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
	f := func(x float64) float64 {
		return math.Cos(math.Exp(2.0 * x))
	}

	// allocate polynomial
	Nvals := []int{5, 6}
	for _, N := range Nvals {
		checkIandIs(tst, N, f, 1e-15, chk.Verbose)
	}
}

// cmpD1che compares D1 matrices with numerical differentiation and with each other
func cmpD1che(tst *testing.T, msg string, o *ChebyInterp, D1, D1trig [][]float64, tolD, tolCmp float64, verb bool) {
	io.Pf("\n\n . . . . . . . . . . . . . %s . . . . . . . . . . . . . . . . . . . . . \n\n", msg)
	m := len(D1)
	for j := 0; j < m; j++ {
		xj := o.X[j]
		for l := 0; l < m; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tolD, D1[j][l], xj, 1e-2, verb, func(t float64) float64 {
				return o.L(l, t)
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
	o := NewChebyInterp(N, false) // Gauss-Lobatto
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// noFlip,noNst
	o.Trig, o.Flip, o.Nst = false, false, false
	o.CalcD1()
	D1 := o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, false, false
	o.CalcD1()
	D1trig := o.D1.GetDeep2()
	cmpD1che(tst, "[noFlip,noNst]", o, D1, D1trig, tolD, tolCmp, verb)

	// flip,noNst
	o.Trig, o.Flip, o.Nst = false, true, false
	o.CalcD1()
	D1 = o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, true, false
	o.CalcD1()
	D1trig = o.D1.GetDeep2()
	cmpD1che(tst, "[flip,noNst]", o, D1, D1trig, tolD, tolCmp, verb)

	// nst
	o.Trig, o.Flip, o.Nst = false, false, true
	o.CalcD1()
	D1 = o.D1.GetDeep2()
	o.Trig, o.Flip, o.Nst = true, false, true
	o.CalcD1()
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

// cmpD1ana compares D1 matrices with numerical differentiation and with each other
func cmpD1ana(tst *testing.T, msg string, o *ChebyInterp, f, dfdxAna Ss, tol float64, verb bool) {
	o.CalcCoefIs(f)
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

	o := NewChebyInterp(N, false) // Gauss-Lobatto

	o.Trig, o.Flip, o.Nst = false, false, false
	o.CalcD1()
	cmpD1ana(tst, "[noTrig, noFlip, noNst]", o, f, dfdxAna, tol, verb)

	o.Trig, o.Flip, o.Nst = true, false, false
	o.CalcD1()
	cmpD1ana(tst, "[  trig, noFlip, noNst]", o, f, dfdxAna, tolTrig, verb)

	o.Trig, o.Flip, o.Nst = false, true, false
	o.CalcD1()
	cmpD1ana(tst, "[noTrig,   flip, noNst]", o, f, dfdxAna, tol, verb)

	o.Trig, o.Flip, o.Nst = true, true, false
	o.CalcD1()
	cmpD1ana(tst, "[  trig,   flip, noNst]", o, f, dfdxAna, tolTrig, verb)

	// nst

	o.Trig, o.Flip, o.Nst = false, false, true
	o.CalcD1()
	cmpD1ana(tst, "[noTrig,  ---  ,   nst]", o, f, dfdxAna, tolNst, verb)

	o.Trig, o.Flip, o.Nst = true, false, true
	o.CalcD1()
	cmpD1ana(tst, "[  trig,  ---  ,   nst]", o, f, dfdxAna, tolNst, verb)
}

func TestChebyInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp05. I(x) versus Is(x). D1 matrices")

	// test function
	f := func(x float64) float64 {
		//return math.Cos(math.Exp(2.0 * x))
		//return 1.0 / (1.0 + 4.0*x*x)
		return math.Pow(x, 8)
	}
	g := func(x float64) float64 {
		//return -2.0 * math.Exp(2.0*x) * math.Sin(math.Exp(2.0*x))
		//return -8.0 * x / math.Pow(1.0+4.0*x*x, 2.0)
		return 8.0 * math.Pow(x, 7)
	}

	// allocate polynomial
	N := 8
	o := NewChebyInterp(N, false) // Gauss-Lobatto

	// check interpolations
	o.CalcCoefI(f)
	o.CalcCoefIs(f)
	xx := utl.LinSpace(-1, 1, 11)
	for _, x := range xx {
		i1 := o.I(x)
		i2 := o.Il(x)
		chk.AnaNum(tst, "I(x) == Is(x)", 1e-14, i1, i2, chk.Verbose)
	}

	// check D1 matrices
	io.Pl()
	checkD1ana(tst, N, f, g, 1e-14, 1e-13, 1e-13, chk.Verbose)
}

func calcD1errorChe(tst *testing.T, N int, f, dfdxAna Ss, trig, flip, nst bool) (maxDiff float64) {

	// allocate polynomial
	o := NewChebyInterp(N, false) // Gauss-Lobatto

	// compute coefficients
	o.CalcCoefIs(f)

	// compute D1 matrix
	o.Trig, o.Flip, o.Nst = trig, flip, nst
	o.CalcD1()

	// compute error
	maxDiff = o.CalcErrorD1(dfdxAna)
	return
}

func calcD2errorChe(tst *testing.T, N int, f, dfdxAna Ss, stdD2 bool) (maxDiff float64) {

	// allocate polynomial
	o := NewChebyInterp(N, false) // Gauss-Lobatto

	// compute coefficients
	o.CalcCoefIs(f)

	// compute D2 matrix
	o.StdD2 = stdD2
	o.CalcD2()

	// compute error
	maxDiff = o.CalcErrorD2(dfdxAna)
	return
}

func checkD2che(tst *testing.T, N int, h, tolD float64, verb bool) {

	// allocate
	o := NewChebyInterp(N, false) // Gauss-Lobatto
	if verb {
		io.Pf("\n\n----------------------------- N = %d -----------------------------------------\n\n", N)
	}

	// check D2 matrix
	hh := h * h
	o.CalcD2()
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

func TestChebyInterp07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp05. Second derivative of ℓ")

	// check D2
	k := 0
	hs := []float64{1e-2, 1e-3, 1e-3, 1e-3}
	tols := []float64{1e-9, 1e-5, 1e-4, 1e-3}
	for N := 3; N <= 6; N++ {
		checkD2che(tst, N, hs[k], tols[k], chk.Verbose)
		k++
	}
}

func TestChebyInterp08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ChebyInterp08. D1 and D2. analytical")

	f := func(x float64) float64 {
		return math.Pow(x, 8)
	}
	g := func(x float64) float64 {
		return 8.0 * math.Pow(x, 7)
	}
	h := func(x float64) float64 {
		return 56.0 * math.Pow(x, 6)
	}

	N := 8

	maxDiff := calcD1errorChe(tst, N, f, g, false, false, false)
	io.Pf("no nst: err(D1{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D1{f})", 1e-14, maxDiff, 0)

	maxDiff = calcD1errorChe(tst, N, f, g, false, false, true)
	io.Pf("   nst: err(D1{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D1{f})", 1e-14, maxDiff, 0)

	maxDiff = calcD2errorChe(tst, N, f, h, true)
	io.Pf("   std: err(D2{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-12, maxDiff, 0)

	maxDiff = calcD2errorChe(tst, N, f, h, false)
	io.Pf("use D1: err(D2{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-12, maxDiff, 0)
}
