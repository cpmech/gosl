// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_bezier01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bezier01. quadratic Bezier. point-distance")

	bez := BezierQuad{
		Q: [][]float64{
			{-1, 1},
			{0.5, -2},
			{2, 4},
		},
	}

	doplot := false

	np := 21
	T := utl.LinSpace(0, 1, np)
	X := make([]float64, np)
	Y := make([]float64, np)
	X2 := utl.LinSpace(-1.0, 2.0, np)
	Y2 := make([]float64, np)
	C := make([]float64, 2)
	for i, t := range T {
		bez.Point(C, t)
		X[i] = C[0]
		Y[i] = C[1]
		Y2[i] = X2[i] * X2[i]
		chk.Scalar(tst, "y=y", 1e-15, Y[i], X[i]*X[i])
	}
	if doplot {
		plt.SetForPng(1, 400, 200)
		plt.Plot(X2, Y2, "'y-', lw=4,label='y=x*x'")
		plt.Plot(X, Y, "'b-', marker='.', label='Bezier'")
		plt.Gll("x", "y", "")
		plt.Equal()
		plt.SaveD("/tmp", "fig_gm_bezier01.png")
	}
}
