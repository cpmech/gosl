// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// ChebyshevPoly defines a structure for efficient computations with Chebyshev polynomials
type ChebyshevPoly struct {

	// input
	N     int  // degree of polynomial
	Gauss bool // use roots (Gauss) or points (Lobatto)?

	// derived
	X     []float64 // points
	Wb    []float64 // weights for Gaussian quadrature
	Gamma []float64 // denominador of coefficients equation ~ ‖p[i]‖²
	CoefI []float64 // coefficients of interpolant
	CoefP []float64 // coefficients of projection (estimated)

	// constants
	EstimationN int // N to use when estimating CoefP [default=128]

	// computed
	C   *la.Matrix // physical to transform space conversion matrix
	Ci  *la.Matrix // transform to physical space conversion matrix
	Psi []float64  // basis functions @ Gauss-Lobatto points
	D1  *la.Matrix // (dψj/dx)(xi)
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
	o.Gauss = gaussChebyshev
	o.Wb = make([]float64, N+1)
	o.Gamma = make([]float64, N+1)
	o.CoefI = make([]float64, N+1)
	o.CoefP = make([]float64, N+1)
	o.EstimationN = 128

	// roots or points
	if gaussChebyshev {
		o.X = ChebyshevXgauss(N)
	} else {
		o.X = ChebyshevXlob(N)
	}

	// set data
	wb, wb0, wbN, gam, gam0, gamN, _, _ := o.gaussData(o.N)
	for i := 0; i < o.N+1; i++ {
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
func (o *ChebyshevPoly) gaussData(N int) (wb, wb0, wbN, gam, gam0, gamN, a float64, b func(i int) float64) {
	if o.Gauss {
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

// P computes p_i(x)=T_i(x) (polynomial function) of Chebyshev polynomial
//func (o *ChebyshevPoly) P(i int, x float64) float64 {
//return ChebyshevT(i, x)
//}

// W computes W(x) (weight function) of Chebyshev polynomial
func (o *ChebyshevPoly) W(x float64) float64 {
	return 1.0 / math.Sqrt(1.0-x*x)
}

// Approx2 computes the approximated projection or interpolation via series approximation
// after computing the coefficients of the interpolant or estimated projection
//
//    Approx(x) = Σ a[i] * φ[i](x)  where  'a' is CoefI or CoefP (if projection==true)
//
func (o *ChebyshevPoly) Approx2(x float64, projection bool) (res float64) {
	a := o.CoefI
	if projection {
		a = o.CoefP
	}
	for i := 0; i < o.N+1; i++ {
		res += a[i] * ChebyshevT(i, x)
	}
	return
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
		res += a[i] * ChebyshevT(i, x)
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
			o.CoefI[i] += fx[j] * ChebyshevT(i, o.X[j]) * o.Wb[j]
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
	wb, wb0, wbN, gam, gam0, gamN, a, b := o.gaussData(nn)

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
			o.CoefP[i] += fx[j] * ChebyshevT(i, xx[j]) * wb
		}
		o.CoefP[i] += fx[0] * ChebyshevT(i, xx[0]) * wb0
		o.CoefP[i] += fx[nn] * ChebyshevT(i, xx[nn]) * wbN
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

// HierarchicalT computes Tn(x) using hierarchical definition (but NOT recursive)
//   NOTE: this function is not as efficient as ChebyshevT and should be used for testing only
func (o *ChebyshevPoly) HierarchicalT(i int, x float64) float64 {
	if i == 0 {
		return 1.0
	}
	if i == 1 {
		return x
	}
	tjm2 := 1.0 // value at step j - 2
	tjm1 := x   // value at step j - 1
	var tj float64
	for j := 2; j <= i; j++ {
		tj = 2*x*tjm1 - tjm2
		tjm2 = tjm1
		tjm1 = tj
	}
	return tjm1
}
