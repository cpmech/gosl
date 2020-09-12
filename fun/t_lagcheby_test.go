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

func compareLambda(tst *testing.T, N int, f Ss, tolU, tolL float64) {

	// allocate Lagrange structure and calculate U
	lag := NewLagrangeInterp(N, "cgl")
	lag.CalcU(f)

	// allocate Chebyshev structure and calculate U
	che := NewChebyInterp(N, false) // Gauss-Lobatto
	che.CalcCoefIs(f)

	// check U values
	io.Pf("\n-------------------------------- N = %d -----------------------------------\n", N)
	cheU := utl.GetReversed(che.CoefIs)
	if N < 9 {
		io.Pforan("lag.U = %+8.4f\n", lag.U)
		io.Pfyel("che.U = %+8.4f\n", cheU)
	}
	chk.Array(tst, "U", tolU, lag.U, cheU)

	// check λ values
	cheL := utl.GetReversed(che.Lam)
	chk.Array(tst, "λ", tolL, lag.Lam, cheL)
}

func TestLagCheby01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby01")

	// test function
	f := func(x float64) float64 {
		return math.Cos(math.Exp(2.0 * x))
	}

	// test
	Nvals := []int{6, 7, 8, 9, 700, 701, 1024, 2048}
	tolsU := []float64{1e-17, 1e-17, 1e-17, 1e-17, 1e-17, 1e-17, 1e-17, 1e-17}
	tolsL := []float64{1e-15, 1e-15, 1e-15, 1e-14, 1e-11, 1e-11, 1e-11, 1e-10}
	for k, N := range Nvals {
		compareLambda(tst, N, f, tolsU[k], tolsL[k])
	}
}

func runD1err(tst *testing.T, fnkey string, Nvals []int, f, g Ss) {
	nn := make([]float64, len(Nvals))
	eeA := make([]float64, len(Nvals))
	eeB := make([]float64, len(Nvals))
	eeC := make([]float64, len(Nvals))
	eeD := make([]float64, len(Nvals))
	dummy := false
	for i, N := range Nvals {
		nn[i] = float64(N)
		eeA[i] = calcD1errorChe(tst, N, f, g, false, dummy, true) // std,nst
		eeB[i] = calcD1errorChe(tst, N, f, g, true, dummy, true)  // tri,nst
		eeC[i] = calcD1errorLag(tst, N, f, g, false)              // lag,---
		eeD[i] = calcD1errorLag(tst, N, f, g, true)               // lag,eta
		io.Pf("%4d: %.2e  %.2e  %.2e  %.2e\n", N, eeA[i], eeB[i], eeC[i], eeD[i])
	}
}

func runD2err(tst *testing.T, fnkey string, Nvals []int, f, h Ss) {
	nn := make([]float64, len(Nvals))
	eeA := make([]float64, len(Nvals))
	eeB := make([]float64, len(Nvals))
	eeC := make([]float64, len(Nvals))
	eeD := make([]float64, len(Nvals))
	for i, N := range Nvals {
		nn[i] = float64(N)
		eeA[i] = calcD2errorChe(tst, N, f, h, false) // che,useD1
		eeB[i] = calcD2errorChe(tst, N, f, h, true)  // che,std
		eeC[i] = calcD2errorLag(tst, N, f, h, false) // lag,lam
		eeD[i] = calcD2errorLag(tst, N, f, h, true)  // lag,eta
		io.Pf("%4d: %.2e  %.2e  %.2e  %.2e\n", N, eeA[i], eeB[i], eeC[i], eeD[i])
	}
}

func TestLagCheby02a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby02a. round-off errors")

	f := func(x float64) float64 {
		return math.Pow(x, 8)
	}
	g := func(x float64) float64 {
		return 8.0 * math.Pow(x, 7)
	}
	if chk.Verbose {
		Nvals := []int{16, 32, 50, 64, 100, 128, 250, 256, 500, 512, 1000, 1024, 2000, 2048}
		runD1err(tst, "lagcheby02a", Nvals, f, g)
	}
}

func TestLagCheby02b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby02b. round-off errors")

	f := func(x float64) float64 {
		return math.Sin(8.0*x) / math.Pow(x+1.1, 1.5)
	}
	g := func(x float64) float64 {
		d := math.Pow(x+1.1, 1.5)
		return (8*math.Cos(8*x))/d - (3*math.Sin(8*x))/(2*(1.1+x)*d)
	}
	if chk.Verbose {
		Nvals := []int{64, 100, 128, 250, 256, 500, 512, 1000, 1024, 2000, 2048}
		runD1err(tst, "lagcheby02b", Nvals, f, g)
	}
}

func TestLagCheby03a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby03a. round-off errors")

	f := func(x float64) float64 {
		return math.Pow(x, 8)
	}
	h := func(x float64) float64 {
		return 56.0 * math.Pow(x, 6)
	}
	if chk.Verbose {
		Nvals := []int{16, 32, 50, 64, 100, 128, 250, 256, 500}
		runD2err(tst, "lagcheby03a", Nvals, f, h)
	}
}

func TestLagCheby03b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LagCheby03b. round-off errors")

	f := func(x float64) float64 {
		return math.Sin(8.0*x) / math.Pow(x+1.1, 1.5)
	}
	h := func(x float64) float64 {
		m := x + 1.1
		d := math.Pow(m, 1.5)
		return -(64*math.Sin(8*x))/d + (3.75*math.Sin(8*x))/(d*m*m) - (24*math.Cos(8*x))/(d*m)
	}
	if chk.Verbose {
		Nvals := []int{64, 100, 128, 250, 256, 500}
		runD2err(tst, "lagcheby03b", Nvals, f, h)
	}
}
