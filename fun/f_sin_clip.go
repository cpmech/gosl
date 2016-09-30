// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Sin_clip implements if y < 0: y(t) = a * sin(b*t) + c
type Sin_clip struct {

	// parameters
	A float64
	B float64
	C float64

	// derived
	b_is_b_div_pi bool
	pos_is_zero   bool
}

// set allocators database
func init() {
	allocators["sin_clip"] = func() Func { return new(Sin_clip) }
}

// Init initialises the function
func (o *Sin_clip) Init(prms Prms) (err error) {
	e := prms.Connect(&o.A, "a", "sin_clip function")
	e += prms.Connect(&o.C, "c", "sin_clip function")
	p := prms.Find("b/pi")
	if p == nil {
		e += prms.Connect(&o.B, "b", "sin_clip function")
	} else {
		e += prms.Connect(&o.B, "b/pi", "sin_clip function")
		o.b_is_b_div_pi = true
	}
	p = prms.Find("pos")
	if p == nil {
		o.pos_is_zero = false
	} else {
		o.pos_is_zero = true
	}
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Sin_clip) F(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.pos_is_zero {
		if o.A*math.Sin(b*t)+o.C <= 0.0 {
			return o.A*math.Sin(b*t) + o.C
		} else {
			return 0.0
		}
	} else {
		if o.A*math.Sin(b*t)+o.C >= 0.0 {
			return o.A*math.Sin(b*t) + o.C
		} else {
			return 0.0
		}
	}
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Sin_clip) G(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.pos_is_zero {
		if o.A*math.Sin(b*t)+o.C <= 0.0 {
			return o.A * b * math.Cos(b*t)
		} else {
			return 0.0
		}
	} else {
		if o.A*math.Sin(b*t)+o.C >= 0.0 {
			return o.A * b * math.Cos(b*t)
		} else {
			return 0.0
		}
	}
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Sin_clip) H(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.pos_is_zero {
		if o.A*math.Sin(b*t)+o.C <= 0.0 {
			return -o.A * b * b * math.Sin(b*t)
		} else {
			return 0.0
		}
	} else {
		if o.A*math.Sin(b*t)+o.C >= 0.0 {
			return -o.A * b * b * math.Sin(b*t)
		} else {
			return 0.0
		}
	}
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Sin_clip) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
