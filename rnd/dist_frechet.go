// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

// DistFrechet implements the Frechet / Type II Extreme Value Distribution (largest value)
type DistFrechet struct {
	L float64 // location
	C float64 // scale
	A float64 // shape
}

// set factory
func init() {
	distallocators[D_Frechet] = func() Distribution { return new(DistFrechet) }
}

// Init initialises lognormal distribution
func (o *DistFrechet) Init(p *VarData) error {
	o.L, o.C, o.A = p.L, p.C, p.A
	if math.Abs(o.C) < ZERO {
		o.C = 1
	}
	return nil
}

// Pdf computes the probability density function @ x
func (o DistFrechet) Pdf(x float64) float64 {
	if x-o.L < ZERO {
		return 0
	}
	z := (x - o.L) / o.C
	return math.Exp(-math.Pow(z, -o.A)) * math.Pow(z, -1.0-o.A) * o.A / o.C
}

// Cdf computes the cumulative probability function @ x
func (o DistFrechet) Cdf(x float64) float64 {
	if x-o.L < ZERO {
		return 0
	}
	z := (x - o.L) / o.C
	return math.Exp(-math.Pow(z, -o.A))
}
