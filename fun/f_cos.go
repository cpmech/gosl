// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cos implements y(t) = a * cos(b*t) + c
type Cos struct {
	a float64
	b float64
	c float64
}

// set allocators database
func init() {
	allocators["cos"] = func() Func { return new(Cos) }
}

// Init initialises the function
func (o *Cos) Init(prms Prms) (err error) {
	for _, p := range prms {
		switch p.N {
		case "a":
			o.a = p.V
		case "b":
			o.b = p.V
		case "c":
			o.c = p.V
		case "b/pi": // b/π => b = b/pi * π
			o.b = p.V * math.Pi
		default:
			return chk.Err("cos: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Cos) F(t float64, x []float64) float64 {
	return o.a*math.Cos(o.b*t) + o.c
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Cos) G(t float64, x []float64) float64 {
	return -o.a * o.b * math.Sin(o.b*t)
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cos) H(t float64, x []float64) float64 {
	return -o.a * o.b * o.b * math.Cos(o.b*t)
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cos) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
