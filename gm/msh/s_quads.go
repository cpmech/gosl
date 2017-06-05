// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

// FuncQua4 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua4
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua4(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   3-----------2
	   |     s     |
	   |     |     |
	   |     +--r  |
	   |           |
	   |           |
	   0-----------1
	*/
	r, s := R[0], R[1]
	S[0] = (1.0 - r - s + r*s) / 4.0
	S[1] = (1.0 + r - s - r*s) / 4.0
	S[2] = (1.0 + r + s + r*s) / 4.0
	S[3] = (1.0 - r + s - r*s) / 4.0

	if !derivs {
		return
	}

	dSdR[0][0] = (-1.0 + s) / 4.0
	dSdR[0][1] = (-1.0 + r) / 4.0
	dSdR[1][0] = (+1.0 - s) / 4.0
	dSdR[1][1] = (-1.0 - r) / 4.0
	dSdR[2][0] = (+1.0 + s) / 4.0
	dSdR[2][1] = (+1.0 + r) / 4.0
	dSdR[3][0] = (-1.0 - s) / 4.0
	dSdR[3][1] = (+1.0 - r) / 4.0
}

// FuncQua8 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua8
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua8(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   3-----6-----2
	   |     s     |
	   |     |     |
	   7     +--r  5
	   |           |
	   |           |
	   0-----4-----1
	*/
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

	dSdR[0][0] = -(1.0 - s) * (-r - r - s) / 4.0
	dSdR[1][0] = (1.0 - s) * (r + r - s) / 4.0
	dSdR[2][0] = (1.0 + s) * (r + r + s) / 4.0
	dSdR[3][0] = -(1.0 + s) * (-r - r + s) / 4.0
	dSdR[4][0] = -(1.0 - s) * r
	dSdR[5][0] = (1.0 - s*s) / 2.0
	dSdR[6][0] = -(1.0 + s) * r
	dSdR[7][0] = -(1.0 - s*s) / 2.0

	dSdR[0][1] = -(1.0 - r) * (-s - s - r) / 4.0
	dSdR[1][1] = -(1.0 + r) * (-s - s + r) / 4.0
	dSdR[2][1] = (1.0 + r) * (s + s + r) / 4.0
	dSdR[3][1] = (1.0 - r) * (s + s - r) / 4.0
	dSdR[4][1] = -(1.0 - r*r) / 2.0
	dSdR[5][1] = -(1.0 + r) * s
	dSdR[6][1] = (1.0 - r*r) / 2.0
	dSdR[7][1] = -(1.0 - r) * s
}

// FuncQua9 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua9
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua9(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   3-----6-----2
	   |     s     |
	   |     |     |
	   7     8--r  5
	   |           |
	   |           |
	   0-----4-----1
	*/
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

	dSdR[0][0] = (r + r - 1.0) * s * (s - 1.0) / 4.0
	dSdR[1][0] = (r + r + 1.0) * s * (s - 1.0) / 4.0
	dSdR[2][0] = (r + r + 1.0) * s * (s + 1.0) / 4.0
	dSdR[3][0] = (r + r - 1.0) * s * (s + 1.0) / 4.0

	dSdR[0][1] = r * (r - 1.0) * (s + s - 1.0) / 4.0
	dSdR[1][1] = r * (r + 1.0) * (s + s - 1.0) / 4.0
	dSdR[2][1] = r * (r + 1.0) * (s + s + 1.0) / 4.0
	dSdR[3][1] = r * (r - 1.0) * (s + s + 1.0) / 4.0

	dSdR[4][0] = -(r + r) * s * (s - 1.0) / 2.0
	dSdR[5][0] = -(r + r + 1.0) * (s*s - 1.0) / 2.0
	dSdR[6][0] = -(r + r) * s * (s + 1.0) / 2.0
	dSdR[7][0] = -(r + r - 1.0) * (s*s - 1.0) / 2.0

	dSdR[4][1] = -(r*r - 1.0) * (s + s - 1.0) / 2.0
	dSdR[5][1] = -r * (r + 1.0) * (s + s) / 2.0
	dSdR[6][1] = -(r*r - 1.0) * (s + s + 1.0) / 2.0
	dSdR[7][1] = -r * (r - 1.0) * (s + s) / 2.0

	dSdR[8][0] = 2.0 * r * (s*s - 1.0)
	dSdR[8][1] = 2.0 * s * (r*r - 1.0)
}

