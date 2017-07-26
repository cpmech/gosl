// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "github.com/cpmech/gosl/num/qpck"

// QuadGen performs automatic integration (quadrature) using the general-purpose
// QUADPACK routine AGSE (Automatic, general-purpose, end-points singularities).
//
//   INPUT:
//     a      -- lower limit of integration
//     b      -- upper limit of integration
//     fid    -- index of goroutine (to avoid race problems)
//     f      -- function defining the integrand
//
//   OUTPUT:          b
//             res = ∫  f(x) dx
//                   a
//
func QuadGen(a, b float64, fid int, f func(x float64) float64) (res float64, err error) {
	id := int32(fid)
	res, _, _, _, err = qpck.Agse(id, f, a, b, 0, 0, nil, nil, nil, nil, nil)
	return
}

// QuadCs performs automatic integration (quadrature) using the cosine or sine weights
// QUADPACK routine AWOE (Automatic with weight, Oscillatory)
//
//   INPUT:
//     a      -- lower limit of integration
//     b      -- upper limit of integration
//     ω      -- omega
//     useSin -- use sin(ω⋅x) instead of cos(ω⋅x)
//     fid    -- index of goroutine (to avoid race problems)
//     f      -- function defining the integrand
//
//   OUTPUT:          b                                     b
//             res = ∫  f(x) ⋅ cos(ω⋅x) dx     or    res = ∫ f(x) ⋅ sin(ω⋅x) dx
//                   a                                     a
//
func QuadCs(a, b, ω float64, useSin bool, fid int, f func(x float64) float64) (res float64, err error) {
	id := int32(fid)
	cs := int32(1) // cos
	if useSin {
		cs = 2 // sin
	}
	res, _, _, _, err = qpck.Awoe(id, f, a, b, ω, cs, 0, 0, 0, 0, nil, nil, nil, nil, nil, nil, 0, nil)
	return
}

// QuadExpIx approximates the integral of f(x) ⋅ exp(i⋅m⋅x) with i = √-1
//
//   INPUT:
//     a      -- lower limit of integration
//     b      -- upper limit of integration
//     m      -- coefficient of x
//     fid    -- index of goroutine (to avoid race problems)
//     f      -- function defining the integrand
//
//   OUTPUT:        b                           b                           b
//           res = ∫  f(x) ⋅ exp(i⋅m⋅x) dx   = ∫  f(x) ⋅ cos(m⋅x) dx + i ⋅ ∫  f(x) ⋅ sin(m⋅x) dx
//                 a                           a                           a
//
func QuadExpIx(a, b, m float64, fid int, f func(x float64) float64) (res complex128, err error) {

	// allocate workspace
	limit := 50
	alist := make([]float64, limit)
	blist := make([]float64, limit)
	rlist := make([]float64, limit)
	elist := make([]float64, limit)
	iord := make([]int32, limit)
	nnlog := make([]int32, limit)

	// set flags
	id := int32(fid)
	var icall int32 = 1  // do not reuse moments
	var maxp1 int32 = 50 // upper bound on the number of Chebyshev moments
	var momcom int32     // 0 => do compute moments

	// allocate Chebyshev moments array
	chebmo := make([]float64, 25*maxp1)

	// perform integration of cos term
	var integr int32 = 1 // w(x) = cos(m*x)
	Icos, _, _, _, err := qpck.Awoe(id, f, a, b, m, integr, 0, 0, icall, maxp1, alist, blist, rlist, elist, iord, nnlog, momcom, chebmo)
	if err != nil {
		return
	}

	// set flags
	icall = 2  // do reuse moments
	momcom = 1 // do not compute moments

	// perform integration of sin term
	integr = 2 // w(x) = sin(m*x)
	Isin, _, _, _, err := qpck.Awoe(id, f, a, b, m, integr, 0, 0, icall, maxp1, alist, blist, rlist, elist, iord, nnlog, momcom, chebmo)
	if err != nil {
		return
	}

	// results
	res = complex(Icos, Isin)
	return
}
