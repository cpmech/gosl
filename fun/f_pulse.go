// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Pulse implements a ramp function
type Pulse struct {
	Ca float64
	Cb float64
	Ta float64
	Tb float64
}

// set allocators database
func init() {
	allocators["pulse"] = func() Func { return new(Pulse) }
}

// Init initialises the function
func (o *Pulse) Init(prms Prms) (err error) {
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
			return chk.Err("pulse: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Pulse) F(t float64, x []float64) float64 {
	if t <= o.Ta {
		return o.Ca
	}
	if t >= 2.0*o.Tb-o.Ta {
		return o.Ca
	}
	return o.Cb + (o.Ca-o.Cb)*(1.0+math.Cos(math.Pi*(t-o.Ta)/(o.Tb-o.Ta)))/2.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Pulse) G(t float64, x []float64) float64 {
	if t <= o.Ta {
		return 0
	}
	if t >= 2.0*o.Tb-o.Ta {
		return 0
	}
	return -(o.Ca - o.Cb) * math.Sin(math.Pi*(t-o.Ta)/(o.Tb-o.Ta)) * math.Pi / (2.0 * (o.Tb - o.Ta))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Pulse) H(t float64, x []float64) float64 {
	if t <= o.Ta {
		return 0
	}
	if t >= 2.0*o.Tb-o.Ta {
		return 0
	}
	return -(o.Ca - o.Cb) * math.Cos(math.Pi*(t-o.Ta)/(o.Tb-o.Ta)) * math.Pi * math.Pi / (2.0 * (o.Tb - o.Ta) * (o.Tb - o.Ta))
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Pulse) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
