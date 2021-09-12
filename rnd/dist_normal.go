// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"math/rand"
)

// Normal returns a random number belonging to a normal distribution
func Normal(μ, σ float64) float64 {
	return μ + σ*rand.NormFloat64()
}

// Stdphi implements φ(x), the standard probability density function
func Stdphi(x float64) float64 {
	return math.Exp(-x*x/2.0) / math.Sqrt2 / math.SqrtPi
}

// StdPhi implements Φ(x), the standard cumulative distribution function
func StdPhi(x float64) float64 {
	return (1.0 + math.Erf(x/math.Sqrt2)) / 2.0
}

// StdInvPhi implements Φ⁻¹(x), the inverse standard cumulative distribution function
func StdInvPhi(x float64) float64 {
	return ltqnorm(x)
}

// DistNormal implements the normal distribution
type DistNormal struct {

	// input
	Mu  float64 // μ: mean
	Sig float64 // σ: std deviation

	// auxiliary
	a float64 // 1 / (σ sqrt(2 π))
	b float64 // -1 / (2 σ²)
}

// set factory
func init() {
	distallocators["N"] = func() Distribution { return new(DistNormal) }
}

// Name returns the name of this probability distribution
func (o *DistNormal) Name() string { return "Normal" }

// CalcDerived compute derived/auxiliary quantities
func (o *DistNormal) CalcDerived() {
	o.a = 1.0 / (o.Sig * math.Sqrt2 * math.SqrtPi)
	o.b = -1.0 / (2.0 * o.Sig * o.Sig)
}

// Init initializes normal distribution
func (o *DistNormal) Init(p *Variable) {
	o.Mu, o.Sig = p.M, p.S
	o.CalcDerived()
}

// Pdf computes the probability density function @ x
func (o DistNormal) Pdf(x float64) float64 {
	return o.a * math.Exp(o.b*math.Pow(x-o.Mu, 2.0))
}

// Cdf computes the cumulative probability function @ x
func (o DistNormal) Cdf(x float64) float64 {
	return (1.0 + math.Erf((x-o.Mu)/(o.Sig*math.Sqrt2))) / 2.0
}
