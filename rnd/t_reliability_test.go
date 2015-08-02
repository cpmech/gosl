// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_reliab01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("reliab01. simply supported beam")

	// Simply supported beam
	// Analyse the max deflection at mid-span of simply supported beam
	// with uniform distributed load q and concentrated load at midspan
	//  Data:
	//   L    -- span
	//   EI   -- Young's modulus times cross-sectional moment of inertia
	//   p    -- x[0] concentrated load at mid-span
	//   q    -- x[1] distributed load
	//   δlim -- max deflection (vertical displacement) at mid-span
	// Reference
	//  Haldar A, Reliability-Based Structura Design, 2005

	// constants
	δlim := 0.0381 // [m] max allowed deflection
	L := 9.144     // [m] span
	EI := 182262.0 // [kN m²] flexural rigidity
	L3 := math.Pow(L, 3.0)

	// statistics of p=x[0] and q=x[1]
	μ := []float64{111.2, 35.03} // mean values
	σ := []float64{11.12, 5.25}  // deviation values
	lrv := []bool{true, false}   // is lognormal random variable?

	// limit state function
	gfcn := func(x []float64, args ...interface{}) (g float64, err error) {
		p, q := x[0], x[1]
		g = δlim - (p*L3/EI/48.0 + 5.0*q*L3*L/EI/384.0)
		return
	}

	// derivative of limit state function
	hfcn := func(dgdx, x []float64, args ...interface{}) (err error) {
		dgdx[0] = -L3 / EI / 48.0            // dg/dp
		dgdx[1] = -5.0 * L3 * L / EI / 384.0 // dg/dq
		return
	}

	// first order reliability method
	var form ReliabFORM
	form.Init(μ, σ, lrv, gfcn, hfcn)
	form.NlsSilent = !chk.Verbose
	form.NlsCheckJ = chk.Verbose
	form.TolA = 0.005
	form.TolB = 0.005
	if chk.Verbose {
		form.PlotFnk = "beam"
	}

	// run FORM
	verbose := chk.Verbose // show messages
	βtrial := 3.0
	β, _, _, _ := form.Run(βtrial, verbose)
	chk.Scalar(tst, "β", 1e-4, β, 3.8754)
}
