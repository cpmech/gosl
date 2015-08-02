// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

type DistNormal struct {
	Mu  float64 // μ: mean
	Sig float64 // σ: std deviation

	// auxiliary
	vari2 float64 // 2 * σ²: 2 times variance
	den   float64 // σ * sqrt(2 * π)
}

// CalcDerived compute derived/auxiliary quantities
func (o *DistNormal) CalcDerived() {
	o.vari2 = 2.0 * o.Sig * o.Sig
	o.den = o.Sig * math.Sqrt2 * math.SqrtPi
}

// Init initialises normal distribution
func (o *DistNormal) Init(μ, σ float64) {
	o.Mu, o.Sig = μ, σ
	o.CalcDerived()
}

// Pdf computes the probability density function @ x
func (o DistNormal) Pdf(x float64) float64 {
	return math.Exp(-math.Pow(x-o.Mu, 2.0)/o.vari2) / o.den
}

// Cdf computes the cumulative probability function @ x
func (o DistNormal) Cdf(x float64) float64 {
	return (1.0 + math.Erf((x-o.Mu)/(o.Sig*math.Sqrt2))) / 2.0
}
