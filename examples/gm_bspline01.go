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

	// order of B-spline
	p := 3

	// knots for clamped curve
	startT := make([]float64, p+1) // p+1 zeros
	endT := utl.Ones(p + 1)        // p+1 ones

	// knots
	T1 := append(append(startT, utl.LinSpace(0.1, 0.9, 9)...), endT...)

	// B-spline
	b1 := gm.NewBspline(T1, p)

	// set control points
	b1.SetControl([][]float64{
		{0.5, 0.5},
		{1.0, 0.5},
		{1.0, 0.5}, // repeated
		{1.0, 0.5}, // repeated => 3x => discontinuity
		{1.0, 0.2},
		{0.7, 0.0},
		{0.3, 0.0},
		{0.0, 0.3},
		{0.0, 0.7},
		{0.3, 1.0},
		{0.7, 1.0},
		{0.9, 0.9},
		{1.0, 0.8},
	})

	// configuration
	withCtrl := true
	argsCurve := &plt.A{C: "r", Lw: 10, L: "curve", NoClip: true}
	argsCtrl := &plt.A{C: "k", M: ".", L: "control", NoClip: true}

	// plot
	np := 101
	plt.Reset(false, nil)
	b1.Draw2d(np, 0, withCtrl, argsCurve, argsCtrl)
	plt.Equal()
	plt.HideAllBorders()
	plt.AxisXmax(1.0)
	plt.Save("/tmp/gosl", "gm_bspline01")
}
