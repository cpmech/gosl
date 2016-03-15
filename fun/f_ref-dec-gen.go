// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// RefDecGen implements the reference decreasing model (general)
//  [1] Pedroso DM, Sheng D, Zhao J. The concept of reference curves for constitutive modelling in soil mechanics, Computers and Geotechnics, 36, 1-2, http://dx.doi.org/10.1016/j.compgeo.2008.01.009
type RefDecGen struct {
	// paramters
	β          float64 // beta coeficient
	a, b, c    float64 // distance function coeffcients
	A, B       float64 // lambda function coeficients
	xini, yini float64 // initial point on curve

	// constants
	c1, c2, c3 float64
}

// set allocators database
func init() {
	allocators["ref-dec-gen"] = func() Func { return new(RefDecGen) }
}

// Init initialises the function
func (o *RefDecGen) Init(prms Prms) (err error) {

	// parameters
	e := prms.Connect(&o.β, "bet", "ref-dec-gen function")
	e += prms.Connect(&o.a, "a", "ref-dec-gen function")
	e += prms.Connect(&o.b, "b", "ref-dec-gen function")
	e += prms.Connect(&o.c, "c", "ref-dec-gen function")
	e += prms.Connect(&o.A, "A", "ref-dec-gen function")
	e += prms.Connect(&o.B, "B", "ref-dec-gen function")
	e += prms.Connect(&o.xini, "xini", "ref-dec-gen function")
	e += prms.Connect(&o.yini, "yini", "ref-dec-gen function")
	if e != "" {
		err = chk.Err("%v\n", e)
		return
	}

	// constants
	o.c1 = o.β * (o.b*o.A - o.a)
	o.c2 = ((o.A - o.B) / (o.A - o.a/o.b)) * math.Exp(-o.β*o.c)
	o.c3 = math.Exp(o.β*o.b*(o.yini+o.A*o.xini)) - o.c2*math.Exp(o.c1*o.xini)
	return
}

// F returns y = F(t, x)
func (o RefDecGen) F(t float64, x []float64) float64 {
	return -o.A*t + (1.0/(o.β*o.b))*math.Log(o.c3+o.c2*math.Exp(o.c1*t))
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o RefDecGen) G(t float64, x []float64) float64 {
	return -o.A + (o.c1*o.c2*math.Exp(o.c1*t))/(o.β*o.b*(o.c3+o.c2*math.Exp(o.c1*t)))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o RefDecGen) H(t float64, x []float64) float64 {
	d := o.c3 + o.c2*math.Exp(o.c1*t)
	return (o.c1 * o.c1 * o.c2 * o.c3 * math.Exp(o.c1*t)) / (o.β * o.b * d * d)
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o RefDecGen) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
