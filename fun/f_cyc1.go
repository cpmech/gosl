// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/chk"

// Cyc1 implements a cycles function # 1 (by Yunpeng Zhang for paper on groundwater settlement)
type Cyc1 struct {
	mx float64 // x scaling factor multiplier. default == 1.0
	my float64 // y scaling factor multiplier. default == 1.0
}

// set allocators database
func init() {
	allocators["cyc1"] = func() Func { return new(Cyc1) }
}

// Init initialises the function
func (o *Cyc1) Init(prms Prms) (err error) {
	o.mx, o.my = 1.0, 1.0
	for _, p := range prms {
		switch p.N {
		case "mx":
			o.mx = p.V
		case "my":
			o.my = p.V
		default:
			return chk.Err("cyc1: parameter named %q is invalid", p.N)
		}
	}
	return
}

// F returns y = F(t, x)
func (o Cyc1) F(t float64, x []float64) float64 {

	t *= o.mx

	if t <= 6.0 {
		return (130.0 - 22.557112*t) * o.my

	} else if t >= 6.0 && t <= 60.0 {
		return (-5.30612) * o.my

	} else if t > 60.0 && t <= 119.0 {
		return (0.000299*t*t*t - 0.1022*t*t + 12.404*t - 443.022) * o.my

	} else if t > 119.0 && t <= 123.729 {
		return (90.3272 + (t-118.9)/(123.729-118.925)*(-5.30612-90.3272)) * o.my

	} else if t >= 123.729 && t <= 178.847 {
		return (-5.30612) * o.my

	} else if t >= 178.847 && t <= 239.2686 {
		return (0.000198*t*t*t - 0.1452*t*t + 36.106*t - 2948) * o.my

	} else if t >= 239.2686 && t <= 249.0 {
		return (90.5818 + (t-239.2686)/(249.0-239.2686)*(-5.30612-90.5818)) * o.my

	} else if t >= 249.0 && t <= 300.540205529 {
		return (-5.30612) * o.my

	} else if t >= 300.540205529 && t <= 358.0986 {
		return (0.0004351*t*t*t - 0.4535*t*t + 158.41*t - 18461) * o.my

	} else if t >= 358.0986 && t <= 369.976214206 {
		return (90.9 + (t-358.0986)/(369.976214206-358.0986)*(-5.30612-90.9)) * o.my

	} else if t >= 369.976214206 && t <= 419.368838935 {
		return (-5.30612) * o.my

	} else if t >= 419.368838935 && t < 477.6417 {
		return (0.00037272*t*t*t - 0.52885*t*t + 250.8*t - 39664) * o.my
	}

	chk.Panic("cyc1: F cannot compute with t=%g", t)
	return 0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Cyc1) G(t float64, x []float64) float64 {
	chk.Panic("cyc1: G is not implemented")
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cyc1) H(t float64, x []float64) float64 {
	chk.Panic("cyc1: H is not implemented")
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cyc1) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}
