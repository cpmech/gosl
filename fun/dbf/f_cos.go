// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cos implements y(t) = a * cos(b*t) + c
// Input:
//  b/pi -- is a flag that says that b is in fact b divided by π
//          thus, the code will multiply b by π internally
type Cos struct {

	// parameters
	A float64
	B float64
	C float64

	// derived
	bIsBdivPi bool
}

// set allocators database
func init() {
	allocators["cos"] = func() T { return new(Cos) }
}

// Init initialises the function
func (o *Cos) Init(prms Params) (err error) {
	e := prms.Connect(&o.A, "a", "cos function")
	e += prms.Connect(&o.C, "c", "cos function")
	p := prms.Find("b/pi")
	if p == nil {
		e += prms.Connect(&o.B, "b", "cos function")
	} else {
		e += prms.Connect(&o.B, "b/pi", "cos function")
		o.bIsBdivPi = true
	}
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Cos) F(t float64, x []float64) float64 {
	b := o.B
	if o.bIsBdivPi {
		b = o.B * math.Pi
	}
	return o.A*math.Cos(b*t) + o.C
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Cos) G(t float64, x []float64) float64 {
	b := o.B
	if o.bIsBdivPi {
		b = o.B * math.Pi
	}
	return -o.A * b * math.Sin(b*t)
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cos) H(t float64, x []float64) float64 {
	b := o.B
	if o.bIsBdivPi {
		b = o.B * math.Pi
	}
	return -o.A * b * b * math.Cos(b*t)
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cos) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
