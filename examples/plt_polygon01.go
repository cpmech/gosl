// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/plt"
)

func main() {

	// point coordinates
	P := [][]float64{
		{-2.5, 0.0},
		{-5.5, 4.0},
		{0.0, 10.0},
		{5.5, 4.0},
		{2.5, 0.0},
	}

	// formatting/styling data
	// Fc: face color, Ec: edge color, Lw: linewidth
	stPlines := &plt.A{Fc: "#c1d7cf", Ec: "#4db38e", Lw: 4.5, Closed: true, NoClip: true}
	stCircles := &plt.A{Fc: "#b2cfa5", Ec: "#5dba35", Z: 1}
	stArrows := &plt.A{Fc: "cyan", Ec: "blue", Z: 2, Scale: 50, Style: "fancy"}

	// clear drawing area, with defaults
	setDefault := true
	plt.Reset(setDefault, nil)

	// draw polyline
	plt.Polyline(P, stPlines)

	// draw circle
	plt.Circle(0, 4, 2.0, stCircles)

	// draw arrow
	plt.Arrow(-4, 2, 4, 7, stArrows)

	// draw arc
	plt.Arc(0, 4, 3, 0, 90, nil)

	// autoscale axes
	plt.AutoScale(P)

	// enforce same scales
	plt.Equal()

	// draw a _posteriori_ legend
	plt.LegendX([]*plt.A{
		{C: "red", M: "o", Ls: "-", Lw: 1, Ms: -1, L: "first", Me: -1},
		{C: "green", M: "s", Ls: "-", Lw: 2, Ms: 0, L: "second", Me: -1},
		{C: "blue", M: "+", Ls: "-", Lw: 3, Ms: 10, L: "third", Me: -1},
	}, nil)

	// save figure (default is PNG)
	plt.Save("/tmp/gosl", "plt_polygon01")
}
