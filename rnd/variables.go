// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
)

// DistType indicates the distribution to which a random variable appears to belong to
type DistType int

const (
	// NormalType defines the Normal distribution type
	NormalType DistType = 1

	// LognormalType defines the Lognormal distribution type
	LognormalType DistType = 2

	// GumbelType defines the Gumbel (Type I Extreme Value) distribution type
	GumbelType DistType = 3

	// FrechetType defines the Frechet (Type II Extreme Value) distribution type
	FrechetType DistType = 4

	// UniformType defines the Uniform distribution type
	UniformType DistType = 5
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

	// input: uniform
	Min float64 // min value
	Max float64 // max value

	// optional
	Key string // auxiliary indentifier
	Prm *dbf.P // parameter connected to this random variable

	// derived
	Distr Distribution // pointer to distribution
}

// Transform transform x into standard normal space
func (o *VarData) Transform(x float64) (y float64, invalid bool) {
	if o.D == NormalType {
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
		return NormalType
	case "lognormal":
		return LognormalType
	case "gumbel":
		return GumbelType
	case "frechet":
		return FrechetType
	case "uniform":
		return UniformType
	default:
		chk.Panic("cannot get distribution named %q", name)
	}
	return NormalType
}

// GetDistrName returns distribution name from ID
func GetDistrName(typ DistType) (name string) {
	switch typ {
	case NormalType:
		return "normal"
	case LognormalType:
		return "lognormal"
	case GumbelType:
		return "gumbel"
	case FrechetType:
		return "frechet"
	case UniformType:
		return "uniform"
	default:
		chk.Panic("cannot get distribution %v", typ)
	}
	return "<unknown>"
}

// GetDistrKey returns distribution key from ID
func GetDistrKey(typ DistType) (name string) {
	switch typ {
	case NormalType:
		return "N"
	case LognormalType:
		return "L"
	case GumbelType:
		return "G"
	case FrechetType:
		return "F"
	case UniformType:
		return "U"
	default:
		chk.Panic("cannot get distribution %v", typ)
	}
	return "<unknown>"
}
