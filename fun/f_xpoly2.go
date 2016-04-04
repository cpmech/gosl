// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

// Xpoly2 implements F(x) as a 2nd order polynomial with the components x[i]
//
//  functions:
//
//    F(x) =  a0 x0    +  a1 x1    +  a2 x2
//         +  b0 x0²   +  b1 x1²   +  b2 x2²
//         + c01 x0 x1 + c12 x1 x2 + c20 x2 x0
//
//  or, if '2D = true':
//
//    F(x) =  a0 x0    +  a1 x1
//         +  b0 x0²   +  b1 x1²
//         + c01 x0 x1
//
type Xpoly2 struct {
	ndim          int
	a0, a1, a2    float64
	b0, b1, b2    float64
	c01, c12, c20 float64
}

// set allocators database
func init() {
	allocators["xpoly2"] = func() Func { return new(Xpoly2) }
}

// Init initialises the function
func (o *Xpoly2) Init(prms Prms) (err error) {
	o.ndim = 3
	for _, p := range prms {
		if p.N == "2D" {
			o.ndim = 2
			break
		}
	}
	prms.Connect(&o.a0, "a0", "xpoly2 function")
	prms.Connect(&o.a1, "a1", "xpoly2 function")
	prms.Connect(&o.b0, "b0", "xpoly2 function")
	prms.Connect(&o.b1, "b1", "xpoly2 function")
	prms.Connect(&o.c01, "c01", "xpoly2 function")
	if o.ndim == 3 {
		prms.Connect(&o.a2, "a2", "xpoly2 function")
		prms.Connect(&o.b2, "b2", "xpoly2 function")
		prms.Connect(&o.c12, "c12", "xpoly2 function")
		prms.Connect(&o.c20, "c20", "xpoly2 function")
	}
	return
}

// F returns y = F(t, x)
func (o Xpoly2) F(t float64, x []float64) float64 {
	if o.ndim == 3 {
		return o.a0*x[0] + o.a1*x[1] + o.a2*x[2] +
			+o.b0*x[0]*x[0] + o.b1*x[1]*x[1] + o.b2*x[2]*x[2] +
			+o.c01*x[0]*x[1] + o.c12*x[1]*x[2] + o.c20*x[2]*x[0]
	}
	return o.a0*x[0] + o.a1*x[1] +
		+o.b0*x[0]*x[0] + o.b1*x[1]*x[1] +
		+o.c01*x[0]*x[1]
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Xpoly2) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Xpoly2) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Xpoly2) Grad(v []float64, t float64, x []float64) {
	if o.ndim == 3 && len(v) == 3 {
		v[0] = o.a0 + 2.0*o.b0*x[0] + o.c01*x[1] + o.c20*x[2]
		v[1] = o.a1 + 2.0*o.b1*x[1] + o.c01*x[0] + o.c12*x[2]
		v[2] = o.a2 + 2.0*o.b2*x[2] + o.c12*x[1] + o.c20*x[0]
		return
	}
	v[0] = o.a0 + 2.0*o.b0*x[0] + o.c01*x[1]
	v[1] = o.a1 + 2.0*o.b1*x[1] + o.c01*x[0]
}
