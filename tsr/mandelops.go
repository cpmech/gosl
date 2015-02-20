// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// M_Norm calculates the norm of a 2nd order tensor represented in Mandel's basis
func M_Norm(a []float64) float64 {
	if len(a) == 4 {
		return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2] + a[3]*a[3])
	}
	return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2] + a[3]*a[3] + a[4]*a[4] + a[5]*a[5])
}

// M_Tr calculates the trace a 2nd order tensor represented in Mandel's basis
func M_Tr(a []float64) float64 {
	return a[0] + a[1] + a[2]
}

// M_Dev calculates the deviator a 2nd order tensor represented in Mandel's basis
func M_Dev(a []float64) (s []float64) {
	s = make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		s[i] = a[i] - (a[0]+a[1]+a[2])*Im[i]/3.0
	}
	return
}

// M_Det calculates the determinant a 2nd order tensor represented in Mandel's basis
func M_Det(a []float64) float64 {
	if len(a) == 4 {
		return a[0]*a[1]*a[2] - a[2]*a[3]*a[3]/2.0
	}
	return a[0]*a[1]*a[2] + a[3]*a[4]*a[5]/SQ2 - a[0]*a[4]*a[4]/2.0 - a[1]*a[5]*a[5]/2.0 - a[2]*a[3]*a[3]/2.0
}

// M_DetDeriv computes the derivative of the determinant of a w.r.t a
//  d := dDet(a)/da == dI3(a)/da
func M_DetDeriv(d, a []float64) {
	aa := make([]float64, len(a)) // aa := a single-dot a
	M_Sq(aa, a)
	I1 := M_Tr(a)
	I2 := (I1*I1 - M_Tr(aa)) / 2.0
	for i := 0; i < len(a); i++ {
		d[i] = aa[i] - I1*a[i] + I2*Im[i]
	}
}

// M_Sq returns the square of a tensor in Mandel's representation
//  b = a² = a single-dot a
func M_Sq(b, a []float64) {
	if len(a) == 4 {
		b[0] = (a[3]*a[3] + 2.0*a[0]*a[0]) / 2.0
		b[1] = (a[3]*a[3] + 2.0*a[1]*a[1]) / 2.0
		b[2] = a[2] * a[2]
		b[3] = (a[1] + a[0]) * a[3]
		return
	}
	b[0] = (a[5]*a[5] + a[3]*a[3] + 2.0*a[0]*a[0]) / 2.0
	b[1] = (a[4]*a[4] + a[3]*a[3] + 2.0*a[1]*a[1]) / 2.0
	b[2] = (a[5]*a[5] + a[4]*a[4] + 2.0*a[2]*a[2]) / 2.0
	b[3] = (a[4]*a[5] + (SQ2*a[1]+SQ2*a[0])*a[3]) / SQ2
	b[4] = (a[3]*a[5] + (SQ2*a[2]+SQ2*a[1])*a[4]) / SQ2
	b[5] = (a[3]*a[4] + (SQ2*a[2]+SQ2*a[0])*a[5]) / SQ2
}

