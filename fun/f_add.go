// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// Add implements the addition of two other functions.
//  F(t, x) := A*Fa(t,x) + B*Fb(t,x)
type Add struct {
	Fa, Fb Func
	A, B   float64
}

// set allocators database
func init() {
	allocators["add"] = func() Func { return new(Add) }
}

// Init initialises the function
func (o *Add) Init(prms Prms) (err error) {
	e := prms.Connect(&o.A, "a", "add function")
	e += prms.Connect(&o.B, "b", "add function")
	for _, p := range prms {
		switch p.N {
		case "fa":
			o.Fa = p.Fcn
		case "fb":
			o.Fb = p.Fcn
		}
	}
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Add) F(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.A*o.Fa.F(t, x) + o.B*o.Fb.F(t, x)
	}
	return 0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Add) G(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.A*o.Fa.G(t, x) + o.B*o.Fb.G(t, x)
	}
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Add) H(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.A*o.Fa.H(t, x) + o.B*o.Fb.H(t, x)
	}
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Add) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
