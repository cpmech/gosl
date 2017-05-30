// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"
import m "math"

// Gauss implements the the gaussian discribution in 2D
// like: F(x,y) = 1 / (2*pi*sig) *exp((x-mux)^2/(2)/(2*sig^2) + (y-muy)^2/(2)/(2*sig^2))
type Gauss struct {
	b1 float64 // mux
	b2 float64 // muy
	c  float64 // sig
}

// set allocators database
func init() {
	allocators["gauss"] = func() Func { return new(Gauss) }
}

// Init initialises the function
func (o *Gauss) Init(prms Prms) (err error) {

	for _, p := range prms {
		switch p.N {
		case "mux":
			o.b1 = p.V
		case "muy":
			o.b2 = p.V
		case "sig":
			o.c = p.V
		default:
			return chk.Err("Gauss: parameter named %q is invalid", p.N)
		}
	}

	return
}

// F returns y = F(t, x)
func (o Gauss) F(t float64, x []float64) float64 {

	b1 := o.b1
	b2 := o.b2
	c := o.c

	a := 1. / (c * m.Sqrt(2.*m.Pi))

	G := a * m.Exp(-(m.Pow(x[0]-b1, 2)+m.Pow(x[1]-b2, 2))/(2.*m.Pow(c, 2)))

	return G
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Gauss) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Gauss) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Gauss) Grad(v []float64, t float64, x []float64) {

	b := o.b1
	c := o.c

	a := 1. / (c * m.Sqrt(2.*m.Pi))

	v[0] = -2. * (x[0] - b) / (2. * m.Pow(c, 2)) * a * m.Exp(-2.*m.Pow(x[0]-b, 2)/(2.*m.Pow(c, 2)))
	v[1] = -2. * (x[0] - b) / (2. * m.Pow(c, 2)) * a * m.Exp(-2.*m.Pow(x[0]-b, 2)/(2.*m.Pow(c, 2)))
	return
}
