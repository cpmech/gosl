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

func TestDft01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dft01. FFT")

	x := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	err := Dft1d(x, false)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	y := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	Y := dft1dslow(y)

	chk.ArrayC(tst, "X", 1e-14, x, Y)
}

func TestDft02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dft02. FFT sinusoid")

	// set sinusoid equation
	T := 1.0 / 5.0      // period [s]
	A0 := 0.0           // mean value
	C1 := 1.0           // amplitude
	θ := -math.Pi / 2.0 // phase shift [rad]
	ss := NewSinusoidEssential(T, A0, C1, θ)

	// discrete data
	N := 16
	dt := 1.0 / float64(N-1)
	tt := make([]float64, N)      // time
	xx := make([]float64, N)      // x[n]
	data := make([]complex128, N) // x[n] to use as input of FFT
	for i := 0; i < N; i++ {
		tt[i] = float64(i) * dt
		xx[i] = ss.Ybasis(tt[i])
		data[i] = complex(xx[i], 0)
	}

	// execute FFT
	err := Dft1d(data, false)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// extract results
	Xr := make([]float64, N) // real(X[n])
	Xi := make([]float64, N) // imag(X[n])
	Rf := make([]float64, N) // |X[n]|/n
	maxRf := 0.0
	for k := 0; k < N; k++ {
		Xr[k] = real(data[k])
		Xi[k] = imag(data[k])
		Rf[k] = math.Sqrt(Xr[k]*Xr[k]+Xi[k]*Xi[k]) / float64(N)
		if Rf[k] > maxRf {
			maxRf = Rf[k]
		}
	}
	io.Pforan("maxRf = %v\n", maxRf)
	chk.Float64(tst, "maxRf", 1e-12, maxRf, 0.383616856748)

	// plot
	if chk.Verbose {
		ts := utl.LinSpace(0, 1, 201)
		xs := make([]float64, len(ts))
		for i := 0; i < len(ts); i++ {
			xs[i] = ss.Ybasis(ts[i])
		}
		fn := utl.LinSpace(0, float64(N), N)

		plt.Reset(true, &plt.A{Prop: 1.2})

		plt.Subplot(3, 1, 1)
		plt.Plot(ts, xs, &plt.A{C: "b", L: "continuous signal", NoClip: true})
		plt.Plot(tt, xx, &plt.A{C: "r", M: ".", L: "discrete signal", NoClip: true})
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Gll("t", "x(t)", &plt.A{LegOut: true, LegNcol: 3})

		plt.Subplot(3, 1, 2)
		plt.Plot(tt, Xr, &plt.A{C: "r", M: ".", L: "real(X)", NoClip: true})
		plt.HideAllBorders()
		plt.Gll("t", "f(t)", &plt.A{LegOut: true, LegNcol: 3})

		plt.Subplot(3, 1, 3)
		plt.Plot(fn, Rf, &plt.A{C: "m", M: ".", NoClip: true})
		plt.HideAllBorders()
		plt.Gll("freq", "|X(f)|/n", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl/fun", "dft02")
	}
}

