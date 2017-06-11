// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/chk"
)

// FourierTransLL (LL:low-level) computes the discrete Fourier transform.
// It replaces data[0..2*n-1] by its discrete Fourier transform, if inverse==false
// or replaces data[0..2*n-1] by its inverse discrete Fourier transform, if inverse==true
//
//   Computes:
//                      N-1         -i 2 π k l / N
//               X[l] =  Σ  x[k] ⋅ e
//                      k=0
//
//   Notes: (a) n=len(data)/2 must be an integer power of 2.
//          (b) the sign of the direct transform is opposite in [1]
//
//   Input:
//     data -- is a complex array stored as a real array of length 2*n. [real,imag, real,imag, ...]
//     inverse -- computes inverse FFT
//
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//
//   NOTE: if possible, use the fun/fftw package that may be up to 5 times faster than this function
//
func FourierTransLL(data []float64, inverse bool) (err error) {

	// check length of data
	ldata := len(data)
	if ldata < 4 || ldata%2 > 0 {
		err = chk.Err("len(data)=2*n must be greater than 4 and must be even. %d is invalid\n", ldata)
		return
	}

	// check for power of two
	n := ldata / 2
	if n < 2 || !IsPowerOfTwo(n) {
		err = chk.Err("n=len(data)/2 must be power of 2. n=%d is invalid\n", n)
		return
	}

	// this is the bit-reversal section of the routine.
	var m int
	nn := n << 1
	j := 1
	for i := 1; i < nn; i += 2 {
		if j > i {
			Swap(&data[j-1], &data[i-1]) // Exchange the two complex numbers.
			Swap(&data[j], &data[i])
		}
		m = n
		for m >= 2 && j > m {
			j -= m
			m >>= 1
		}
		j += m
	}

	// set isign
	isign := -1.0 // direct transform. note that this is opposite than what's used in [1]
	if inverse {
		isign = 1.0
	}

	// here begins the Danielson-Lanczos section of the routine.
	var istep int
	var wtemp, wr, wpr, wpi, wi, theta, tempr, tempi float64
	mmax := 2
	for nn > mmax { // outer loop executed log2(n) times.
		istep = mmax << 1
		theta = isign * (2.0 * math.Pi / float64(mmax)) // initialize the trigonometric recurrence.
		wtemp = math.Sin(0.5 * theta)
		wpr = -2.0 * wtemp * wtemp
		wpi = math.Sin(theta)
		wr = 1.0
		wi = 0.0
		for m = 1; m < mmax; m += 2 { // here are the two nested inner loops.
			for i := m; i <= nn; i += istep {
				j = i + mmax // this is the Danielson-Lanczos formula:
				tempr = wr*data[j-1] - wi*data[j]
				tempi = wr*data[j] + wi*data[j-1]
				data[j-1] = data[i-1] - tempr
				data[j] = data[i] - tempi
				data[i-1] += tempr
				data[i] += tempi
			}
			wtemp = wr
			wr = wr*wpr - wi*wpi + wr // trigonometric recurrence.
			wi = wi*wpr + wtemp*wpi + wi
		}
		mmax = istep
	}

	// fix inverse results
	if inverse {
		mul := 1.0 * float64(n)
		for i := 0; i < n; i++ {
			data[i] *= mul
		}
	}
	return
}

// DftSlow computes the discrete Fourier transform of x (complex) by using the "slow" method; i.e.
// by directly computing the summation with N² operations
// NOTE: This function is mostly useful for verifications only.
func DftSlow(x []complex128) (X []complex128) {
	N := len(x)
	X = make([]complex128, N)
	for n := 0; n < N; n++ {
		for k := 0; k < N; k++ {
			a := 2.0 * math.Pi * float64(k*n) / float64(N)
			X[n] += x[k] * cmplx.Exp(-1i*complex(a, 0))
		}
	}
	return
}
