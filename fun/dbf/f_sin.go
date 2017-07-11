// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Sin implements y(t) = a * sin(b*t) + c
// Input:
//  b/pi -- is a flag that says that b is in fact b divided by π
//          thus, the code will multiply b by π internally
type Sin struct {

	// parameters
	A float64
	B float64
	C float64

	// derived
	bDivPi bool
}

// set allocators database
func init() {
	allocators["sin"] = func() T { return new(Sin) }
}

// Init initialises the function
func (o *Sin) Init(prms Params) (err error) {
	e := prms.Connect(&o.A, "a", "sin function")
	e += prms.Connect(&o.C, "c", "sin function")
	p := prms.Find("b/pi")
	if p == nil {
		e += prms.Connect(&o.B, "b", "sin function")
	} else {
		e += prms.Connect(&o.B, "b/pi", "sin function")
		o.bDivPi = true
	}
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Sin) F(t float64, x []float64) float64 {
	b := o.B
	if o.bDivPi {
		b = o.B * math.Pi
	}
	return o.A*math.Sin(b*t) + o.C
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Sin) G(t float64, x []float64) float64 {
	b := o.B
	if o.bDivPi {
		b = o.B * math.Pi
	}
	return o.A * b * math.Cos(b*t)
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Sin) H(t float64, x []float64) float64 {
	b := o.B
	if o.bDivPi {
		b = o.B * math.Pi
	}
	return -o.A * b * b * math.Sin(b*t)
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Sin) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
