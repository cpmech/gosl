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
	p := float64(n)
	if x < -1 {
		if (n & 1) == 0 { // n is even
			return math.Cosh(p * math.Acosh(-x))
		}
		return -math.Cosh(p * math.Acosh(-x))
	}
	if x > 1 {
		return math.Cosh(p * math.Acosh(x))
	}
	return math.Cos(p * math.Acos(x))
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
		l := (N+3)/2 - 1
		for i := 0; i < l; i++ {
			X[N-i] = math.Cos(float64(2*i+1) * math.Pi / d)
			if i < l {
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
		l := (N+3)/2 - 1
		for i := 1; i < l; i++ {
			X[N-i] = math.Sin(π * (n - 2.0*float64(i)) / d)
			if i < l {
				X[i] = -X[N-i]
			}
		}
	}
	return
}

// ChebyshevTdiff1 computes the first derivative of the Chebyshev function Tn(x)
//
//       dTn
//       ————
//        dx
//
func ChebyshevTdiff1(n int, x float64) float64 {
	p := float64(n)
	if x > -1 && x < +1 {
		t := math.Acos(x)
		d1 := -p * math.Sin(p*t) // derivatives of cos(n⋅t) with respect to t
		s := math.Sin(t)
		return -d1 / s
	}
	if x == +1 {
		return p * p
	}
	if x == -1 {
		if (n & 1) == 0 { // n is even ⇒ n+1 is odd
			return -p * p
		}
		return p * p // n is odd ⇒ n+1 is even
	}
	if x < -1 {
		return -math.Pow(-1, p) * (p * math.Sinh(p*math.Acosh(-x))) / math.Sqrt(x*x-1.0)
	}
	// x > +1
	return (p * math.Sinh(p*math.Acosh(x))) / math.Sqrt(x*x-1.0)
}

// ChebyshevTdiff2 computes the second derivative of the Chebyshev function Tn(x)
//
//       d²Tn
//       —————
//        dx²
//
func ChebyshevTdiff2(n int, x float64) float64 {
	p := float64(n)
	pp := p * p
	t := math.Acos(x)
	if x > -1 && x < +1 {
		d1 := -p * math.Sin(p*t) // derivatives of cos(n⋅t) with respect to t
		d2 := -pp * math.Cos(p*t)
		c := math.Cos(t)
		s := math.Sin(t)
		return (s*d2 - c*d1) / (s * s * s)
	}
	if x == -1 {
		if (n & 1) != 0 { // n is odd
			return -(pp*pp - pp) / 3.0
		}
		return (pp*pp - pp) / 3.0
	}
	if x == +1 {
		return (pp*pp - pp) / 3.0
	}
	d := x*x - 1.0
	if x < -1 {
		r := (pp*math.Cosh(p*math.Acosh(-x)))/d + (p*x*math.Sinh(p*math.Acosh(-x)))/math.Pow(d, 1.5)
		if (n & 1) == 0 { // n is even
			return r
		}
		return -r
	}
	// x > +1
	return (pp*math.Cosh(p*math.Acosh(x)))/d - (p*x*math.Sinh(p*math.Acosh(x)))/math.Pow(d, 1.5)
}
