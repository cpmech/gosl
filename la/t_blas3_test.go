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
	chk.PrintTitle("Blas3tst01. (real) MatMat mul functions")

	// allocate data
	a := NewMatrixDeep2([][]float64{ // 2 x 3
		{1.0, 2.00, 3.0},
		{0.5, 0.75, 1.5},
	})
	b := NewMatrixDeep2([][]float64{ // 3 x 4
		{0.1, 0.5, 0.5, 0.75},
		{0.2, 2.0, 2.0, 2.0},
		{0.3, 0.5, 0.5, 0.5},
	})
	cref := [][]float64{
		{2.80, 12.0, 12.0, 12.50},
		{1.30, +5.0, +5.0, +5.25},
	}

	// tranpose matrices
	atrans := a.GetTranspose()
	btrans := b.GetTranspose()

	// MatMatMul
	c := NewMatrix(2, 4)
	c.Set(0, 0, 12333) // noise
	c.Set(1, 3, 111)   // noise
	MatMatMul(c, 2, a, b)
	chk.Deep2(tst, "c := 2⋅a⋅b", 1e-15, c.GetDeep2(), cref)

	// MatMatTrMul
	c.Set(0, 0, 12333) // noise
	c.Set(1, 3, 111)   // noise
	MatTrMatMul(c, 2, atrans, b)
	chk.Deep2(tst, "c := 2⋅aᵀ⋅b", 1e-15, c.GetDeep2(), cref)

	// MatMatTrMul
	c.Set(0, 0, 12333) // noise
	c.Set(1, 3, 111)   // noise
	MatMatTrMul(c, 2, a, btrans)
	chk.Deep2(tst, "c := 2⋅a⋅bᵀ", 1e-15, c.GetDeep2(), cref)

	// MatTrMatTrMul
	c.Set(0, 0, 12333) // noise
	c.Set(1, 3, 111)   // noise
	MatTrMatTrMul(c, 2, atrans, btrans)
	chk.Deep2(tst, "c := 2⋅aᵀ⋅bᵀ", 1e-15, c.GetDeep2(), cref)
}

func TestBlas3tst02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Blas3tst02. (real) MatMat mul + add functions")

	// allocate data
	a := NewMatrixDeep2([][]float64{ // 2 x 3
		{1.0, 2.00, 3.0},
		{0.5, 0.75, 1.5},
	})
	b := NewMatrixDeep2([][]float64{ // 3 x 4
		{0.1, 0.5, 0.5, 0.75},
		{0.2, 2.0, 2.0, 2.0},
		{0.3, 0.5, 0.5, 0.5},
	})
	cref := [][]float64{
		{102.80, 112.0, 112.0, 112.50},
		{101.30, 105.0, 105.0, 105.25},
	}

	// tranpose matrices
	atrans := a.GetTranspose()
	btrans := b.GetTranspose()

	// MatMatMulAdd
	c := NewMatrix(2, 4)
	c.Fill(100)
	MatMatMulAdd(c, 2, a, b)
	chk.Deep2(tst, "c := 2⋅a⋅b + 100", 1e-17, c.GetDeep2(), cref)

	// MatMatTrMulAdd
	c.Fill(100)
	MatTrMatMulAdd(c, 2, atrans, b)
	chk.Deep2(tst, "c := 2⋅aᵀ⋅b + 100", 1e-17, c.GetDeep2(), cref)

	// MatMatTrMulAdd
	c.Fill(100)
	MatMatTrMulAdd(c, 2, a, btrans)
	chk.Deep2(tst, "c := 2⋅a⋅bᵀ + 100", 1e-17, c.GetDeep2(), cref)

	// MatTrMatTrMulAdd
	c.Fill(100)
	MatTrMatTrMulAdd(c, 2, atrans, btrans)
	chk.Deep2(tst, "c := 2⋅aᵀ⋅bᵀ + 100", 1e-17, c.GetDeep2(), cref)
}
