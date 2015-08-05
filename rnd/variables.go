// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// DistType indicates the distribution to which a random variable appears to belong to
type DistType int

const (
	D_Normal DistType = iota
	D_Log
	D_Beta
)

// VarData implements data defining one random variable
type VarData struct {

	// input
	D   DistType // type of distribution: "nrm", "log", "beta"
	M   float64  // mean-type parameter
	S   float64  // deviation-type parameter
	Std bool     // mean and deviation are standard values; e.g μ and σ instead of m and s
	A   float64  // lower limit
	B   float64  // upper limit
	Q   float64  // q-parameter (e.g. for Beta distribution)
	R   float64  // r-parameter (e.g. for Beta distribution)

	// derived
	distr Distribution // pointer to distribution
	m, s  float64      // auxiliary: parameters of lognormal distribution
}

// CalcEquiv computes equivalent normal parameters at check/design point
func (o *VarData) CalcEquiv(x float64) (μN, σN float64, err error) {
	switch o.D {
	case D_Normal:
		return o.M, o.S, nil
	case D_Log:
		if x < TOLMINLOG {
			err = chk.Err("x must be ≥ 0 for lognormal distribution. %g is invalid", x)
			return
		}
		σN = o.s * x
		μN = (1.0 - math.Log(x) + o.m) * x
	default:
		chk.Panic("cannot handle %+v distribution yet", o.D)
	}
	return
}

// Map maps x into a standard/normal space
func (o *VarData) Map(x float64) (y float64, err error) {
	μ, σ, err := o.CalcEquiv(x)
	if err != nil {
		return
	}
	y = (x - μ) / σ
	return
}

// InvMap maps y from a standard/normal space back to x-space
func (o *VarData) InvMap(y float64) (x float64, err error) {
	chk.Panic("InvMap is not implemented yet")
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
