// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
)

// SmpCalcμ computes μ=q/p to satisfy Mohr-Coulomb criterion @ compression
func SmpCalcμ(φ, b float64) (μ float64) {
	sφ := math.Sin(φ * math.Pi / 180.0)
	R := (1.0 + sφ) / (1.0 - sφ)
	return SQ2 * (R - 1.0) * math.Pow(R, b) / (2.0*math.Pow(R, 2.0*b) + R)
}

// SmpCalcμNum computes μ=q/p to satisfy Mohr-Coulomb criterion @ compression
func SmpCalcμNum(φ, b float64) (μ float64) {
	sφ := math.Sin(φ * math.Pi / 180.0)
	R := (1.0 + sφ) / (1.0 - sφ)
	σ := []float64{R, 1, 1}
	n := make([]float64, 3)
	SmpUnitDirector(n, σ, b)
	p, q, _ := GenInvs(σ, n, 1)
	return q / p
}

/// SMP director ////////////////////////////////////////////////////////////////////////////////////

// SmpDirector computes the director (normal vector) of the spatially mobilised plane
// Note: σ are the shifted and positive eigenvalues of the tensor, i.e:
// σ = {-λ0+σc, -λ1+σc, -λ2+σc}
func SmpDirector(N, σ []float64, b float64) {
	for i := 0; i < 3; i++ {
		N[i] = 1.0 / math.Pow(σ[i], b)
	}
}

// SmpDirectorDeriv1 computes the first order derivative of the SMP director
func SmpDirectorDeriv1(dNdσ [][]float64, σ []float64, b float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dNdσ[i][j] = 0
			if i == j {
				dNdσ[i][j] = -b * math.Pow(σ[i], -(b+1.0))
			}
		}
	}
}

// SmpDirectorDeriv2 computes the second order derivative of the SMP director
// d²N[i]/dσ[j]dσ[k]
func SmpDirectorDeriv2(i, j, k int, σ []float64, b float64) (res float64) {
	if i == j && j == k {
		res = b * (b + 1.0) * math.Pow(σ[i], -(b+2.0))
	}
	return
}

/// norm of SMP director /////////////////////////////////////////////////////////////////////////////

// SmpNormDirectorDeriv1 computes the first derivative of the norm of the SMP director
func SmpNormDirectorDeriv1(dmdσ []float64, σ []float64, b float64) (m float64) {
	m = math.Sqrt(math.Pow(σ[0], -2.0*b) + math.Pow(σ[1], -2.0*b) + math.Pow(σ[2], -2.0*b))
	for i := 0; i < 3; i++ {
		dmdσ[i] = -b * math.Pow(σ[i], -(2.0*b+1.0)) / m
	}
	return
}

// SmpNormDirectorDeriv2 computes the second order derivative of the norm of the SMP director
func SmpNormDirectorDeriv2(d2mdσdσ [][]float64, σ []float64, b, m float64, dmdσ []float64) {
	m2 := m * m
	var ci float64
	for i := 0; i < 3; i++ {
		ci = math.Pow(σ[i], -(2.0*b + 1.0))
		for j := 0; j < 3; j++ {
			d2mdσdσ[i][j] = b * ci * dmdσ[j] / m2
			if i == j {
				d2mdσdσ[i][j] += b * (2.0*b + 1.0) * math.Pow(σ[i], -2.0*(b+1.0)) / m
			}
		}
	}
}

/// unit SMP director ///////////////////////////////////////////////////////////////////////////////

// SmpUnitDirector computed the unit normal of the SMP
// Note: σ are the shifted and positive eigenvalues of the tensor, i.e:
// σ = {-λ0+σc, -λ1+σc, -λ2+σc}
// m = norm(N)
func SmpUnitDirector(n, σ []float64, b float64) (m float64) {
	m = math.Sqrt(math.Pow(σ[0], -2.0*b) + math.Pow(σ[1], -2.0*b) + math.Pow(σ[2], -2.0*b))
	for i := 0; i < 3; i++ {
		n[i] = math.Pow(σ[i], -b) / m
	}
	return
}

// SmpUnitDirectorDeriv1 computes the first derivative of the SMP unit normal
func SmpUnitDirectorDeriv1(dndσ [][]float64, σ, n []float64, b, m float64, dmdσ []float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dndσ[i][j] = -n[i] * dmdσ[j] / m
			if i == j {
				dndσ[i][j] -= b * math.Pow(σ[i], -(b+1.0)) / m
			}
		}
	}
}

// SmpUnitDirectorDeriv2 computes the second order derivative of the unit director of the SMP
// d²n[i]/dσ[j]dσ[k]
func SmpUnitDirectorDeriv2(d2ndσdσ [][][]float64, σ, n, dmdσ []float64, b, m float64, d2mdσdσ, dndσ [][]float64) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				d2ndσdσ[i][j][k] = (SmpDirectorDeriv2(i, j, k, σ, b) - n[i]*d2mdσdσ[j][k] - dmdσ[j]*dndσ[i][k] - dndσ[i][j]*dmdσ[k]) / m
			}
		}
	}
}

/// auxiliary functions /////////////////////////////////////////////////////////////////////////////

// SmpDerivs computes the SMP director and its first order derivative. In addition, the
// norm of the SMP director (m) and its first and second order derivatives are computed.
func SmpDerivs(d2mdσdσ, dndσ [][]float64, dmdσ, n, σ []float64, b float64) (m float64) {
	m = math.Sqrt(math.Pow(σ[0], -2.0*b) + math.Pow(σ[1], -2.0*b) + math.Pow(σ[2], -2.0*b))
	e1 := -2.0*b - 1.0
	for i := 0; i < 3; i++ {
		n[i] = math.Pow(σ[i], -b) / m
		dmdσ[i] = -b * math.Pow(σ[i], e1) / m
	}
	e2 := -b - 1.0
	e3 := -2.0*b - 2.0
	m2 := m * m
	var ci float64
	for i := 0; i < 3; i++ {
		ci = math.Pow(σ[i], e1)
		for j := 0; j < 3; j++ {
			dndσ[i][j] = -n[i] * dmdσ[j] / m
			d2mdσdσ[i][j] = b * ci * dmdσ[j] / m2
			if i == j {
				dndσ[i][j] -= b * math.Pow(σ[i], e2) / m
				d2mdσdσ[i][j] -= b * e1 * math.Pow(σ[i], e3) / m
			}
		}
	}
	return
}
