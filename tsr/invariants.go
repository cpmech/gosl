// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
)

///////////////////////////////// Mandel's representation functions //////////////////////////////////

// M_p returns the hydrostatic pressure == negative of the mean pressure == - tr(σ) / 3
func M_p(σ []float64) float64 {
	return -(σ[0] + σ[1] + σ[2]) / 3.0
}

// M_q returns von Mises' equivalent stress
func M_q(σ []float64) float64 {
	var m float64
	if len(σ) == 6 { // 3D
		m = σ[4]*σ[4] + σ[5]*σ[5]
	}
	return math.Sqrt(((σ[0]-σ[1])*(σ[0]-σ[1]) + (σ[1]-σ[2])*(σ[1]-σ[2]) + (σ[2]-σ[0])*(σ[2]-σ[0]) + 3.0*(σ[3]*σ[3]+m)) / 2.0)
}

// M_εv returns the volumetric strain
func M_εv(ε []float64) float64 {
	return ε[0] + ε[1] + ε[2]
}

// M_εd returns the deviatoric strain
func M_εd(ε []float64) float64 {
	var m float64
	if len(ε) == 6 { // 3D
		m = ε[4]*ε[4] + ε[5]*ε[5]
	}
	return math.Sqrt(((ε[0]-ε[1])*(ε[0]-ε[1])+(ε[1]-ε[2])*(ε[1]-ε[2])+(ε[2]-ε[0])*(ε[2]-ε[0])+3.0*(ε[3]*ε[3]+m))*2.0) / 3.0
}

// L_strains compute strain invariants given principal values
//  ε -- principal values [3]
func L_strains(ε []float64) (εv, εd float64) {
	εv = ε[0] + ε[1] + ε[2]
	εd = math.Sqrt(((ε[0]-ε[1])*(ε[0]-ε[1])+(ε[1]-ε[2])*(ε[1]-ε[2])+(ε[2]-ε[0])*(ε[2]-ε[0]))*2.0) / 3.0
	return
}

// M_devσ returns the deviator of σ (s := dev(σ)), the norm of the deviator (sno) and the p, q invariants
func M_devσ(s, σ []float64) (sno, p, q float64) {
	p = -(σ[0] + σ[1] + σ[2]) / 3.0
	for i := 0; i < len(σ); i++ {
		s[i] = σ[i] + p*Im[i]
		sno += s[i] * s[i]
	}
	sno = math.Sqrt(sno)
	q = SQ3by2 * sno
	return
}

// M_devε returns the deviator of ε (e := dev(ε)), the norm of the deviator (eno) and the εv, εd invariants
func M_devε(e, ε []float64) (eno, εv, εd float64) {
	εv = ε[0] + ε[1] + ε[2]
	for i := 0; i < len(ε); i++ {
		e[i] = ε[i] - εv*Im[i]/3.0
		eno += e[i] * e[i]
	}
	eno = math.Sqrt(eno)
	εd = SQ2by3 * eno
	return
}

// M_w returns the Lode invariant -1 ≤ w := sin(3θ) ≤ 1
func M_w(a []float64) (w float64) {
	s := make([]float64, len(a))
	ns := 0.0 // norm of s
	for i := 0; i < len(a); i++ {
		s[i] = a[i] - (a[0]+a[1]+a[2])*Im[i]/3.0
		ns += s[i] * s[i]
	}
	ns = math.Sqrt(ns)
	w = 1.0
	if SQ3by2*ns > QMIN {
		ds := M_Det(s)
		w = -3.0 * SQ6 * ds / (ns * ns * ns)
		if w < -1.0 {
			w = -1.0
		}
		if w > 1.0 {
			w = 1.0
		}
	}
	return
}

// M_θ returns the Lode invariant -30° ≤ θ := asin(w) / 3 ≤ 30°
func M_θ(a []float64) (θdeg float64) {
	return math.Asin(M_w(a)) * 180.0 / (3.0 * math.Pi)
}

// M_pqw returns p, q and w invariants
func M_pqw(a []float64) (p, q, w float64) {
	s := make([]float64, len(a))
	ns := 0.0 // norm of s
	p = -(a[0] + a[1] + a[2]) / 3.0
	for i := 0; i < len(a); i++ {
		s[i] = a[i] + p*Im[i]
		ns += s[i] * s[i]
	}
	ns = math.Sqrt(ns)
	q = SQ3by2 * ns
	w = 1.0
	if q > QMIN {
		ds := M_Det(s)
		w = -3.0 * SQ6 * ds / (ns * ns * ns)
		if w < -1.0 {
			w = -1.0
		}
		if w > 1.0 {
			w = 1.0
		}
	}
	return
}

