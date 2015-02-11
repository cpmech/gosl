// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"github.com/cpmech/gosl/utl"
)

// Lin implements a linear function w.r.t t
//  y = m * (t - ts)
type Lin struct {
	M  float64 // slope
	Ts float64 // shift
}

// set allocators database
func init() {
	allocators["lin"] = func() Func { return new(Lin) }
}

// Init initialises the function
func (o *Lin) Init(prms Prms) (err error) {
	for _, p := range prms {
		switch p.N {
		case "m":
			o.M = p.V
		case "ts":
			o.Ts = p.V
		default:
			return utl.Err("lin: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Lin) F(t float64, x []float64) float64 {
	return o.M * (t - o.Ts)
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Lin) G(t float64, x []float64) float64 {
	return o.M
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Lin) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Lin) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
