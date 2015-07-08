// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// GenInvs returns the SMP invariants
//  Note: L are the eigenvalues (shifted or not)
func GenInvs(L, n []float64, a float64) (p, q float64, err error) {
	p = a * (L[0]*n[0]*n[0] + L[1]*n[1]*n[1] + L[2]*n[2]*n[2])
	d := L[0]*L[0]*n[0]*n[0] + L[1]*L[1]*n[1]*n[1] + L[2]*L[2]*n[2]*n[2] - p*p
	if d < 0 {
		if math.Abs(d) > SMPINVSTOL {
			err = chk.Err("difference==%g (>%g) is negative and cannot be used to computed q=sqrt(d). L=%v", d, SMPINVSTOL, L)
			return
		}
		d = 0
	}
	q = math.Sqrt(d)
	return
}

// GenTvec computes the t vector (stress vector on SMP via Cauchy's rule: t = L dot n)
//  Note: L are the eigenvalues (shifted or not)
func GenTvec(t, L, n []float64) {
	for i := 0; i < 3; i++ {
		t[i] = L[i] * n[i]
	}
}

// GenTvecDeriv1 computes the first derivative dt/dL
//  Note: L are the eigenvalues (shifted or not)
func GenTvecDeriv1(dtdL [][]float64, L, n []float64, dndL [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dtdL[i][j] = L[i] * dndL[i][j]
			if i == j {
				dtdL[i][j] += n[i]
			}
		}
	}
}

// GenTvecDeriv2 computes the second derivative d²t/dLdL
//  Note: L are the eigenvalues (shifted or not)
func GenTvecDeriv2(i, j, k int, L []float64, dndL [][]float64, d2ndLdL_ijk float64) (res float64) {
	res = L[i] * d2ndLdL_ijk
	if i == j {
		res += dndL[i][k]
	}
	if i == k {
		res += dndL[i][j]
	}
	return
}

// GenInvsDeriv1 computes the first order derivatives of p and q w.r.t L (shifted eigenvalues)
//  Note: L are the eigenvalues (shifted or not)
func GenInvsDeriv1(dpdL, dqdL []float64, L, n []float64, dndL [][]float64, a float64) (p, q float64, err error) {
	// invariants
	p = a * (L[0]*n[0]*n[0] + L[1]*n[1]*n[1] + L[2]*n[2]*n[2])
	d := L[0]*L[0]*n[0]*n[0] + L[1]*L[1]*n[1]*n[1] + L[2]*L[2]*n[2]*n[2] - p*p
	if d < 0 {
		if math.Abs(d) > SMPINVSTOL {
			err = chk.Err("difference==%g (>%g) is negative and cannot be used to computed q=sqrt(d). L=%v", d, SMPINVSTOL, L)
			return
		}
		d = 0
	}
	q = math.Sqrt(d)
	// derivatives
	var t_k, dtdL_ki float64
	for i := 0; i < 3; i++ {
		dpdL[i] = a * n[i] * n[i]
		dqdL[i] = 0
		for k := 0; k < 3; k++ {
			t_k = L[k] * n[k]
			dtdL_ki = L[k] * dndL[k][i]
			if k == i {
				dtdL_ki += n[k]
			}
			dpdL[i] += 2.0 * a * t_k * dndL[k][i]
			dqdL[i] += t_k * dtdL_ki
		}
		if q > 0 {
			dqdL[i] = (dqdL[i] - p*dpdL[i]) / q
		}
	}
	return
}

// GenInvsDeriv2 computes the second order derivatives of p and q w.r.t L (shifted eigenvalues)
//  Note: L are the eigenvalues (shifted or not)
func GenInvsDeriv2(d2pdLdL, d2qdLdL [][]float64, L, n, dpdL, dqdL []float64, p, q float64, dndL [][]float64, d2ndLdL [][][]float64, a float64) {
	var t_k, dtdL_ki, dtdL_kj, d2tdLdL_kij float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			d2pdLdL[i][j] = 2.0 * a * n[i] * dndL[i][j]
			d2qdLdL[i][j] = -dpdL[i]*dpdL[j] - dqdL[i]*dqdL[j]
			for k := 0; k < 3; k++ {
				t_k = L[k] * n[k]
				dtdL_ki = L[k] * dndL[k][i]
				dtdL_kj = L[k] * dndL[k][j]
				d2tdLdL_kij = L[k] * d2ndLdL[k][i][j]
				if k == i {
					dtdL_ki += n[k]
					d2tdLdL_kij += dndL[k][j]
				}
				if k == j {
					dtdL_kj += n[k]
					d2tdLdL_kij += dndL[k][i]
				}
				d2pdLdL[i][j] += 2.0*a*dndL[k][i]*dtdL_kj + 2.0*a*t_k*d2ndLdL[k][i][j]
				d2qdLdL[i][j] += dtdL_ki*dtdL_kj + t_k*d2tdLdL_kij
			}
			if q > 0 {
				d2qdLdL[i][j] = (d2qdLdL[i][j] - p*d2pdLdL[i][j]) / q
			}
		}
	}
}

// SMPinvs computes the SMP invariants, after the internal computation of the SMP unit director
//  Note: internal variables are created => not efficient
func SMPinvs(L []float64, a, b, β, ϵ float64) (p, q float64, err error) {
	W := make([]float64, 3)                   // workspace
	m := SmpDirector(W, L, a, b, β, ϵ)        // W := N
	W[0], W[1], W[2] = W[0]/m, W[1]/m, W[2]/m // W := n
	p, q, err = GenInvs(L, W, a)
	return
}

// SMPderivs1 computes the 1st order derivatives of SMP invariants
//  Note: internal variables are created => not efficient
func SMPderivs1(dpdL, dqdL, L []float64, a, b, β, ϵ float64) (p, q float64, err error) {
	dndL := la.MatAlloc(3, 3)
	dNdL := make([]float64, 3)
	Frmp := make([]float64, 3)
	Grmp := make([]float64, 3)
	W := make([]float64, 3)                                   // workspace
	m := SmpDerivs1(dndL, dNdL, W, Frmp, Grmp, L, a, b, β, ϵ) // W := N
	W[0], W[1], W[2] = W[0]/m, W[1]/m, W[2]/m                 // W := n
	p, q, err = GenInvsDeriv1(dpdL, dqdL, L, W, dndL, a)
	return
}
