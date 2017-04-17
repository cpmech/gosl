// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// Zero implements an specialisation of Cte function that always returns zero
var Zero Cte

// One implements an specialisation of Cte function that always returns one
var One = Cte{1}

// Cte implements a constant function
type Cte struct {
	C float64
}

// set allocators database
func init() {
	allocators["cte"] = func() TimeSpace { return new(Cte) }
}

// Init initialises the function
func (o *Cte) Init(prms Prms) (err error) {
	e := prms.Connect(&o.C, "c", "cte function")
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Cte) F(t float64, x []float64) float64 {
	return o.C
}

// first derivative
// G returns ∂y/∂t_cteX = G(t, x)
func (o Cte) G(t float64, x []float64) float64 {
	return 0.0
}

// second derivative
// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cte) H(t float64, x []float64) float64 {
	return 0.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cte) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
