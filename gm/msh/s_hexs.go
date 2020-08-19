// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import "gosl/la"

// FuncHex8 calculates the shape functions (S) and derivatives of shape functions (dSdR) of hex8
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//              4________________7
//            ,'|              ,'|
//          ,'  |            ,'  |
//        ,'    |          ,'    |
//      ,'      |        ,'      |
//    5'===============6'        |
//    |         |      |         |
//    |         |      |         |
//    |         0_____ | ________3
//    |       ,'       |       ,'
//    |     ,'         |     ,'
//    |   ,'           |   ,'
//    | ,'             | ,'
//    1________________2'
//
func FuncHex8(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s, t := R[0], R[1], R[2]
	S[0] = (1.0 - r - s + r*s - t + s*t + r*t - r*s*t) / 8.0
	S[1] = (1.0 + r - s - r*s - t + s*t - r*t + r*s*t) / 8.0
	S[2] = (1.0 + r + s + r*s - t - s*t - r*t - r*s*t) / 8.0
	S[3] = (1.0 - r + s - r*s - t - s*t + r*t + r*s*t) / 8.0
	S[4] = (1.0 - r - s + r*s + t - s*t - r*t + r*s*t) / 8.0
	S[5] = (1.0 + r - s - r*s + t - s*t + r*t - r*s*t) / 8.0
	S[6] = (1.0 + r + s + r*s + t + s*t + r*t + r*s*t) / 8.0
	S[7] = (1.0 - r + s - r*s + t + s*t - r*t - r*s*t) / 8.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, (-1.0+s+t-s*t)/8.0)
	dSdR.Set(0, 1, (-1.0+r+t-r*t)/8.0)
	dSdR.Set(0, 2, (-1.0+r+s-r*s)/8.0)

	dSdR.Set(1, 0, (+1.0-s-t+s*t)/8.0)
	dSdR.Set(1, 1, (-1.0-r+t+r*t)/8.0)
	dSdR.Set(1, 2, (-1.0-r+s+r*s)/8.0)

	dSdR.Set(2, 0, (+1.0+s-t-s*t)/8.0)
	dSdR.Set(2, 1, (+1.0+r-t-r*t)/8.0)
	dSdR.Set(2, 2, (-1.0-r-s-r*s)/8.0)

	dSdR.Set(3, 0, (-1.0-s+t+s*t)/8.0)
	dSdR.Set(3, 1, (+1.0-r-t+r*t)/8.0)
	dSdR.Set(3, 2, (-1.0+r-s+r*s)/8.0)

	dSdR.Set(4, 0, (-1.0+s-t+s*t)/8.0)
	dSdR.Set(4, 1, (-1.0+r-t+r*t)/8.0)
	dSdR.Set(4, 2, (+1.0-r-s+r*s)/8.0)

	dSdR.Set(5, 0, (+1.0-s+t-s*t)/8.0)
	dSdR.Set(5, 1, (-1.0-r-t-r*t)/8.0)
	dSdR.Set(5, 2, (+1.0+r-s-r*s)/8.0)

	dSdR.Set(6, 0, (+1.0+s+t+s*t)/8.0)
	dSdR.Set(6, 1, (+1.0+r+t+r*t)/8.0)
	dSdR.Set(6, 2, (+1.0+r+s+r*s)/8.0)

	dSdR.Set(7, 0, (-1.0-s-t-s*t)/8.0)
	dSdR.Set(7, 1, (+1.0-r+t-r*t)/8.0)
	dSdR.Set(7, 2, (+1.0-r+s-r*s)/8.0)
}