// FuncQua12 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua12
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua12(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	    3      10       6        2
	      @-----@-------@------@
	      |               (1,1)|
	      |       s ^          |
	    7 @         |          @ 9
	      |         |          |
	      |         +----> r   |
	      |       (0,0)        |
	   11 @                    @ 5
	      |                    |
	      |(-1,-1)             |
	      @-----@-------@------@
	    0       4       8        1
	*/
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

	dSdR[0][0] = sm * (9.0*(2.0*r-3.0*r*r-s*s) + 10.0) / 32.0
	dSdR[1][0] = sm * (9.0*(2.0*r+3.0*r*r+s*s) - 10.0) / 32.0
	dSdR[2][0] = sp * (9.0*(2.0*r+3.0*r*r+s*s) - 10.0) / 32.0
	dSdR[3][0] = sp * (9.0*(2.0*r-3.0*r*r-s*s) + 10.0) / 32.0
	dSdR[4][0] = 9.0 * sm * (9.0*r*r - 2.0*r - 3.0) / 32.0
	dSdR[5][0] = 9.0 * (1.0 - s*s) * (1.0 - 3.0*s) / 32.0
	dSdR[6][0] = 9.0 * sp * (-9.0*r*r - 2.0*r + 3.0) / 32.0
	dSdR[7][0] = -9.0 * (1.0 - s*s) * (1.0 + 3.0*s) / 32.0
	dSdR[8][0] = 9.0 * sm * (-9.0*r*r - 2.0*r + 3.0) / 32.0
	dSdR[9][0] = 9.0 * (1.0 - s*s) * (1.0 + 3.0*s) / 32.0
	dSdR[10][0] = 9.0 * sp * (9.0*r*r - 2.0*r - 3.0) / 32.0
	dSdR[11][0] = -9.0 * (1.0 - s*s) * (1.0 - 3.0*s) / 32.0

	dSdR[0][1] = rm * (9.0*(2.0*s-3.0*s*s-r*r) + 10.0) / 32.0
	dSdR[1][1] = rp * (9.0*(2.0*s-3.0*s*s-r*r) + 10.0) / 32.0
	dSdR[2][1] = rp * (9.0*(2.0*s+3.0*s*s+r*r) - 10.0) / 32.0
	dSdR[3][1] = rm * (9.0*(2.0*s+3.0*s*s+r*r) - 10.0) / 32.0
	dSdR[4][1] = -9.0 * (1.0 - r*r) * (1.0 - 3.0*r) / 32.0
	dSdR[5][1] = 9.0 * rp * (9.0*s*s - 2.0*s - 3.0) / 32.0
	dSdR[6][1] = 9.0 * (1.0 - r*r) * (1.0 + 3.0*r) / 32.0
	dSdR[7][1] = 9.0 * rm * (-9.0*s*s - 2.0*s + 3.0) / 32.0
	dSdR[8][1] = -9.0 * (1.0 - r*r) * (1.0 + 3.0*r) / 32.0
	dSdR[9][1] = 9.0 * rp * (-9.0*s*s - 2.0*s + 3.0) / 32.0
	dSdR[10][1] = 9.0 * (1.0 - r*r) * (1.0 - 3.0*r) / 32.0
	dSdR[11][1] = 9.0 * rm * (9.0*s*s - 2.0*s - 3.0) / 32.0
}

// FuncQua16 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua16
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua16(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	    3      10       6        2
	      @-----@-------@------@
	      |               (1,1)|
	      |       s ^          |
	    7 @   15@   |    @14   @ 9
	      |         |          |
	      |         +----> r   |
	      |       (0,0)        |
	   11 @   12@       @13    @ 5
	      |                    |
	      |(-1,-1)             |
	      @-----@-------@------@
	    0       4       8        1
	*/
	r, s := R[0], R[1]
	sr, ss := make([]float64, 4), make([]float64, 4)
	var dr, ds [][]float64
	if derivs {
		dr, ds = make([][]float64, 4), make([][]float64, 4)
		for i := 0; i < 4; i++ {
			dr[i], ds[i] = make([]float64, 1), make([]float64, 1)
		}
	}

	FuncLin4(sr, dr, []float64{r}, derivs)
	FuncLin4(ss, ds, []float64{s}, derivs)

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

	dSdR[0][0] = dr[0][0] * ss[0]
	dSdR[1][0] = dr[1][0] * ss[0]
	dSdR[2][0] = dr[1][0] * ss[1]
	dSdR[3][0] = dr[0][0] * ss[1]
	dSdR[4][0] = dr[2][0] * ss[0]
	dSdR[5][0] = dr[1][0] * ss[2]
	dSdR[6][0] = dr[3][0] * ss[1]
	dSdR[7][0] = dr[0][0] * ss[3]
	dSdR[8][0] = dr[3][0] * ss[0]
	dSdR[9][0] = dr[1][0] * ss[3]
	dSdR[10][0] = dr[2][0] * ss[1]
	dSdR[11][0] = dr[0][0] * ss[2]
	dSdR[12][0] = dr[2][0] * ss[2]
	dSdR[13][0] = dr[3][0] * ss[2]
	dSdR[14][0] = dr[3][0] * ss[3]
	dSdR[15][0] = dr[2][0] * ss[3]

	dSdR[0][1] = sr[0] * ds[0][0]
	dSdR[1][1] = sr[1] * ds[0][0]
	dSdR[2][1] = sr[1] * ds[1][0]
	dSdR[3][1] = sr[0] * ds[1][0]
	dSdR[4][1] = sr[2] * ds[0][0]
	dSdR[5][1] = sr[1] * ds[2][0]
	dSdR[6][1] = sr[3] * ds[1][0]
	dSdR[7][1] = sr[0] * ds[3][0]
	dSdR[8][1] = sr[3] * ds[0][0]
	dSdR[9][1] = sr[1] * ds[3][0]
	dSdR[10][1] = sr[2] * ds[1][0]
	dSdR[11][1] = sr[0] * ds[2][0]
	dSdR[12][1] = sr[2] * ds[2][0]
	dSdR[13][1] = sr[3] * ds[2][0]
	dSdR[14][1] = sr[3] * ds[3][0]
	dSdR[15][1] = sr[2] * ds[3][0]
}

