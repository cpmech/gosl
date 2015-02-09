// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"github.com/cpmech/gosl/utl"
)

// Zero implements an specialisation of Cte function that always returns zero
var Zero Cte

// One implements an specialisation of Cte function that always returns one
var One = Cte{1}

// Cte implements a constant function
type Cte struct {
	C float64
}

// set allocators database
func init() {
	allocators["cte"] = func() Func { return new(Cte) }
}

// Init initialises the function
func (o *Cte) Init(prms Prms) {
	for _, p := range prms {
		switch p.N {
		case "C", "c":
			o.C = p.V
		default:
			utl.Panic("parameter named %q is incorrect", p.N)
		}
	}
}

// F returns y = F(t, x)
func (o Cte) F(t float64, x []float64) float64 {
	return o.C
}

// first derivative
// G returns ∂y/∂t_cteX = G(t, x)
func (o Cte) G(t float64, x []float64) float64 {
	return 0.0
}

// second derivative
// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cte) H(t float64, x []float64) float64 {
	return 0.0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cte) Grad(v []float64, t float64, x []float64) {
	utl.Panic("not implemented")
}
