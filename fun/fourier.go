// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/fftw"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Dft1d computes the discrete Fourier transform (DFT) in 1D.
// It replaces data by its discrete Fourier transform, if inverse==false
// or replaces data by its inverse discrete Fourier transform, if inverse==true
//
//   Computes:
//                      N-1         -i 2 π k l / N
//               X[l] =  Σ  x[k] ⋅ e
//                      k=0
//
//   NOTE: (1) ideally, N=len(data) is an integer power of 2.
//         (2) using FFTW
//
func Dft1d(data []complex128, inverse bool) (err error) {
	plan, err := fftw.NewPlan1d(data, inverse, false)
	if err != nil {
		return
	}
	defer plan.Free()
	plan.Execute()
	return
}

// dft1dslow computes the discrete Fourier transform of x (complex) by using the "slow" method; i.e.
// by directly computing the DFT summation using N² operations.
//   NOTE: This function is useful for verifications (testing) only.
func dft1dslow(x []complex128) (X []complex128) {
	N := len(x)
	X = make([]complex128, N)
	for n := 0; n < N; n++ {
		for k := 0; k < N; k++ {
			a := 2.0 * math.Pi * float64(k*n) / float64(N)
			X[n] += x[k] * ExpMix(a) // x[k]⋅exp(-a)
		}
	}
	return
}

// FourierInterp performs interpolation using truncated Fourier series
//
//                  N/2 - 1
//                   ————          +i k X[j]
//        f(x[j]) =  \     A[k] ⋅ e              with    X[j] = 2 π j / N
//                   /
//                   ————
//                  k = -N/2
//
//   x ϵ [0, 2π]
//
//   Equation (2.1.27) of [1]. Note that f=u in [1] and A[k] is the tilde(u[k]) of [1]
//
//   Reference:
//     [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
//         Single Domains. Springer. 563p
//
type FourierInterp struct {
	N int          // number of terms. must be power of 2; i.e. N = 2ⁿ
	F Ss           // f(x)
	A []complex128 // coefficients for interpolation. from FFT
	X []float64    // grid (len=N): X[j] = 2 π j / N (excluding last point ⇒ periodic)
}

// NewFourierInterp allocates a new FourierInterp object
//
//                    N - 1
//                1   ————             -i k X[j]
//        A[k] = ———  \     f(x[j]) ⋅ e              with    X[j] = 2 π j / N
//                N   /
//                    ————
//                   j = 0
//
//   x ϵ [0, 2π]
//
//   Equation (2.1.25) of [1]. Note that f=u in [1] and A[k] is the tilde(u[k]) of [1]
//
//   Reference:
//     [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
//         Single Domains. Springer. 563p
//
func NewFourierInterp(N int, f Ss) (o *FourierInterp, err error) {

	// check
	if N%2 != 0 {
		err = chk.Err("N must be even. N=%d is invalid\n", N)
		return
	}

	// allocate
	o = new(FourierInterp)
	o.N = N
	o.F = f
	o.A = make([]complex128, o.N)
	o.X = make([]float64, o.N)

	// compute grid coordinates and F(X[i])
	π := math.Pi
	n := float64(o.N)
	for i := 0; i < o.N; i++ {
		o.X[i] = 2.0 * π * float64(i) / n
		fx, e := o.F(o.X[i])
		if e != nil {
			err = e
			return
		}
		o.A[i] = complex(fx/n, 0)
	}

	// perform Fourier transform to find A
	err = Dft1d(o.A, false)
	return
}

// I computes the interpolation
//
//                  N/2 - 1
//                    ————          +i k x
//        I {f}(x) =  \     A[k] ⋅ e
//         N          /
//                    ————
//                   k = -N/2
//
//   x ϵ [0, 2π]
//
//   Equation (2.1.28) of [1]. Note that f=u in [1] and A[k] is the tilde(u[k]) of [1]
//
//   Reference:
//     [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
//         Single Domains. Springer. 563p
//
func (o *FourierInterp) I(x float64) float64 {
	var res complex128
	for i := 0; i < o.N/2; i++ {

		// first half
		k := i
		res += o.A[i] * cmplx.Exp(complex(0, float64(k)*x))

		// second half
		ii := o.N/2 + i
		kk := ii - o.N
		res += o.A[ii] * cmplx.Exp(complex(0, float64(kk)*x))

	}
	return real(res)
}

// DI computes the p-derivative of the interpolation
//
//                   p       N/2 - 1
//        p         d(I{f})    ————       p           +i k x
//       DI{f}(x) = ——————— =  \     (i⋅k)  ⋅ A[k] ⋅ e
//        N             p      /
//                    dx       ————
//                            k = -N/2
//
//   x ϵ [0, 2π]
//
func (o *FourierInterp) DI(p int, x float64) float64 {
	var res complex128
	for i := 0; i < o.N/2; i++ {

		// first half
		k := i
		kf := float64(k)
		cf := cmplx.Pow(complex(0, kf), complex(float64(p), 0)) // TODO: simplify this
		res += cf * o.A[i] * cmplx.Exp(complex(0, kf*x))

		// second half
		ii := o.N/2 + i
		kk := ii - o.N
		kkf := float64(kk)
		ccf := cmplx.Pow(complex(0, kkf), complex(float64(p), 0)) // TODO: simplify this
		res += ccf * o.A[ii] * cmplx.Exp(complex(0, kkf*x))
	}
	return real(res)
}