func TestDft03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Dft03. FFT and inverse FFT")

	// function
	π := math.Pi
	f := func(x float64) float64 { return math.Sin(x / 2.0) }

	// data
	N := 4 // number of terms
	U := make([]complex128, N)
	Ucopy := make([]complex128, N)

	// run with 3 places for performing normalisation
	for place := 1; place <= 3; place++ {

		// message
		io.Pf("\n\n~~~~~~~~~~~~~~~~~~~~ place = %v ~~~~~~~~~~~~~~~~~~~~~~~~\n", place)

		// f @ points
		for i := 0; i < N; i++ {
			x := 2.0 * π * float64(i) / float64(N)
			U[i] = complex(f(x), 0)
			Ucopy[i] = U[i]
		}
		io.Pf("before: U = %.3f\n", U)

		switch place {

		// normalise at the beginning
		case 1:

			// normalise
			for i := 0; i < N; i++ {
				U[i] /= complex(float64(N), 0)
			}
			io.Pfblue2("normalised\n")

			// execute FFT
			err := Dft1d(U, false)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pforan("FFT(U) = %.3f\n", U)

			// execute inverse FFT
			err = Dft1d(U, true)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pf("invFFT(U) = %.3f\n", U)
			chk.ArrayC(tst, "U", 1e-15, U, Ucopy)

		// normalise after direct FFT
		case 2:

			// execute FFT
			err := Dft1d(U, false)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pforan("FFT(U) = %.3f\n", U)

			// normalise
			for i := 0; i < N; i++ {
				U[i] /= complex(float64(N), 0)
			}
			io.Pfblue2("normalised\n")

			// execute inverse FFT
			err = Dft1d(U, true)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pf("invFFT(U) = %.3f\n", U)
			chk.ArrayC(tst, "U", 1e-15, U, Ucopy)

		// normalise after inverse FFT
		case 3:

			// execute FFT
			err := Dft1d(U, false)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pforan("FFT(U) = %.3f\n", U)

			// execute inverse FFT
			err = Dft1d(U, true)
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}
			io.Pf("invFFT(U) = %.3f\n", U)

			// normalise
			for i := 0; i < N; i++ {
				U[i] /= complex(float64(N), 0)
			}
			io.Pfblue2("normalised\n")

			// check
			chk.ArrayC(tst, "U", 1e-15, U, Ucopy)
		}
	}
}

func TestFourierInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp01. check k(j)")

	// constants
	N := 8
	fou, err := NewFourierInterp(N, 0)
	chk.EP(err)

	// check k
	kvals := make([]float64, N)
	for j := 0; j < N; j++ {
		kvals[j] = fou.K(j)
	}
	chk.Array(tst, "k(j)", 1e-17, kvals, []float64{0, 1, 2, 3, -4, -3, -2, -1})
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
	fou, err := NewFourierInterp(N, SmoNoneKind)
	chk.EP(err)

	// compute A[k]
	err = fou.CalcA(f)
	chk.EP(err)

	// check interpolation @ nodes
	n := float64(N)
	for j := 0; j < N; j++ {
		xj := 2.0 * math.Pi * float64(j) / n
		fx, err := f(xj)
		chk.EP(err)
		chk.AnaNum(tst, io.Sf("I{f}(%5.3f)", xj), 1e-15, fx, fou.I(xj), chk.Verbose)
	}

	// check first derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.DI(1, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.I(t), nil
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-10, fou.DI(2, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.DI(1, t), nil
		})
	}

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

		fou.Plot(3, 3, f, dfdx, d2fdx2, nil, nil, nil, nil)
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
	fou, err := NewFourierInterp(N, SmoLanczosKind)
	chk.EP(err)

	// compute A[k]
	err = fou.CalcA(f)
	chk.EP(err)

	// check first derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D1I{f}(%5.3f)", x), 1e-10, fou.DI(1, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.I(t), nil
		})
	}

	// check second derivative of interpolation
	io.Pl()
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("D2I{f}(%5.3f)", x), 1e-9, fou.DI(2, x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.DI(1, t), nil
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
			fou, err := NewFourierInterp(N, SmoLanczosKind)
			chk.EP(err)
			err = fou.CalcA(f)
			chk.EP(err)

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
			fou, err := NewFourierInterp(N, SmoRcosKind)
			chk.EP(err)
			err = fou.CalcA(f)
			chk.EP(err)

			ll := ""
			if k == 2 {
				ll = "Rcos"
			}
			fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ll}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""})
		}

		for k, p := range []uint64{2, 3, 4, 5} {

			N := 1 << p
			fou, err := NewFourierInterp(N, SmoCesaroKind)
			chk.EP(err)
			err = fou.CalcA(f)
			chk.EP(err)

			ll := ""
			if k == 2 {
				ll = "Cesaro"
			}
			fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ll}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""})
		}

		plt.Save("/tmp/gosl/fun", "fourierinterp03")
	}
}
