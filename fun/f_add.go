// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
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
	o.A, o.B = 1, 1
	for _, p := range prms {
		switch p.N {
		case "A", "a":
			o.A = p.V
		case "B", "b":
			o.B = p.V
		case "Fa", "fa":
			o.Fa = p.Fcn
		case "Fb", "fb":
			o.Fb = p.Fcn
		default:
			return chk.Err("add: parameter named %q is invalid", p.N)
		}
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
