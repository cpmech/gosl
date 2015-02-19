// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/fun"
)

// NewSmpCalcμ computes μ=q/p to satisfy Mohr-Coulomb criterion @ compression
func NewSmpCalcμ(φ, a, b, β, ϵ float64) (μ float64) {
	sφ := math.Sin(φ * math.Pi / 180.0)
	R := (1.0 + sφ) / (1.0 - sφ)
	λ := []float64{a * R, a, a}
	N := make([]float64, 3)
	n := make([]float64, 3)
	m := NewSmpDirector(N, λ, a, b, β, ϵ)
	NewSmpUnitDirector(n, m, N)
	p, q, _ := GenInvs(λ, n, a)
	return q / p
}

/// SMP director ////////////////////////////////////////////////////////////////////////////////////

// NewSmpDirector computes the director (normal vector) of the spatially mobilised plane
//  Note: the norm of N is returned => m := norm(N)
func NewSmpDirector(N, λ []float64, a, b, β, ϵ float64) (m float64) {
	N[0] = a / math.Pow(ϵ+fun.Sramp(a*λ[0], β), b)
	N[1] = a / math.Pow(ϵ+fun.Sramp(a*λ[1], β), b)
	N[2] = a / math.Pow(ϵ+fun.Sramp(a*λ[2], β), b)
	m = math.Sqrt(N[0]*N[0] + N[1]*N[1] + N[2]*N[2])
	return
}

// NewSmpDirectorDeriv1 computes the first order derivative of the SMP director
//  Notes: Only non-zero components are returned; i.e. dNdλ[i] := dNdλ[i][i]
func NewSmpDirectorDeriv1(dNdλ []float64, λ []float64, a, b, β, ϵ float64) {
	dNdλ[0] = -b * fun.SrampD1(a*λ[0], β) * math.Pow(ϵ+fun.Sramp(a*λ[0], β), -b-1.0)
	dNdλ[1] = -b * fun.SrampD1(a*λ[1], β) * math.Pow(ϵ+fun.Sramp(a*λ[1], β), -b-1.0)
	dNdλ[2] = -b * fun.SrampD1(a*λ[2], β) * math.Pow(ϵ+fun.Sramp(a*λ[2], β), -b-1.0)
}

// NewSmpDirectorDeriv2 computes the second order derivative of the SMP director
//  Notes: Only the non-zero components are returned; i.e.: d²Ndλ2[i] := d²N[i]/dλ[i]dλ[i]
func NewSmpDirectorDeriv2(d2Ndλ2 []float64, λ []float64, a, b, β, ϵ float64) {
	var F_i, G_i, H_i float64
	for i := 0; i < 3; i++ {
		F_i = fun.Sramp(a*λ[i], β)
		G_i = fun.SrampD1(a*λ[i], β)
		H_i = fun.SrampD2(a*λ[i], β)
		d2Ndλ2[i] = a * b * ((b+1.0)*G_i*G_i - (ϵ+F_i)*H_i) * math.Pow(ϵ+F_i, -b-2.0)
	}
}

/// norm of SMP director /////////////////////////////////////////////////////////////////////////////

// NewSmpNormDirectorDeriv1 computes the first derivative of the norm of the SMP director
//  Note: m, N and dNdλ are input
func NewSmpNormDirectorDeriv1(dmdλ []float64, m float64, N, dNdλ []float64) {
	dmdλ[0] = N[0] * dNdλ[0] / m
	dmdλ[1] = N[1] * dNdλ[1] / m
	dmdλ[2] = N[2] * dNdλ[2] / m
}

// NewSmpNormDirectorDeriv2 computes the second order derivative of the norm of the SMP director
//  Note: m, N, dNdλ, d2Ndλ2 and dmdλ are input
func NewSmpNormDirectorDeriv2(d2mdλdλ [][]float64, λ []float64, a, b, β, ϵ, m float64, N, dNdλ, d2Ndλ2, dmdλ []float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			d2mdλdλ[i][j] = -N[i] * dNdλ[i] * dmdλ[j] / (m * m)
			if i == j {
				d2mdλdλ[i][j] += (N[i]*d2Ndλ2[i] + dNdλ[i]*dNdλ[i]) / m
			}
		}
	}
}

/// unit SMP director ///////////////////////////////////////////////////////////////////////////////

// NewSmpUnitDirector computed the unit normal of the SMP
//  Note: m and N are input
func NewSmpUnitDirector(n []float64, m float64, N []float64) {
	n[0] = N[0] / m
	n[1] = N[1] / m
	n[2] = N[2] / m
}

