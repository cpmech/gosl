// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import "math"

// EqCubicSolveReal solves a cubic equation, ignoring the complex answers.
//  The equation is specified by:
//   x³ + a x² + b x + c = 0
//  Notes:
//   1) every cubic equation with real coefficients has at least one solution
//      x among the real numbers
//   2) from Numerical Recipes 2007, page 228
//  Output:
//   x[i] -- roots
//   nx   -- number of real roots: 1, 2 or 3
func EqCubicSolveReal(a, b, c float64) (x1, x2, x3 float64, nx int) { //, err error) {

	// tolerance
	ϵ := 1e-14

	// auxiliary
	aa := a * a
	Q := (aa - 3.0*b) / 9.0
	R := (2.0*a*aa - 9.0*a*b + 27.0*c) / 54.0
	QQQ := Q * Q * Q

	// three real roots
	if R*R < QQQ {
		//if QQQ < ϵ {
		//return chk.Err("cannot compute roots of cubic equation because Q³ = %g < %g", QQQ, ϵ)
		//}
		θ := math.Acos(R / math.Sqrt(QQQ))
		c := -2.0 * math.Sqrt(Q)
		x1 = c*math.Cos(θ/3.0) - a/3.0
		x2 = c*math.Cos((θ+2.0*math.Pi)/3.0) - a/3.0
		x3 = c*math.Cos((θ-2.0*math.Pi)/3.0) - a/3.0
		nx = 3
		return
	}

	// auxiliary
	A := -sign(R) * math.Pow(math.Abs(R)+math.Sqrt(R*R-QQQ), 1.0/3.0)
	B := 0.0
	if math.Abs(A) > ϵ {
		B = Q / A
	}

	// one root
	if math.Abs(A) < ϵ && math.Abs(B) < ϵ {
		x1 = -a / 3.0
		nx = 1
		return
	}

	// two roots
	if math.Abs(A-B) < ϵ {
		x1 = (A + B) - a/3.0
		x2 = -(A+B)/2.0 - a/3.0
		nx = 2
		return
	}

	// one real root
	x1 = (A + B) - a/3.0
	nx = 1
	return
}
