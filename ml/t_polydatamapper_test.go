// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestPolyDataMapper01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("PolyDataMapper01. simple mapping")

	// raw data
	XYraw := [][]float64{ //                ↓       ↓
		{10, 1, 0.1, -1, 100, -10}, // x0, x1, x2, x3, x4, y
		{20, 2, 0.2, -2, 200, -20},
		{30, 3, 0.3, -3, 300, -30},
		{40, 4, 0.4, -4, 400, -40},
		{50, 5, 0.5, -5, 500, -50},
	}

	// mapper
	nOriFeatures := len(XYraw[0]) - 1
	iFeature := 1
	jFeature := 3
	degree := 2
	mapper := NewPolyDataMapper(nOriFeatures, iFeature, jFeature, degree)

	// data
	data := mapper.GetMapped(XYraw, true)

	// check
	chk.Deep2(tst, "X", 1e-15, data.X.GetDeep2(), [][]float64{
		{10, 1, 0.1, -1, 100, +1, -1., +1}, // x0, x1, x2, x3, x4, (x1)², (x1)*(x3), (x3)²
		{20, 2, 0.2, -2, 200, +4, -4., +4},
		{30, 3, 0.3, -3, 300, +9, -9., +9},
		{40, 4, 0.4, -4, 400, 16, -16, 16},
		{50, 5, 0.5, -5, 500, 25, -25, 25},
	})
	chk.Array(tst, "y", 1e-15, data.Y, []float64{-10, -20, -30, -40, -50})
}
