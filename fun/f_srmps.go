// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Increasing or decreasing smooth-ramp-smooth function
type Srmps struct {
	Ca, Cb float64
	Ta, Tb float64
}

// set allocators databse
func init() {
	allocators["srmps"] = func() Func { return new(Srmps) }
}

// Init initialises the function
func (o *Srmps) Init(prms Prms) (err error) {
	e := prms.Connect(&o.Ca, "ca", "srmps function")
	e += prms.Connect(&o.Cb, "cb", "srmps function")
	e += prms.Connect(&o.Ta, "ta", "srmps function")
	e += prms.Connect(&o.Tb, "tb", "srmps function")
	if e != "" {
		err = chk.Err("%v\n", e)
	}
	return
}

// F returns y = F(t, x)
func (o Srmps) F(t float64, x []float64) float64 {
	if t < o.Ta {
		return o.Ca
	}
	if t > o.Tb {
		return o.Cb
	}
	return o.Cb + (o.Ca-o.Cb)*(math.Cos(math.Pi*(t-o.Ta)/(o.Tb-o.Ta))+1.0)/2.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Srmps) G(t float64, x []float64) float64 {
	if t < o.Ta {
		return 0
	}
	if t > o.Tb {
		return 0
	}
	return -math.Pi * (o.Ca - o.Cb) * math.Sin(math.Pi*(t-o.Ta)/(o.Tb-o.Ta)) / (2.0 * (o.Tb - o.Ta))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Srmps) H(t float64, x []float64) float64 {
	if t < o.Ta {
		return 0
	}
	if t > o.Tb {
		return 0
	}
	return -math.Pi * math.Pi * (o.Ca - o.Cb) * math.Cos(math.Pi*(t-o.Ta)/(o.Tb-o.Ta)) / (2.0 * (o.Tb - o.Ta) * (o.Tb - o.Ta))
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Srmps) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
