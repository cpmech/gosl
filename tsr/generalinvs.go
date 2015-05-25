// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
)

// ShiftedEigenvs computes the positive and shifted eigenvalues of a second order tensor
// Note: λbar are the eigenvalues of the tensor
//       λc is a positive value
func ShiftedEigenvs(λ, λbar []float64, λc, tol float64) (err error) {
	for i := 0; i < 3; i++ {
		λ[i] = -λbar[i] + λc
		if λ[i] < tol {
			err = chk.Err(_geninvs_err1, λ, λbar, λc, tol)
			return
		}
	}
	return
}

// GenInvs returns the SMP invariants
//  Note: λ are the eigenvalues (shifted or not)
func GenInvs(λ, n []float64, a float64) (p, q float64, err error) {
	p = a * (λ[0]*n[0]*n[0] + λ[1]*n[1]*n[1] + λ[2]*n[2]*n[2])
	d := λ[0]*λ[0]*n[0]*n[0] + λ[1]*λ[1]*n[1]*n[1] + λ[2]*λ[2]*n[2]*n[2] - p*p
	if d < 0 {
		if math.Abs(d) > SMPINVSTOL {
			err = chk.Err(_geninvs_err2, "GenInvs", d, SMPINVSTOL, λ)
			return
		}
		d = 0
	}
	q = math.Sqrt(d)
	return
}

// GenTvec computes the t vector (stress vector on SMP via Cauchy's rule: t = λ dot n)
//  Note: λ are the eigenvalues (shifted or not)
func GenTvec(t, λ, n []float64) {
	for i := 0; i < 3; i++ {
		t[i] = λ[i] * n[i]
	}
}

// GenTvecDeriv1 computes the first derivative dt/dλ
//  Note: λ are the eigenvalues (shifted or not)
func GenTvecDeriv1(dtdλ [][]float64, λ, n []float64, dndλ [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dtdλ[i][j] = λ[i] * dndλ[i][j]
			if i == j {
				dtdλ[i][j] += n[i]
			}
		}
	}
}

// GenTvecDeriv2 computes the second derivative d²t/dλdλ
//  Note: λ are the eigenvalues (shifted or not)
func GenTvecDeriv2(i, j, k int, λ []float64, dndλ [][]float64, d2ndλdλ_ijk float64) (res float64) {
	res = λ[i] * d2ndλdλ_ijk
	if i == j {
		res += dndλ[i][k]
	}
	if i == k {
		res += dndλ[i][j]
	}
	return
}

// GenInvsDeriv1 computes the first order derivatives of p and q w.r.t λ (shifted eigenvalues)
//  Note: λ are the eigenvalues (shifted or not)
func GenInvsDeriv1(dpdλ, dqdλ []float64, λ, n []float64, dndλ [][]float64, a float64) (p, q float64, err error) {
	// invariants
	p = a * (λ[0]*n[0]*n[0] + λ[1]*n[1]*n[1] + λ[2]*n[2]*n[2])
	d := λ[0]*λ[0]*n[0]*n[0] + λ[1]*λ[1]*n[1]*n[1] + λ[2]*λ[2]*n[2]*n[2] - p*p
	if d < 0 {
		if math.Abs(d) > SMPINVSTOL {
			err = chk.Err(_geninvs_err2, "GenInvsDeriv1", d, SMPINVSTOL, λ)
			return
		}
		d = 0
	}
	q = math.Sqrt(d)
	// derivatives
	var t_k, dtdλ_ki float64
	for i := 0; i < 3; i++ {
		dpdλ[i] = a * n[i] * n[i]
		dqdλ[i] = 0
		for k := 0; k < 3; k++ {
			t_k = λ[k] * n[k]
			dtdλ_ki = λ[k] * dndλ[k][i]
			if k == i {
				dtdλ_ki += n[k]
			}
			dpdλ[i] += 2.0 * a * t_k * dndλ[k][i]
			dqdλ[i] += t_k * dtdλ_ki
		}
		if q > 0 {
			dqdλ[i] = (dqdλ[i] - p*dpdλ[i]) / q
		}
	}
	return
}

// GenInvsDeriv2 computes the second order derivatives of p and q w.r.t λ (shifted eigenvalues)
//  Note: λ are the eigenvalues (shifted or not)
func GenInvsDeriv2(d2pdλdλ, d2qdλdλ [][]float64, λ, n, dpdλ, dqdλ []float64, p, q float64, dndλ [][]float64, d2ndλdλ [][][]float64, a float64) {
	var t_k, dtdλ_ki, dtdλ_kj, d2tdλdλ_kij float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			d2pdλdλ[i][j] = 2.0 * a * n[i] * dndλ[i][j]
			d2qdλdλ[i][j] = -dpdλ[i]*dpdλ[j] - dqdλ[i]*dqdλ[j]
			for k := 0; k < 3; k++ {
				t_k = λ[k] * n[k]
				dtdλ_ki = λ[k] * dndλ[k][i]
				dtdλ_kj = λ[k] * dndλ[k][j]
				d2tdλdλ_kij = λ[k] * d2ndλdλ[k][i][j]
				if k == i {
					dtdλ_ki += n[k]
					d2tdλdλ_kij += dndλ[k][j]
				}
				if k == j {
					dtdλ_kj += n[k]
					d2tdλdλ_kij += dndλ[k][i]
				}
				d2pdλdλ[i][j] += 2.0*a*dndλ[k][i]*dtdλ_kj + 2.0*a*t_k*d2ndλdλ[k][i][j]
				d2qdλdλ[i][j] += dtdλ_ki*dtdλ_kj + t_k*d2tdλdλ_kij
			}
			if q > 0 {
				d2qdλdλ[i][j] = (d2qdλdλ[i][j] - p*d2pdλdλ[i][j]) / q
			}
		}
	}
}

// SMPinvs computes the SMP invariants, after the internal computation of the SMP unit director
//  Note: internal variables are created => not efficient
func SMPinvs(λ []float64, a, b, β, ϵ float64) (p, q float64, err error) {
	W := make([]float64, 3)                   // workspace
	m := SmpDirector(W, λ, a, b, β, ϵ)        // W := N
	W[0], W[1], W[2] = W[0]/m, W[1]/m, W[2]/m // W := n
	p, q, err = GenInvs(λ, W, a)
	return
}

// SMPderivs1 computes the 1st order derivatives of SMP invariants
//  Note: internal variables are created => not efficient
func SMPderivs1(dpdλ, dqdλ, λ []float64, a, b, β, ϵ float64) (p, q float64, err error) {
	dndλ := la.MatAlloc(3, 3)
	dNdλ := make([]float64, 3)
	Frmp := make([]float64, 3)
	Grmp := make([]float64, 3)
	W := make([]float64, 3)                                   // workspace
	m := SmpDerivs1(dndλ, dNdλ, W, Frmp, Grmp, λ, a, b, β, ϵ) // W := N
	W[0], W[1], W[2] = W[0]/m, W[1]/m, W[2]/m                 // W := n
	p, q, err = GenInvsDeriv1(dpdλ, dqdλ, λ, W, dndλ, a)
	return
}

// error messages
var (
	_geninvs_err1 = "generalinvs.go: ShiftedEigenvs: shifted eigenvalues are negative:\n  λ(out)=%v\n  λbar(in)=%v\n  (λc=%v, tol=%v)"
	_geninvs_err2 = "generalinvs.go: %s: difference==%g (>%g) is negative and cannot be used to computed q=sqrt(d). λ=%v"
)
