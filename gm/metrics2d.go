// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// Metrics2d holds metrics data for 2d grids
//
//   Surface             B[2](r)
//   functions:         ┌───────┐
//                      │       │
//               B[3](s)│       │B[1](s)
//                      │       │
//                      └───────┘
//                       B[0](r)
//
type Metrics2d struct {
	Nr int           // number of points along r direction
	Ns int           // number of points along s direction
	X  [][]float64   // [ns][nr] grid x coordinates
	Y  [][]float64   // [ns][nr] grid y coordinates
	Xr [][]float64   // [ns][nr] grid x derivatives w.r.t. r
	Xs [][]float64   // [ns][nr] grid x derivatives w.r.t. s
	Yr [][]float64   // [ns][nr] grid y derivatives w.r.t. r
	Ys [][]float64   // [ns][nr] grid y derivatives w.r.t. s
	J  [][]float64   // [ns][nr] grid Jacobian
	Xf [][]float64   // [4][nr or ns] surface x coordinates
	Yf [][]float64   // [4][nr or ns] surface y coordinates
	Jf [][]float64   // [4][nr or ns] surface Jacobian
	Sf [][]float64   // [4][nr or ns] surface scaling factor = ‖dXsdu‖ with u being either r or s @ surface
	Nf [][]la.Vector // [4][nr or ns] surface outward normals
}

// Allocate allocates all arrays in Metrics2d
func (o *Metrics2d) Allocate(nr, ns int) {

	o.Nr = nr
	o.Ns = ns

	o.X = make([][]float64, o.Ns)
	o.Y = make([][]float64, o.Ns)
	o.Xr = make([][]float64, o.Ns)
	o.Xs = make([][]float64, o.Ns)
	o.Yr = make([][]float64, o.Ns)
	o.Ys = make([][]float64, o.Ns)
	o.J = make([][]float64, o.Ns)

	for j := 0; j < o.Ns; j++ {
		o.X[j] = make([]float64, o.Nr)
		o.Y[j] = make([]float64, o.Nr)
		o.Xr[j] = make([]float64, o.Nr)
		o.Xs[j] = make([]float64, o.Nr)
		o.Yr[j] = make([]float64, o.Nr)
		o.Ys[j] = make([]float64, o.Nr)
		o.J[j] = make([]float64, o.Nr)
	}

	o.Xf = make([][]float64, 4)
	o.Yf = make([][]float64, 4)
	o.Jf = make([][]float64, 4)
	o.Sf = make([][]float64, 4)
	o.Nf = make([][]la.Vector, 4)

	o.Xf[0] = make([]float64, o.Nr)
	o.Yf[0] = make([]float64, o.Nr)
	o.Jf[0] = make([]float64, o.Nr)
	o.Sf[0] = make([]float64, o.Nr)
	o.Nf[0] = make([]la.Vector, o.Nr)

	o.Xf[1] = make([]float64, o.Ns)
	o.Yf[1] = make([]float64, o.Ns)
	o.Jf[1] = make([]float64, o.Ns)
	o.Sf[1] = make([]float64, o.Ns)
	o.Nf[1] = make([]la.Vector, o.Ns)

	o.Xf[2] = make([]float64, o.Nr)
	o.Yf[2] = make([]float64, o.Nr)
	o.Jf[2] = make([]float64, o.Nr)
	o.Sf[2] = make([]float64, o.Nr)
	o.Nf[2] = make([]la.Vector, o.Nr)

	o.Xf[3] = make([]float64, o.Ns)
	o.Yf[3] = make([]float64, o.Ns)
	o.Jf[3] = make([]float64, o.Ns)
	o.Sf[3] = make([]float64, o.Ns)
	o.Nf[3] = make([]la.Vector, o.Ns)

	for i := 0; i < o.Nr; i++ {
		o.Nf[0][i] = la.NewVector(2)
		o.Nf[2][i] = la.NewVector(2)
	}

	for j := 0; j < o.Ns; j++ {
		o.Nf[1][j] = la.NewVector(2)
		o.Nf[3][j] = la.NewVector(2)
	}
}

// Draw draw metric vectors
func (o *Metrics2d) Draw(sf float64, argsU, argsV *plt.A) {
	if argsU == nil {
		argsU = &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10, NoClip: true}
	}
	if argsV == nil {
		argsV = &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10, NoClip: true}
	}
	c := la.NewVector(2)
	v := la.NewVector(2)
	for j := 0; j < o.Ns; j++ {
		for i := 0; i < o.Nr; i++ {
			c[0] = o.X[j][i]
			c[1] = o.Y[j][i]
			plt.PlotOne(c[0], c[1], &plt.A{C: "k", M: ".", Ms: 3, NoClip: true})
			if true {
				v[0] = o.Xr[j][i]
				v[1] = o.Yr[j][i]
				DrawArrow2d(c, v, true, sf, argsU)
				v[0] = o.Xs[j][i]
				v[1] = o.Ys[j][i]
				DrawArrow2d(c, v, true, sf, argsV)
			} else {
				plt.Quiver(o.X, o.Y, o.Xr, o.Yr, argsU)
				plt.Quiver(o.X, o.Y, o.Xs, o.Ys, argsV)
			}
		}
	}
	for i := 0; i < o.Nr; i++ {
		plt.PlotOne(o.Xf[0][i], o.Yf[0][i], &plt.A{C: plt.C(0, 2), M: "s", Ms: 4, NoClip: true})
		plt.PlotOne(o.Xf[2][i], o.Yf[2][i], &plt.A{C: plt.C(0, 2), M: "^", Ms: 5, NoClip: true})
		c[0] = o.Xf[0][i]
		c[1] = o.Yf[0][i]
		DrawArrow2d(c, o.Nf[0][i], true, sf, &plt.A{C: plt.C(2, 0), NoClip: true})
		c[0] = o.Xf[2][i]
		c[1] = o.Yf[2][i]
		DrawArrow2d(c, o.Nf[2][i], true, sf, &plt.A{C: plt.C(3, 0), NoClip: true})
	}
	for j := 0; j < o.Ns; j++ {
		plt.PlotOne(o.Xf[1][j], o.Yf[1][j], &plt.A{C: plt.C(1, 2), M: "s", Ms: 4, NoClip: true})
		plt.PlotOne(o.Xf[3][j], o.Yf[3][j], &plt.A{C: plt.C(1, 2), M: "^", Ms: 5, NoClip: true})
		c[0] = o.Xf[1][j]
		c[1] = o.Yf[1][j]
		DrawArrow2d(c, o.Nf[1][j], true, sf, &plt.A{C: plt.C(4, 0), NoClip: true})
		c[0] = o.Xf[3][j]
		c[1] = o.Yf[3][j]
		DrawArrow2d(c, o.Nf[3][j], true, sf, &plt.A{C: plt.C(5, 0), NoClip: true})
	}
}

