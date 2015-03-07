// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// Mul implements the multiplication of two other functions.
//  F(t, x) := Fa(t,x) * Fb(t,x)
type Mul struct {
	Fa, Fb Func
}

// set allocators database
func init() {
	allocators["mul"] = func() Func { return new(Mul) }
}

// Init initialises the function
func (o *Mul) Init(prms Prms) (err error) {
	for _, p := range prms {
		switch p.N {
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
func (o Mul) F(t float64, x []float64) float64 {
	if o.Fa != nil && o.Fb != nil {
		return o.Fa.F(t, x) * o.Fb.F(t, x)
	}
	return 0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Mul) G(t float64, x []float64) float64 {
	panic("f_mul.go: G: test needed")
	if o.Fa != nil && o.Fb != nil {
		return o.Fa.G(t, x)*o.Fb.F(t, x) + o.Fa.F(t, x)*o.Fb.G(t, x)
	}
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Mul) H(t float64, x []float64) float64 {
	panic("f_mul.go: H is not implemented yet")
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Mul) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
