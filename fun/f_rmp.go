// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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
	for _, p := range prms {
		switch p.N {
		case "Ca", "ca":
			o.Ca = p.V
		case "Cb", "cb":
			o.Cb = p.V
		case "Ta", "ta":
			o.Ta = p.V
		case "Tb", "tb":
			o.Tb = p.V
		default:
			return chk.Err("rmp: parameter named %q is invalid", p.N)
		}
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
