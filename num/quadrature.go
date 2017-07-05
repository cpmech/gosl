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
