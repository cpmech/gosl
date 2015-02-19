// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// RefIncRL1 implements a specialisation of the reference increasing model
// Reference concept model: dydx growth but solved from right to left
// with right-most initial point @ x,y=1,1 with (0 <= x <= 1)
// Flipped model is also available if λ1 < λ0 ( b == 1 )
type RefIncRL1 struct {

	// paramters
	λ0 float64 // slope @ right side
	λ1 float64 // slope @ left side
	α  float64 // minimum y @ left side
	β  float64 // transition coefficient

	// constants
	c1, c2, c3 float64 // constants

	// auxiliary
	b float64 // constant indicating flipped model or not (-1 => not flipped)
}

// set allocators database
func init() {
	allocators["ref-inc-rl1"] = func() Func { return new(RefIncRL1) }
}

// Init initialises the model
func (o *RefIncRL1) Init(prms Prms) (err error) {

	// parameters
	for _, p := range prms {
		switch p.N {
		case "lam0":
			o.λ0 = p.V
		case "lam1":
			o.λ1 = p.V
		case "alp":
			o.α = p.V
		case "bet":
			o.β = p.V
		default:
			return chk.Err("ref-inc-rl1: parameter named %q is invalid", p.N)
		}
	}

	// set b
	o.b = -1.0 // not flipped
	if o.λ1 < o.λ0 {
		o.b = 1.0 // flipped
	}

	// constants
	o.c1 = o.β * o.b * (o.λ1 - o.λ0)
	o.c2 = math.Exp(o.β * o.b * o.α)
	o.c3 = math.Exp(o.β*o.b*(1.0-o.λ0)) - o.c2*math.Exp(o.c1)
	return
}

// F returns y = F(t, x)
func (o RefIncRL1) F(t float64, x []float64) float64 {
	return o.λ0*t + math.Log(o.c3+o.c2*math.Exp(o.c1*t))/(o.β*o.b)
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o RefIncRL1) G(t float64, x []float64) float64 {
	return o.λ0 + o.c1*o.c2*math.Exp(o.c1*t)/(o.β*o.b*(o.c3+o.c2*math.Exp(o.c1*t)))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o RefIncRL1) H(t float64, x []float64) float64 {
	return o.c1 * o.c1 * o.c2 * o.c3 * math.Exp(o.c1*t) / (o.β * o.b * math.Pow(o.c3+o.c2*math.Exp(o.c1*t), 2.0))
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o RefIncRL1) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
