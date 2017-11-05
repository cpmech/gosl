// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func TestVecDot(tst *testing.T) {
	verbose()
	chk.PrintTitle("NaiveVecDot")
	u := []float64{1, 2, 3}
	v := []float64{3, 2, 1}
	chk.Float64(tst, "u・v", 1e-17, NaiveVecDot(u, v), 10.0)
}

func TestVecAdd(tst *testing.T) {
	verbose()
	chk.PrintTitle("NaiveVecAdd")
	u := []float64{1, 2, 3}
	v := []float64{3, 2, 1}
	w := make([]float64, len(u))
	NaiveVecAdd(w, 1, u, -2, v)
	chk.Array(tst, "w := 1⋅u - 2⋅v", 1e-17, w, []float64{-5, -2, 1})
}

func TestMatVecMul(tst *testing.T) {
	verbose()
	chk.PrintTitle("NaiveMatVecMul")
	a := [][]float64{
		{2, +3, +0, 0, 0},
		{3, +0, +4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, +0, +1, 0, 0},
		{0, +4, +2, 0, 1},
	}
	x := []float64{1, 2, 3, 4, 5}
	r := make([]float64, len(x))
	NaiveMatVecMul(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x", 1e-17, r, []float64{8, 45, -3, 3, 19})
}

func TestMatMatMul(tst *testing.T) {
	verbose()
	chk.PrintTitle("NaiveMatMatMul")
	a := [][]float64{ // 2 x 3
		{1.0, 2.00, 3.0},
		{0.5, 0.75, 1.5},
	}
	b := [][]float64{ // 3 x 4
		{0.1, 0.5, 0.5, 0.75},
		{0.2, 2.0, 2.0, 2.0},
		{0.3, 0.5, 0.5, 0.5},
	}
	cref := [][]float64{
		{2.80, 12.0, 12.0, 12.50},
		{1.30, +5.0, +5.0, +5.25},
	}
	c := utl.Alloc(2, 4)
	NaiveMatMatMul(c, 2, a, b)
	chk.Deep2(tst, "c := 2⋅a⋅b", 1e-15, c, cref)
}
