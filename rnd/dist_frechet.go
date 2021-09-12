// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
)

// DistFrechet implements the Frechet / Type II Extreme Value Distribution (largest value)
type DistFrechet struct {
	L float64 // location. default = 0
	C float64 // scale. default = 1
	A float64 // shape
}

// set factory
func init() {
	distallocators["F"] = func() Distribution { return new(DistFrechet) }
}

// Name returns the name of this probability distribution
func (o *DistFrechet) Name() string { return "Frechet" }

// Init initializes Frechet distribution
func (o *DistFrechet) Init(p *Variable) {
	o.L, o.C, o.A = p.L, p.C, p.A
	if math.Abs(o.C) < 1e-15 {
		o.C = 1
	}
	p.M = o.Mean()
	p.S = math.Sqrt(o.Variance())
}

// Pdf computes the probability density function @ x
func (o DistFrechet) Pdf(x float64) float64 {
	if x-o.L < 1e-15 {
		return 0
	}
	z := (x - o.L) / o.C
	return math.Exp(-math.Pow(z, -o.A)) * math.Pow(z, -1.0-o.A) * o.A / o.C
}

// Cdf computes the cumulative probability function @ x
func (o DistFrechet) Cdf(x float64) float64 {
	if x-o.L < 1e-15 {
		return 0
	}
	z := (x - o.L) / o.C
	return math.Exp(-math.Pow(z, -o.A))
}

// Mean returns the expected value
func (o DistFrechet) Mean() float64 {
	if o.A > 1.0 {
		return o.L + o.C*math.Gamma(1.0-1.0/o.A)
	}
	return math.Inf(1)
}

// Variance returns the variance
func (o DistFrechet) Variance() float64 {
	if o.A > 2.0 {
		return o.C * o.C * (math.Gamma(1.0-2.0/o.A) - math.Pow(math.Gamma(1.0-1.0/o.A), 2.0))
	}
	return math.Inf(1)
}
