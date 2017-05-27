// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// CarlsonRf computes Carlson's elliptic integral of the first kind according to [1]. See also [2]
// Computes Rf(x,y,z) where x,y,z must be non-negative and at most one can be zero.
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//   [2] Carlson BC (1977) Elliptic Integrals of the First Kind, SIAM Journal on Mathematical
//       Analysis, vol. 8, pp. 231–242.
func CarlsonRf(x, y, z float64) float64 {
	ERRTOL := 0.0025
	THIRD := 1.0 / 3.0
	C1 := 1.0 / 24.0
	C2 := 0.1
	C3 := 3.0 / 44.0
	C4 := 1.0 / 14.0
	TINY := 5.0 * math.SmallestNonzeroFloat64
	BIG := 0.2 * math.MaxFloat64
	if min(min(x, y), z) < 0.0 || min(min(x+y, x+z), y+z) < TINY || max(max(x, y), z) > BIG {
		chk.Panic("cannot compute elliptic integral of first kind in Rf with x=%g, y=%g, z=%g. All values must be non-negative and at most one can be zero", x, y, z)
	}
	xt := x
	yt := y
	zt := z
	var sqrtx, sqrty, sqrtz, alamb, ave, delx, dely, delz float64
	it, MAXIT := 0, 100
	for it = 0; it < MAXIT; it++ {
		sqrtx = math.Sqrt(xt)
		sqrty = math.Sqrt(yt)
		sqrtz = math.Sqrt(zt)
		alamb = sqrtx*(sqrty+sqrtz) + sqrty*sqrtz
		xt = 0.25 * (xt + alamb)
		yt = 0.25 * (yt + alamb)
		zt = 0.25 * (zt + alamb)
		ave = THIRD * (xt + yt + zt)
		delx = (ave - xt) / ave
		dely = (ave - yt) / ave
		delz = (ave - zt) / ave
		if max(max(math.Abs(delx), math.Abs(dely)), math.Abs(delz)) < ERRTOL {
			break
		}
	}
	if it == MAXIT {
		chk.Panic("CarlsonRf failed to converge after %d iterations", it)
	}
	e2 := delx*dely - delz*delz
	e3 := delx * dely * delz
	return (1.0 + (C1*e2-C2-C3*e3)*e2 + C4*e3) / math.Sqrt(ave)
}

// Elliptic1 computes Legendre elliptic integral of the first kind F(φ,k),
// evaluated using Carlson’s function Rf.
// The argument ranges are 0 ≤ φ ≤ π/2  and  0 ≤ k·sin(φ) ≤ 1
//
//   Computes:
//                      φ
//                     ⌠         dt
//         F(φ, k)  =  │  _________________
//                     │    _______________
//                     ⌡  \╱ 1 - k² sin²(t)
//                    0
//   where:
//            0 ≤ φ ≤ π/2
//            0 ≤ k·sin(φ) ≤ 1
//
func Elliptic1(φ, k float64) float64 {
	if φ < 0 || k < 0 {
		chk.Panic("φ and k must be non-negative. φ=%g, k=%g is invalid", φ, k)
	}
	if φ > math.Pi/2.0+1e-15 {
		chk.Panic("φ must be in 0 ≤ φ ≤ π/2. φ=%g is invalid", φ)
	}
	if φ < math.SmallestNonzeroFloat64 {
		return 0
	}
	if k < math.SmallestNonzeroFloat64 {
		return φ
	}
	s := math.Sin(φ)
	if math.Abs(k*s-1.0) < 1e-15 {
		return math.Inf(1)
	}
	return s * CarlsonRf(math.Pow(math.Cos(φ), 2.0), (1.0-s*k)*(1.0+s*k), 1.0)
}
