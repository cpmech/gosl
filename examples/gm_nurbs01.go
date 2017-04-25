// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func main() {

	// set of available vertices => to define control points
	verts := [][]float64{
		{0.0, 0.0, 0, 1},
		{1.0, 0.2, 0, 1},
		{0.5, 1.5, 0, 1},
		{2.5, 2.0, 0, 1},
		{2.0, 0.4, 0, 1},
		{3.0, 0.0, 0, 1},
	}

	// NURBS knots
	knots := [][]float64{
		{0, 0, 0, 0, 0.3, 0.7, 1, 1, 1, 1},
	}

	// NURBS curve
	gnd := 1           // geometry number of dimensions: 1=>curve, 2=>surface, 3=>volume
	orders := []int{3} // 3rd order along the only dimension
	curve := new(gm.Nurbs)
	curve.Init(gnd, orders, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))

	// refine NURBS
	refined := curve.Krefine([][]float64{
		{0.15, 0.5, 0.85}, // new knots along the only dimension
	})

	// configuration
	argsCtrlA := &plt.A{C: "k", Ls: "--", L: "control"}
	argsCtrlB := &plt.A{C: "green", L: "refined: control"}
	argsElemsA := &plt.A{C: "b", L: "curve"}
	argsElemsB := &plt.A{C: "orange", Ls: "none", M: "*", Me: 20, L: "refined: curve"}
	argsIdsA := &plt.A{C: "k", Fsz: 7}
	argsIdsB := &plt.A{C: "green", Fsz: 7}

	// plot
	ndim := 2
	npts := 41
	plt.Reset(true, &plt.A{WidthPt: 400})
	curve.DrawCtrl(ndim, true, argsCtrlA, argsIdsA)
	curve.DrawElems(ndim, npts, true, argsElemsA, nil)
	refined.DrawCtrl(ndim, true, argsCtrlB, argsIdsB)
	refined.DrawElems(ndim, npts, false, argsElemsB, nil)
	plt.AxisOff()
	plt.Equal()
	plt.LegendX([]*plt.A{argsCtrlA, argsCtrlB, argsElemsA, argsElemsB}, &plt.A{LegOut: true, LegNcol: 2})
	plt.Save("/tmp/gosl", "gm_nurbs01")
}
