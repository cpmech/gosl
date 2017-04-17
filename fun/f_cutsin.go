// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cutsin implements
// 	if find["cps"]: # means cut_positive is True
//		if y < 0: y(t) = a * sin(b*t) + c
//	 	else: y(t) = 0
//	else:			# means cut_positive is False so cut negative values
//		if y > 0: y(t) = a * sin(b*t) + c
//	 	else: y(t) = 0
// Input:
//  b/pi -- is a flag that says that b is in fact b divided by π
//          thus, the code will multiply b by π internally
type CutSin struct {

	// parameters
	A float64
	B float64
	C float64

	// derived
	b_is_b_div_pi bool
	cut_positive  bool
}

// set allocators database
func init() {
	allocators["cut-sin"] = func() TimeSpace { return new(CutSin) }
}

// Init initialises the function
func (o *CutSin) Init(prms Prms) (err error) {
	e := prms.Connect(&o.A, "a", "cut-sin function")
	e += prms.Connect(&o.C, "c", "cut-sin function")
	p := prms.Find("b/pi")
	if p == nil {
		e += prms.Connect(&o.B, "b", "cut-sin function")
	} else {
		e += prms.Connect(&o.B, "b/pi", "cut-sin function")
		o.b_is_b_div_pi = true
	}
	p = prms.Find("cps")
	if p == nil {
		o.cut_positive = false
	} else {
		o.cut_positive = true
	}
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o CutSin) F(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.cut_positive {
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
func (o CutSin) G(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.cut_positive {
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
func (o CutSin) H(t float64, x []float64) float64 {
	b := o.B
	if o.b_is_b_div_pi {
		b = o.B * math.Pi
	}
	if o.cut_positive {
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
func (o CutSin) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
