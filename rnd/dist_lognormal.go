// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"math/rand"
)

// Lognormal returns a random number belonging to a lognormal distribution
func Lognormal(μ, σ float64) float64 {
	δ := σ / μ
	v := math.Log(1.0 + δ*δ)
	z := math.Sqrt(v)
	n := math.Log(μ) - v/2.0
	return math.Exp(n + z*rand.NormFloat64())
}

// DistLogNormal implements the lognormal distribution
type DistLogNormal struct {

	// input
	N float64 // mean of log(x)
	Z float64 // standard deviation of log(x)

	// auxiliary
	A float64 // 1 / (z sqrt(2 π))
	B float64 // -1 / (2 z²)
}

// set factory
func init() {
	distallocators["L"] = func() Distribution { return new(DistLogNormal) }
}

// Name returns the name of this probability distribution
func (o *DistLogNormal) Name() string { return "Lognormal" }

// CalcDerived computes derived/auxiliary quantities
func (o *DistLogNormal) CalcDerived() {
	o.A = 1.0 / (o.Z * math.Sqrt2 * math.SqrtPi)
	o.B = -1.0 / (2.0 * o.Z * o.Z)
}

// Init initializes lognormal distribution
func (o *DistLogNormal) Init(p *Variable) {
	μ, σ := p.M, p.S
	δ := σ / μ
	v := math.Log(1.0 + δ*δ)
	o.Z = math.Sqrt(v)
	o.N = math.Log(μ) - v/2.0
	o.CalcDerived()
}

// Pdf computes the probability density function @ x
func (o DistLogNormal) Pdf(x float64) float64 {
	if x < 1e-15 {
		return 0
	}
	return o.A * math.Exp(o.B*math.Pow(math.Log(x)-o.N, 2.0)) / x
}

// Cdf computes the cumulative probability function @ x
func (o DistLogNormal) Cdf(x float64) float64 {
	if x < 1e-15 {
		return 0
	}
	return (1.0 + math.Erf((math.Log(x)-o.N)/(o.Z*math.Sqrt2))) / 2.0
}
