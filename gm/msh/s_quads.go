// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import "gosl/la"

// FuncQua4 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua4
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    3-----------2
//    |     s     |
//    |     |     |
//    |     +--r  |
//    |           |
//    |           |
//    0-----------1
//
func FuncQua4(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	S[0] = (1.0 - r - s + r*s) / 4.0
	S[1] = (1.0 + r - s - r*s) / 4.0
	S[2] = (1.0 + r + s + r*s) / 4.0
	S[3] = (1.0 - r + s - r*s) / 4.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, (-1.0+s)/4.0)
	dSdR.Set(0, 1, (-1.0+r)/4.0)
	dSdR.Set(1, 0, (+1.0-s)/4.0)
	dSdR.Set(1, 1, (-1.0-r)/4.0)
	dSdR.Set(2, 0, (+1.0+s)/4.0)
	dSdR.Set(2, 1, (+1.0+r)/4.0)
	dSdR.Set(3, 0, (-1.0-s)/4.0)
	dSdR.Set(3, 1, (+1.0-r)/4.0)
}

// FuncQua8 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua8
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    3-----6-----2
//    |     s     |
//    |     |     |
//    7     +--r  5
//    |           |
//    |           |
//    0-----4-----1
//
func FuncQua8(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	S[0] = (1.0 - r) * (1.0 - s) * (-r - s - 1.0) / 4.0
	S[1] = (1.0 + r) * (1.0 - s) * (r - s - 1.0) / 4.0
	S[2] = (1.0 + r) * (1.0 + s) * (r + s - 1.0) / 4.0
	S[3] = (1.0 - r) * (1.0 + s) * (-r + s - 1.0) / 4.0
	S[4] = (1.0 - s) * (1.0 - r*r) / 2.0
	S[5] = (1.0 + r) * (1.0 - s*s) / 2.0
	S[6] = (1.0 + s) * (1.0 - r*r) / 2.0
	S[7] = (1.0 - r) * (1.0 - s*s) / 2.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, -(1.0-s)*(-r-r-s)/4.0)
	dSdR.Set(1, 0, (1.0-s)*(r+r-s)/4.0)
	dSdR.Set(2, 0, (1.0+s)*(r+r+s)/4.0)
	dSdR.Set(3, 0, -(1.0+s)*(-r-r+s)/4.0)
	dSdR.Set(4, 0, -(1.0-s)*r)
	dSdR.Set(5, 0, (1.0-s*s)/2.0)
	dSdR.Set(6, 0, -(1.0+s)*r)
	dSdR.Set(7, 0, -(1.0-s*s)/2.0)

	dSdR.Set(0, 1, -(1.0-r)*(-s-s-r)/4.0)
	dSdR.Set(1, 1, -(1.0+r)*(-s-s+r)/4.0)
	dSdR.Set(2, 1, (1.0+r)*(s+s+r)/4.0)
	dSdR.Set(3, 1, (1.0-r)*(s+s-r)/4.0)
	dSdR.Set(4, 1, -(1.0-r*r)/2.0)
	dSdR.Set(5, 1, -(1.0+r)*s)
	dSdR.Set(6, 1, (1.0-r*r)/2.0)
	dSdR.Set(7, 1, -(1.0-r)*s)
}