// M_pqθ returns p, q and θ invariants
func M_pqθ(a []float64) (p, q, θ float64) {
	s := make([]float64, len(a))
	ns := 0.0 // norm of s
	p = -(a[0] + a[1] + a[2]) / 3.0
	for i := 0; i < len(a); i++ {
		s[i] = a[i] + p*Im[i]
		ns += s[i] * s[i]
	}
	ns = math.Sqrt(ns)
	q = SQ3by2 * ns
	w := 1.0
	if q > QMIN {
		ds := M_Det(s)
		w = -3.0 * SQ6 * ds / (ns * ns * ns)
		if w < -1.0 {
			w = -1.0
		}
		if w > 1.0 {
			w = 1.0
		}
	}
	θ = math.Asin(w) * 180.0 / (3.0 * math.Pi)
	return
}

// M_pqws returns p, q, w invariants and the deviatoric stress s := dev(σ)
func M_pqws(s, a []float64) (p, q, w float64) {
	ns := 0.0 // norm of s
	p = -(a[0] + a[1] + a[2]) / 3.0
	for i := 0; i < len(a); i++ {
		s[i] = a[i] + p*Im[i]
		ns += s[i] * s[i]
	}
	ns = math.Sqrt(ns)
	q = SQ3by2 * ns
	w = 1.0
	if q > QMIN {
		ds := M_Det(s)
		w = -3.0 * SQ6 * ds / (ns * ns * ns)
		if w < -1.0 {
			w = -1.0
		}
		if w > 1.0 {
			w = 1.0
		}
	}
	return
}

// M_LodeDeriv1 computes the first derivative of w w.r.t σ
//  Note: only dwdσ is output
func M_LodeDeriv1(dwdσ, σ, s []float64, p, q, w float64) {
	nσ := len(σ)
	if q > QMIN {
		n := SQ2by3 * q
		M_Sq(dwdσ, s)                       // dwdσ := s²
		trs2 := dwdσ[0] + dwdσ[1] + dwdσ[2] // tr(s²)
		for i := 0; i < nσ; i++ {
			dwdσ[i] = -3.0 * (SQ6*(dwdσ[i]-trs2*Im[i]/3.0)/n + w*s[i]) / (n * n)
		}
	} else {
		for i := 0; i < nσ; i++ {
			dwdσ[i] = 0
		}
	}
}

// M_LodeDeriv2 computes the first and second derivatives of w w.r.t. σ
//  Note: d2wdσdσ and dwdσ output
func M_LodeDeriv2(d2wdσdσ [][]float64, dwdσ, σ, s []float64, p, q, w float64) {
	nσ := len(σ)
	if q > QMIN {
		n := SQ2by3 * q
		Ss := make([]float64, nσ)     // Ss := dev(s²)
		M_Sq(Ss, s)                   // Ss := s²
		trs2 := Ss[0] + Ss[1] + Ss[2] // tr(s²)
		for i := 0; i < nσ; i++ {
			Ss[i] = Ss[i] - trs2*Im[i]/3.0
			dwdσ[i] = -3.0 * (SQ6*Ss[i]/n + w*s[i]) / (n * n)
		}
		n2 := n * n
		n3 := n2 * n
		n4 := n2 * n2
		ds := M_Det(s)
		M_Ts(d2wdσdσ, s) // d2wdσdσ := Ts
		for i := 0; i < nσ; i++ {
			for j := 0; j < nσ; j++ {
				d2wdσdσ[i][j] = 9.0 * SQ6 * (s[i]*Ss[j]/n2 + Ss[i]*s[j]/n2 - 5.0*ds*s[i]*s[j]/n4 + ds*Psd[i][j]/n2 - d2wdσdσ[i][j]/3.0) / n3
			}
		}
	} else {
		for i := 0; i < nσ; i++ {
			for j := 0; j < nσ; j++ {
				d2wdσdσ[i][j] = 0
			}
			dwdσ[i] = 0
		}
	}
}

