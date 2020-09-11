// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "gosl/utl"

// Variable holds all data defining a single random variable including information about a
// probability distribution that bests represents this variable
//
//   Some distributions:
//      "N" : Normal
//      "L" : Lognormal
//      "G" : Gumbel (Type I Extreme Value)
//      "F" : Frechet (Type II Extreme Value)
//      "U" : Uniform
//
type Variable struct {

	// input: required by many distributions
	D string  // [required] type of distribution
	M float64 // [optional] mean
	S float64 // [optional] standard deviation

	// input: Frechet
	L float64 // [Frechet] location
	C float64 // [Frechet] scale
	A float64 // [Frechet] shape

	// input: limits
	Min float64 // [optional] min value
	Max float64 // [optional] max value

	// optional
	Key string // [optional] auxiliary indentifier
	Prm *utl.P // [optional] parameter connected to this random variable

	// derived
	Normal bool         // [derived] is normal distribution
	Distr  Distribution // [derived] pointer to distribution
}

// Variables implements a set of random variables
type Variables []*Variable

// SetDistribution sets the implementation of Distribution in VarData
func (o *Variable) SetDistribution(dtype string) {
	o.Normal = dtype == "N"
	o.Distr = GetDistrib(dtype)
	o.Distr.Init(o)
}

// Transform transform x into standard normal space
func (o *Variable) Transform(x float64) (y float64, invalid bool) {
	if o.Normal {
		y = (x - o.M) / o.S
		return
	}
	F := o.Distr.Cdf(x)
	if F == 0 || F == 1 { // y = Φ⁻¹(F) → -∞ or +∞
		invalid = true
		return
	}
	y = StdInvPhi(F)
	return
}

// Init initialises distributions in Variables
func (o *Variables) Init() {
	for _, d := range *o {
		d.SetDistribution(d.D)
	}
}

// Transform transforms all variables
func (o Variables) Transform(x []float64) (y []float64, invalid bool) {
	y = make([]float64, len(x))
	for i, d := range o {
		y[i], invalid = d.Transform(x[i])
		if invalid {
			return
		}
	}
	return
}
