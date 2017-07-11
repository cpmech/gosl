// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "math"

// ChebyshevT directly computes the Chebyshev polynomial of first kind Tn(x) using the trigonometric
// functions.
//
//           │ (-1)ⁿ cosh[n⋅acosh(-x)]   if x < -1
//   Tn(x) = ┤       cosh[n⋅acosh( x)]   if x > 1
//           │       cos [n⋅acos ( x)]   if |x| ≤ 1
//
func ChebyshevT(n int, x float64) float64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	nn := float64(n)
	if x < -1 {
		return math.Pow(-1, nn) * math.Cosh(nn*math.Acosh(-x))
	}
	if x > 1 {
		return math.Cosh(nn * math.Acosh(x))
	}
	return math.Cos(nn * math.Acos(x))
}

// ChebyshevXgauss computes Chebyshev-Gauss roots considering symmetry
//
//                       /  (2i+1)⋅π  \
//           X[i] = -cos | —————————— |       i = 0 ... N
//                       \   2N + 2   /
//
func ChebyshevXgauss(N int) (X []float64) {
	X = make([]float64, N+1)
	n := float64(N)
	d := 2.0*n + 2.0
	if (N & 1) == 0 { // even number of segments
		for i := 0; i < N/2; i++ {
			X[N-i] = math.Cos(float64(2*i+1) * math.Pi / d)
			X[i] = -X[N-i]
		}
	} else { // odd number of segments
		l := (N + 3) / 2
		for i := 0; i < l; i++ {
			X[N-i] = math.Cos(float64(2*i+1) * math.Pi / d)
			if i < l-1 {
				X[i] = -X[N-i]
			}
		}
	}
	return
}

// ChebyshevXlob computes Chebyshev-Gauss-Lobatto points using the sin function and
// considering symmetry
//
//                       /  π⋅(N-2i)  \
//           X[i] = -sin | —————————— |       i = 0 ... N
//                       \    2⋅N     /
//     or
//                       /  i⋅π  \
//           X[i] = -cos | ————— |       i = 0 ... N
//                       \   N   /
//
//   Reference:
//   [1] Baltensperger R and Trummer MR (2003) Spectral differencing with a twist, SIAM J. Sci.
//       Comput. 24(5):1465-1487
//
func ChebyshevXlob(N int) (X []float64) {
	X = make([]float64, N+1)
	X[0] = -1
	X[N] = +1
	if N < 3 {
		return
	}
	n := float64(N)
	d := 2.0 * n
	if (N & 1) == 0 { // even number of segments
		for i := 1; i < N/2; i++ {
			X[N-i] = math.Sin(π * (n - 2.0*float64(i)) / d)
			X[i] = -X[N-i]
		}
	} else { // odd number of segments
		l := (N + 3) / 2
		for i := 1; i < l; i++ {
			X[N-i] = math.Sin(π * (n - 2.0*float64(i)) / d)
			if i < l-1 {
				X[i] = -X[N-i]
			}
		}
	}
	return
}
