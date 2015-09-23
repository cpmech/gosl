// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"github.com/cpmech/gosl/chk"
)

// sign-x gives calculates the sign function at every coordiante
// depending on the x coordinate x0. If x >= x0 ->a and if x < x0 -> -a
type XSign struct {
	x0 float64 // x coordinate where sign changes
	a  float64 // "amplitude of sign value. e.g. 6 -> 6, -6"
}

// set allocators database
func init() {
	allocators["sign-x0"] = func() Func { return new(XSign) }
}

// Init initialises the function
func (o *XSign) Init(prms Prms) (err error) {

	for _, p := range prms {
		switch p.N {
		case "x":
			o.x0 = p.V
		case "a":
			o.a = p.V
		default:
			return chk.Err("XSign: parameter named %q is invalid", p.N)
		}
	}

	return
}

// F returns y = F(t, x)
func (o XSign) F(t float64, x []float64) float64 {

	if x[0]-o.x0 < 10.e-13 {
		return -o.a
	} else {
		return o.a
	}
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o XSign) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o XSign) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o XSign) Grad(v []float64, t float64, x []float64) {

	return
}