// NewSmpUnitDirectorDeriv1 computes the first derivative of the SMP unit normal
//  Note: m, N, dNdλ and dmdλ are input
func NewSmpUnitDirectorDeriv1(dndλ [][]float64, m float64, N, dNdλ, dmdλ []float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dndλ[i][j] = -N[i] * dmdλ[j] / (m * m)
			if i == j {
				dndλ[i][j] += dNdλ[i] / m
			}
		}
	}
}

// NewSmpUnitDirectorDeriv2 computes the second order derivative of the unit director of the SMP
// d²n[i]/dλ[j]dλ[k]
//  Note: m, N, dNdλ, d2Ndλ2, dmdλ, n, d2mdλdλ and dndλ are input
func NewSmpUnitDirectorDeriv2(d2ndλdλ [][][]float64, m float64, N, dNdλ, d2Ndλ2, dmdλ, n []float64, d2mdλdλ, dndλ [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				if i == j && j == k {
					d2ndλdλ[i][j][k] = d2Ndλ2[i] / m
				}
				d2ndλdλ[i][j][k] -= (n[i]*d2mdλdλ[j][k] + dmdλ[j]*dndλ[i][k] + dndλ[i][j]*dmdλ[k]) / m
			}
		}
	}
}

/// auxiliary functions /////////////////////////////////////////////////////////////////////////////

// NewSmpDerivs1 computes the first derivative and other variables
//  Note: m, dNdλ, N, F and G are output
func NewSmpDerivs1(dndλ [][]float64, dNdλ, N, F, G []float64, λ []float64, a, b, β, ϵ float64) (m float64) {
	F[0] = fun.Sramp(a*λ[0], β)
	F[1] = fun.Sramp(a*λ[1], β)
	F[2] = fun.Sramp(a*λ[2], β)
	G[0] = fun.SrampD1(a*λ[0], β)
	G[1] = fun.SrampD1(a*λ[1], β)
	G[2] = fun.SrampD1(a*λ[2], β)
	N[0] = a / math.Pow(ϵ+F[0], b)
	N[1] = a / math.Pow(ϵ+F[1], b)
	N[2] = a / math.Pow(ϵ+F[2], b)
	m = math.Sqrt(N[0]*N[0] + N[1]*N[1] + N[2]*N[2])
	dNdλ[0] = -b * G[0] * math.Pow(ϵ+F[0], -b-1.0)
	dNdλ[1] = -b * G[1] * math.Pow(ϵ+F[1], -b-1.0)
	dNdλ[2] = -b * G[2] * math.Pow(ϵ+F[2], -b-1.0)
	var dmdλ_j float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dmdλ_j = N[j] * dNdλ[j] / m
			dndλ[i][j] = -N[i] * dmdλ_j / (m * m)
			if i == j {
				dndλ[i][j] += dNdλ[i] / m
			}
		}
	}
	return
}

// NewSmpDerivs2 computes the second order derivative
//  Note: m, N, F, G, dNdλ and dndλ are input
func NewSmpDerivs2(d2ndλdλ [][][]float64, λ []float64, a, b, β, ϵ, m float64, N, F, G, dNdλ []float64, dndλ [][]float64) {
	H := []float64{
		fun.SrampD2(a*λ[0], β),
		fun.SrampD2(a*λ[1], β),
		fun.SrampD2(a*λ[2], β),
	}
	var dmdλ_k, dmdλ_j, d2mdλdλ_jk, d2Ndλ2_jj, d2Ndλdλ_ijk float64
	for k := 0; k < 3; k++ {
		dmdλ_k = N[k] * dNdλ[k] / m
		for j := 0; j < 3; j++ {
			dmdλ_j = N[j] * dNdλ[j] / m
			d2mdλdλ_jk = -N[j] * dNdλ[j] * dmdλ_k / (m * m)
			if j == k {
				d2Ndλ2_jj = a * b * ((b+1.0)*G[j]*G[j] - (ϵ+F[j])*H[j]) * math.Pow(ϵ+F[j], -b-2.0)
				d2mdλdλ_jk += (N[j]*d2Ndλ2_jj + dNdλ[j]*dNdλ[j]) / m
			}
			for i := 0; i < 3; i++ {
				d2Ndλdλ_ijk = 0
				if i == j && j == k {
					d2Ndλdλ_ijk = a * b * ((b+1.0)*G[i]*G[i] - (ϵ+F[i])*H[i]) * math.Pow(ϵ+F[i], -b-2.0)
				}
				d2ndλdλ[i][j][k] = (d2Ndλdλ_ijk - (N[i]/m)*d2mdλdλ_jk - dmdλ_j*dndλ[i][k] - dndλ[i][j]*dmdλ_k) / m
			}
		}
	}
}
