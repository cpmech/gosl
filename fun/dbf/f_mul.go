// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import "github.com/cpmech/gosl/chk"

// Mul implements the multiplication of two other functions.
//  F(t, x) := fa(t,x) * fb(t,x)
type Mul struct {
	Fa, Fb T
}

// set allocators database
func init() {
	allocators["mul"] = func() T { return new(Mul) }
}

// Init initialises the function
func (o *Mul) Init(prms Params) {
	for _, p := range prms {
		switch p.N {
		case "fa":
			o.Fa = p.Fcn
		case "fb":
			o.Fb = p.Fcn
		default:
			chk.Panic("mul: parameter named %q is invalid", p.N)
		}
	}
}

// F returns y = F(t, x)
func (o Mul) F(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.Fa.F(t, x) * o.Fb.F(t, x)
	}
	chk.Panic("mul: fa and fb functions are <nil>\n")
	return 0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Mul) G(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.Fa.F(t, x)*o.Fb.G(t, x) + o.Fb.F(t, x)*o.Fa.G(t, x)
	}
	chk.Panic("mul: fa and fb functions are <nil>\n")
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Mul) H(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.Fa.F(t, x)*o.Fb.H(t, x) + 2.0*o.Fa.G(t, x)*o.Fb.G(t, x) + o.Fb.F(t, x)*o.Fa.H(t, x)
	}
	chk.Panic("mul: fa and fb functions are <nil>\n")
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Mul) Grad(v []float64, t float64, x []float64) {
	chk.Panic("mul: Grad is not implemented yet")
	setvzero(v)
	return
}
