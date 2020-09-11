// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
)

func TestLagCardinal01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCardinal01. Lagrange cardinal polynomials")

	// allocate structure
	N := 5
	kind := "uni"
	o := NewLagrangeInterp(N, kind)
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
}

func TestLagInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp01. Lagrange interpolation")

	// cos-exp function
	f := func(x float64) float64 {
		return math.Cos(math.Exp(2.0 * x))
	}

	// allocate structure and calculate U
	N := 5
	kind := "uni"
	o := NewLagrangeInterp(N, kind)
	o.CalcU(f)

	// check interpolation
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()
}

func TestLagInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp02. Lagrange interp. Runge problem")

	// Runge function
	f := func(x float64) float64 {
		return 1.0 / (1.0 + 16.0*x*x)
	}

	// allocate structure and calculate U
	N := 8
	kind := "uni"
	o := NewLagrangeInterp(N, kind)
	o.CalcU(f)

	// check interpolation
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()
}

func TestLagInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp03. Chebyshev-Gauss. Runge problem")

	// Runge function
	f := func(x float64) float64 {
		return 1.0 / (1.0 + 16.0*x*x)
	}

	// allocate structure and calculate U
	N := 8
	kind := "cg"
	o := NewLagrangeInterp(N, kind)
	o.CalcU(f)

	// check interpolation
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// check Lebesgue constants and compute max error
	ΛN := []float64{1.988854381999833e+00, 2.361856787767076e+00, 3.011792612349363e+00}
	for i, n := range []int{4, 8, 24} {
		p := NewLagrangeInterp(n, kind)
		chk.Float64(tst, "ΛN (Lebesgue constant)", 1e-13, p.EstimateLebesgue(), ΛN[i])
	}
}

func TestLagInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp04. Chebyshev-Gauss-Lobatto. Runge problem")

	// Runge function
	f := func(x float64) float64 {
		return 1.0 / (1.0 + 16.0*x*x)
	}

	// allocate structure and calculate U
	N := 8
	kind := "cgl"
	o := NewLagrangeInterp(N, kind)
	o.CalcU(f)

	// check interpolation
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, chk.Verbose)
	}
	io.Pl()

	// check Lebesgue constants and compute max error
	ΛN := []float64{1.798761778849085e+00, 2.274730699116020e+00, 2.984443326362511e+00}
	for i, n := range []int{4, 8, 24} {
		p := NewLagrangeInterp(n, kind)
		chk.Float64(tst, "ΛN (Lebesgue constant)", 1e-14, p.EstimateLebesgue(), ΛN[i])
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
		chk.AnaNum(tst, io.Sf("λ%d", i), tol, o.Lam[i], 1.0/d/m, chk.Verbose)
	}
}

func checkIandLam(tst *testing.T, N int, tolLam float64, f Ss) {

	// allocate structure and calculate U
	o := NewLagrangeInterp(N, "cgl")
	o.CalcU(f)

	// check interpolation (std)
	o.Bary = false
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// check λ
	checkLam(tst, o, tolLam)

	// check interpolation (barycentric)
	o.Bary = true
	for i, x := range o.X {
		ynum := o.I(x)
		yana := f(x)
		chk.AnaNum(tst, io.Sf("I(X[%d])", i), 1e-17, ynum, yana, false)
	}

	// compare std and barycentric
	xx := utl.LinSpace(-1, 1, 14)
	for _, x := range xx {
		for i := 0; i < o.N+1; i++ {
			o.Bary = false
			i1 := o.I(x)
			o.Bary = true
			i2 := o.I(x)
			chk.AnaNum(tst, io.Sf("I%d", i), 1e-15, i1, i2, false)
		}
	}
}

func TestLagInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestLagInterp05. Barycentric formulae")

	// Runge function
	f := func(x float64) float64 {
		return math.Cos(math.Exp(2.0 * x))
	}

	// test
	Nvals := []int{3, 4, 5, 6, 7, 8}
	tolsL := []float64{1e-15, 1e-15, 1e-15, 1e-15, 1e-14, 1e-14}
	for k, N := range Nvals {
		io.Pf("\n\n-------------------------------- N = %d -----------------------------------------------\n\n", N)
		checkIandLam(tst, N, tolsL[k], f)
	}
}

func cmpD1lag(tst *testing.T, N int, tol float64) {

	// allocate structure
	o := NewLagrangeInterp(N, "cgl")

	// calc and check D1
	o.CalcD1()
	for j := 0; j < N+1; j++ {
		xj := o.X[j]
		for l := 0; l < N+1; l++ {
			chk.DerivScaSca(tst, io.Sf("D1[%d,%d](%+.3f)", j, l, xj), tol, o.D1.Get(j, l), xj, 1e-2, chk.Verbose, func(t float64) float64 {
				return o.L(l, t)
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
		cmpD1lag(tst, N, tols[k])
	}
}

func checkD2lag(tst *testing.T, N int, h, tolD float64, verb bool) {

	// allocate
	o := NewLagrangeInterp(N, "cgl")
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

func TestLagInterp07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagInterp07. Second derivative of ℓ")

	// check D2
	Nvals := []int{3, 4, 5, 6}
	hs := []float64{1e-2, 1e-3, 1e-3, 1e-3}
	tols := []float64{1e-11, 1e-5, 1e-4, 1e-3}
	for k, N := range Nvals {
		checkD2lag(tst, N, hs[k], tols[k], chk.Verbose)
	}
}

func calcD1errorLag(tst *testing.T, N int, f, dfdxAna Ss, useEta bool) (maxDiff float64) {

	// allocate polynomial
	o := NewLagrangeInterp(N, "cgl")

	// compute coefficients
	o.CalcU(f)

	// compute D1 matrix
	o.UseEta = useEta
	o.CalcD1()

	// compute error
	maxDiff = o.CalcErrorD1(dfdxAna)
	return
}

func calcD2errorLag(tst *testing.T, N int, f, d2fdx2Ana Ss, useEta bool) (maxDiff float64) {

	// allocate polynomial
	o := NewLagrangeInterp(N, "cgl")

	// compute coefficients
	o.CalcU(f)

	// compute D2 matrix
	o.UseEta = useEta
	o.CalcD2()

	// compute error
	maxDiff = o.CalcErrorD2(d2fdx2Ana)
	return
}

func TestLagInterp08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagInterp08. D1 and D2. analytical")

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

	maxDiff := calcD1errorLag(tst, N, f, g, true)
	io.Pf("useEta: err(D1{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-13, maxDiff, 0)

	maxDiff = calcD2errorLag(tst, N, f, h, true)
	io.Pf("useEta: err(D2{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-12, maxDiff, 0)

	maxDiff = calcD1errorLag(tst, N, f, g, false)
	io.Pf("no eta: err(D1{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-13, maxDiff, 0)

	maxDiff = calcD2errorLag(tst, N, f, h, false)
	io.Pf("no eta: err(D2{f}) = %v\n", maxDiff)
	chk.Float64(tst, "err(D2{f})", 1e-12, maxDiff, 0)
}
