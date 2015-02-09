// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// Tr returns the trace of a second order tensor
func Tr(a [][]float64) float64 {
	return a[0][0] + a[1][1] + a[2][2]
}

// Dev returns the second order deviatoric tensor
func Dev(a [][]float64) (deva [][]float64) {
	p := -Tr(a) / 3.0
	deva = [][]float64{
		{a[0][0] + p, a[0][1], a[0][2]},
		{a[1][0], a[1][1] + p, a[1][2]},
		{a[2][0], a[2][1], a[2][2] + p},
	}
	return
}

// Det computes the determinant of a second order tensor
func Det(a [][]float64) float64 {
	return a[0][0]*(a[1][1]*a[2][2]-a[1][2]*a[2][1]) - a[0][1]*(a[1][0]*a[2][2]-a[1][2]*a[2][0]) + a[0][2]*(a[1][0]*a[2][1]-a[1][1]*a[2][0])
}

// Inv computes the inverse of a second order tensor
//  ai := Inv(a)
func Inv(ai, a [][]float64) (det float64, err error) {
	det = a[0][0]*(a[1][1]*a[2][2]-a[1][2]*a[2][1]) - a[0][1]*(a[1][0]*a[2][2]-a[1][2]*a[2][0]) + a[0][2]*(a[1][0]*a[2][1]-a[1][1]*a[2][0])
	if math.Abs(det) > MINDET {

		ai[0][0] = (a[1][1]*a[2][2] - a[1][2]*a[2][1]) / det
		ai[0][1] = (a[0][2]*a[2][1] - a[0][1]*a[2][2]) / det
		ai[0][2] = (a[0][1]*a[1][2] - a[0][2]*a[1][1]) / det

		ai[1][0] = (a[1][2]*a[2][0] - a[1][0]*a[2][2]) / det
		ai[1][1] = (a[0][0]*a[2][2] - a[0][2]*a[2][0]) / det
		ai[1][2] = (a[0][2]*a[1][0] - a[0][0]*a[1][2]) / det

		ai[2][0] = (a[1][0]*a[2][1] - a[1][1]*a[2][0]) / det
		ai[2][1] = (a[0][1]*a[2][0] - a[0][0]*a[2][1]) / det
		ai[2][2] = (a[0][0]*a[1][1] - a[0][1]*a[1][0]) / det
		return
	}
	err = utl.Err(_tsr_inv_1, det)
	return
}

// Add adds two second order tensors according to:
//  u := α*a + β*b
func Add(u [][]float64, α float64, a [][]float64, β float64, b [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			u[i][j] = α*a[i][j] + β*b[i][j]
		}
	}
}

// errors
var (
	_tsr_inv_1 = "Inv: determinant of tensor is near zero (%g)\n"
)
