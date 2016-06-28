// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

// Xpoly1 implements F(x) as a 1st order polynomial with the components x[i]
//
//  functions:
//
//    F(x) =  a0 x0  +  a1 x1  +  a2 x2
//
//  or, if '2D = true':
//
//    F(x) =  a0 x0  +  a1 x1
//
type Xpoly1 struct {
	ndim       int
	a0, a1, a2 float64
}

// set allocators database
func init() {
	allocators["xpoly1"] = func() Func { return new(Xpoly1) }
}

// Init initialises the function
func (o *Xpoly1) Init(prms Prms) (err error) {
	o.ndim = 3
	for _, p := range prms {
		if p.N == "2D" {
			o.ndim = 2
			break
		}
	}
	prms.Connect(&o.a0, "a0", "xpoly1 function")
	prms.Connect(&o.a1, "a1", "xpoly1 function")
	if o.ndim == 3 {
		prms.Connect(&o.a2, "a2", "xpoly1 function")
	}
	return
}

// F returns y = F(t, x)
func (o Xpoly1) F(t float64, x []float64) float64 {
	if o.ndim == 3 {
		return o.a0*x[0] + o.a1*x[1] + o.a2*x[2]
	}
	return o.a0*x[0] + o.a1*x[1]
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Xpoly1) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Xpoly1) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Xpoly1) Grad(v []float64, t float64, x []float64) {
	if o.ndim == 3 && len(v) == 3 {
		v[0] = o.a0
		v[1] = o.a1
		v[2] = o.a2
		return
	}
	v[0] = o.a0
	v[1] = o.a1
}