// FuncHex20 calculates the shape functions (S) and derivatives of shape functions (dSdR) of hex20
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//                4_______15_______7
//              ,'|              ,'|
//           12'  |            ,'  |
//          ,'    16         ,14   |
//        ,'      |        ,'      19
//      5'=====13========6'        |
//      |         |      |         |
//      |         |      |         |
//      |         0_____ | _11_____3
//     17       ,'       |       ,'
//      |     8'        18     ,'
//      |   ,'           |   ,10
//      | ,'             | ,'
//      1_______9________2'
//
func FuncHex20(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s, t := R[0], R[1], R[2]
	rp1 := 1.0 + r
	rm1 := 1.0 - r
	sp1 := 1.0 + s
	sm1 := 1.0 - s
	tp1 := 1.0 + t
	tm1 := 1.0 - t

	S[0] = rm1 * sm1 * tm1 * (-r - s - t - 2) / 8.0
	S[1] = rp1 * sm1 * tm1 * (r - s - t - 2) / 8.0
	S[2] = rp1 * sp1 * tm1 * (r + s - t - 2) / 8.0
	S[3] = rm1 * sp1 * tm1 * (-r + s - t - 2) / 8.0
	S[4] = rm1 * sm1 * tp1 * (-r - s + t - 2) / 8.0
	S[5] = rp1 * sm1 * tp1 * (r - s + t - 2) / 8.0
	S[6] = rp1 * sp1 * tp1 * (r + s + t - 2) / 8.0
	S[7] = rm1 * sp1 * tp1 * (-r + s + t - 2) / 8.0
	S[8] = (1.0 - r*r) * sm1 * tm1 / 4.0
	S[9] = rp1 * (1.0 - s*s) * tm1 / 4.0
	S[10] = (1.0 - r*r) * sp1 * tm1 / 4.0
	S[11] = rm1 * (1.0 - s*s) * tm1 / 4.0
	S[12] = (1.0 - r*r) * sm1 * tp1 / 4.0
	S[13] = rp1 * (1.0 - s*s) * tp1 / 4.0
	S[14] = (1.0 - r*r) * sp1 * tp1 / 4.0
	S[15] = rm1 * (1.0 - s*s) * tp1 / 4.0
	S[16] = rm1 * sm1 * (1.0 - t*t) / 4.0
	S[17] = rp1 * sm1 * (1.0 - t*t) / 4.0
	S[18] = rp1 * sp1 * (1.0 - t*t) / 4.0
	S[19] = rm1 * sp1 * (1.0 - t*t) / 4.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, -0.125*sm1*tm1*(-r-s-t-2.0)-0.125*rm1*sm1*tm1)
	dSdR.Set(1, 0, 0.125*sm1*tm1*(r-s-t-2.0)+0.125*rp1*sm1*tm1)
	dSdR.Set(2, 0, 0.125*sp1*tm1*(r+s-t-2.0)+0.125*rp1*sp1*tm1)
	dSdR.Set(3, 0, -0.125*sp1*tm1*(-r+s-t-2.0)-0.125*rm1*sp1*tm1)
	dSdR.Set(4, 0, -0.125*sm1*tp1*(-r-s+t-2.0)-0.125*rm1*sm1*tp1)
	dSdR.Set(5, 0, 0.125*sm1*tp1*(r-s+t-2.0)+0.125*rp1*sm1*tp1)
	dSdR.Set(6, 0, 0.125*sp1*tp1*(r+s+t-2.0)+0.125*rp1*sp1*tp1)
	dSdR.Set(7, 0, -0.125*sp1*tp1*(-r+s+t-2.0)-0.125*rm1*sp1*tp1)
	dSdR.Set(8, 0, -0.5*r*sm1*tm1)
	dSdR.Set(9, 0, 0.25*(1.0-s*s)*tm1)
	dSdR.Set(10, 0, -0.5*r*sp1*tm1)
	dSdR.Set(11, 0, -0.25*(1.0-s*s)*tm1)
	dSdR.Set(12, 0, -0.5*r*sm1*tp1)
	dSdR.Set(13, 0, 0.25*(1.0-s*s)*tp1)
	dSdR.Set(14, 0, -0.5*r*sp1*tp1)
	dSdR.Set(15, 0, -0.25*(1.0-s*s)*tp1)
	dSdR.Set(16, 0, -0.25*sm1*(1.0-t*t))
	dSdR.Set(17, 0, 0.25*sm1*(1.0-t*t))
	dSdR.Set(18, 0, 0.25*sp1*(1.0-t*t))
	dSdR.Set(19, 0, -0.25*sp1*(1.0-t*t))

	dSdR.Set(0, 1, -0.125*rm1*tm1*(-r-s-t-2.0)-0.125*rm1*sm1*tm1)
	dSdR.Set(1, 1, -0.125*rp1*tm1*(r-s-t-2.0)-0.125*rp1*sm1*tm1)
	dSdR.Set(2, 1, 0.125*rp1*tm1*(r+s-t-2.0)+0.125*rp1*sp1*tm1)
	dSdR.Set(3, 1, 0.125*rm1*tm1*(-r+s-t-2.0)+0.125*rm1*sp1*tm1)
	dSdR.Set(4, 1, -0.125*rm1*tp1*(-r-s+t-2.0)-0.125*rm1*sm1*tp1)
	dSdR.Set(5, 1, -0.125*rp1*tp1*(r-s+t-2.0)-0.125*rp1*sm1*tp1)
	dSdR.Set(6, 1, 0.125*rp1*tp1*(r+s+t-2.0)+0.125*rp1*sp1*tp1)
	dSdR.Set(7, 1, 0.125*rm1*tp1*(-r+s+t-2.0)+0.125*rm1*sp1*tp1)
	dSdR.Set(8, 1, -0.25*(1.0-r*r)*tm1)
	dSdR.Set(9, 1, -0.5*s*rp1*tm1)
	dSdR.Set(10, 1, 0.25*(1.0-r*r)*tm1)
	dSdR.Set(11, 1, -0.5*s*rm1*tm1)
	dSdR.Set(12, 1, -0.25*(1.0-r*r)*tp1)
	dSdR.Set(13, 1, -0.5*s*rp1*tp1)
	dSdR.Set(14, 1, 0.25*(1.0-r*r)*tp1)
	dSdR.Set(15, 1, -0.5*s*rm1*tp1)
	dSdR.Set(16, 1, -0.25*rm1*(1.0-t*t))
	dSdR.Set(17, 1, -0.25*rp1*(1.0-t*t))
	dSdR.Set(18, 1, 0.25*rp1*(1.0-t*t))
	dSdR.Set(19, 1, 0.25*rm1*(1.0-t*t))

	dSdR.Set(0, 2, -0.125*rm1*sm1*(-r-s-t-2.0)-0.125*rm1*sm1*tm1)
	dSdR.Set(1, 2, -0.125*rp1*sm1*(r-s-t-2.0)-0.125*rp1*sm1*tm1)
	dSdR.Set(2, 2, -0.125*rp1*sp1*(r+s-t-2.0)-0.125*rp1*sp1*tm1)
	dSdR.Set(3, 2, -0.125*rm1*sp1*(-r+s-t-2.0)-0.125*rm1*sp1*tm1)
	dSdR.Set(4, 2, 0.125*rm1*sm1*(-r-s+t-2.0)+0.125*rm1*sm1*tp1)
	dSdR.Set(5, 2, 0.125*rp1*sm1*(r-s+t-2.0)+0.125*rp1*sm1*tp1)
	dSdR.Set(6, 2, 0.125*rp1*sp1*(r+s+t-2.0)+0.125*rp1*sp1*tp1)
	dSdR.Set(7, 2, 0.125*rm1*sp1*(-r+s+t-2.0)+0.125*rm1*sp1*tp1)
	dSdR.Set(8, 2, -0.25*(1.0-r*r)*sm1)
	dSdR.Set(9, 2, -0.25*rp1*(1.0-s*s))
	dSdR.Set(10, 2, -0.25*(1.0-r*r)*sp1)
	dSdR.Set(11, 2, -0.25*rm1*(1.0-s*s))
	dSdR.Set(12, 2, 0.25*(1.0-r*r)*sm1)
	dSdR.Set(13, 2, 0.25*rp1*(1.0-s*s))
	dSdR.Set(14, 2, 0.25*(1.0-r*r)*sp1)
	dSdR.Set(15, 2, 0.25*rm1*(1.0-s*s))
	dSdR.Set(16, 2, -0.5*t*rm1*sm1)
	dSdR.Set(17, 2, -0.5*t*rp1*sm1)
	dSdR.Set(18, 2, -0.5*t*rp1*sp1)
	dSdR.Set(19, 2, -0.5*t*rm1*sp1)
}
