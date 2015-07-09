// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/fun"
)

// SmpCalcμ computes μ=q/p to satisfy Mohr-Coulomb criterion @ compression
func SmpCalcμ(φ, a, b, β, ϵ float64) (μ float64) {
	sφ := math.Sin(φ * math.Pi / 180.0)
	R := (1.0 + sφ) / (1.0 - sφ)
	L := []float64{a * R, a, a}
	N := make([]float64, 3)
	n := make([]float64, 3)
	m := SmpDirector(N, L, a, b, β, ϵ)
	SmpUnitDirector(n, m, N)
	p, q, _ := GenInvs(L, n, a)
	return q / p
}

/// SMP director ////////////////////////////////////////////////////////////////////////////////////

// SmpDirector computes the director (normal vector) of the spatially mobilised plane
//  Notes:
//    1) the norm of N is returned => m := norm(N)
//    2) if !SMPUSESRAMP, β==eps and must be a small quantity
func SmpDirector(N, L []float64, a, b, β, ϵ float64) (m float64) {
	if SMPUSESRAMP {
		N[0] = a / math.Pow(ϵ+fun.Sramp(a*L[0], β), b)
		N[1] = a / math.Pow(ϵ+fun.Sramp(a*L[1], β), b)
		N[2] = a / math.Pow(ϵ+fun.Sramp(a*L[2], β), b)
	} else {
		eps := β
		N[0] = a / math.Pow(ϵ+fun.Sabs(a*L[0], eps), b)
		N[1] = a / math.Pow(ϵ+fun.Sabs(a*L[1], eps), b)
		N[2] = a / math.Pow(ϵ+fun.Sabs(a*L[2], eps), b)
	}
	m = math.Sqrt(N[0]*N[0] + N[1]*N[1] + N[2]*N[2])
	return
}

// SmpDirectorDeriv1 computes the first order derivative of the SMP director
//  Notes: Only non-zero components are returned; i.e. dNdL[i] := dNdL[i][i]
func SmpDirectorDeriv1(dNdL []float64, L []float64, a, b, β, ϵ float64) {
	if SMPUSESRAMP {
		dNdL[0] = -b * fun.SrampD1(a*L[0], β) * math.Pow(ϵ+fun.Sramp(a*L[0], β), -b-1.0)
		dNdL[1] = -b * fun.SrampD1(a*L[1], β) * math.Pow(ϵ+fun.Sramp(a*L[1], β), -b-1.0)
		dNdL[2] = -b * fun.SrampD1(a*L[2], β) * math.Pow(ϵ+fun.Sramp(a*L[2], β), -b-1.0)
	} else {
		eps := β
		dNdL[0] = -b * fun.SabsD1(a*L[0], eps) * math.Pow(ϵ+fun.Sabs(a*L[0], eps), -b-1.0)
		dNdL[1] = -b * fun.SabsD1(a*L[1], eps) * math.Pow(ϵ+fun.Sabs(a*L[1], eps), -b-1.0)
		dNdL[2] = -b * fun.SabsD1(a*L[2], eps) * math.Pow(ϵ+fun.Sabs(a*L[2], eps), -b-1.0)
	}
}

// SmpDirectorDeriv2 computes the second order derivative of the SMP director
//  Notes: Only the non-zero components are returned; i.e.: d²NdL2[i] := d²N[i]/dL[i]dL[i]
func SmpDirectorDeriv2(d2NdL2 []float64, L []float64, a, b, β, ϵ float64) {
	var F_i, G_i, H_i float64
	for i := 0; i < 3; i++ {
		if SMPUSESRAMP {
			F_i = fun.Sramp(a*L[i], β)
			G_i = fun.SrampD1(a*L[i], β)
			H_i = fun.SrampD2(a*L[i], β)
		} else {
			eps := β
			F_i = fun.Sabs(a*L[i], eps)
			G_i = fun.SabsD1(a*L[i], eps)
			H_i = fun.SabsD2(a*L[i], eps)
		}
		d2NdL2[i] = a * b * ((b+1.0)*G_i*G_i - (ϵ+F_i)*H_i) * math.Pow(ϵ+F_i, -b-2.0)
	}
}

/// norm of SMP director /////////////////////////////////////////////////////////////////////////////

// SmpNormDirectorDeriv1 computes the first derivative of the norm of the SMP director
//  Note: m, N and dNdL are input
func SmpNormDirectorDeriv1(dmdL []float64, m float64, N, dNdL []float64) {
	dmdL[0] = N[0] * dNdL[0] / m
	dmdL[1] = N[1] * dNdL[1] / m
	dmdL[2] = N[2] * dNdL[2] / m
}

// SmpNormDirectorDeriv2 computes the second order derivative of the norm of the SMP director
//  Note: m, N, dNdL, d2NdL2 and dmdL are input
func SmpNormDirectorDeriv2(d2mdLdL [][]float64, L []float64, a, b, β, ϵ, m float64, N, dNdL, d2NdL2, dmdL []float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			d2mdLdL[i][j] = -N[i] * dNdL[i] * dmdL[j] / (m * m)
			if i == j {
				d2mdLdL[i][j] += (N[i]*d2NdL2[i] + dNdL[i]*dNdL[i]) / m
			}
		}
	}
}

/// unit SMP director ///////////////////////////////////////////////////////////////////////////////

