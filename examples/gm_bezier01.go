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

	// quadratic Bezier
	var bez gm.BezierQuad

	// control points
	bez.Q = [][]float64{
		{-1.0, +1.0},
		{+0.5, -2.0},
		{+2.0, +4.0},
	}

	// points on Bezier curve
	np := 11
	Xb, Yb, _ := bez.GetPoints(utl.LinSpace(0, 1, np))

	// quadratic curve
	Xc := utl.LinSpace(-1, 2, np)
	Yc := utl.GetMapped(Xc, func(x float64) float64 { return x * x })

	// control points
	Xq, Yq, _ := bez.GetControlCoords()

	// plot
	plt.Reset(true, &plt.A{WidthPt: 300})
	plt.Plot(Xq, Yq, &plt.A{C: "k", M: "*", NoClip: true, L: "control"})
	plt.Plot(Xc, Yc, &plt.A{C: "r", M: "o", Void: true, Ms: 10, Ls: "none", L: "y=x*x", NoClip: true})
	plt.Plot(Xb, Yb, &plt.A{C: "b", Ls: "-", M: ".", L: "Bezier", NoClip: true})
	plt.HideAllBorders()
	plt.Gll("x", "y", &plt.A{LegLoc: "upper left"})
	plt.Equal()
	plt.Save("/tmp/gosl", "gm_bezier01")
}
