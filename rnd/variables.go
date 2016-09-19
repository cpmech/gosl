// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
)

// DistType indicates the distribution to which a random variable appears to belong to
type DistType int

const (
	D_Normal    DistType = iota + 1 // normal
	D_Lognormal                     // lognormal
	D_Gumbel                        // Type I Extreme Value
	D_Frechet                       // Type II Extreme Value
	D_Uniform                       // uniform
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

	// optional
	Min float64  // min value
	Max float64  // max value
	Prm *fun.Prm // parameter connected to this random variable

	// derived
	Distr Distribution // pointer to distribution
}

// Transform transform x into standard normal space
func (o *VarData) Transform(x float64) (y float64, invalid bool) {
	if o.D == D_Normal {
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

// Variables implements a set of random variables
type Variables []*VarData

// Init initialises distributions in Variables
func (o *Variables) Init() (err error) {
	for _, d := range *o {
		d.Distr, err = GetDistrib(d.D)
		if err != nil {
			chk.Err("cannot find distribution:\n%v", err)
			return
		}
		err = d.Distr.Init(d)
		if err != nil {
			chk.Err("cannot initialise variables:\n%v", err)
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

// GetDistribution returns distribution ID from name
func GetDistribution(name string) DistType {
	switch name {
	case "normal":
		return D_Normal
	case "lognormal":
		return D_Lognormal
	case "gumbel":
		return D_Gumbel
	case "frechet":
		return D_Frechet
	default:
		chk.Panic("cannot get distribution named %q", name)
	}
	return D_Normal
}

// GetDistrName returns distribution name from ID
func GetDistrName(typ DistType) (name string) {
	switch typ {
	case D_Normal:
		return "normal"
	case D_Lognormal:
		return "lognormal"
	case D_Gumbel:
		return "gumbel"
	case D_Frechet:
		return "frechet"
	default:
		chk.Panic("cannot get distribution %v", typ)
	}
	return "<unknown>"
}