// M_SqDeriv (Mandel) derivative of square of a tensor
//  d = derivative of a² w.r.t a
func M_SqDeriv(d [][]float64, a []float64) {
	if len(a) == 4 {
		d[0][0] = 2 * a[0]
		d[0][1] = 0
		d[0][2] = 0
		d[0][3] = a[3]
		d[1][0] = 0
		d[1][1] = 2 * a[1]
		d[1][2] = 0
		d[1][3] = a[3]
		d[2][0] = 0
		d[2][1] = 0
		d[2][2] = 2 * a[2]
		d[2][3] = 0
		d[3][0] = a[3]
		d[3][1] = a[3]
		d[3][2] = 0
		d[3][3] = a[1] + a[0]
		return
	}
	d[0][0] = 2 * a[0]
	d[0][1] = 0
	d[0][2] = 0
	d[0][3] = a[3]
	d[0][4] = 0
	d[0][5] = a[5]
	d[1][0] = 0
	d[1][1] = 2 * a[1]
	d[1][2] = 0
	d[1][3] = a[3]
	d[1][4] = a[4]
	d[1][5] = 0
	d[2][0] = 0
	d[2][1] = 0
	d[2][2] = 2 * a[2]
	d[2][3] = 0
	d[2][4] = a[4]
	d[2][5] = a[5]
	d[3][0] = a[3]
	d[3][1] = a[3]
	d[3][2] = 0
	d[3][3] = a[1] + a[0]
	d[3][4] = a[5] / SQ2
	d[3][5] = a[4] / SQ2
	d[4][0] = 0
	d[4][1] = a[4]
	d[4][2] = a[4]
	d[4][3] = a[5] / SQ2
	d[4][4] = a[2] + a[1]
	d[4][5] = a[3] / SQ2
	d[5][0] = a[5]
	d[5][1] = 0
	d[5][2] = a[5]
	d[5][3] = a[4] / SQ2
	d[5][4] = a[3] / SQ2
	d[5][5] = a[2] + a[0]
}

// M_Dy returns the dyadic product between a and b
//  c := a dy b
func M_Dy(a, b []float64) (c [][]float64) {
	if len(a) == 4 {
		return [][]float64{
			{a[0] * b[0], a[0] * b[1], a[0] * b[2], a[0] * b[3]},
			{a[1] * b[0], a[1] * b[1], a[1] * b[2], a[1] * b[3]},
			{a[2] * b[0], a[2] * b[1], a[2] * b[2], a[2] * b[3]},
			{a[3] * b[0], a[3] * b[1], a[3] * b[2], a[3] * b[3]},
		}
	}
	return [][]float64{
		{a[0] * b[0], a[0] * b[1], a[0] * b[2], a[0] * b[3], a[0] * b[4], a[0] * b[5]},
		{a[1] * b[0], a[1] * b[1], a[1] * b[2], a[1] * b[3], a[1] * b[4], a[1] * b[5]},
		{a[2] * b[0], a[2] * b[1], a[2] * b[2], a[2] * b[3], a[2] * b[4], a[2] * b[5]},
		{a[3] * b[0], a[3] * b[1], a[3] * b[2], a[3] * b[3], a[3] * b[4], a[3] * b[5]},
		{a[4] * b[0], a[4] * b[1], a[4] * b[2], a[4] * b[3], a[4] * b[4], a[4] * b[5]},
		{a[5] * b[0], a[5] * b[1], a[5] * b[2], a[5] * b[3], a[5] * b[4], a[5] * b[5]},
	}
}

// M_DyAdd adds the dyadic product between a and b scaled by s
//  c += s * a dy b
func M_DyAdd(c [][]float64, s float64, a, b []float64) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			c[i][j] += s * a[i] * b[j]
		}
	}
}

