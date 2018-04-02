// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// InterpQuad computes a quadratic polynomial to perform interpolation either using 3 points
// or 2 points and a known derivative
type InterpQuad struct {
	A, B, C float64 // coefficients of polynomial
	TolDen  float64 // tolerance to avoid zero denominator
}

// NewInterpQuad returns a new object
func NewInterpQuad() (o *InterpQuad) {
	o = new(InterpQuad)
	o.TolDen = 1e-15
	return
}

// F computes y = f(x) curve
func (o *InterpQuad) F(x float64) float64 {
	return o.A*x*x + o.B*x + o.C
}

// G computes y' = df/x|(x) curve
func (o *InterpQuad) G(x float64) float64 {
	return 2.0*o.A*x + o.B
}

// Optimum returns the minimum or maximum point; i.e. the point with zero derivative
//   xopt -- x @ optimum
//   fopt -- f(xopt) = y @ optimum
func (o *InterpQuad) Optimum() (xopt, fopt float64) {
	if math.Abs(o.A) < o.TolDen {
		chk.Panic("cannot compute optimum because zero A=%g\n", o.A)
	}
	xopt = -0.5 * o.B / o.A
	fopt = o.F(xopt)
	return
}

// Fit3points fits polynomial to 3 points
//   (x0, y0) -- first point
//   (x1, y1) -- second point
//   (x2, y2) -- third point
func (o *InterpQuad) Fit3points(x0, y0, x1, y1, x2, y2 float64) (err error) {
	z0, z1, z2 := x0*x0, x1*x1, x2*x2
	den := x0*(z2-z1) - x1*z2 + x2*z1 + (x1-x2)*z0
	if math.Abs(den) < o.TolDen {
		return chk.Err("Cannot fit 3 points because denominator=%g is near zero.\n\t(x0,y0)=(%g,%g)\t(x1,y1)=(%g,%g)\t(x2,y2)=(%g,%g)\n", x0, y0, x1, y1, x2, y2)
	}
	o.A = ((x1-x2)*y0 + x2*y1 - x1*y2 + x0*(y2-y1)) / den
	o.B = ((y1-y2)*z0 + y2*z1 - y1*z2 + y0*(z2-z1)) / den
	o.C = -((x2*y1-x1*y2)*z0 + y0*(x1*z2-x2*z1) + x0*(y2*z1-y1*z2)) / den
	return
}

// Fit2pointsD fits polynomial to 2 points and known derivative
//   (x0, y0) -- first point
//   (x1, y1) -- second point
//   (x2, d2) -- derivative @ x2
func (o *InterpQuad) Fit2pointsD(x0, y0, x1, y1, x2, d2 float64) (err error) {
	z0, z1 := x0*x0, x1*x1
	den := -z1 + z0 + 2*x1*x2 - 2*x0*x2
	if math.Abs(den) < o.TolDen {
		return chk.Err("Cannot fit 2 points + deriv because denominator=%g is near zero.\n\t(x0,y0)=(%g,%g)\t(x1,y1)=(%g,%g)\t(x2,d2)=(%g,%g)\n", x0, y0, x1, y1, x2, d2)
	}
	o.A = (-d2*x0 + d2*x1 + y0 - y1) / den
	o.B = (-2*x2*y0 + 2*x2*y1 + d2*z0 - d2*z1) / den
	o.C = ((y1-d2*x1)*z0 + y0*(2*x1*x2-z1) + x0*(d2*z1-2*x2*y1)) / den
	return
}
