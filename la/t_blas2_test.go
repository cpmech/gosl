// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestBlas2tst01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas2tst01. (real) Blas2 functions")

	// allocate data
	a := NewMatrixSlice([][]float64{
		{2, +3, +0, 0, 0},
		{3, +0, +4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, +0, +1, 0, 0},
		{0, +4, +2, 0, 1},
	})
	x := Vector([]float64{1, 2, 3, 4, 5})

	// MatVecMul
	r := NewVector(5)
	MatVecMul(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x", 1e-17, r, []float64{8, 45, -3, 3, 19})

	// MatVecMul
	r.Fill(11234)
	MatVecMul(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x (again)", 1e-17, r, []float64{8, 45, -3, 3, 19})

	// MatVecMulAdd
	r.Fill(0)
	MatVecMulAdd(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x + 0", 1e-17, r, []float64{8, 45, -3, 3, 19})

	// CopyInto
	r.Apply(-1, []float64{8, 45, -3, 3, 19}) // r := -b
	chk.Array(tst, "r := -b", 1e-17, r, []float64{-8, -45, 3, -3, -19})

	// MatVecMulAdd
	MatVecMulAdd(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x - b", 1e-17, r, nil)
}

func TestBlas2tst02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas2tst02. (complex) Blas2 functions")

	// allocate data
	a := NewMatrixSliceC([][]complex128{
		{2, +3, +0, 0, 0},
		{3, +0, +4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, +0, +1, 0, 0},
		{0, +4, +2, 0, 1},
	})
	x := VectorC([]complex128{1, 2, 3, 4, 5})

	// MatVecMul
	r := NewVectorC(5)
	MatVecMulC(r, 1, a, x)
	chk.ArrayC(tst, "r = 1⋅a⋅x", 1e-17, r, []complex128{8, 45, -3, 3, 19})

	// MatVecMul
	r.Fill(11234)
	MatVecMulC(r, 1, a, x)
	chk.ArrayC(tst, "r = 1⋅a⋅x (again)", 1e-17, r, []complex128{8, 45, -3, 3, 19})

	// MatVecMulAdd
	r.Fill(0)
	MatVecMulAddC(r, 1, a, x)
	chk.ArrayC(tst, "r = 1⋅a⋅x + 0", 1e-17, r, []complex128{8, 45, -3, 3, 19})

	// CopyInto
	r.Apply(-1, []complex128{8, 45, -3, 3, 19}) // r := -b
	chk.ArrayC(tst, "r := -b", 1e-17, r, []complex128{-8, -45, 3, -3, -19})

	// MatVecMulAdd
	MatVecMulAddC(r, 1, a, x)
	chk.ArrayC(tst, "r = 1⋅a⋅x - b", 1e-17, r, nil)
}
