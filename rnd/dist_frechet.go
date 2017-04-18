// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// DistFrechet implements the Frechet / Type II Extreme Value Distribution (largest value)
type DistFrechet struct {
	L float64 // location. default = 0
	C float64 // scale. default = 1
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
	p.M = o.Mean()
	p.S = math.Sqrt(o.Variance())
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

// FrechetPlotCoef plots coefficients for Frechet parameter's estimation
func FrechetPlotCoef(dirout, fn string, amin, amax float64) {
	np := 201
	A := utl.LinSpace(amin, amax, np)
	X := make([]float64, np)
	Y := make([]float64, np)
	var dist DistFrechet
	for i := 0; i < np; i++ {
		dist.Init(&VarData{L: 0, A: A[i]})
		X[i] = 1.0 / A[i]
		μ := dist.Mean()
		σ2 := dist.Variance()
		δ2 := σ2 / (μ * μ)
		Y[i] = 1.0 + δ2
	}
	k := np - 1
	plt.Plot(X, Y, nil)
	plt.Text(X[k], Y[k], io.Sf("(%.4f,%.4f)", X[k], Y[k]), nil)
	plt.Text(X[0], Y[0], io.Sf("(%.4f,%.4f)", X[0], Y[0]), nil)
	plt.Gll("$1/\\alpha$", "$1+\\delta^2$", nil)
	plt.SaveD(dirout, fn)
}