// GetMetrics2d computes all metrics data
func (o *Transfinite) GetMetrics2d(rvals, svals []float64) (m *Metrics2d) {

	// allocate
	m = new(Metrics2d)
	nr := len(rvals)
	ns := len(svals)
	m.Allocate(nr, ns)
	dxdu := la.NewMatrix(2, 2)
	x := la.NewVector(2)
	u := la.NewVector(2)

	// grid
	var Xr, Xs, Yr, Ys, cf float64
	for j, s := range svals {
		for i, r := range rvals {
			u[0] = r
			u[1] = s
			o.Derivs(dxdu, x, u)
			Xr = dxdu.Get(0, 0)
			Xs = dxdu.Get(0, 1)
			Yr = dxdu.Get(1, 0)
			Ys = dxdu.Get(1, 1)
			m.X[j][i] = x[0]
			m.Y[j][i] = x[1]
			m.Xr[j][i] = Xr
			m.Xs[j][i] = Xs
			m.Yr[j][i] = Yr
			m.Ys[j][i] = Ys
			m.J[j][i] = Xr*Ys - Xs*Yr
		}
	}

	// surface
	for i, r := range rvals {

		// B0
		u[0] = r
		u[1] = -1
		o.Derivs(dxdu, x, u)
		Xr = dxdu.Get(0, 0)
		Xs = dxdu.Get(0, 1)
		Yr = dxdu.Get(1, 0)
		Ys = dxdu.Get(1, 1)
		m.Xf[0][i] = x[0]
		m.Yf[0][i] = x[1]
		m.Jf[0][i] = Xr*Ys - Xs*Yr
		m.Sf[0][i] = math.Sqrt(Xr*Xr + Yr*Yr)
		cf = -fun.Sign(m.Jf[0][i]) / m.Sf[0][i]
		m.Nf[0][i][0] = -Yr * cf
		m.Nf[0][i][1] = +Xr * cf

		// B2
		u[0] = r
		u[1] = +1
		o.Derivs(dxdu, x, u)
		Xr = dxdu.Get(0, 0)
		Xs = dxdu.Get(0, 1)
		Yr = dxdu.Get(1, 0)
		Ys = dxdu.Get(1, 1)
		m.Xf[2][i] = x[0]
		m.Yf[2][i] = x[1]
		m.Jf[2][i] = Xr*Ys - Xs*Yr
		m.Sf[2][i] = math.Sqrt(Xr*Xr + Yr*Yr)
		cf = fun.Sign(m.Jf[2][i]) / m.Sf[2][i]
		m.Nf[2][i][0] = -Yr * cf
		m.Nf[2][i][1] = +Xr * cf
	}

	for j, s := range svals {

		// B1
		u[0] = +1
		u[1] = s
		o.Derivs(dxdu, x, u)
		Xr = dxdu.Get(0, 0)
		Xs = dxdu.Get(0, 1)
		Yr = dxdu.Get(1, 0)
		Ys = dxdu.Get(1, 1)
		m.Xf[1][j] = x[0]
		m.Yf[1][j] = x[1]
		m.Jf[1][j] = Xr*Ys - Xs*Yr
		m.Sf[1][j] = math.Sqrt(Xs*Xs + Ys*Ys)
		cf = fun.Sign(m.Jf[1][j]) / m.Sf[1][j]
		m.Nf[1][j][0] = +Ys * cf
		m.Nf[1][j][1] = -Xs * cf

		// B3
		u[0] = -1
		u[1] = s
		o.Derivs(dxdu, x, u)
		Xr = dxdu.Get(0, 0)
		Xs = dxdu.Get(0, 1)
		Yr = dxdu.Get(1, 0)
		Ys = dxdu.Get(1, 1)
		m.Xf[3][j] = x[0]
		m.Yf[3][j] = x[1]
		m.Jf[3][j] = Xr*Ys - Xs*Yr
		m.Sf[3][j] = math.Sqrt(Xs*Xs + Ys*Ys)
		cf = -fun.Sign(m.Jf[3][j]) / m.Sf[3][j]
		m.Nf[3][j][0] = +Ys * cf
		m.Nf[3][j][1] = -Xs * cf
	}
	return
}