// M_Dot multiplies two second order symmetric tensors (the result may be non-symmetric)
// An error is returned in case the result is non-symmetric
//  c = a dot b  =>  cij = aik * bkj
func M_Dot(c []float64, a, b []float64, nonsymTol float64) (err error) {
	if len(a) == 4 {
		c[0] = (a[3]*b[3] + 2.0*a[0]*b[0]) / 2.0
		c[1] = (a[3]*b[3] + 2.0*a[1]*b[1]) / 2.0
		c[2] = a[2] * b[2]
		c[3] = ((a[1]+a[0])*b[3] + a[3]*b[1] + a[3]*b[0]) / 2.0
		c6 := ((a[1]-a[0])*b[3] - a[3]*b[1] + a[3]*b[0]) / 2.0
		if math.Abs(c6) > nonsymTol {
			err = chk.Err(_mandelops_err2, "2D", c6)
			return
		}
		return
	}
	c[0] = (a[5]*b[5] + a[3]*b[3] + 2*a[0]*b[0]) / 2.0
	c[1] = (a[4]*b[4] + a[3]*b[3] + 2*a[1]*b[1]) / 2.0
	c[2] = (a[5]*b[5] + a[4]*b[4] + 2*a[2]*b[2]) / 2.0
	c[3] = (a[4]*b[5] + a[5]*b[4] + (SQ2*a[1]+SQ2*a[0])*b[3] + SQ2*a[3]*b[1] + SQ2*a[3]*b[0]) / TWOSQ2
	c[4] = (a[3]*b[5] + (SQ2*a[2]+SQ2*a[1])*b[4] + a[5]*b[3] + SQ2*a[4]*b[2] + SQ2*a[4]*b[1]) / TWOSQ2
	c[5] = ((SQ2*a[2]+SQ2*a[0])*b[5] + a[3]*b[4] + a[4]*b[3] + SQ2*a[5]*b[2] + SQ2*a[5]*b[0]) / TWOSQ2
	c6 := (a[4]*b[5] - a[5]*b[4] + (SQ2*a[1]-SQ2*a[0])*b[3] - SQ2*a[3]*b[1] + SQ2*a[3]*b[0]) / TWOSQ2
	c7 := (a[5]*b[3] - a[3]*b[5] - (SQ2*a[1]-SQ2*a[2])*b[4] - SQ2*a[4]*b[2] + SQ2*a[4]*b[1]) / TWOSQ2
	c8 := ((SQ2*a[2]-SQ2*a[0])*b[5] - a[3]*b[4] + a[4]*b[3] - SQ2*a[5]*b[2] + SQ2*a[5]*b[0]) / TWOSQ2
	if math.Abs(c6) > nonsymTol {
		err = chk.Err(_mandelops_err2, "3D", c6)
		return
	}
	if math.Abs(c7) > nonsymTol {
		err = chk.Err(_mandelops_err2, "3D", c7)
		return
	}
	if math.Abs(c8) > nonsymTol {
		err = chk.Err(_mandelops_err2, "3D", c8)
		return
	}
	return
}

// M_Inv computes the inverse of a 2nd order symmetric tensor 'a'
func M_Inv(ai, a []float64, tol float64) (det float64, err error) {
	if len(a) == 4 {
		det = a[0]*a[1]*a[2] - a[2]*a[3]*a[3]/2.0
		if math.Abs(det) < tol {
			return 0, chk.Err(_mandelops_err1, a, det, tol)
		}
		ai[0] = a[1] * a[2] / det
		ai[1] = a[0] * a[2] / det
		ai[2] = (a[0]*a[1] - a[3]*a[3]/2.0) / det
		ai[3] = -a[2] * a[3] / det
		return
	}
	det = a[0]*a[1]*a[2] + a[3]*a[4]*a[5]/SQ2 - a[0]*a[4]*a[4]/2.0 - a[1]*a[5]*a[5]/2.0 - a[2]*a[3]*a[3]/2.0
	if math.Abs(det) < tol {
		return 0, chk.Err(_mandelops_err1, a, det, tol)
	}
	ai[0] = (a[1]*a[2] - a[4]*a[4]/2.0) / det
	ai[1] = (a[0]*a[2] - a[5]*a[5]/2.0) / det
	ai[2] = (a[0]*a[1] - a[3]*a[3]/2.0) / det
	ai[3] = (a[4]*a[5]/SQ2 - a[2]*a[3]) / det
	ai[4] = (a[3]*a[5]/SQ2 - a[0]*a[4]) / det
	ai[5] = (a[3]*a[4]/SQ2 - a[1]*a[5]) / det
	return
}

