// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Elliptic1 computes Legendre elliptic integral of the first kind F(φ,k),
// evaluated using Carlson’s function Rf [1].
// The argument ranges are  0 ≤ φ ≤ π/2  and  0 ≤ k·sin(φ) ≤ 1
//
//   Computes:
//                      φ
//                     ⌠          dt
//         F(φ, k)  =  │  ___________________
//                     │     _______________
//                     ⌡   \╱ 1 - k² sin²(t)
//                    0
//   where:
//            0 ≤ φ ≤ π/2
//            0 ≤ k·sin(φ) ≤ 1
//
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
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

// Elliptic2 computes Legendre elliptic integral of the second kind E(φ,k),
// evaluated using Carlson's functions Rf and Rd [1].
// The argument ranges are  0 ≤ φ ≤ π/2  and  0 ≤ k⋅sin(φ) ≤ 1
//
//   Computes:
//                      φ
//                     ⌠     _______________
//         E(φ, k)  =  │   \╱ 1 - k² sin²(t)  dt
//                     ⌡
//                    0
//   where:
//            0 ≤ φ ≤ π/2
//            0 ≤ k·sin(φ) ≤ 1
//
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func Elliptic2(φ, k float64) float64 {
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
	cc := math.Pow(math.Cos(φ), 2.0)
	q := (1.0 - s*k) * (1.0 + s*k)
	return s * (CarlsonRf(cc, q, 1.0) - (math.Pow(s*k, 2.0))*CarlsonRd(cc, q, 1.0)/3.0)
}

// Elliptic3 computes Legendre elliptic integral of the third kind Π(φ,n,k),
// evaluated using Carlson's functions Rf and Rj.
// NOTE that the sign convention on n is opposite to that of Abramowitz and Stegun [2]
// The argument ranges are  0 ≤ φ ≤ π/2  and  0 ≤ k⋅sin(φ) ≤ 1
//
//   Computes:
//                         φ
//                        ⌠                  dt
//         Π(φ, n, k)  =  │  ___________________________________
//                        │                     _______________
//                        ⌡   (1 + n sin²(t)) \╱ 1 - k² sin²(t)
//                       0
//   where:
//            0 ≤ φ ≤ π/2
//            0 ≤ k·sin(φ) ≤ 1
//
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//   [2] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
//       and Mathematical Tables. U.S. Department of Commerce, NIST
func Elliptic3(φ, n, k float64) float64 {
	s := math.Sin(φ)
	t := n * s * s
	cc := math.Pow(math.Cos(φ), 2.0)
	q := (1.0 - s*k) * (1.0 + s*k)
	return s * (CarlsonRf(cc, q, 1.0) - t*CarlsonRj(cc, q, 1.0, 1.0+t)/3.0)
}

