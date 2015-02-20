// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"sort"

	"github.com/cpmech/gosl/chk"
)

// Pts is a function based on a linear interpolation over a set of points
type Pts struct {
	p    points
	tmin float64
	tmax float64
}

// set allocators databse
func init() {
	allocators["pts"] = func() Func { return new(Pts) }
}

// Init initialises the function
func (o *Pts) Init(prms Prms) (err error) {
	var T, Y []float64
	for _, p := range prms {
		if len(p.N) < 2 {
			return chk.Err(_pts_err01, p.N)
		}
		switch p.N[:1] {
		case "t":
			T = append(T, p.V)
		case "y":
			Y = append(Y, p.V)
		default:
			return chk.Err("pts: parameter named %q is invalid", p.N)
		}
	}
	if len(T) != len(Y) {
		return chk.Err(_pts_err04, len(T), len(Y))
	}
	for i, t := range T {
		o.p = append(o.p, &point{t, Y[i]})
	}
	sort.Sort(o.p)
	for i := 1; i < len(o.p); i++ {
		if math.Abs(o.p[i].t-o.p[i-1].t) < 1e-7 {
			return chk.Err(_pts_err02, o.p[i].t, 1e-7)
		}
	}
	o.tmin = o.p[0].t
	o.tmax = o.p[len(o.p)-1].t
	return
}

// F returns y = F(t, x)
func (o Pts) F(t float64, x []float64) float64 {
	if t < o.tmin {
		return o.p[0].y
	}
	if t > o.tmax {
		return o.p[len(o.p)-1].y
	}
	for i := 1; i < len(o.p); i++ {
		if t <= o.p[i].t {
			return o.p[i-1].y + (t-o.p[i-1].t)*(o.p[i].y-o.p[i-1].y)/(o.p[i].t-o.p[i-1].t)
		}
	}
	return 0
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Pts) G(t float64, x []float64) float64 {
	if t < o.tmin || t > o.tmax {
		return 0
	}
	for i := 1; i < len(o.p); i++ {
		if t <= o.p[i].t {
			return (o.p[i].y - o.p[i-1].y) / (o.p[i].t - o.p[i-1].t)
		}
	}
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Pts) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Pts) Grad(v []float64, t float64, x []float64) {
	setvzero(v)
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// point holds point data
type point struct {
	t, y float64
}

// points is a set of points
type points []*point

// Len the length of points
func (o points) Len() int {
	return len(o)
}

// Swap swaps two points
func (o points) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less compares points considering the t-coordinate
func (o points) Less(i, j int) bool {
	return o[i].t < o[j].t
}

// error messages
var (
	_pts_err01 = "parameter name must start with 't' or 'y'; e.g. t0, t1, t2, ... y0, y1, y2. %q is incorrect"
	_pts_err02 = "points must not have coincident t coordinates (t=%g; tol=%g)"
	_pts_err03 = "t=%v is out of range. tmin=%g tmax=%g"
	_pts_err04 = "number of 't' parameters must be the same as the number of 'y' parameters. len(T)=%d != len(Y)=%d"
)