// M_Ts computes Ts = Psd:(ds²/ds):Psd
func M_Ts(Ts [][]float64, s []float64) {
	if len(s) == 4 {
		Ts[0][0] = (2 * s[0]) / 3
		Ts[0][1] = -(2*s[1] + 2*s[0]) / 3
		Ts[0][2] = -(2*s[2] + 2*s[0]) / 3
		Ts[0][3] = s[3] / 3
		Ts[1][0] = -(2*s[1] + 2*s[0]) / 3
		Ts[1][1] = (2 * s[1]) / 3
		Ts[1][2] = -(2*s[2] + 2*s[1]) / 3
		Ts[1][3] = s[3] / 3
		Ts[2][0] = -(2*s[2] + 2*s[0]) / 3
		Ts[2][1] = -(2*s[2] + 2*s[1]) / 3
		Ts[2][2] = (2 * s[2]) / 3
		Ts[2][3] = -(2 * s[3]) / 3
		Ts[3][0] = s[3] / 3
		Ts[3][1] = s[3] / 3
		Ts[3][2] = -(2 * s[3]) / 3
		Ts[3][3] = s[1] + s[0]
		return
	}
	Ts[0][0] = (2 * s[0]) / 3
	Ts[0][1] = -(2*s[1] + 2*s[0]) / 3
	Ts[0][2] = -(2*s[2] + 2*s[0]) / 3
	Ts[0][3] = s[3] / 3
	Ts[0][4] = -(2 * s[4]) / 3
	Ts[0][5] = s[5] / 3
	Ts[1][0] = -(2*s[1] + 2*s[0]) / 3
	Ts[1][1] = (2 * s[1]) / 3
	Ts[1][2] = -(2*s[2] + 2*s[1]) / 3
	Ts[1][3] = s[3] / 3
	Ts[1][4] = s[4] / 3
	Ts[1][5] = -(2 * s[5]) / 3
	Ts[2][0] = -(2*s[2] + 2*s[0]) / 3
	Ts[2][1] = -(2*s[2] + 2*s[1]) / 3
	Ts[2][2] = (2 * s[2]) / 3
	Ts[2][3] = -(2 * s[3]) / 3
	Ts[2][4] = s[4] / 3
	Ts[2][5] = s[5] / 3
	Ts[3][0] = s[3] / 3
	Ts[3][1] = s[3] / 3
	Ts[3][2] = -(2 * s[3]) / 3
	Ts[3][3] = s[1] + s[0]
	Ts[3][4] = s[5] / SQ2
	Ts[3][5] = s[4] / SQ2
	Ts[4][0] = -(2 * s[4]) / 3
	Ts[4][1] = s[4] / 3
	Ts[4][2] = s[4] / 3
	Ts[4][3] = s[5] / SQ2
	Ts[4][4] = s[2] + s[1]
	Ts[4][5] = s[3] / SQ2
	Ts[5][0] = s[5] / 3
	Ts[5][1] = -(2 * s[5]) / 3
	Ts[5][2] = s[5] / 3
	Ts[5][3] = s[4] / SQ2
	Ts[5][4] = s[3] / SQ2
	Ts[5][5] = s[2] + s[0]
}

// M_CharInvs computes the characteristic invariants of a 2nd order symmetric tensor
func M_CharInvs(a []float64) (I1, I2, I3 float64) {
	I1 = a[0] + a[1] + a[2]
	I2 = a[0]*a[1] + a[1]*a[2] + a[2]*a[0] - a[3]*a[3]/2.0
	I3 = a[0]*a[1]*a[2] - a[2]*a[3]*a[3]/2.0
	if len(a) > 4 {
		I2 += (-a[4]*a[4]/2.0 - a[5]*a[5]/2.0)
		I3 += (a[3]*a[4]*a[5]/SQ2 - a[0]*a[4]*a[4]/2.0 - a[1]*a[5]*a[5]/2.0)
	}
	return
}

// M_CharInvsAndDerivs computes the characteristic invariants of a
// 2nd order symmetric and their derivatives
func M_CharInvsAndDerivs(a []float64) (I1, I2, I3 float64, dI1da, dI2da, dI3da []float64) {
	I1 = a[0] + a[1] + a[2]
	I2 = a[0]*a[1] + a[1]*a[2] + a[2]*a[0] - a[3]*a[3]/2.0
	I3 = a[0]*a[1]*a[2] - a[2]*a[3]*a[3]/2.0
	n := len(a)
	dI2da = make([]float64, n)
	dI3da = make([]float64, n)
	if n > 4 {
		I2 += (-a[4]*a[4]/2.0 - a[5]*a[5]/2.0)
		I3 += (a[3]*a[4]*a[5]/SQ2 - a[0]*a[4]*a[4]/2.0 - a[1]*a[5]*a[5]/2.0)
		dI1da = []float64{1, 1, 1, 0, 0, 0}
	} else {
		dI1da = []float64{1, 1, 1, 0}
	}
	M_Sq(dI3da, a) // dI3da := a²
	for i := 0; i < n; i++ {
		dI2da[i] = I1*dI1da[i] - a[i]
		dI3da[i] += I2*dI1da[i] - I1*a[i]
	}
	return
}

// M_pq_smp computes p and q SMP invariants of 2nd order symmetric tensor (Mandel components)
//  Note: 1) σ is a 2D or 3D symmetric tensor (len(σ)==4 or 6)
//        2) this function creates a number of local arrays => not efficient
func M_pq_smp(σ []float64, a, b, β, ϵ float64) (p, q float64, err error) {
	λ := make([]float64, 3)
	err = M_EigenValsNum(λ, σ)
	if err != nil {
		return
	}
	W := make([]float64, 3)
	m := NewSmpDirector(W, λ, a, b, β, ϵ)     // W := N
	W[0], W[1], W[2] = W[0]/m, W[1]/m, W[2]/m // W := n
	p, q, err = GenInvs(λ, W, a)
	return
}
