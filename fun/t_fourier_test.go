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

func TestFourierInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp01. check k(j) and j(k)")

	// constants
	N := 8
	fou := NewFourierInterp(N, "")
	defer fou.Free()

	// check k
	chk.Array(tst, "k[j]", 1e-17, fou.K, []float64{0, 1, 2, 3, -4, -3, -2, -1})

	// check j
	jvals := make([]int, N)
	i := 0
	for k := -N / 2; k <= N/2-1; k++ {
		jvals[i] = fou.CalcJ(float64(k))
		i++
	}
	io.Pf("jvals = %v\n", jvals)
	chk.Ints(tst, "j[k]", jvals, []int{4, 5, 6, 7, 0, 1, 2, 3})
}

// check interpolation @ nodes
func fouCheckI(tst *testing.T, fou *FourierInterp, f Ss) {

	n := float64(fou.N)
	for j := 0; j < fou.N; j++ {
		xj := 2.0 * math.Pi * float64(j) / n
		fx := f(xj)
		chk.AnaNum(tst, io.Sf("I{f}(%5.3f)", xj), 1e-15, fx, fou.I(xj), chk.Verbose)
	}
}

// check derivatives @ notes
func fouCheckD1andD2(tst *testing.T, fou *FourierInterp, f Ss) {

	// check first derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.Idiff(1, x), x, 1e-3, chk.Verbose, func(t float64) float64 {
			return fou.I(t)
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-9, fou.Idiff(2, x), x, 1e-3, chk.Verbose, func(t float64) float64 {
			return fou.Idiff(1, t)
		})
	}
}

func TestFourierInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp02. interpolation using DFT")

	// function and analytic derivative
	f := func(x float64) float64 { return math.Sin(x / 2.0) }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou := NewFourierInterp(N, "")
	defer fou.Free()

	// compute A using 3/2-rule
	fou.CalcAwithAliasRemoval(f)
	io.Pf("\n............3/2-rule.............\n")
	fouCheckD1andD2(tst, fou, f)

	// compute U and A (standard)
	fou.CalcU(f)
	fou.CalcA()
	io.Pf("\n.......no aliasing removal.......\n")
	fouCheckI(tst, fou, f)
	fouCheckD1andD2(tst, fou, f)
}

func TestFourierInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp03. square wave")

	// function and analytic derivative
	f := func(x float64) float64 { return Boxcar(x-math.Pi/2, 0, math.Pi) }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou := NewFourierInterp(N, "lanc")
	defer fou.Free()

	// compute U and A
	fou.CalcU(f)
	fou.CalcA()

	// check first derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.Idiff(1, x), x, 1e-3, chk.Verbose, func(t float64) float64 {
			return fou.I(t)
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-9, fou.Idiff(2, x), x, 1e-3, chk.Verbose, func(t float64) float64 {
			return fou.Idiff(1, t)
		})
	}
}

func TestFourierInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp04. setting U externally")

	// function
	f := func(x float64) float64 { return math.Sin(x / 2.0) }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou := NewFourierInterp(N, "")
	defer fou.Free()

	// calc fvals and set U
	fvals := make([]float64, N)
	for j := 0; j < N; j++ {
		fvals[j] = f(fou.X[j])
	}
	fou.U = fvals

	// compute A[k] using fvals
	fou.CalcA()
	Afvals := fou.A.GetCopy()

	// compute A[k] using f(x)
	fou.CalcU(f)
	fou.CalcA()

	// check
	chk.ArrayC(tst, "A", 1e-17, Afvals, fou.A)
}

func TestFourierInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp05. Derivative @ grid points ⇒ DFT")

	// function
	f := func(x float64) float64 { return math.Sin(x / 2.0) }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou := NewFourierInterp(N, "")
	defer fou.Free()

	// compute U and A
	fou.CalcU(f)
	fou.CalcA()

	// compute 1st derivatives @ grid points
	fou.CalcD1()
	fou.CalcD(1)
	chk.Array(tst, "D1=D(1)", 1e-17, fou.Du1, fou.Du)

	// compute 2nd derivatives @ grid points
	fou.CalcD2()
	fou.CalcD(2)
	chk.Array(tst, "D2=D(2)", 1e-17, fou.Du2, fou.Du)

	// check
	for j, x := range fou.X {
		d1 := fou.Idiff(1, x)
		d2 := fou.Idiff(2, x)
		chk.AnaNum(tst, "d1", 1e-15, fou.Du1[j], d1, chk.Verbose)
		chk.AnaNum(tst, "d2", 1e-15, fou.Du2[j], d2, chk.Verbose)
	}
}
