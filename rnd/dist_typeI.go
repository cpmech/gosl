// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

const (
	EULER_CTE = 0.577215664901532860606512090082402431042159335939923598805767234884867726777664670936947063291746749
)

// DistTypeI implements the Type I Extreme Value Distribution (largest value)
type DistTypeI struct {
	U float64 // location: characteristic largest value
	B float64 // scale: measure of dispersion of the largest value
}

// set factory
func init() {
	distallocators[D_TypeI] = func() Distribution { return new(DistTypeI) }
}

// Init initialises lognormal distribution
func (o *DistTypeI) Init(p *VarData) error {
	μ, σ := p.M, p.S
	o.B = σ * math.Sqrt(6.0) / math.Pi
	o.U = μ - EULER_CTE*o.B
	return nil
}

// Pdf computes the probability density function @ x
func (o DistTypeI) Pdf(x float64) float64 {
	mz := (o.U - x) / o.B
	return math.Exp(mz) * math.Exp(-math.Exp(mz)) / o.B
}

// Cdf computes the cumulative probability function @ x
func (o DistTypeI) Cdf(x float64) float64 {
	mz := (o.U - x) / o.B
	return math.Exp(-math.Exp(mz))
}
