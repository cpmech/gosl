// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

// DistGev implements the generalised extreme value (GEV) distribution
//  Gumbel, Fréchet and Weibull families also known as type I, II and III
//  extreme value distributions.
//  The shape parameter ξ governs the tail behaviour of the distribution.
//  The sub-families defined by ξ=0, ξ>0 and ξ<0 correspond, respectively,
//  to the Gumbel (ξ=0), Fréchet (ξ>0) and Weibull (ξ<0) families. In summary:
//   Type-I   == Gumbel:  ξ = 0
//   Type-II  == Fréchet: ξ > 0
//   Type-III == Weibull: ξ < 0
type DistGev struct {
	M float64 // location: μ
	S float64 // scale: σ (or β in R-lang)
	K float64 // shape: ξ. If ξ==0, the distribution is of type Gumbel
}

// set factory
func init() {
	distallocators[D_Gev] = func() Distribution { return new(DistGev) }
}

// Init initialises lognormal distribution
func (o *DistGev) Init(p *VarData) error {
	o.M, o.S, o.K = p.M, p.S, p.K
	return nil
}

// Pdf computes the probability density function @ x
func (o DistGev) Pdf(x float64) float64 {
	if o.K > 0 {
		if x > o.M-o.S/o.K {
			t := math.Pow(1.0+o.K*(x-o.M)/o.S, -1.0/o.K)
			return math.Pow(t, o.K+1) * math.Exp(-t) / o.S
		}
		return 0
	}
	if o.K < 0 {
		if x < o.M-o.S/o.K {
			t := math.Pow(1.0+o.K*(x-o.M)/o.S, -1.0/o.K)
			return math.Pow(t, o.K+1) * math.Exp(-t) / o.S
		}
		return 0
	}
	t := math.Exp((o.M - x) / o.S)
	return math.Pow(t, o.K+1) * math.Exp(-t) / o.S
}

// Cdf computes the cumulative probability function @ x
func (o DistGev) Cdf(x float64) float64 {
	if o.K > 0 {
		if x > o.M-o.S/o.K {
			t := math.Pow(1.0+o.K*(x-o.M)/o.S, -1.0/o.K)
			return math.Exp(-t)
		}
		return 0
	}
	if o.K < 0 {
		if x < o.M-o.S/o.K {
			t := math.Pow(1.0+o.K*(x-o.M)/o.S, -1.0/o.K)
			return math.Exp(-t)
		}
		return 1
	}
	t := math.Exp((o.M - x) / o.S)
	return math.Exp(-t)
}
