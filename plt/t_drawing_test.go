// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"gosl/chk"
)

func Test_draw01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("draw01. 2d polygon/polyline")

	P := [][]float64{
		{-2.5, 0.0},
		{-5.5, 4.0},
		{0.0, 10.0},
		{5.5, 4.0},
		{2.5, 0.0},
	}

	Reset(true, nil)
	Polyline(P, &A{Fc: "#c1d7cf", Ec: "#4db38e", Lw: 4.5, Closed: true, NoClip: true})
	Circle(0, 4, 2.0, &A{Fc: "#b2cfa5", Ec: "#5dba35", Z: 1})
	Arrow(-4, 2, 4, 7, &A{Fc: "cyan", Ec: "blue", Z: 2, Scale: 50, Style: "fancy"})
	Arc(0, 4, 3, 0, 90, nil)
	AutoScale(P)
	Equal()

	if chk.Verbose {
		Save("/tmp/gosl", "t_draw01")
	}
}

func Test_draw02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("draw02. 3d polygon")

	P := [][]float64{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
	}

	Q := [][]float64{
		{1, 1, 0},
		{0, 1, 0},
		{0, 0, 0},
	}

	Reset(true, nil)
	Polygon3d(Q, &A{Fc: "#ace3ba", Ec: "#8700c6", Lw: 2})
	Polygon3d(P, nil)

	if chk.Verbose {
		Default3dView(0, 1, 0, 1, 0, 1, true)
		Save("/tmp/gosl", "t_draw02")
	}
}

func Test_draw03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("draw02. box and 3d points")

	Reset(true, nil)
	Plot3dPoint(0, -1.5, -0.9, &A{C: "grey", M: "o", Ms: 3000, Mec: "orange"})
	Plot3dPoint(-2, 0.5, -0.9, &A{C: "r", M: "*", Ms: 5000, Mec: "green", Void: true})
	Plot3dPoint(2.5, 1.5, -0.9, &A{C: "r", M: "s", Ms: 1000, Mec: "k", Void: false})
	Box(-0.5, 1, -1, 2, -3, 0, &A{A: 0.5, Lw: 3, Fc: "#5294ed", Ec: "#ffec4f", Wire: true})
	if chk.Verbose {
		Triad(1.5, "x", "y", "z", nil, nil)
		Default3dView(-1, 1.5, -1.5, 2.5, -3.5, 0.5, true)
		Save("/tmp/gosl", "t_draw03")
	}
}
