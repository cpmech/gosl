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
	D_Nrm DistType = iota + 1
	D_Log
	D_Bet
	D_Gev
)

// VarData implements data defining one random variable
type VarData struct {

	// input
	D     DistType // type of distribution: "nrm", "log", "beta"
	M     float64  // mean-type parameter
	S     float64  // deviation-type parameter
	MSlog bool     // M and S parameters are m and s of log distribution and not μ and σ from standard (default)
	K     float64  // shape parameter
	A     float64  // lower limit
	B     float64  // upper limit
	Q     float64  // q-parameter (e.g. for Beta distribution)
	R     float64  // r-parameter (e.g. for Beta distribution)

	// derived
	distr Distribution // pointer to distribution
	m, s  float64      // auxiliary: parameters of lognormal distribution
}

// CalcEquiv computes equivalent normal parameters at check/design point
func (o *VarData) CalcEquiv(x float64) (μN, σN float64, invalid bool) {
	switch o.D {
	case D_Nrm:
		return o.M, o.S, false
	case D_Log:
		if x < TOLMINLOG {
			//err = chk.Err("x must be ≥ 0 for lognormal distribution. %g is invalid", x)
			invalid = true
			return
		}
		σN = o.s * x
		μN = (1.0 - math.Log(x) + o.m) * x
		//io.Pforan("s=%v m=%v σN=%v μN=%v\n", o.s, o.m, σN, μN)
		if μN < 0 { // TODO: check this
			F := o.distr.Cdf(x)
			if F == 0 || F == 1 { // z = Φ⁻¹(F) → -∞ or +∞   TODO: check this
				//io.Pfred("cannot compute equivalent normal parameters @ %g because F=%g", x, F)
				μN = 0
				σN = o.S
				return
			}
			z := StdInvPhi(F)
			μN = 0
			σN = x / z
		}
	case D_Gev:
		// using algorithm from:
		//  Rackwitz R, Fiessler B. An algorithm for calculation of structural reliability
		//  under combined loading. Berichte zur Sicherheitstheorie der Bauwerke,
		//  Lab. f. Konstr. Ingb. Munich, Germany; 1977
		F := o.distr.Cdf(x)
		if F == 0 || F == 1 { // z = Φ⁻¹(F) → -∞ or +∞   TODO: check this
			//io.Pfred("cannot compute equivalent normal parameters @ %g because F=%g", x, F)
			μN = 0
			σN = o.S
			return
		}
		f := o.distr.Pdf(x)
		z := StdInvPhi(F)
		σN = Stdphi(z) / f
		μN = x - σN*z
		if μN < 0 { // TODO: check if this is neccessary for all GEVs
			//io.Pfred("warning: fixing μN and σN to avoid μN<0\n")
			μN = 0
			σN = x / z
		}
	default:
		chk.Panic("cannot handle <%v> distribution yet", o.D)
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