// FuncQua17 calculates the shape functions (S) and derivatives of shape functions (dSdR) of qua17
// (serendipity) elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func FuncQua17(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
		    3      14    10     6     2
		      @-----@-----@-----@-----@
		      |                  (1,1)|
		      |                       |
			7 @                       @ 13
		      |         s ^           |
		      |           |           |
		   11 @           |16         @ 9
		      |           @----> r    |
		      |         (0,0)         |
		   15 @                       @ 5
		      |                       |
		      |(-1,-1)                |
		      @-----@-----@-----@-----@
		    0       4     8    12       1
	*/

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

	dSdR[0][0] = b * sm * (16.0*rrr - 12.0*rr - 6.0*r*s - 8.0*r + 4.0*sss - s + 4.0)
	dSdR[1][0] = b * sm * (16.0*rrr + 12.0*rr - 6.0*r*s - 8.0*r - 4.0*sss + s - 4.0)
	dSdR[2][0] = b * sp * (16.0*rrr + 12.0*rr + 6.0*r*s - 8.0*r + 4.0*sss - s - 4.0)
	dSdR[3][0] = b * sp * (16.0*rrr - 12.0*rr + 6.0*r*s - 8.0*r - 4.0*sss + s + 4.0)
	dSdR[4][0] = -a * (1.0 - 4.0*r - 3.0*rr + 8.0*rrr) * sm
	dSdR[5][0] = a * s * sm * sp * (-1.0 + 2.0*s)
	dSdR[6][0] = -a * (-1.0 - 4.0*r + 3.0*rr + 8.0*rrr) * sp
	dSdR[7][0] = -a * s * sm * sp * (1.0 + 2.0*s)
	dSdR[8][0] = r * sm * (8.0*rr + s - 4.0)
	dSdR[9][0] = 0.5 * sm * sp * (2.0*r - 4.0*ss + 1.0)
	dSdR[10][0] = r * sp * (8.0*rr - s - 4.0)
	dSdR[11][0] = 0.5 * sm * sp * (2.0*r - 1.0 + 4.0*ss)
	dSdR[12][0] = a * (1.0 + 4.0*r - 3.0*rr - 8.0*rrr) * sm
	dSdR[13][0] = a * s * sm * sp * (1.0 + 2.0*s)
	dSdR[14][0] = -a * (1.0 - 4.0*r - 3.0*rr + 8.0*rrr) * sp
	dSdR[15][0] = a * s * sm * sp * (1.0 - 2.0*s)
	dSdR[16][0] = -2.0 * r * sm * sp

	dSdR[0][1] = b * rm * (16.0*sss - 12.0*ss - 6.0*r*s - 8.0*s + 4.0*rrr - r + 4.0)
	dSdR[1][1] = -b * rp * (-16.0*sss + 12.0*ss - 6.0*r*s + 8.0*s + 4.0*rrr - r - 4.0)
	dSdR[2][1] = b * rp * (16.0*sss + 12.0*ss + 6.0*r*s - 8.0*s + 4.0*rrr - r - 4.0)
	dSdR[3][1] = b * r1 * (-16.0*sss - 12.0*ss + 6.0*r*s + 8.0*s + 4.0*rrr - r + 4.0)
	dSdR[4][1] = a * r * r1 * rp * (2.0*r - 1.0)
	dSdR[5][1] = -a * (1.0 - 4.0*s - 3.0*ss + 8.0*sss) * rp
	dSdR[6][1] = -a * r * r1 * rp * (1.0 + 2.0*r)
	dSdR[7][1] = a * (-1.0 - 4.0*s + 3.0*ss + 8.0*sss) * r1
	dSdR[8][1] = -0.5 * r1 * rp * (2.0*s - 1.0 + 4.0*rr)
	dSdR[9][1] = -s * rp * (-8.0*ss + r + 4.0)
	dSdR[10][1] = 0.5 * r1 * rp * (-2.0*s + 4.0*rr - 1.0)
	dSdR[11][1] = -s * r1 * (8.0*ss + r - 4.0)
	dSdR[12][1] = a * r * r1 * rp * (1.0 + 2.0*r)
	dSdR[13][1] = -a * (-1.0 - 4.0*s + 3.0*ss + 8.0*sss) * rp
	dSdR[14][1] = -a * r * r1 * rp * (2.0*r - 1.0)
	dSdR[15][1] = a * (1.0 - 4.0*s - 3.0*ss + 8.0*sss) * r1
	dSdR[16][1] = 2.0 * s * r1 * rp
}
