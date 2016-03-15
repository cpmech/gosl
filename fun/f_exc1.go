// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Excitation #1 y(t) = a * (1 - cos(b*π*t)) / 2
type Exc1 struct {
	A, B float64
}

// set allocators databse
func init() {
	allocators["exc1"] = func() Func { return new(Exc1) }
}

// Init initialises the function
func (o *Exc1) Init(prms Prms) (err error) {
	e := prms.Connect(&o.A, "a", "exc1 function")
	e += prms.Connect(&o.B, "b", "exc1 function")
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Exc1) F(t float64, x []float64) float64 {
	return o.A * (1.0 - math.Cos(o.B*math.Pi*t)) / 2.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Exc1) G(t float64, x []float64) float64 {
	return o.A * o.B * math.Pi * math.Sin(o.B*math.Pi*t) / 2.0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Exc1) H(t float64, x []float64) float64 {
	return o.A * o.B * o.B * math.Pi * math.Pi * math.Cos(o.B*math.Pi*t) / 2.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Exc1) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
