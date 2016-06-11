// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

func init() {
	Functions["lin2"] = Lin2
	Functions["lin3"] = Lin3
	Functions["lin4"] = Lin4
	Functions["lin5"] = Lin5
}

// Lin2 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin2
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func Lin2(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   -1     0    +1
	    0-----------1-->r
	*/
	r := R[0]
	S[0] = 0.5 * (1.0 - r)
	S[1] = 0.5 * (1.0 + r)

	if !derivs {
		return
	}

	dSdR[0][0] = -0.5
	dSdR[1][0] = 0.5
}

// Lin3 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin3
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func Lin3(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   -1     0    +1
	    0-----2-----1-->r
	*/
	r := R[0]
	S[0] = 0.5 * (r*r - r)
	S[1] = 0.5 * (r*r + r)
	S[2] = 1.0 - r*r

	if !derivs {
		return
	}

	dSdR[0][0] = r - 0.5
	dSdR[1][0] = r + 0.5
	dSdR[2][0] = -2.0 * r
}

// Lin4 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin4
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func Lin4(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	   -1                  +1
	   @------@-----@------@  --> r
	   0      2     3      1
	*/
	r := R[0]
	S[0] = (-9.0*r*r*r + 9.0*r*r + r - 1.0) / 16.0
	S[1] = (9.0*r*r*r + 9.0*r*r - r - 1.0) / 16.0
	S[2] = (27.0*r*r*r - 9.0*r*r - 27.0*r + 9.0) / 16.0
	S[3] = (-27.0*r*r*r - 9.0*r*r + 27.0*r + 9.0) / 16.0

	if !derivs {
		return
	}

	dSdR[0][0] = 1.0 / 16.0 * (-27.0*r*r + 18.0*r + 1.0)
	dSdR[1][0] = 1.0 / 16.0 * (27.0*r*r + 18.0*r - 1.0)
	dSdR[2][0] = 1.0 / 16.0 * (81.0*r*r - 18.0*r - 27.0)
	dSdR[3][0] = 1.0 / 16.0 * (-81.0*r*r - 18.0*r + 27.0)
}

// Lin5 calculates the shape functions (S) and derivatives of shape functions (dSdR) of lin5
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
func Lin5(S []float64, dSdR [][]float64, R []float64, derivs bool) {
	/*
	    @-----@-----@-----@-----@-> r
	    0     3     2     4     1
	    |           |           |
	   r=-1  -1/2   r=0  1/2   r=+1
	*/
	r := R[0]
	S[0] = (r - 1.0) * (1.0 - 2.0*r) * r * (-1.0 - 2.0*r) / 6.0
	S[1] = (1.0 - 2.0*r) * r * (-1.0 - 2.0*r) * (1.0 + r) / 6.0
	S[2] = (1.0 - r) * (1.0 - 2.0*r) * (-1.0 - 2.0*r) * (-1.0 - r)
	S[3] = 4.0 * (1.0 - r) * (1.0 - 2.0*r) * r * (-1.0 - r) / 3.0
	S[4] = 4.0 * (1.0 - r) * r * (-1.0 - 2.0*r) * (-1.0 - r) / 3.0

	if !derivs {
		return
	}

	dSdR[0][0] = -((1.0-2.0*r)*(r-1.0)*r)/3.0 - ((-2.0*r-1.0)*(r-1.0)*r)/3.0 + ((-2.0*r-1.0)*(1.0-2.0*r)*r)/6.0 + ((-2.0*r-1.0)*(1.0-2.0*r)*(r-1.0))/6.0
	dSdR[1][0] = -((1.0-2.0*r)*r*(r+1.0))/3.0 - ((-2.0*r-1.0)*r*(r+1.0))/3.0 + ((-2.0*r-1.0)*(1.0-2.0*r)*(r+1.0))/6.0 + ((-2.0*r-1.0)*(1.0-2.0*r)*r)/6.0
	dSdR[2][0] = -2.0*(1.0-2.0*r)*(-r-1.0)*(1.0-r) - 2.0*(-2.0*r-1.0)*(-r-1.0)*(1.0-r) - (-2.0*r-1.0)*(1.0-2.0*r)*(1.0-r) - (-2.0*r-1.0)*(1.0-2.0*r)*(-r-1.0)
	dSdR[3][0] = -(8.0*(-r-1.0)*(1.0-r)*r)/3.0 - (4.0*(1.0-2.0*r)*(1.0-r)*r)/3.0 - (4.0*(1.0-2.0*r)*(-r-1.0)*r)/3.0 + (4.0*(1.0-2.0*r)*(-r-1.0)*(1.0-r))/3.0
	dSdR[4][0] = -(8.0*(-r-1.0)*(1.0-r)*r)/3.0 - (4.0*(-2.0*r-1.0)*(1.0-r)*r)/3.0 - (4.0*(-2.0*r-1.0)*(-r-1.0)*r)/3.0 + (4.0*(-2.0*r-1.0)*(-r-1.0)*(1.0-r))/3.0
}
