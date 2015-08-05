// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

const TOLMINLOG = 1e-16

// DistLogNormal implements the lognormal distribution
type DistLogNormal struct {

	// input
	M float64 // location
	S float64 // scale

	// auxiliary
	A float64 // 1 / (s sqrt(2 π))
	B float64 // -1 / (2 s²)
}

// set factory
func init() {
	distallocators["log"] = func() Distribution { return new(DistLogNormal) }
}

// CalcDerived computes derived/auxiliary quantities
func (o *DistLogNormal) CalcDerived() {
	o.A = 1.0 / (o.S * math.Sqrt2 * math.SqrtPi)
	o.B = -1.0 / (2.0 * o.S * o.S)
}

// Init initialises lognormal distribution
func (o *DistLogNormal) Init(p *VarData) error {
	o.M, o.S = p.M, p.S
	if p.Std {
		μ, σ := p.M, p.S
		δ := σ / μ
		o.S = math.Sqrt(math.Log(1.0 + δ*δ))
		o.M = math.Log(μ) - o.S*o.S/2.0
	}
	o.CalcDerived()
	return nil
}

// Pdf computes the probability density function @ x
func (o DistLogNormal) Pdf(x float64) float64 {
	if x < TOLMINLOG {
		return 0
	}
	return o.A * math.Exp(o.B*math.Pow(math.Log(x)-o.M, 2.0)) / x
}

// Cdf computes the cumulative probability function @ x
func (o DistLogNormal) Cdf(x float64) float64 {
	if x < TOLMINLOG {
		return 0
	}
	return (1.0 + math.Erf((math.Log(x)-o.M)/(o.S*math.Sqrt2))) / 2.0
}