// CarlsonRf computes Carlson's elliptic integral of the first kind according to [1]. See also [2]
// Computes Rf(x,y,z) where x,y,z must be non-negative and at most one can be zero.
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
//   [2] Carlson BC (1977) Elliptic Integrals of the First Kind, SIAM Journal on Mathematical
//       Analysis, vol. 8, pp. 231-242.
func CarlsonRf(x, y, z float64) float64 {
	ERRTOL := 0.0025 // a value of 0.0025 for the error tolerance parameter gives full double precision (16 sig ant digits) [1]
	THIRD := 1.0 / 3.0
	C1 := 1.0 / 24.0
	C2 := 0.1
	C3 := 3.0 / 44.0
	C4 := 1.0 / 14.0
	TINY := 5.0 * math.SmallestNonzeroFloat64
	BIG := 0.2 * math.MaxFloat64
	if min(min(x, y), z) < 0.0 || min(min(x+y, x+z), y+z) < TINY || max(max(x, y), z) > BIG {
		chk.Panic("cannot compute Carlson's Rf function with x=%g, y=%g, z=%g. All values must be non-negative and at most one can be zero", x, y, z)
	}
	xt := x
	yt := y
	zt := z
	var sqrtx, sqrty, sqrtz, alamb, ave, delx, dely, delz float64
	it, MAXIT := 0, 11
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

// CarlsonRd computes Carlson’s elliptic integral of the second kind according to [1]
// Computes Rf(x,y,z) where x,y must be non-negative and at most one can be zero.
// z must be positive.
//   References:
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//       Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func CarlsonRd(x, y, z float64) float64 {
	ERRTOL := 0.0015
	C1 := 3.0 / 14.0
	C2 := 1.0 / 6.0
	C3 := 9.0 / 22.0
	C4 := 3.0 / 26.0
	C5 := 0.25 * C3
	C6 := 1.5 * C4
	TINY := 2.0 * math.Pow(math.MaxFloat64, -2.0/3.0)
	BIG := 0.1 * ERRTOL * math.Pow(math.SmallestNonzeroFloat64, -2.0/3.0)
	if min(x, y) < 0.0 || min(x+y, z) < TINY || max(max(x, y), z) > BIG {
		chk.Panic("cannot compute Carlson's Rd function with x=%g, y=%g, z=%g. x,y must be non-negative and at most one can be zero. z must be positive", x, y, z)
	}
	xt := x
	yt := y
	zt := z
	sum := 0.0
	fac := 1.0
	var alamb, ave, delx, dely, delz, ea, eb, ec, ed, ee, sqrtx, sqrty, sqrtz float64
	it, MAXIT := 0, 11
	for it = 0; it < MAXIT; it++ {
		sqrtx = math.Sqrt(xt)
		sqrty = math.Sqrt(yt)
		sqrtz = math.Sqrt(zt)
		alamb = sqrtx*(sqrty+sqrtz) + sqrty*sqrtz
		sum += fac / (sqrtz * (zt + alamb))
		fac = 0.25 * fac
		xt = 0.25 * (xt + alamb)
		yt = 0.25 * (yt + alamb)
		zt = 0.25 * (zt + alamb)
		ave = 0.2 * (xt + yt + 3.0*zt)
		delx = (ave - xt) / ave
		dely = (ave - yt) / ave
		delz = (ave - zt) / ave
		if max(max(math.Abs(delx), math.Abs(dely)), math.Abs(delz)) < ERRTOL {
			break
		}
	}
	if it == MAXIT {
		chk.Panic("CarlsonRd failed to converge after %d iterations", it)
	}
	ea = delx * dely
	eb = delz * delz
	ec = ea - eb
	ed = ea - 6.0*eb
	ee = ed + ec + ec
	return 3.0*sum + fac*(1.0+ed*(-C1+C5*ed-C6*delz*ee)+delz*(C2*ee+delz*(-C3*ec+delz*C4*ea)))/(ave*math.Sqrt(ave))
}

