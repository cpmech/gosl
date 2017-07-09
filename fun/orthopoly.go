// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// ChebyshevPoly defines a structure for efficient computations with orthogonal polynomials
type ChebyshevPoly struct {

	// input
	Kind io.Enum // type of orthogonal polynomial
	N    int     // degree of polynomial

	// derived
	X     []float64 // points
	Wb    []float64 // weights for Gaussian quadrature
	Gamma []float64 // denominador of coefficients equation ~ ‖p[i]‖²
	CoefI []float64 // coefficients of interpolant
	CoefP []float64 // coefficients of projection (estimated)

	// constants
	EstimationN int // N to use when estimating CoefP [default=128]
}

// NewChebyshevPoly returns a new ChebyshevPoly structure
//
//   gaussChebyshev == true:
//
//                       /  (2i+1)⋅π  \
//           X[i] = -cos | —————————— |       i = 0 ... N
//                       \   2N + 2   /
//
//   gaussChebyshev == false: (Gauss-Lobatto-Chebyshev)
//
//                       /  i⋅π  \
//           X[i] = -cos | ————— |       i = 0 ... N
//                       \   N   /
//
func NewChebyshevPoly(N int, gaussChebyshev bool) (o *ChebyshevPoly, err error) {

	// allocate
	o = new(ChebyshevPoly)
	o.N = N
	o.X = make([]float64, N+1)
	o.Wb = make([]float64, N+1)
	o.Gamma = make([]float64, N+1)
	o.CoefI = make([]float64, N+1)
	o.CoefP = make([]float64, N+1)
	o.EstimationN = 128

	// set data
	wb, wb0, wbN, gam, gam0, gamN, a, b := o.gaussData(o.N, gaussChebyshev)
	for i := 0; i < o.N+1; i++ {
		o.X[i] = -math.Cos(a * b(i))
		o.Wb[i] = wb
		o.Gamma[i] = gam
	}
	o.Wb[0] = wb0
	o.Wb[o.N] = wbN
	o.Gamma[0] = gam0
	o.Gamma[o.N] = gamN
	return
}

// gaussData returns quadrature data for either Gauss-Chebyshev or Gauss-Lobatto
func (o *ChebyshevPoly) gaussData(N int, gCheby bool) (wb, wb0, wbN, gam, gam0, gamN, a float64, b func(i int) float64) {
	if gCheby {
		wb = π / float64(N+1)
		wb0 = wb
		wbN = wb
		gam = π / 2.0
		gam0 = π
		gamN = π / 2.0
		a = π / float64(2*N+2)
		b = func(i int) float64 { return float64(2*i + 1) }
	} else {
		wb = π / float64(N)
		wb0 = π / float64(2*N)
		wbN = wb0
		gam = π / 2.0
		gam0 = π
		gamN = π
		a = π / float64(N)
		b = func(i int) float64 { return float64(i) }
	}
	return
}

// Phi computes φi(x) (polynomial function) of Chebyshev polynomial
func (o *ChebyshevPoly) Phi(i int, x float64) float64 {
	return ChebyshevT(i, x)
}

// W computes W(x) (weight function) of Chebyshev polynomial
func (o *ChebyshevPoly) W(x float64) float64 {
	return 1.0 / math.Sqrt(1.0-x*x)
}

// Approx computes the approximated projection or interpolation via series approximation
// after computing the coefficients of the interpolant or estimated projection
//
//    Approx(x) = Σ a[i] * φ[i](x)  where  'a' is CoefI or CoefP (if projection==true)
//
func (o *ChebyshevPoly) Approx(x float64, projection bool) (res float64) {
	a := o.CoefI
	if projection {
		a = o.CoefP
	}
	for i := 0; i < o.N+1; i++ {
		res += a[i] * o.Phi(i, x)
	}
	return
}

// EstimateMaxErr estimates the maximum error using 10000 stations along [-1,1]
// This function also returns the location (xloc) of the estimated max error
//
//    maxerr = max(|f - I{f}|)  or  maxerr = max(|f - Π{f}|)
//
func (o *ChebyshevPoly) EstimateMaxErr(f Ss, projection bool) (maxerr, xloc float64) {
	nsta := 10000 // generate several points along [-1,1]
	xloc = -1
	for i := 0; i < nsta; i++ {
		x := -1.0 + 2.0*float64(i)/float64(nsta-1)
		fx, err := f(x)
		if err != nil {
			chk.Panic("f(x) failed:%v\n", err)
		}
		fa := o.Approx(x, projection)
		e := math.Abs(fx - fa)
		if e > maxerr {
			maxerr = e
			xloc = x
		}
	}
	return
}

// CoefInterpolantSlow computes the coefficients of the interpolant by (slow) formula
//
//   A[i] = Σ f(x[i]) ⋅ φi(x[i]) ⋅ wb[i]
//
//   NOTE: the results will be stored in o.CoefI
func (o *ChebyshevPoly) CoefInterpolantSlow(f Ss) (err error) {

	// evaluate function at all points
	fx := make([]float64, o.N+1)
	for i := 0; i < o.N+1; i++ {
		fx[i], err = f(o.X[i])
		if err != nil {
			return
		}
	}

	// computation of coefficients
	for i := 0; i < o.N+1; i++ {
		o.CoefI[i] = 0
		for j := 0; j < o.N+1; j++ {
			o.CoefI[i] += fx[j] * o.Phi(i, o.X[j]) * o.Wb[j]
		}
		o.CoefI[i] /= o.Gamma[i]
	}
	return
}

// EstimateCoefProjection computes the coefficients of the projection (slow)
// using GaussChebyshev quadrature and EstimationN + 1 points
//   NOTE: the results will be stored in o.CoefP
func (o *ChebyshevPoly) EstimateCoefProjection(f Ss) (err error) {

	// quadrature data
	nn := o.EstimationN
	gaussChebyshev := true
	wb, wb0, wbN, gam, gam0, gamN, a, b := o.gaussData(nn, gaussChebyshev)

	// evaluate function at many points
	xx := make([]float64, nn+1)
	fx := make([]float64, nn+1)
	for i := 0; i < nn+1; i++ {
		xx[i] = -math.Cos(a * b(i))
		fx[i], err = f(xx[i])
		if err != nil {
			return
		}
	}

	// computation of coefficients using GaussChebyshev quadrature
	for i := 0; i < o.N+1; i++ {
		o.CoefP[i] = 0
		for j := 1; j < nn; j++ {
			o.CoefP[i] += fx[j] * o.Phi(i, xx[j]) * wb
		}
		o.CoefP[i] += fx[0] * o.Phi(i, xx[0]) * wb0
		o.CoefP[i] += fx[nn] * o.Phi(i, xx[nn]) * wbN
		if i == 0 {
			o.CoefP[i] /= gam0
		} else if i == nn {
			o.CoefP[i] /= gamN
		} else {
			o.CoefP[i] /= gam
		}
	}
	return
}
