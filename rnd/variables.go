// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/fun/dbf"
)

// VarData implements data defining one random variable
//
//   example of distributions
//      "N" : Normal
//      "L" : Lognormal
//      "G" : Gumbel (Type I Extreme Value)
//      "F" : Frechet (Type II Extreme Value)
//      "U" : Uniform
//
type VarData struct {

	// input
	D string  // type of distribution
	M float64 // mean
	S float64 // standard deviation

	// input: Frechet
	L float64 // location
	C float64 // scale
	A float64 // shape

	// input: limits
	Min float64 // min value
	Max float64 // max value

	// optional
	Key string // auxiliary indentifier
	Prm *dbf.P // parameter connected to this random variable

	// derived
	Normal bool         // is normal distribution
	Distr  Distribution // pointer to distribution
}

// Variables implements a set of random variables
type Variables []*VarData

// SetDistribution sets the implementation of Distribution in VarData
func (o *VarData) SetDistribution(dtype string) (err error) {
	o.Normal = dtype == "N"
	o.Distr, err = GetDistrib(dtype)
	if err != nil {
		return
	}
	err = o.Distr.Init(o)
	return
}

// Transform transform x into standard normal space
func (o *VarData) Transform(x float64) (y float64, invalid bool) {
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
func (o *Variables) Init() (err error) {
	for _, d := range *o {
		err = d.SetDistribution(d.D)
		if err != nil {
			return
		}
	}
	return
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
