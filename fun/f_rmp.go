// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// Rmp implements a ramp function
type Rmp struct {
	Ca float64
	Cb float64
	Ta float64
	Tb float64
}

// set allocators database
func init() {
	allocators["rmp"] = func() Func { return new(Rmp) }
}

// Init initialises the function
func (o *Rmp) Init(prms Prms) (err error) {
	e := prms.Connect(&o.Ca, "ca", "rmp function")
	e += prms.Connect(&o.Cb, "cb", "rmp function")
	e += prms.Connect(&o.Ta, "ta", "rmp function")
	e += prms.Connect(&o.Tb, "tb", "rmp function")
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Rmp) F(t float64, x []float64) float64 {
	if t < o.Ta {
		return o.Ca
	}
	if t < o.Tb {
		return o.Ca - (t-o.Ta)*(o.Ca-o.Cb)/(o.Tb-o.Ta)
	}
	return o.Cb
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Rmp) G(t float64, x []float64) float64 {
	if t < o.Ta {
		return 0
	}
	if t < o.Tb {
		return -(o.Ca - o.Cb) / (o.Tb - o.Ta)
	}
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Rmp) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Rmp) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