// CarlsonRj computes Carlson’s elliptic integral of the third kind according to [1]
// Computes Rj(x,y,z,p) where x,y,z must be nonnegative, and at most one can be zero.
// p must be nonzero. If p < 0, the Cauchy principal value is returned.
func CarlsonRj(x, y, z, p float64) float64 {
	ERRTOL := 0.0015
	C1 := 3.0 / 14.0
	C2 := 1.0 / 3.0
	C3 := 3.0 / 22.0
	C4 := 3.0 / 26.0
	C5 := 0.75 * C3
	C6 := 1.5 * C4
	C7 := 0.5 * C2
	C8 := C3 + C3
	TINY := math.Pow(5.0*math.SmallestNonzeroFloat64, 1.0/3.0)
	BIG := 0.3 * math.Pow(0.2*math.MaxFloat64, 1.0/3.0)
	var a, alamb, alpha, ans, ave, b, beta, delp, delx, dely, delz, ea, eb, ec, ed, ee, pt, rcx, rho, sqrtx, sqrty, sqrtz, tau, xt, yt, zt float64
	if min(min(x, y), z) < 0.0 || min(min(x+y, x+z), min(y+z, math.Abs(p))) < TINY || max(max(x, y), max(z, math.Abs(p))) > BIG {
		chk.Panic("cannot compute Carlson's Rj function with x=%g, y=%g, z=%g, p=%g. x,y,z must be non-negative and at most one can be zero. p must be nonzero", x, y, z, p)
	}
	sum := 0.0
	fac := 1.0
	if p > 0.0 {
		xt = x
		yt = y
		zt = z
		pt = p
	} else {
		xt = min(min(x, y), z)
		zt = max(max(x, y), z)
		yt = x + y + z - xt - zt
		a = 1.0 / (yt - p)
		b = a * (zt - yt) * (yt - xt)
		pt = yt + b
		rho = xt * zt / yt
		tau = p * pt / yt
		rcx = CarlsonRc(rho, tau)
	}
	it, MAXIT := 0, 11
	for it = 0; it < MAXIT; it++ {
		sqrtx = math.Sqrt(xt)
		sqrty = math.Sqrt(yt)
		sqrtz = math.Sqrt(zt)
		alamb = sqrtx*(sqrty+sqrtz) + sqrty*sqrtz
		alpha = math.Pow(pt*(sqrtx+sqrty+sqrtz)+sqrtx*sqrty*sqrtz, 2.0)
		beta = pt * math.Pow(pt+alamb, 2.0)
		sum += fac * CarlsonRc(alpha, beta)
		fac = 0.25 * fac
		xt = 0.25 * (xt + alamb)
		yt = 0.25 * (yt + alamb)

		zt = 0.25 * (zt + alamb)
		pt = 0.25 * (pt + alamb)
		ave = 0.2 * (xt + yt + zt + pt + pt)
		delx = (ave - xt) / ave
		dely = (ave - yt) / ave
		delz = (ave - zt) / ave
		delp = (ave - pt) / ave
		if max(max(math.Abs(delx), math.Abs(dely)), max(math.Abs(delz), math.Abs(delp))) < ERRTOL {
			break
		}
	}
	if it == MAXIT {
		chk.Panic("CarlsonRj failed to converge after %d iterations", it)
	}
	ea = delx*(dely+delz) + dely*delz
	eb = delx * dely * delz
	ec = delp * delp
	ed = ea - 3.0*ec
	ee = eb + 2.0*delp*(ea-ec)
	ans = 3.0*sum + fac*(1.0+ed*(-C1+C5*ed-C6*ee)+eb*(C7+delp*(-C8+delp*C4))+delp*ea*(C2-delp*C3)-C2*delp*ec)/(ave*math.Sqrt(ave))
	if p <= 0.0 {
		ans = a * (b*ans + 3.0*(rcx-CarlsonRf(xt, yt, zt)))
	}
	return ans
}

// CarlsonRc computes Carlson’s degenerate elliptic integral according to [1]
// Computes Rc(x,y) where x must be nonnegative and y must be nonzero.
// If y < 0, the Cauchy principal value is returned.
func CarlsonRc(x, y float64) float64 {
	ERRTOL := 0.0012
	THIRD := 1.0 / 3.0
	C1 := 0.3
	C2 := 1.0 / 7.0
	C3 := 0.375
	C4 := 9.0 / 22.0
	TINY := 5.0 * math.SmallestNonzeroFloat64
	BIG := 0.2 * math.MaxFloat64
	COMP1 := 2.236 / math.Sqrt(TINY)
	COMP2 := math.Pow(TINY*BIG, 2.0) / 25.0
	var alamb, ave, s, w, xt, yt float64
	if x < 0.0 || y == 0.0 || (x+math.Abs(y)) < TINY || (x+math.Abs(y)) > BIG || (y < -COMP1 && x > 0.0 && x < COMP2) {
		chk.Panic("cannot compute Carlson's Rc function with x=%g, y=%g. x,y must be non-negative", x, y)
	}
	if y > 0.0 {
		xt = x
		yt = y
		w = 1.0
	} else {
		xt = x - y
		yt = -y
		w = math.Sqrt(x) / math.Sqrt(xt)
	}
	it, MAXIT := 0, 11
	for it = 0; it < MAXIT; it++ {
		alamb = 2.0*math.Sqrt(xt)*math.Sqrt(yt) + yt
		xt = 0.25 * (xt + alamb)
		yt = 0.25 * (yt + alamb)
		ave = THIRD * (xt + yt + yt)
		s = (yt - ave) / ave
		if math.Abs(s) < ERRTOL {
			break
		}
	}
	return w * (1.0 + s*s*(C1+s*(C2+s*(C3+s*C4)))) / math.Sqrt(ave)
}