// FuncQua9 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua9
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//    3-----6-----2
//    |     s     |
//    |     |     |
//    7     8--r  5
//    |           |
//    |           |
//    0-----4-----1
//
func FuncQua9(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	S[0] = r * (r - 1.0) * s * (s - 1.0) / 4.0
	S[1] = r * (r + 1.0) * s * (s - 1.0) / 4.0
	S[2] = r * (r + 1.0) * s * (s + 1.0) / 4.0
	S[3] = r * (r - 1.0) * s * (s + 1.0) / 4.0

	S[4] = -(r*r - 1.0) * s * (s - 1.0) / 2.0
	S[5] = -r * (r + 1.0) * (s*s - 1.0) / 2.0
	S[6] = -(r*r - 1.0) * s * (s + 1.0) / 2.0
	S[7] = -r * (r - 1.0) * (s*s - 1.0) / 2.0

	S[8] = (r*r - 1.0) * (s*s - 1.0)

	if !derivs {
		return
	}

	dSdR.Set(0, 0, (r+r-1.0)*s*(s-1.0)/4.0)
	dSdR.Set(1, 0, (r+r+1.0)*s*(s-1.0)/4.0)
	dSdR.Set(2, 0, (r+r+1.0)*s*(s+1.0)/4.0)
	dSdR.Set(3, 0, (r+r-1.0)*s*(s+1.0)/4.0)

	dSdR.Set(0, 1, r*(r-1.0)*(s+s-1.0)/4.0)
	dSdR.Set(1, 1, r*(r+1.0)*(s+s-1.0)/4.0)
	dSdR.Set(2, 1, r*(r+1.0)*(s+s+1.0)/4.0)
	dSdR.Set(3, 1, r*(r-1.0)*(s+s+1.0)/4.0)

	dSdR.Set(4, 0, -(r+r)*s*(s-1.0)/2.0)
	dSdR.Set(5, 0, -(r+r+1.0)*(s*s-1.0)/2.0)
	dSdR.Set(6, 0, -(r+r)*s*(s+1.0)/2.0)
	dSdR.Set(7, 0, -(r+r-1.0)*(s*s-1.0)/2.0)

	dSdR.Set(4, 1, -(r*r-1.0)*(s+s-1.0)/2.0)
	dSdR.Set(5, 1, -r*(r+1.0)*(s+s)/2.0)
	dSdR.Set(6, 1, -(r*r-1.0)*(s+s+1.0)/2.0)
	dSdR.Set(7, 1, -r*(r-1.0)*(s+s)/2.0)

	dSdR.Set(8, 0, 2.0*r*(s*s-1.0))
	dSdR.Set(8, 1, 2.0*s*(r*r-1.0))
}

// FuncQua12 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua12
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//     3      10       6        2
//       @-----@-------@------@
//       |               (1,1)|
//       |       s ^          |
//     7 @         |          @ 9
//       |         |          |
//       |         +----> r   |
//       |       (0,0)        |
//    11 @                    @ 5
//       |                    |
//       |(-1,-1)             |
//       @-----@-------@------@
//     0       4       8        1
//
func FuncQua12(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	rm := 1.0 - r
	rp := 1.0 + r
	sm := 1.0 - s
	sp := 1.0 + s

	S[0] = rm * sm * (9.0*(r*r+s*s) - 10.0) / 32.0
	S[1] = rp * sm * (9.0*(r*r+s*s) - 10.0) / 32.0
	S[2] = rp * sp * (9.0*(r*r+s*s) - 10.0) / 32.0
	S[3] = rm * sp * (9.0*(r*r+s*s) - 10.0) / 32.0
	S[4] = 9.0 * (1.0 - r*r) * (1.0 - 3.0*r) * sm / 32.0
	S[5] = 9.0 * (1.0 - s*s) * (1.0 - 3.0*s) * rp / 32.0
	S[6] = 9.0 * (1.0 - r*r) * (1.0 + 3.0*r) * sp / 32.0
	S[7] = 9.0 * (1.0 - s*s) * (1.0 + 3.0*s) * rm / 32.0
	S[8] = 9.0 * (1.0 - r*r) * (1.0 + 3.0*r) * sm / 32.0
	S[9] = 9.0 * (1.0 - s*s) * (1.0 + 3.0*s) * rp / 32.0
	S[10] = 9.0 * (1.0 - r*r) * (1.0 - 3.0*r) * sp / 32.0
	S[11] = 9.0 * (1.0 - s*s) * (1.0 - 3.0*s) * rm / 32.0

	if !derivs {
		return
	}

	dSdR.Set(0, 0, sm*(9.0*(2.0*r-3.0*r*r-s*s)+10.0)/32.0)
	dSdR.Set(1, 0, sm*(9.0*(2.0*r+3.0*r*r+s*s)-10.0)/32.0)
	dSdR.Set(2, 0, sp*(9.0*(2.0*r+3.0*r*r+s*s)-10.0)/32.0)
	dSdR.Set(3, 0, sp*(9.0*(2.0*r-3.0*r*r-s*s)+10.0)/32.0)
	dSdR.Set(4, 0, 9.0*sm*(9.0*r*r-2.0*r-3.0)/32.0)
	dSdR.Set(5, 0, 9.0*(1.0-s*s)*(1.0-3.0*s)/32.0)
	dSdR.Set(6, 0, 9.0*sp*(-9.0*r*r-2.0*r+3.0)/32.0)
	dSdR.Set(7, 0, -9.0*(1.0-s*s)*(1.0+3.0*s)/32.0)
	dSdR.Set(8, 0, 9.0*sm*(-9.0*r*r-2.0*r+3.0)/32.0)
	dSdR.Set(9, 0, 9.0*(1.0-s*s)*(1.0+3.0*s)/32.0)
	dSdR.Set(10, 0, 9.0*sp*(9.0*r*r-2.0*r-3.0)/32.0)
	dSdR.Set(11, 0, -9.0*(1.0-s*s)*(1.0-3.0*s)/32.0)

	dSdR.Set(0, 1, rm*(9.0*(2.0*s-3.0*s*s-r*r)+10.0)/32.0)
	dSdR.Set(1, 1, rp*(9.0*(2.0*s-3.0*s*s-r*r)+10.0)/32.0)
	dSdR.Set(2, 1, rp*(9.0*(2.0*s+3.0*s*s+r*r)-10.0)/32.0)
	dSdR.Set(3, 1, rm*(9.0*(2.0*s+3.0*s*s+r*r)-10.0)/32.0)
	dSdR.Set(4, 1, -9.0*(1.0-r*r)*(1.0-3.0*r)/32.0)
	dSdR.Set(5, 1, 9.0*rp*(9.0*s*s-2.0*s-3.0)/32.0)
	dSdR.Set(6, 1, 9.0*(1.0-r*r)*(1.0+3.0*r)/32.0)
	dSdR.Set(7, 1, 9.0*rm*(-9.0*s*s-2.0*s+3.0)/32.0)
	dSdR.Set(8, 1, -9.0*(1.0-r*r)*(1.0+3.0*r)/32.0)
	dSdR.Set(9, 1, 9.0*rp*(-9.0*s*s-2.0*s+3.0)/32.0)
	dSdR.Set(10, 1, 9.0*(1.0-r*r)*(1.0-3.0*r)/32.0)
	dSdR.Set(11, 1, 9.0*rm*(9.0*s*s-2.0*s-3.0)/32.0)
}

