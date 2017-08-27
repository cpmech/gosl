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

func TestFourierInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp01. check k(j) and j(k)")

	// constants
	N := 8
	fou, err := NewFourierInterp(N, "")
	status(tst, err)
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
		fx, err := f(xj)
		status(tst, err)
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
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.Idiff(1, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.I(t), nil
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-9, fou.Idiff(2, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.Idiff(1, t), nil
		})
	}
}

func TestFourierInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp02. interpolation using DFT")

	// function and analytic derivative
	f := func(x float64) (float64, error) { return math.Sin(x / 2.0), nil }
	dfdx := func(x float64) (float64, error) { return math.Cos(x/2.0) / 2.0, nil }
	d2fdx2 := func(x float64) (float64, error) { return -math.Cos(x/2.0) / 4.0, nil }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou, err := NewFourierInterp(N, "")
	status(tst, err)
	defer fou.Free()

	// compute A using 3/2-rule
	status(tst, fou.CalcAwithAliasRemoval(f))
	A32 := fou.A.GetCopy()
	io.Pf("\n............3/2-rule.............\n")
	fouCheckD1andD2(tst, fou, f)

	// compute U and A (standard)
	status(tst, fou.CalcU(f))
	fou.CalcA()
	io.Pf("\n.......no aliasing removal.......\n")
	fouCheckI(tst, fou, f)
	fouCheckD1andD2(tst, fou, f)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 1.7})
		plt.SplotGap(0.0, 0.3)

		plt.Subplot(3, 1, 1)
		plt.Title(io.Sf("f(x) and interpolation. N=%d", N), &plt.A{Fsz: 9})

		plt.Subplot(3, 1, 2)
		plt.Title(io.Sf("df/dx(x) and derivative of interpolation. N=%d", N), &plt.A{Fsz: 9})

		plt.Subplot(3, 1, 3)
		plt.Title(io.Sf("d2f/dx2(x) and second deriv interpolation. N=%d", N), &plt.A{Fsz: 9})

		fou.Plot(3, 3, f, dfdx, d2fdx2, nil,
			&plt.A{L: "noAliasRem", C: plt.C(1, 1), Ls: "--", NoClip: true},
			&plt.A{L: "noAliasRem", C: plt.C(2, 1), Ls: "--", NoClip: true},
			&plt.A{L: "noAliasRem", C: plt.C(3, 1), Ls: "--", NoClip: true})
		copy(fou.A, A32)
		fou.Plot(3, 3, nil, nil, nil, nil, nil, nil, nil)
		plt.Save("/tmp/gosl/fun", "fourierinterp02")
	}
}

func TestFourierInterp03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp03. square wave")

	// function and analytic derivative
	f := func(x float64) (float64, error) { return Boxcar(x-math.Pi/2, 0, math.Pi), nil }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou, err := NewFourierInterp(N, "lanc")
	status(tst, err)
	defer fou.Free()

	// compute U and A
	status(tst, fou.CalcU(f))
	fou.CalcA()

	// check first derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.Idiff(1, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.I(t), nil
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-9, fou.Idiff(2, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.Idiff(1, t), nil
		})
	}

	// plot
	if chk.Verbose {

		plt.Reset(true, &plt.A{Prop: 1.7})
		plt.SplotGap(0.0, 0.3)

		plt.Subplot(3, 1, 1)
		plt.Title("f(x) and interpolation", &plt.A{Fsz: 9})

		plt.Subplot(3, 1, 2)
		plt.Title("df/dx(x) and derivative of interpolation", &plt.A{Fsz: 9})

		plt.Subplot(3, 1, 3)
		plt.Title("d2f/dx2(x) and second deriv interpolation", &plt.A{Fsz: 9})

		for k, p := range []uint64{2, 3, 4, 5} {

			N := 1 << p
			fou, err := NewFourierInterp(N, "lanc")
			status(tst, err)
			defer fou.Free()

			status(tst, fou.CalcU(f))
			fou.CalcA()

			ff := f
			ll := ""
			if k == 2 {
				ff = nil
				ll = "Lanczos"
			}
			l := io.Sf("%d", N)
			fou.Plot(3, 3, ff, nil, nil, &plt.A{C: "k", L: ""}, &plt.A{C: plt.C(k, 0), L: l}, &plt.A{C: plt.C(k, 0), L: ll}, &plt.A{C: plt.C(k, 0), L: l})
		}

		for k, p := range []uint64{2, 3, 4, 5} {

			N := 1 << p
			fou, err := NewFourierInterp(N, "rcos")
			status(tst, err)
			defer fou.Free()

			status(tst, fou.CalcU(f))
			fou.CalcA()

			ll := ""
			if k == 2 {
				ll = "Rcos"
			}
			fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ll}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""})
		}

		for k, p := range []uint64{2, 3, 4, 5} {

			N := 1 << p
			fou, err := NewFourierInterp(N, "ces")
			status(tst, err)
			defer fou.Free()

			status(tst, fou.CalcU(f))
			fou.CalcA()

			ll := ""
			if k == 2 {
				ll = "Cesaro"
			}
			fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ll}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""})
		}

		plt.Save("/tmp/gosl/fun", "fourierinterp03")
	}
}

func TestFourierInterp04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp04. setting U externally")

	// function
	f := func(x float64) (float64, error) { return math.Sin(x / 2.0), nil }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou, err := NewFourierInterp(N, "")
	status(tst, err)
	defer fou.Free()

	// calc fvals and set U
	fvals := make([]float64, N)
	for j := 0; j < N; j++ {
		fvals[j], _ = f(fou.X[j])
	}
	fou.U = fvals

	// compute A[k] using fvals
	fou.CalcA()
	Afvals := fou.A.GetCopy()

	// compute A[k] using f(x)
	status(tst, fou.CalcU(f))
	fou.CalcA()

	// check
	chk.ArrayC(tst, "A", 1e-17, Afvals, fou.A)
}

func TestFourierInterp05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp05. Derivative @ grid points ⇒ DFT")

	// function
	f := func(x float64) (float64, error) { return math.Sin(x / 2.0), nil }
	dfdx := func(x float64) (float64, error) { return math.Cos(x/2.0) / 2.0, nil }
	d2fdx2 := func(x float64) (float64, error) { return -math.Cos(x/2.0) / 4.0, nil }

	// constants
	var p uint64 = 3 // exponent of 2ⁿ
	N := 1 << p      // 2ⁿ (n=p): number of terms
	fou, err := NewFourierInterp(N, "")
	status(tst, err)
	defer fou.Free()

	// compute U and A
	status(tst, fou.CalcU(f))
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

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 1.7})
		plt.SplotGap(0.0, 0.3)

		plt.Subplot(3, 1, 1)
		plt.Title(io.Sf("f(x) and interpolation. N=%d", N), &plt.A{Fsz: 9})

		plt.Subplot(3, 1, 2)
		plt.Title(io.Sf("df/dx(x) and derivative of interpolation. N=%d", N), &plt.A{Fsz: 9})
		plt.Plot(fou.X, fou.Du1, &plt.A{C: "k", M: ".", Ls: "none", Void: true, NoClip: true})

		plt.Subplot(3, 1, 3)
		plt.Title(io.Sf("d2f/dx2(x) and second deriv interpolation. N=%d", N), &plt.A{Fsz: 9})
		plt.Plot(fou.X, fou.Du2, &plt.A{C: "k", M: ".", Ls: "none", Void: true, NoClip: true})

		fou.Plot(3, 3, f, dfdx, d2fdx2, nil, nil, nil, nil)
		plt.Save("/tmp/gosl/fun", "fourierinterp05")
	}
}
