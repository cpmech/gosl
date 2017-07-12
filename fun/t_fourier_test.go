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

	chk.VectorC(tst, "X", 1e-14, x, Y)
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
	chk.Scalar(tst, "maxRf", 1e-12, maxRf, 0.383616856748)

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
			chk.VectorC(tst, "U", 1e-15, U, Ucopy)

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
			chk.VectorC(tst, "U", 1e-15, U, Ucopy)

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
			chk.VectorC(tst, "U", 1e-15, U, Ucopy)
		}
	}
}

func TestFourierInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("FourierInterp01. interpolation using DFT")

	// function and analytic derivative
	f := func(x float64) (float64, error) { return math.Sin(x / 2.0), nil }
	dfdx := func(x float64) (float64, error) { return math.Cos(x/2.0) / 2.0, nil }

	// constants
	var p uint64 = 2 // exponent of 2ⁿ
	N := 2 << p      // 2ⁿ (n=p): number of terms
	fou, err := NewFourierInterp(N, f)
	chk.EP(err)

	// check interpolation @ nodes
	for i := 0; i < N; i++ {
		x := fou.X[i]
		fx, err := f(x)
		chk.EP(err)
		chk.AnaNum(tst, io.Sf("I{f}(%5.3f)", x), 1e-15, fx, fou.I(x), chk.Verbose)
	}

	// check derivative of interpolation
	io.Pl()
	xx := utl.LinSpace(0, 2*math.Pi, 11)
	for i := 0; i < len(xx); i++ {
		x := xx[i]
		chk.DerivScaSca(tst, io.Sf("DI{f}(%5.3f)", x), 1e-10, fou.DI(x), x, 1e-3, chk.Verbose, func(t float64) (float64, error) {
			return fou.I(t), nil
		})
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{Prop: 1.5})
		plt.Subplot(2, 1, 1)
		plt.Title(io.Sf("f(x) and interpolation. N=%d", N), &plt.A{Fsz: 9})
		plt.Subplot(2, 1, 2)
		plt.SplotGap(0.0, 0.25)
		plt.Title(io.Sf("df/dx(x) and derivative of interpolation. N=%d", N), &plt.A{Fsz: 9})
		fou.Plot(2, dfdx, nil, nil, nil)
		plt.Save("/tmp/gosl/fun", "fourierinterp01")
	}
}
