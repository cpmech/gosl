// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import "github.com/cpmech/gosl/la"

// FuncLin2 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin2
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    -1     0    +1
//     0-----------1-->r
//
func FuncLin2(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r := R[0]
	S[0] = 0.5 * (1.0 - r)
	S[1] = 0.5 * (1.0 + r)

	if !derivs {
		return
	}

	dSdR.Set(0, 0, -0.5)
	dSdR.Set(1, 0, 0.5)
}

// FuncLin3 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin3
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    -1     0    +1
//     0-----2-----1-->r
//
func FuncLin3(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r := R[0]
	S[0] = 0.5 * (r*r - r)
	S[1] = 0.5 * (r*r + r)
	S[2] = 1.0 - r*r

	if !derivs {
		return
	}

	dSdR.Set(0, 0, r-0.5)
	dSdR.Set(1, 0, r+0.5)
	dSdR.Set(2, 0, -2.0*r)
}

// FuncLin4 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin4
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    -1                  +1
//     @------@-----@------@  --> r
//     0      2     3      1
//
func FuncLin4(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r := R[0]
	S[0] = (-9.0*r*r*r + 9.0*r*r + r - 1.0) / 16.0
	S[1] = (9.0*r*r*r + 9.0*r*r - r - 1.0) / 16.0
	S[2] = (27.0*r*r*r - 9.0*r*r - 27.0*r + 9.0) / 16.0
	S[3] = (-27.0*r*r*r - 9.0*r*r + 27.0*r + 9.0) / 16.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, 1.0/16.0*(-27.0*r*r+18.0*r+1.0))
	dSdR.Set(1, 0, 1.0/16.0*(27.0*r*r+18.0*r-1.0))
	dSdR.Set(2, 0, 1.0/16.0*(81.0*r*r-18.0*r-27.0))
	dSdR.Set(3, 0, 1.0/16.0*(-81.0*r*r-18.0*r+27.0))
}

// FuncLin5 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin5
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//      @-----@-----@-----@-----@-> r
//      0     3     2     4     1
//      |           |           |
//     r=-1  -1/2   r=0  1/2   r=+1
//
func FuncLin5(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r := R[0]
	S[0] = (r - 1.0) * (1.0 - 2.0*r) * r * (-1.0 - 2.0*r) / 6.0
	S[1] = (1.0 - 2.0*r) * r * (-1.0 - 2.0*r) * (1.0 + r) / 6.0
	S[2] = (1.0 - r) * (1.0 - 2.0*r) * (-1.0 - 2.0*r) * (-1.0 - r)
	S[3] = 4.0 * (1.0 - r) * (1.0 - 2.0*r) * r * (-1.0 - r) / 3.0
	S[4] = 4.0 * (1.0 - r) * r * (-1.0 - 2.0*r) * (-1.0 - r) / 3.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, -((1.0-2.0*r)*(r-1.0)*r)/3.0-((-2.0*r-1.0)*(r-1.0)*r)/3.0+((-2.0*r-1.0)*(1.0-2.0*r)*r)/6.0+((-2.0*r-1.0)*(1.0-2.0*r)*(r-1.0))/6.0)
	dSdR.Set(1, 0, -((1.0-2.0*r)*r*(r+1.0))/3.0-((-2.0*r-1.0)*r*(r+1.0))/3.0+((-2.0*r-1.0)*(1.0-2.0*r)*(r+1.0))/6.0+((-2.0*r-1.0)*(1.0-2.0*r)*r)/6.0)
	dSdR.Set(2, 0, -2.0*(1.0-2.0*r)*(-r-1.0)*(1.0-r)-2.0*(-2.0*r-1.0)*(-r-1.0)*(1.0-r)-(-2.0*r-1.0)*(1.0-2.0*r)*(1.0-r)-(-2.0*r-1.0)*(1.0-2.0*r)*(-r-1.0))
	dSdR.Set(3, 0, -(8.0*(-r-1.0)*(1.0-r)*r)/3.0-(4.0*(1.0-2.0*r)*(1.0-r)*r)/3.0-(4.0*(1.0-2.0*r)*(-r-1.0)*r)/3.0+(4.0*(1.0-2.0*r)*(-r-1.0)*(1.0-r))/3.0)
	dSdR.Set(4, 0, -(8.0*(-r-1.0)*(1.0-r)*r)/3.0-(4.0*(-2.0*r-1.0)*(1.0-r)*r)/3.0-(4.0*(-2.0*r-1.0)*(-r-1.0)*r)/3.0+(4.0*(-2.0*r-1.0)*(-r-1.0)*(1.0-r))/3.0)
}
