// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestBlas3tst01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas3tst01. (real) Blas3 functions")

	// allocate data
	a := NewMatrixSlice([][]float64{
		{1.0, 2.00, 3.0},
		{0.5, 0.75, 1.5},
	})
	b := NewMatrixSlice([][]float64{
		{0.1, 0.5, 0.5, 0.75},
		{0.2, 2.0, 2.0, 2.0},
		{0.3, 0.5, 0.5, 0.5},
	})

	// MatMatMul
	c := NewMatrix(2, 4)
	c.Set(0, 0, 12333)
	c.Set(1, 3, 111)
	MatMatMul(c, 1, a, b)
	chk.Matrix(tst, "c := 1⋅a⋅b", 1e-15, c.GetSlice(), [][]float64{
		{1.40, 6.0, 6.0, 6.250},
		{0.65, 2.5, 2.5, 2.625},
	})

	// MatMatMulAdd
	c.Fill(100)
	MatMatMulAdd(c, 1, a, b)
	chk.Matrix(tst, "c := 1⋅a⋅b", 1e-17, c.GetSlice(), [][]float64{
		{101.40, 106.0, 106.0, 106.250},
		{100.65, 102.5, 102.5, 102.625},
	})
}
