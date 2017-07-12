// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/fftw"
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

	// perform Fourier transform
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

// Plot plots interpolated curve
func (o *FourierInterp) Plot(argsF, argsI, argsX *plt.A) {
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
	y1 := make([]float64, npts)
	y2 := make([]float64, npts)
	for i := 0; i < npts; i++ {
		x := xx[i]
		fx, err := o.F(x)
		if err != nil {
			chk.Panic("f(x) failed:\n%v\n", err)
		}
		y1[i] = fx
		y2[i] = o.I(x)
	}
	plt.Plot(o.X, yX, argsX)
	plt.Plot(xx, y1, argsF)
	plt.Plot(xx, y2, argsI)
	plt.Gll("$x$", "$f(x)$", nil)
}
