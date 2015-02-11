// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// Excitation #1 y(t) = A * (1 - cos(b*π*t)) / 2
type Exc1 struct {
	A, b float64
}

// set allocators databse
func init() {
	allocators["exc1"] = func() Func { return new(Exc1) }
}

// Init initialises the function
func (o *Exc1) Init(prms Prms) {
	for _, p := range prms {
		switch p.N {
		case "A":
			o.A = p.V
		case "b":
			o.b = p.V
		}
	}
}

// F returns y = F(t, x)
func (o Exc1) F(t float64, x []float64) float64 {
	return o.A * (1.0 - math.Cos(o.b*math.Pi*t)) / 2.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Exc1) G(t float64, x []float64) float64 {
	return o.A * o.b * math.Pi * math.Sin(o.b*math.Pi*t) / 2.0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Exc1) H(t float64, x []float64) float64 {
	return o.A * o.b * o.b * math.Pi * math.Pi * math.Cos(o.b*math.Pi*t) / 2.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Exc1) Grad(v []float64, t float64, x []float64) {
	utl.PfRed("not implemented")
}
