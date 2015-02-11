// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// Excitation # 2: y(t) = if t < ta A*sin(b*π*t), else 0
type Exc2 struct {
	ta float64
	A  float64
	b  float64
}

// set allocators databse
func init() {
	allocators["exc2"] = func() Func { return new(Exc2) }
}

// Init initialises the function
func (o *Exc2) Init(prms Prms) (err error) {
	for _, p := range prms {
		switch p.N {
		case "ta":
			o.ta = p.V
		case "A":
			o.A = p.V
		case "b":
			o.b = p.V
		default:
			return utl.Err("exc2: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Exc2) F(t float64, x []float64) float64 {
	if t < o.ta {
		return o.A * math.Sin(o.b*math.Pi*t)
	}
	return 0.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Exc2) G(t float64, x []float64) float64 {
	if t < o.ta {
		return o.A * o.b * math.Pi * math.Cos(o.b*math.Pi*t)
	}
	return 0.0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Exc2) H(t float64, x []float64) float64 {
	if t < o.ta {
		return -o.A * o.b * o.b * math.Pi * math.Pi * math.Sin(o.b*math.Pi*t)
	}
	return 0.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Exc2) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