// M_InvDeriv computes the derivative of the inverse of a tensor with respect to itself
//  ai := inv(a)
//  d  := dai/da
func M_InvDeriv(d [][]float64, ai []float64) {
	if len(ai) == 4 {
		d[0][0] = -ai[0] * ai[0]
		d[0][1] = -ai[3] * ai[3] / 2.0
		d[0][2] = 0
		d[0][3] = -ai[0] * ai[3]
		d[1][0] = -ai[3] * ai[3] / 2.0
		d[1][1] = -ai[1] * ai[1]
		d[1][2] = 0
		d[1][3] = -ai[1] * ai[3]
		d[2][0] = 0
		d[2][1] = 0
		d[2][2] = -ai[2] * ai[2]
		d[2][3] = 0
		d[3][0] = -ai[0] * ai[3]
		d[3][1] = -ai[1] * ai[3]
		d[3][2] = 0
		d[3][3] = -(ai[3]*ai[3] + 2.0*ai[0]*ai[1]) / 2.0
		return
	}
	d[0][0] = -ai[0] * ai[0]
	d[0][1] = -ai[3] * ai[3] / 2.0
	d[0][2] = -ai[5] * ai[5] / 2.0
	d[0][3] = -ai[0] * ai[3]
	d[0][4] = -ai[3] * ai[5] / SQ2
	d[0][5] = -ai[0] * ai[5]
	d[1][0] = -ai[3] * ai[3] / 2.0
	d[1][1] = -ai[1] * ai[1]
	d[1][2] = -ai[4] * ai[4] / 2.0
	d[1][3] = -ai[1] * ai[3]
	d[1][4] = -ai[1] * ai[4]
	d[1][5] = -ai[3] * ai[4] / SQ2
	d[2][0] = -ai[5] * ai[5] / 2.0
	d[2][1] = -ai[4] * ai[4] / 2.0
	d[2][2] = -ai[2] * ai[2]
	d[2][3] = -ai[4] * ai[5] / SQ2
	d[2][4] = -ai[2] * ai[4]
	d[2][5] = -ai[2] * ai[5]
	d[3][0] = -ai[0] * ai[3]
	d[3][1] = -ai[1] * ai[3]
	d[3][2] = -ai[4] * ai[5] / SQ2
	d[3][3] = -(ai[3]*ai[3] + 2.0*ai[0]*ai[1]) / 2.0
	d[3][4] = -(2.0*ai[1]*ai[5] + SQ2*ai[3]*ai[4]) / TWOSQ2
	d[3][5] = -(SQ2*ai[3]*ai[5] + 2.0*ai[0]*ai[4]) / TWOSQ2
	d[4][0] = -ai[3] * ai[5] / SQ2
	d[4][1] = -ai[1] * ai[4]
	d[4][2] = -ai[2] * ai[4]
	d[4][3] = -(2.0*ai[1]*ai[5] + SQ2*ai[3]*ai[4]) / TWOSQ2
	d[4][4] = -(ai[4]*ai[4] + 2.0*ai[1]*ai[2]) / 2.0
	d[4][5] = -(SQ2*ai[4]*ai[5] + 2.0*ai[2]*ai[3]) / TWOSQ2
	d[5][0] = -ai[0] * ai[5]
	d[5][1] = -ai[3] * ai[4] / SQ2
	d[5][2] = -ai[2] * ai[5]
	d[5][3] = -(SQ2*ai[3]*ai[5] + 2.0*ai[0]*ai[4]) / TWOSQ2
	d[5][4] = -(SQ2*ai[4]*ai[5] + 2.0*ai[2]*ai[3]) / TWOSQ2
	d[5][5] = -(ai[5]*ai[5] + 2.0*ai[0]*ai[2]) / 2.0
}

// error messages
var (
	_mandelops_err1 = "mandelops.go: M_Inv: inverse of 2nd order symmetric tensor %v failed with zero determinant (det = %g < %g)"
	_mandelops_err2 = "mandelops.go: M_Dot: (%s) dot product between two 2nd order symmetric tensor resulted in a non-symmetric tensor. err = %g"
)
