// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"gosl/chk"
)

func TestInterp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Interp01. linear interpolation")

	xx := []float64{0, 1, 2, 3, 4, 5}
	yy := []float64{0.50, 0.20, 0.20, 0.05, 0.01, 0.00}

	o := NewDataInterp("lin", 1, xx, yy)

	for i, x := range xx {
		chk.Float64(tst, "P(xi)", 1e-17, o.P(x), yy[i])
	}

	xref := []float64{1.0 / 3.0, 2.5, 2.0 / 3.0, 1.1, 1.5, 3.5, 4.5}
	yref := []float64{0.4, 0.125, 0.3, 0.2, 0.2, 0.03, 0.005}
	for i, x := range xref {
		chk.Float64(tst, "P(xref)", 1e-16, o.P(x), yref[i])
	}
}

func TestInterp02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Interp02. polynomial interpolation")

	xx := []float64{0, 1, 2, 3, 4, 5}
	yy := []float64{0.50, 0.20, 0.20, 0.05, 0.01, 0.00}

	for _, p := range []int{1, 2, 3} {

		o := NewDataInterp("poly", p, xx, yy)

		for i, x := range xx {
			chk.Float64(tst, "P(xi)", 1e-17, o.P(x), yy[i])
		}

		if o.m == 2 {
			xref := []float64{1.0 / 3.0, 2.5, 2.0 / 3.0, 1.1, 1.5, 3.5, 4.5}
			yref := []float64{0.4, 0.125, 0.3, 0.2, 0.2, 0.03, 0.005}
			for i, x := range xref {
				chk.Float64(tst, "P(xref)", 1e-16, o.P(x), yref[i])
			}
		}
	}
}
