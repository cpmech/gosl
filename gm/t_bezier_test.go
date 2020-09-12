// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"gosl/chk"
	"gosl/utl"
)

func Test_bezier01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bezier01. quadratic Bezier.")

	bez := BezierQuad{
		Q: [][]float64{
			{-1, 1},
			{0.5, -2},
			{2, 4},
		},
	}

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
		chk.Float64(tst, "y=y", 1e-15, Y[i], X[i]*X[i])
	}

	XX, YY, _ := bez.GetPoints(T)
	chk.Array(tst, "X", 1e-15, X, XX)
	chk.Array(tst, "Y", 1e-15, Y, YY)

	Xq, Yq, _ := bez.GetControlCoords()
	chk.Array(tst, "Xq", 1e-15, Xq, []float64{-1, 0.5, 2})
	chk.Array(tst, "Yq", 1e-15, Yq, []float64{1, -2, 4})
}
