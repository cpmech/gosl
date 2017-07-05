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
