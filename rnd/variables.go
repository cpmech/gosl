// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "github.com/cpmech/gosl/chk"

// DistType indicates the distribution to which a random variable appears to belong to
type DistType int

const (
	D_Normal    DistType = iota + 1 // normal
	D_Lognormal                     // lognormal
	D_Gumbel                        // Type I Extreme Value
	D_Frechet                       // Type II Extreme Value
)

// VarData implements data defining one random variable
type VarData struct {

	// input
	D DistType // type of distribution
	M float64  // mean
	S float64  // standard deviation

	// input: Frechet
	L float64 // location
	C float64 // scale
	A float64 // shape

	// derived
	distr Distribution // pointer to distribution
}

// Transform transform x into standard normal space
func (o *VarData) Transform(x float64) (y float64, invalid bool) {
	if o.D == D_Normal {
		y = (x - o.M) / o.S
		return
	}
	F := o.distr.Cdf(x)
	if F == 0 || F == 1 { // y = Φ⁻¹(F) → -∞ or +∞
		invalid = true
		return
	}
	y = StdInvPhi(F)
	return
}

// Variables implements a set of random variables
type Variables []*VarData

// Init initialises distributions in Variables
func (o *Variables) Init() (err error) {
	for _, d := range *o {
		d.distr, err = GetDistrib(d.D)
		if err != nil {
			chk.Err("cannot find distribution:\n%v", err)
			return
		}
		err = d.distr.Init(d)
		if err != nil {
			chk.Err("cannot initialise variables:\n%v", err)
			return
		}
	}
	return
}