// Plot plots interpolated curve
//   option -- 1: plot only f(x)
//             2: plot both f(x) and df/dx(x)
//             3: plot all f(x), df/dx(x) and d^2f/dx^2
//             4: plot only df/dx(x)
//             5: plot only d^2f/dx^2(x)
//             6: plot df^p/dx^p
//   p      -- order of the derivative to plot if option == 6
//   dfdx   -- is the analytic df/dx(x) [optional]
//   d2fdx2 -- is the analytic d^2f/dx^2(x) [optional]
func (o *FourierInterp) Plot(option, p int, dfdx, d2fdx2 Ss, argsF, argsI, argsX *plt.A) {
	if argsF == nil {
		argsF = &plt.A{L: "f(x)", C: "#0034ab", NoClip: true}
	}
	if argsI == nil {
		argsI = &plt.A{L: "I{f}(x)", C: "r", NoClip: true}
	}
	if argsX == nil {
		argsX = &plt.A{C: "k", M: "o", Ms: 3, Ls: "none", Void: true, NoClip: true}
	}
	npts := 2001
	xx := utl.LinSpace(0, 2.0*math.Pi, npts)
	yX := make([]float64, len(o.X))
	withF := option == 1 || option == 2 || option == 3
	firstD := option == 2 || option == 3 || option == 4
	secondD := option == 3 || option == 5
	var y1, y2 []float64
	if withF {
		y1 = make([]float64, npts)
		y2 = make([]float64, npts)
	}
	var y3, y4 []float64
	if firstD {
		y3 = make([]float64, npts)
		y4 = make([]float64, npts)
	}
	var y5, y6 []float64
	if secondD {
		y5 = make([]float64, npts)
		y6 = make([]float64, npts)
	}
	var y7 []float64
	if option == 6 {
		y7 = make([]float64, npts)
	}
	for i := 0; i < npts; i++ {
		x := xx[i]
		if withF {
			fx, err := o.F(x)
			if err != nil {
				chk.Panic("f(x) failed:\n%v\n", err)
			}
			y1[i] = fx
			y2[i] = o.I(x)
		}
		if firstD {
			if dfdx != nil {
				dfx, err := dfdx(x)
				if err != nil {
					chk.Panic("df/dx(x) failed:\n%v\n", err)
				}
				y3[i] = dfx
			}
			y4[i] = o.DI(1, x)
		}
		if secondD {
			if d2fdx2 != nil {
				ddfx, err := d2fdx2(x)
				if err != nil {
					chk.Panic("d2f/dx2(x) failed:\n%v\n", err)
				}
				y5[i] = ddfx
			}
			y6[i] = o.DI(2, x)
		}
		if option == 6 {
			y7[i] = o.DI(p, x)
		}
	}
	if option == 2 {
		plt.Subplot(2, 1, 1)
	}
	if option == 3 {
		plt.Subplot(3, 1, 1)
	}
	if withF {
		plt.Plot(o.X, yX, argsX)
		plt.Plot(xx, y1, argsF)
		plt.Plot(xx, y2, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", "$f(x)$", nil)
	}
	if option == 2 {
		plt.Subplot(2, 1, 2)
	}
	if option == 3 {
		plt.Subplot(3, 1, 2)
	}
	if firstD {
		argsF.L = "D1"
		argsI.L = "D1"
		plt.Plot(o.X, yX, argsX)
		if dfdx != nil {
			plt.Plot(xx, y3, argsF)
		}
		plt.Plot(xx, y4, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", "$\\frac{\\mathrm{d}f(x)}{\\mathrm{d}x}$", nil)
	}
	if option == 3 {
		plt.Subplot(3, 1, 3)
	}
	if secondD {
		argsF.L = "D2"
		argsI.L = "D2"
		plt.Plot(o.X, yX, argsX)
		if d2fdx2 != nil {
			plt.Plot(xx, y5, argsF)
		}
		plt.Plot(xx, y6, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", "$\\frac{\\mathrm{d}^2f(x)}{\\mathrm{d}x^2}$", nil)
	}
	if option == 6 {
		argsI.L = io.Sf("D%d", p)
		plt.Plot(o.X, yX, argsX)
		plt.Plot(xx, y7, argsI)
		plt.HideTRborders()
		plt.Gll("$x$", io.Sf("$\\frac{\\mathrm{d}^{%d}f(x)}{\\mathrm{d}x^{%d}}$", p, p), nil)
	}
}
