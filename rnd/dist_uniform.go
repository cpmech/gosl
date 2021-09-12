// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math/rand"

// Uniform returns a random number belonging to a uniform distribution
func Uniform(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// DistUniform implements the normal distribution
type DistUniform struct {

	// input
	A float64 // min value
	B float64 // max value
}

// set factory
func init() {
	distallocators["U"] = func() Distribution { return new(DistUniform) }
}

// Name returns the name of this probability distribution
func (o *DistUniform) Name() string { return "Uniform" }

// Init initializes uniform distribution
func (o *DistUniform) Init(p *Variable) {
	o.A, o.B = p.Min, p.Max
}

// Pdf computes the probability density function @ x
func (o DistUniform) Pdf(x float64) float64 {
	if x < o.A {
		return 0.0
	}
	if x > o.B {
		return 0.0
	}
	return 1.0 / (o.B - o.A)
}

// Cdf computes the cumulative probability function @ x
func (o DistUniform) Cdf(x float64) float64 {
	if x < o.A {
		return 0.0
	}
	if x > o.B {
		return 1.0
	}
	return (x - o.A) / (o.B - o.A)
}
