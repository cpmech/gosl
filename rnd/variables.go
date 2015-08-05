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
func (o *VarData) CalcEquiv(x float64) (μN, σN float64, invalid bool) {
	switch o.D {
	case D_Normal:
		return o.M, o.S, false
	case D_Log:
		if x < TOLMINLOG {
			//err = chk.Err("x must be ≥ 0 for lognormal distribution. %g is invalid", x)
			invalid = true
			return
		}
		σN = o.s * x
		μN = (1.0 - math.Log(x) + o.m) * x
	default:
		chk.Panic("cannot handle %+v distribution yet", o.D)
	}
	return
}

// Transform transform x into standard normal space
func (o *VarData) Transform(x float64) (y float64, invalid bool) {
	μ, σ, invalid := o.CalcEquiv(x)
	if invalid {
		return
	}
	y = (x - μ) / σ
	return
}

// InvTransform transform y from standard normal space into original space
func (o *VarData) InvTransform(y float64) (x float64, invalid bool) {
	chk.Panic("InvTransform is not implemented yet")
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

// Transform transform all variables into standard normal space
func (o *Variables) Transform(x []float64) (y []float64, invalid bool) {
	y = make([]float64, len(x))
	for i, d := range *o {
		y[i], invalid = d.Transform(x[i])
		if invalid {
			return
		}
	}
	return
}
