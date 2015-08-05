// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/num"
)

// DistType indicates the distribution to which a random variable appears to belong to
//type DistType int
//
//const (
//D_Normal DistType = iota
//D_Lognormal
//)

// VarData implements data defining one random variable
type VarData struct {
	D   string  // type of distribution: "nrm", "log", "beta"
	M   float64 // mean-type parameter
	S   float64 // deviation-type parameter
	Std bool    // mean and deviation are standard values
	A   float64 // lower limit
	B   float64 // upper limit
	Q   float64 // q-parameter (e.g. for Beta distribution)
	R   float64 // r-parameter (e.g. for Beta distribution)
}

// Variables implements a set of random variables
type Variables struct {
	Data []*VarData     // information of variables
	Dist []Distribution // distributions
}

// Init initialises distributions in Variables
func (o *Variables) Init() (err error) {
	o.Dist = make([]Distribution, len(o.Data))
	for i, d := range o.Data {
		o.Dist[i], err = GetDistrib(d.D)
		if err != nil {
			chk.Err("cannot find distribution:\n%v", err)
			return
		}
		err = o.Dist[i].Init(d)
		if err != nil {
			chk.Err("cannot initialise variables:\n%v", err)
			return
		}
	}
	return
}

// RecalcPrms recomputes distribution parameters
func (o *Variables) RecalcPrms(x float64) {
}

func lognormal_calc_equiv_prms(μ, σ, x float64, μσ_are_ms bool) (μN, σN float64) {
	if x < TOLMINLOG {
		chk.Panic("cannot compute μN and σN because x<0. x=%g", x)
	}
	m, s := μ, σ
	if !μσ_are_ms { // compute lognormal variables from statistics of x
		δ := σ / μ
		s = math.Sqrt(math.Log(1 + δ*δ))
		m = math.Log(μ) - s*s/2
	}
	σN = s * x
	μN = (1 - math.Log(x) + m) * x
	//io.Pfpink("μ=%v σ=%v x=%v\n", μ, σ, x)
	//io.Pfpink("δ=%v s=%v m=%v\n", δ, s, m)
	//io.Pfpink("μN=%v σN=%v\n", μN, σN)
	return
}

// calc_norm_vars computes normal variables from original variables
func calc_norm_vars(ds []string, μ, σ, x []float64) (y []float64, invalid bool) {
	y = make([]float64, len(x))
	copy(y, x)
	if len(ds) == 0 { // all standard normal variables
		return
	}
	for i, typ := range ds {
		switch typ {
		case "nrm":
			y[i] = (x[i] - μ[i]) / σ[i]
		case "logMS", "logMuSig":
			if x[i] < TOLMINLOG {
				return nil, true
			}
			μN, σN := lognormal_calc_equiv_prms(μ[i], σ[i], x[i], typ == "logMS")
			y[i] = (x[i] - μN) / σN
		default:
			chk.Panic("distribution %q is not available", typ)
		}
	}
	return
}

// calc_orig_vars computes original variables from normal variables
func calc_orig_vars(ds []string, μ, σ, y []float64) (x []float64, invalid bool) {
	x = make([]float64, len(y))
	copy(x, y)
	if len(ds) == 0 { // all standard normal variables
		return
	}
	for i, typ := range ds {
		switch typ {
		case "nrm":
			x[i] = μ[i] + σ[i]*y[i]
		case "logMS", "logMuSig":
			if y[i] < TOLMINLOG {
				return nil, true
			}

			// nonlinear problem for lognormal variable u[0] = x[i]
			var nls num.NlSolver
			nls.Init(1, func(fu, u []float64) error {
				μNtmp, σNtmp := lognormal_calc_equiv_prms(μ[i], σ[i], u[0], typ == "logMS")
				fu[0] = u[0] - (μNtmp + σNtmp*y[i])
				return nil
			}, nil, nil, false, true, nil)
			u := []float64{μ[i] + σ[i]*y[i]}
			nls.SetTols(1e-4, 1e-4, 1e-4, 1e-15)
			nls.Lsearch = false
			err := nls.Solve(u, true)
			if err != nil {
				chk.Panic("nonlinear solver failed:\n%v", err)
			}
			x[i] = u[0]

		default:
			chk.Panic("distribution %q is not available", typ)
		}
	}
	return
}
