// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Increasing or decreasing smooth-ramp-smooth function
type Srmps struct {
	ca, cb float64
	ta, tb float64
}

// set allocators databse
func init() {
	allocators["srmps"] = func() Func { return new(Srmps) }
}

// Init initialises the function
func (o *Srmps) Init(prms Prms) (err error) {
	for _, p := range prms {
		switch p.N {
		case "ca":
			o.ca = p.V
		case "cb":
			o.cb = p.V
		case "ta":
			o.ta = p.V
		case "tb":
			o.tb = p.V
		default:
			return chk.Err("srmps: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Srmps) F(t float64, x []float64) float64 {
	if t < o.ta {
		return o.ca
	}
	if t > o.tb {
		return o.cb
	}
	return o.cb + (o.ca-o.cb)*(math.Cos(math.Pi*(t-o.ta)/(o.tb-o.ta))+1.0)/2.0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Srmps) G(t float64, x []float64) float64 {
	if t < o.ta {
		return 0
	}
	if t > o.tb {
		return 0
	}
	return -math.Pi * (o.ca - o.cb) * math.Sin(math.Pi*(t-o.ta)/(o.tb-o.ta)) / (2.0 * (o.tb - o.ta))
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Srmps) H(t float64, x []float64) float64 {
	if t < o.ta {
		return 0
	}
	if t > o.tb {
		return 0
	}
	return -math.Pi * math.Pi * (o.ca - o.cb) * math.Cos(math.Pi*(t-o.ta)/(o.tb-o.ta)) / (2.0 * (o.tb - o.ta) * (o.tb - o.ta))
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Srmps) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