// SmpUnitDirector computed the unit normal of the SMP
//  Note: m and N are input
func SmpUnitDirector(n []float64, m float64, N []float64) {
	n[0] = N[0] / m
	n[1] = N[1] / m
	n[2] = N[2] / m
}

// SmpUnitDirectorDeriv1 computes the first derivative of the SMP unit normal
//  Note: m, N, dNdL and dmdL are input
func SmpUnitDirectorDeriv1(dndL [][]float64, m float64, N, dNdL, dmdL []float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dndL[i][j] = -N[i] * dmdL[j] / (m * m)
			if i == j {
				dndL[i][j] += dNdL[i] / m
			}
		}
	}
}

// SmpUnitDirectorDeriv2 computes the second order derivative of the unit director of the SMP
// d²n[i]/dL[j]dL[k]
//  Note: m, N, dNdL, d2NdL2, dmdL, n, d2mdLdL and dndL are input
func SmpUnitDirectorDeriv2(d2ndLdL [][][]float64, m float64, N, dNdL, d2NdL2, dmdL, n []float64, d2mdLdL, dndL [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				if i == j && j == k {
					d2ndLdL[i][j][k] = d2NdL2[i] / m
				}
				d2ndLdL[i][j][k] -= (n[i]*d2mdLdL[j][k] + dmdL[j]*dndL[i][k] + dndL[i][j]*dmdL[k]) / m
			}
		}
	}
}

/// auxiliary functions /////////////////////////////////////////////////////////////////////////////

// SmpDerivs1 computes the first derivative and other variables
//  Note: m, dNdL, N, F and G are output
func SmpDerivs1(dndL [][]float64, dNdL, N, F, G []float64, L []float64, a, b, β, ϵ float64) (m float64) {
	if SMPUSESRAMP {
		F[0] = fun.Sramp(a*L[0], β)
		F[1] = fun.Sramp(a*L[1], β)
		F[2] = fun.Sramp(a*L[2], β)
		G[0] = fun.SrampD1(a*L[0], β)
		G[1] = fun.SrampD1(a*L[1], β)
		G[2] = fun.SrampD1(a*L[2], β)
	} else {
		c := β
		F[0] = fun.Sabs(a*L[0], c)
		F[1] = fun.Sabs(a*L[1], c)
		F[2] = fun.Sabs(a*L[2], c)
		G[0] = fun.SabsD1(a*L[0], c)
		G[1] = fun.SabsD1(a*L[1], c)
		G[2] = fun.SabsD1(a*L[2], c)
	}
	N[0] = a / math.Pow(ϵ+F[0], b)
	N[1] = a / math.Pow(ϵ+F[1], b)
	N[2] = a / math.Pow(ϵ+F[2], b)
	m = math.Sqrt(N[0]*N[0] + N[1]*N[1] + N[2]*N[2])
	dNdL[0] = -b * G[0] * math.Pow(ϵ+F[0], -b-1.0)
	dNdL[1] = -b * G[1] * math.Pow(ϵ+F[1], -b-1.0)
	dNdL[2] = -b * G[2] * math.Pow(ϵ+F[2], -b-1.0)
	var dmdL_j float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dmdL_j = N[j] * dNdL[j] / m
			dndL[i][j] = -N[i] * dmdL_j / (m * m)
			if i == j {
				dndL[i][j] += dNdL[i] / m
			}
		}
	}
	return
}

// SmpDerivs2 computes the second order derivative
//  Note: m, N, F, G, dNdL and dndL are input
func SmpDerivs2(d2ndLdL [][][]float64, L []float64, a, b, β, ϵ, m float64, N, F, G, dNdL []float64, dndL [][]float64) {
	var H []float64
	if SMPUSESRAMP {
		H = []float64{
			fun.SrampD2(a*L[0], β),
			fun.SrampD2(a*L[1], β),
			fun.SrampD2(a*L[2], β),
		}
	} else {
		c := β
		H = []float64{
			fun.SabsD2(a*L[0], c),
			fun.SabsD2(a*L[1], c),
			fun.SabsD2(a*L[2], c),
		}
	}
	var dmdL_k, dmdL_j, d2mdLdL_jk, d2NdL2_jj, d2NdLdL_ijk float64
	for k := 0; k < 3; k++ {
		dmdL_k = N[k] * dNdL[k] / m
		for j := 0; j < 3; j++ {
			dmdL_j = N[j] * dNdL[j] / m
			d2mdLdL_jk = -N[j] * dNdL[j] * dmdL_k / (m * m)
			if j == k {
				d2NdL2_jj = a * b * ((b+1.0)*G[j]*G[j] - (ϵ+F[j])*H[j]) * math.Pow(ϵ+F[j], -b-2.0)
				d2mdLdL_jk += (N[j]*d2NdL2_jj + dNdL[j]*dNdL[j]) / m
			}
			for i := 0; i < 3; i++ {
				d2NdLdL_ijk = 0
				if i == j && j == k {
					d2NdLdL_ijk = a * b * ((b+1.0)*G[i]*G[i] - (ϵ+F[i])*H[i]) * math.Pow(ϵ+F[i], -b-2.0)
				}
				d2ndLdL[i][j][k] = (d2NdLdL_ijk - (N[i]/m)*d2mdLdL_jk - dmdL_j*dndL[i][k] - dndL[i][j]*dmdL_k) / m
			}
		}
	}
}
