// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// InterpCubic computes a cubic polynomial to perform interpolation either using 4 points
// or 3 points and a known derivative
type InterpCubic struct {
	A, B, C, D float64 // coefficients of polynomial
	TolDen     float64 // tolerance to avoid zero denominator
}

// NewInterpCubic returns a new object
func NewInterpCubic() (o *InterpCubic) {
	o = new(InterpCubic)
	o.TolDen = 1e-15
	return
}

// F computes y = f(x) curve
func (o *InterpCubic) F(x float64) float64 {
	return o.A*x*x*x + o.B*x*x + o.C*x + o.D
}

// G computes y' = df/x|(x) curve
func (o *InterpCubic) G(x float64) float64 {
	return 3.0*o.A*x*x + 2.0*o.B*x + o.C
}

// Critical returns the critical points
//
//	xmin -- x @ min and y(xmin)
//	xmax -- x @ max and y(xmax)
//	xifl -- x @ inflection point and y(ifl)
//	hasMin, hasMax, hasIfl -- flags telling what is available
func (o *InterpCubic) Critical() (xmin, xmax, xifl float64, hasMin, hasMax, hasIfl bool) {
	delBy4 := o.B*o.B - 3.0*o.A*o.C
	if delBy4 < 0 {
		return // cubic function is strictly monotonic
	}
	den := 3.0 * o.A
	xifl = -o.B / den
	hasIfl = true
	if delBy4 != 0.0 {
		xmin = (-o.B + math.Sqrt(delBy4)) / den
		xmax = (-o.B - math.Sqrt(delBy4)) / den
		if o.F(xmin) > o.F(xmax) {
			xmin, xmax = xmax, xmin
		}
		hasMin = true
		hasMax = true
	}
	return
}

// Fit4points fits polynomial to 4 points
//
//	(x0, y0) -- first point
//	(x1, y1) -- second point
//	(x2, y2) -- third point
//	(x3, y3) -- fourth point
func (o *InterpCubic) Fit4points(x0, y0, x1, y1, x2, y2, x3, y3 float64) (err error) {
	z0, z1, z2, z3 := x0*x0, x1*x1, x2*x2, x3*x3
	w0, w1, w2, w3 := z0*x0, z1*x1, z2*x2, z3*x3
	den := w0*((x2-x3)*z1+x3*z2-x2*z3+x1*(z3-z2)) + w1*(x2*z3-x3*z2) + x0*((w3-w2)*z1-w3*z2+w1*(z2-z3)+w2*z3) + x1*(w3*z2-w2*z3) + (w2*x3-w3*x2)*z1 + ((w2-w3)*x1+w3*x2-w2*x3+w1*(x3-x2))*z0
	if math.Abs(den) < o.TolDen {
		return chk.Err("Cannot fit 4 points because denominator=%g is near zero.\n\t(x0,y0)=(%g,%g)\n\t(x1,y1)=(%g,%g)\n\t(x2,y2)=(%g,%g)\n\t(x3,y3)=(%g,%g)\n", den, x0, y0, x1, y1, x2, y2, x3, y3)
	}
	o.A = -((x1*(y3-y2)-x2*y3+x3*y2+(x2-x3)*y1)*z0 + (x2*y3-x3*y2)*z1 + y1*(x3*z2-x2*z3) + y0*(x2*z3+x1*(z2-z3)-x3*z2+(x3-x2)*z1) + x1*(y2*z3-y3*z2) + x0*(y1*(z3-z2)-y2*z3+y3*z2+(y2-y3)*z1)) / den
	o.B = ((w1*(x3-x2)-w2*x3+w3*x2+(w2-w3)*x1)*y0 + (w2*x3-w3*x2)*y1 + x1*(w3*y2-w2*y3) + x0*(w2*y3+w1*(y2-y3)-w3*y2+(w3-w2)*y1) + w1*(x2*y3-x3*y2) + w0*(x1*(y3-y2)-x2*y3+x3*y2+(x2-x3)*y1)) / den
	o.C = ((w1*(y3-y2)-w2*y3+w3*y2+(w2-w3)*y1)*z0 + (w2*y3-w3*y2)*z1 + y1*(w3*z2-w2*z3) + y0*(w2*z3+w1*(z2-z3)-w3*z2+(w3-w2)*z1) + w1*(y2*z3-y3*z2) + w0*(y1*(z3-z2)-y2*z3+y3*z2+(y2-y3)*z1)) / den
	o.D = ((w1*(x3*y2-x2*y3)+x1*(w2*y3-w3*y2)+(w3*x2-w2*x3)*y1)*z0 + y0*(w1*(x2*z3-x3*z2)+x1*(w3*z2-w2*z3)+(w2*x3-w3*x2)*z1) + x0*(w1*(y3*z2-y2*z3)+y1*(w2*z3-w3*z2)+(w3*y2-w2*y3)*z1) + w0*(x1*(y2*z3-y3*z2)+y1*(x3*z2-x2*z3)+(x2*y3-x3*y2)*z1)) / den
	return
}

// Fit3pointsD fits polynomial to 3 points and known derivative
//
//	(x0, y0) -- first point
//	(x1, y1) -- second point
//	(x2, y2) -- third point
//	(x3, d3) -- derivative @ x3
func (o *InterpCubic) Fit3pointsD(x0, y0, x1, y1, x2, y2, x3, d3 float64) (err error) {
	z0, z1, z2, z3 := x0*x0, x1*x1, x2*x2, x3*x3
	w0, w1, w2 := z0*x0, z1*x1, z2*x2
	den := x0*(2*w1*x3-2*w2*x3-3*z1*z3+3*z2*z3) + x1*(2*w2*x3-3*z2*z3) + z1*(3*x2*z3-w2) + z0*(-w1+w2+3*x1*z3-3*x2*z3) + w1*(z2-2*x2*x3) + w0*(-2*x1*x3+2*x2*x3+z1-z2)
	if math.Abs(den) < o.TolDen {
		return chk.Err("Cannot fit 3 points and known derivative because denominator=%g is near zero.\n\t(x0,y0)=(%g,%g)\n\t(x1,y1)=(%g,%g)\n\t(x2,y2)=(%g,%g)\n\t(x3,d3)=(%g,%g)\n", den, x0, y0, x1, y1, x2, y2, x3, d3)
	}
	o.A = -(-2*x1*x3*y2 + x0*(2*x3*y2-2*x3*y1) + (y1-y2)*z0 + y2*z1 + y1*(2*x2*x3-z2) + y0*(z2-z1-2*x2*x3+2*x1*x3)) / den
	o.B = (w0*(y1-y2) + w1*y2 - 3*x1*y2*z3 + y0*(-3*x2*z3+3*x1*z3+w2-w1) + y1*(3*x2*z3-w2) + x0*(3*y2*z3-3*y1*z3)) / den
	o.C = (-2*w1*x3*y2 + w0*(2*x3*y2-2*x3*y1) + 3*y2*z1*z3 + z0*(3*y1*z3-3*y2*z3) + y1*(2*w2*x3-3*z2*z3) + y0*(3*z2*z3-3*z1*z3-2*w2*x3+2*w1*x3)) / den
	o.D = -(w0*(y1*(z2-2*x2*x3)-y2*z1+2*x1*x3*y2) + z0*(y1*(3*x2*z3-w2)-3*x1*y2*z3+w1*y2) + x0*(y1*(2*w2*x3-3*z2*z3)+3*y2*z1*z3-2*w1*x3*y2) + y0*(x1*(3*z2*z3-2*w2*x3)+z1*(w2-3*x2*z3)+w1*(2*x2*x3-z2))) / den
	return
}
