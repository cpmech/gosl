// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"

	"gosl/chk"
)

// RefDecSp1 implements a specialisation of the reference decreasing model
//  [1] Pedroso DM, Sheng D, Zhao J. The concept of reference curves for constitutive modelling in soil mechanics, Computers and Geotechnics, 36, 1-2, http://dx.doi.org/10.1016/j.compgeo.2008.01.009
//
//                y ^
//                  |
//               ya o
//                  | \
//                  |
//                  |   \
//                  |
//      ------------+-------------------------> x
//                  |
//               yb o--._ -\-------- λ0 = 0
//                  |    `.
//                  |      \\
//                  |       .
//                  |        \
//                  |         \---
//                  |          \ | λ1
//                  |           \|
//                  |            \
//                  |
//                  |
type RefDecSp1 struct {
	// parameters
	β      float64 // beta coefficient
	λ1     float64 // slope
	ya, yb float64 // points on curve

	// constants
	c1, c2, c3  float64
	c1timestmax float64
}

// set allocators database
func init() {
	allocators["ref-dec-sp1"] = func() T { return new(RefDecSp1) }
}

// Init initialises the model
func (o *RefDecSp1) Init(prms Params) {

	// parameters
	e := prms.Connect(&o.β, "bet", "ref-dec-sp1 function")
	e += prms.Connect(&o.λ1, "lam1", "ref-dec-sp1 function")
	e += prms.Connect(&o.ya, "ya", "ref-dec-sp1 function")
	e += prms.Connect(&o.yb, "yb", "ref-dec-sp1 function")
	if e != "" {
		chk.Panic("%v\n", e)
	}

	// check
	if o.yb >= o.ya {
		chk.Panic("yb(%g) must be smaller than ya(%g)", o.yb, o.ya)
	}

	// constants
	o.c1 = o.β * o.λ1
	o.c2 = math.Exp(-o.β * o.ya)
	o.c3 = math.Exp(-o.β*o.yb) - o.c2
	o.c1timestmax = 400

	// check
	if math.IsInf(o.c2, 0) || math.IsInf(o.c3, 0) {
		chk.Panic("β*ya or β*yb is too large:\n β=%v, ya=%v, yb=%v\n c1=%v, c2=%v, c3=%v", o.β, o.ya, o.yb, o.c1, o.c2, o.c3)
	}
}

// F returns y = F(t, x)
func (o RefDecSp1) F(t float64, x []float64) float64 {
	if o.c1*t > o.c1timestmax {
		return o.ya - o.λ1*t
	}
	return -math.Log(o.c3+o.c2*math.Exp(o.c1*t)) / o.β
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o RefDecSp1) G(t float64, x []float64) float64 {
	if o.c1*t > o.c1timestmax {
		return -o.λ1
	}
	return -(o.c1 * o.c2 * math.Exp(o.c1*t)) / (o.β * (o.c3 + o.c2*math.Exp(o.c1*t)))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o RefDecSp1) H(t float64, x []float64) float64 {
	if o.c1*t > o.c1timestmax {
		return 0.0
	}
	d := o.c3 + o.c2*math.Exp(o.c1*t)
	return -(o.c1 * o.c1 * o.c2 * o.c3 * math.Exp(o.c1*t)) / (o.β * d * d)
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o RefDecSp1) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