// FuncQua16 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua16
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//     3      10       6        2
//       @-----@-------@------@
//       |               (1,1)|
//       |       s ^          |
//     7 @   15@   |    @14   @ 9
//       |         |          |
//       |         +----> r   |
//       |       (0,0)        |
//    11 @   12@       @13    @ 5
//       |                    |
//       |(-1,-1)             |
//       @-----@-------@------@
//     0       4       8        1
//
func FuncQua16(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	sr, ss := make(la.Vector, 4), make(la.Vector, 4)
	var dr, ds *la.Matrix
	if derivs {
		dr, ds = la.NewMatrix(4, 1), la.NewMatrix(4, 1)
	}

	FuncLin4(sr, dr, la.Vector{r}, derivs)
	FuncLin4(ss, ds, la.Vector{s}, derivs)

	S[0] = sr[0] * ss[0]
	S[1] = sr[1] * ss[0]
	S[2] = sr[1] * ss[1]
	S[3] = sr[0] * ss[1]
	S[4] = sr[2] * ss[0]
	S[5] = sr[1] * ss[2]
	S[6] = sr[3] * ss[1]
	S[7] = sr[0] * ss[3]
	S[8] = sr[3] * ss[0]
	S[9] = sr[1] * ss[3]
	S[10] = sr[2] * ss[1]
	S[11] = sr[0] * ss[2]
	S[12] = sr[2] * ss[2]
	S[13] = sr[3] * ss[2]
	S[14] = sr[3] * ss[3]
	S[15] = sr[2] * ss[3]

	if !derivs {
		return
	}

	dSdR.Set(0, 0, dr.Get(0, 0)*ss[0])
	dSdR.Set(1, 0, dr.Get(1, 0)*ss[0])
	dSdR.Set(2, 0, dr.Get(1, 0)*ss[1])
	dSdR.Set(3, 0, dr.Get(0, 0)*ss[1])
	dSdR.Set(4, 0, dr.Get(2, 0)*ss[0])
	dSdR.Set(5, 0, dr.Get(1, 0)*ss[2])
	dSdR.Set(6, 0, dr.Get(3, 0)*ss[1])
	dSdR.Set(7, 0, dr.Get(0, 0)*ss[3])
	dSdR.Set(8, 0, dr.Get(3, 0)*ss[0])
	dSdR.Set(9, 0, dr.Get(1, 0)*ss[3])
	dSdR.Set(10, 0, dr.Get(2, 0)*ss[1])
	dSdR.Set(11, 0, dr.Get(0, 0)*ss[2])
	dSdR.Set(12, 0, dr.Get(2, 0)*ss[2])
	dSdR.Set(13, 0, dr.Get(3, 0)*ss[2])
	dSdR.Set(14, 0, dr.Get(3, 0)*ss[3])
	dSdR.Set(15, 0, dr.Get(2, 0)*ss[3])

	dSdR.Set(0, 1, sr[0]*ds.Get(0, 0))
	dSdR.Set(1, 1, sr[1]*ds.Get(0, 0))
	dSdR.Set(2, 1, sr[1]*ds.Get(1, 0))
	dSdR.Set(3, 1, sr[0]*ds.Get(1, 0))
	dSdR.Set(4, 1, sr[2]*ds.Get(0, 0))
	dSdR.Set(5, 1, sr[1]*ds.Get(2, 0))
	dSdR.Set(6, 1, sr[3]*ds.Get(1, 0))
	dSdR.Set(7, 1, sr[0]*ds.Get(3, 0))
	dSdR.Set(8, 1, sr[3]*ds.Get(0, 0))
	dSdR.Set(9, 1, sr[1]*ds.Get(3, 0))
	dSdR.Set(10, 1, sr[2]*ds.Get(1, 0))
	dSdR.Set(11, 1, sr[0]*ds.Get(2, 0))
	dSdR.Set(12, 1, sr[2]*ds.Get(2, 0))
	dSdR.Set(13, 1, sr[3]*ds.Get(2, 0))
	dSdR.Set(14, 1, sr[3]*ds.Get(3, 0))
	dSdR.Set(15, 1, sr[2]*ds.Get(3, 0))
}

