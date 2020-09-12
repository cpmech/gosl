// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"gosl/chk"
)

func TestBlas2tst01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas2tst01. (real) Blas2 functions")

	// allocate data
	a := NewMatrixDeep2([][]float64{
		{2, +3, +0, 0, 0},
		{3, +0, +4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, +0, +1, 0, 0},
		{0, +4, +2, 0, 1},
	})
	x := NewVectorSlice([]float64{1, 2, 3, 4, 5})

	// MatVecMul
	r := NewVector(5)
	MatVecMul(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x", 1e-17, r, []float64{8, 45, -3, 3, 19})

	// MatVecMul
	r.Fill(11234)
	MatVecMul(r, 1, a, x)
	chk.Array(tst, "r = 1⋅a⋅x (again)", 1e-17, r, []float64{8, 45, -3, 3, 19})

	// MatVecMul (larger)
	A := NewMatrix(10, 10)
	A.SetDiag(2)
	U := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	V := NewVector(10)
	MatVecMul(V, 0.5, A, U)
	chk.Array(tst, "V = 1⋅A⋅U (larger)", 1e-17, V, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	MatTrVecMul(V, 0.5, A, U)
	chk.Array(tst, "V = 1⋅Aᵀ⋅U (larger)", 1e-17, V, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

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

	// VecVecTrMul
	b := NewMatrix(5, 2)
	u := []float64{10, 20, 30, 40, 50}
	v := []float64{1, 2}
	VecVecTrMul(b, 0.5, u, v)
	chk.Deep2(tst, "b = 0.5⋅u⋅vᵀ", 1e-17, b.GetDeep2(), [][]float64{
		{+5, 10},
		{10, 20},
		{15, 30},
		{20, 40},
		{25, 50},
	})

	// VecVecTrMul (larger)
	B := NewMatrix(10, 2)
	VecVecTrMul(B, 2.0, U, v)
	chk.Deep2(tst, "B = 2⋅U⋅vᵀ", 1e-17, B.GetDeep2(), [][]float64{
		{2, 4},
		{4, 8},
		{6, 12},
		{8, 16},
		{10, 20},
		{12, 24},
		{14, 28},
		{16, 32},
		{18, 36},
		{20, 40},
	})
}

func TestBlas2tst02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas2tst02. (complex) Blas2 functions")

	// allocate data
	a := NewMatrixDeep2c([][]complex128{
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

func TestBlas2tst03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas2tst03. (real) Blas2 functions (tranpose)")

	// allocate data
	a := NewMatrixDeep2([][]float64{
		{1, 2, +0, 1, -1},
		{2, 3, -1, 1, +1},
		{1, 2, +0, 4, -1},
		{4, 0, +3, 1, +1},
	})
	u := NewVectorSlice([]float64{1, 2, 3, 4, 5})
	v := NewVectorSlice([]float64{1, 2, 3, 4})

	// MatVecMul
	x := NewVector(4)
	MatVecMul(x, 0.5, a, u)
	chk.Array(tst, "0.5⋅a⋅u", 1e-17, x, []float64{2, 7, 8, 11})

	// MatTrVecMul
	y := NewVector(5)
	MatTrVecMul(y, 0.5, a, v)
	chk.Array(tst, "0.5⋅aᵀ⋅v", 1e-17, y, []float64{12, 7, 5, 9.5, 1})
}
