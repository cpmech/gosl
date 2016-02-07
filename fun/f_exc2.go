// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Excitation # 2: y(t) = if t < ta a*sin(b*π*t), else 0
type Exc2 struct {
	Ta float64
	A  float64
	B  float64
}

// set allocators databse
func init() {
	allocators["exc2"] = func() Func { return new(Exc2) }
}

// Init initialises the function
func (o *Exc2) Init(prms Prms) (err error) {
	e := prms.Connect(&o.Ta, "ta")
	e += prms.Connect(&o.A, "a")
	e += prms.Connect(&o.B, "b")
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Exc2) F(t float64, x []float64) float64 {
	if t < o.Ta {
		return o.A * math.Sin(o.B*math.Pi*t)
	}
	return 0.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Exc2) G(t float64, x []float64) float64 {
	if t < o.Ta {
		return o.A * o.B * math.Pi * math.Cos(o.B*math.Pi*t)
	}
	return 0.0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Exc2) H(t float64, x []float64) float64 {
	if t < o.Ta {
		return -o.A * o.B * o.B * math.Pi * math.Pi * math.Sin(o.B*math.Pi*t)
	}
	return 0.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Exc2) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