// FuncQua17 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua17
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//     3      14    10     6     2
//       @-----@-----@-----@-----@
//       |                  (1,1)|
//       |                       |
// 	7 @                       @ 13
//       |         s ^           |
//       |           |           |
//    11 @           |16         @ 9
//       |           @----> r    |
//       |         (0,0)         |
//    15 @                       @ 5
//       |                       |
//       |(-1,-1)                |
//       @-----@-----@-----@-----@
//     0       4     8    12       1
//
func FuncQua17(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool) {

	r, s := R[0], R[1]
	a := 2.0 / 3.0
	rr := r * r
	ss := s * s
	rs := r * s
	rp := 1.0 + r
	rm := 1.0 - r
	sp := 1.0 + s
	sm := 1.0 - s

	S[0] = rm * sm * (-4.0*r*(rr-1.0) - 4.0*s*(ss-1.0) + 3.0*rs) / 12.0
	S[1] = rp * sm * (4.0*r*(rr-1.0) - 4.0*s*(ss-1.0) - 3.0*rs) / 12.0
	S[2] = rp * sp * (4.0*r*(rr-1.0) + 4.0*s*(ss-1.0) + 3.0*rs) / 12.0
	S[3] = rm * sp * (-4.0*r*(rr-1.0) + 4.0*s*(ss-1.0) - 3.0*rs) / 12.0
	S[4] = -a * r * sm * rm * rp * (1.0 - 2.0*r)
	S[5] = -a * s * rp * sm * sp * (1.0 - 2.0*s)
	S[6] = a * r * sp * rm * rp * (1.0 + 2.0*r)
	S[7] = a * s * sm * sp * (1.0 + 2.0*s) * rm
	S[8] = 0.5 * rm * rp * (-s - 4.0*rr) * sm
	S[9] = 0.5 * sm * sp * (r - 4.0*ss) * rp
	S[10] = 0.5 * rm * rp * (s - 4.0*rr) * sp
	S[11] = 0.5 * sm * sp * (-r - 4.0*ss) * rm
	S[12] = a * r * sm * rm * rp * (1.0 + 2.0*r)
	S[13] = a * s * rp * sm * sp * (1.0 + 2.0*s)
	S[14] = -a * r * sp * rm * rp * (1.0 - 2.0*r)
	S[15] = -a * s * rm * sm * sp * (1.0 - 2.0*s)
	S[16] = rm * rp * sm * sp

	if !derivs {
		return
	}

	b := 1.0 / 12.0
	r1 := r - 1.0
	rrr := rr * r
	sss := ss * s

	dSdR.Set(0, 0, b*sm*(16.0*rrr-12.0*rr-6.0*r*s-8.0*r+4.0*sss-s+4.0))
	dSdR.Set(1, 0, b*sm*(16.0*rrr+12.0*rr-6.0*r*s-8.0*r-4.0*sss+s-4.0))
	dSdR.Set(2, 0, b*sp*(16.0*rrr+12.0*rr+6.0*r*s-8.0*r+4.0*sss-s-4.0))
	dSdR.Set(3, 0, b*sp*(16.0*rrr-12.0*rr+6.0*r*s-8.0*r-4.0*sss+s+4.0))
	dSdR.Set(4, 0, -a*(1.0-4.0*r-3.0*rr+8.0*rrr)*sm)
	dSdR.Set(5, 0, a*s*sm*sp*(-1.0+2.0*s))
	dSdR.Set(6, 0, -a*(-1.0-4.0*r+3.0*rr+8.0*rrr)*sp)
	dSdR.Set(7, 0, -a*s*sm*sp*(1.0+2.0*s))
	dSdR.Set(8, 0, r*sm*(8.0*rr+s-4.0))
	dSdR.Set(9, 0, 0.5*sm*sp*(2.0*r-4.0*ss+1.0))
	dSdR.Set(10, 0, r*sp*(8.0*rr-s-4.0))
	dSdR.Set(11, 0, 0.5*sm*sp*(2.0*r-1.0+4.0*ss))
	dSdR.Set(12, 0, a*(1.0+4.0*r-3.0*rr-8.0*rrr)*sm)
	dSdR.Set(13, 0, a*s*sm*sp*(1.0+2.0*s))
	dSdR.Set(14, 0, -a*(1.0-4.0*r-3.0*rr+8.0*rrr)*sp)
	dSdR.Set(15, 0, a*s*sm*sp*(1.0-2.0*s))
	dSdR.Set(16, 0, -2.0*r*sm*sp)

	dSdR.Set(0, 1, b*rm*(16.0*sss-12.0*ss-6.0*r*s-8.0*s+4.0*rrr-r+4.0))
	dSdR.Set(1, 1, -b*rp*(-16.0*sss+12.0*ss-6.0*r*s+8.0*s+4.0*rrr-r-4.0))
	dSdR.Set(2, 1, b*rp*(16.0*sss+12.0*ss+6.0*r*s-8.0*s+4.0*rrr-r-4.0))
	dSdR.Set(3, 1, b*r1*(-16.0*sss-12.0*ss+6.0*r*s+8.0*s+4.0*rrr-r+4.0))
	dSdR.Set(4, 1, a*r*r1*rp*(2.0*r-1.0))
	dSdR.Set(5, 1, -a*(1.0-4.0*s-3.0*ss+8.0*sss)*rp)
	dSdR.Set(6, 1, -a*r*r1*rp*(1.0+2.0*r))
	dSdR.Set(7, 1, a*(-1.0-4.0*s+3.0*ss+8.0*sss)*r1)
	dSdR.Set(8, 1, -0.5*r1*rp*(2.0*s-1.0+4.0*rr))
	dSdR.Set(9, 1, -s*rp*(-8.0*ss+r+4.0))
	dSdR.Set(10, 1, 0.5*r1*rp*(-2.0*s+4.0*rr-1.0))
	dSdR.Set(11, 1, -s*r1*(8.0*ss+r-4.0))
	dSdR.Set(12, 1, a*r*r1*rp*(1.0+2.0*r))
	dSdR.Set(13, 1, -a*(-1.0-4.0*s+3.0*ss+8.0*sss)*rp)
	dSdR.Set(14, 1, -a*r*r1*rp*(2.0*r-1.0))
	dSdR.Set(15, 1, a*(1.0-4.0*s-3.0*ss+8.0*sss)*r1)
	dSdR.Set(16, 1, 2.0*s*r1*rp)
}
