// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cdist implements the distance from point to a circle (2D) or a sphere (3D)
// where the circle/sphere is implicitly defined by means of F(x) = 0.
//  where
//    F(x) = sqrt((x-xc) dot (x-xc)) - r
//  with r being the radius and xc the coordinates of the centre.
//  Thus F > 0 is outside and F < 0 is inside the circle/sphere
type Cdist struct {
	xc []float64 // centre; len(xc) = 2 or 3 => must enter "xc", "yc" and "zc" as parameters
	r  float64   // radius
}

// set allocators database
func init() {
	allocators["cdist"] = func() T { return new(Cdist) }
}

// Init initialises the function
func (o *Cdist) Init(prms Params) (err error) {
	ndim := 2
	for _, p := range prms {
		if p.N == "zc" {
			ndim = 3
			break
		}
	}
	o.xc = make([]float64, ndim)
	e := prms.Connect(&o.r, "r", "cdist function")
	e += prms.Connect(&o.xc[0], "xc", "cdist function")
	e += prms.Connect(&o.xc[1], "yc", "cdist function")
	if ndim == 3 {
		e += prms.Connect(&o.xc[2], "zc", "cdist function")
	}
	if e != "" {
		err = chk.Err("%v\n", e)
		return
	}
	rtol := 1e-10
	if o.r < rtol {
		return chk.Err("cdist: radius must be greater than %g", rtol)
	}
	return
}

// F returns y = F(t, x)
func (o Cdist) F(t float64, x []float64) float64 {
	f := (x[0]-o.xc[0])*(x[0]-o.xc[0]) + (x[1]-o.xc[1])*(x[1]-o.xc[1])
	if len(x) == 3 {
		f += (x[2] - o.xc[2]) * (x[2] - o.xc[2])
	}
	return math.Sqrt(f) - o.r
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Cdist) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cdist) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cdist) Grad(v []float64, t float64, x []float64) {
	d := (x[0]-o.xc[0])*(x[0]-o.xc[0]) + (x[1]-o.xc[1])*(x[1]-o.xc[1])
	if len(x) == 3 {
		d += (x[2] - o.xc[2]) * (x[2] - o.xc[2])
	}
	d = math.Sqrt(d)
	v[0] = (x[0] - o.xc[0]) / d
	v[1] = (x[1] - o.xc[1]) / d
	if len(x) == 3 {
		v[2] = (x[2] - o.xc[1]) / d
	}
	return
}
