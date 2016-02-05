// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

// DistGumbel implements the Gumbel / Type I Extreme Value Distribution (largest value)
type DistGumbel struct {
	U float64 // location: characteristic largest value
	B float64 // scale: measure of dispersion of the largest value
}

// set factory
func init() {
	distallocators[D_Gumbel] = func() Distribution { return new(DistGumbel) }
}

// Init initialises lognormal distribution
func (o *DistGumbel) Init(p *VarData) error {
	μ, σ := p.M, p.S
	o.B = σ * math.Sqrt(6.0) / math.Pi
	o.U = μ - EULER*o.B
	return nil
}

// Pdf computes the probability density function @ x
func (o DistGumbel) Pdf(x float64) float64 {
	mz := (o.U - x) / o.B
	return math.Exp(mz) * math.Exp(-math.Exp(mz)) / o.B
}

// Cdf computes the cumulative probability function @ x
func (o DistGumbel) Cdf(x float64) float64 {
	mz := (o.U - x) / o.B
	return math.Exp(-math.Exp(mz))
}
