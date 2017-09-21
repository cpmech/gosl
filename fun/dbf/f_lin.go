// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import "github.com/cpmech/gosl/chk"

// Lin implements a linear function w.r.t t
//  y = m * (t - ts)
type Lin struct {
	M  float64 // slope
	Ts float64 // shift
}

// set allocators database
func init() {
	allocators["lin"] = func() T { return new(Lin) }
}

// Init initialises the function
func (o *Lin) Init(prms Params) {
	e := prms.Connect(&o.M, "m", "lin function")
	prms.Connect(&o.Ts, "ts", "lin function")
	if e != "" {
		chk.Panic("%v\n", e)
	}
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
