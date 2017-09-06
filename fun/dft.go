// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/fun/fftw"
)

// Dft1d computes the discrete Fourier transform (DFT) in 1D.
// It replaces data by its discrete Fourier transform, if inverse==false
// or replaces data by its inverse discrete Fourier transform, if inverse==true
//
//   Computes:
//                      N-1         -i 2 π j k / N                 __
//     forward:  X[k] =  Σ  x[j] ⋅ e                     with i = √-1
//                      j=0
//
//                      N-1         +i 2 π j k / N
//     inverse:  Y[k] =  Σ  y[j] ⋅ e                     thus x[k] = Y[k] / N
//                      j=0
//
//   NOTE: (1) the inverse operation does not divide by N
//         (2) ideally, N=len(data) is an integer power of 2.
//         (3) using FFTW: http://fftw.org/fftw3_doc/What-FFTW-Really-Computes.html
//
func Dft1d(data []complex128, inverse bool) {
	plan := fftw.NewPlan1d(data, inverse, false)
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
