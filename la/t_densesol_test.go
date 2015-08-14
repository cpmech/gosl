// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func calc_LLt(L [][]float64) (LLt [][]float64) {
	LLt = MatAlloc(len(L), len(L))
	for i := 0; i < len(L); i++ {
		for j := 0; j < len(L[0]); j++ {
			for k := 0; k < len(L); k++ {
				LLt[i][j] += L[i][k] * L[j][k]
			}
		}
	}
	return
}

func TestCholesky01(tst *testing.T) {

	chk.PrintTitle("TestCholesky 01")

	a := [][]float64{
		{25.0, 15.0, -5.0},
		{15.0, 18.0, 0.0},
		{-5.0, 0.0, 11.0},
	}

	L := MatAlloc(3, 3)
	Cholesky(L, a) // L is such as A = L * transp(L)
	PrintMat("a", a, "%6g", false)
	PrintMat("L", L, "%6g", false)
	LLt := calc_LLt(L)
	chk.Matrix(tst, "a = LLt", 1e-17, LLt, a)
	chk.Matrix(tst, "L", 1e-17, L, [][]float64{
		{5, 0, -0},
		{3, 3, 0},
		{-1, 1, 3},
	})
}

func TestCholesky02(tst *testing.T) {

	chk.PrintTitle("TestCholesky 02")

	a := [][]float64{
		{2, 1, 1, 3, 2},
		{1, 2, 2, 1, 1},
		{1, 2, 9, 1, 5},
		{3, 1, 1, 7, 1},
		{2, 1, 5, 1, 8},
	}
	L := MatAlloc(5, 5)
	Cholesky(L, a)
	PrintMat("a", a, "%6g", false)
	PrintMat("L", L, "%10.6f", false)
	chk.Matrix(tst, "a = LLt", 1e-15, calc_LLt(L), a)
	chk.Matrix(tst, "L", 1e-15, L, [][]float64{
		{math.Sqrt2, 0, 0, 0, 0},
		{1.0 / math.Sqrt2, math.Sqrt(3.0 / 2.0), 0, 0, 0},
		{1.0 / math.Sqrt2, math.Sqrt(3.0 / 2.0), math.Sqrt(7.0), 0, 0},
		{3.0 / math.Sqrt2, -1.0 / math.Sqrt(6.0), 0, math.Sqrt(7.0 / 3.0), 0},
		{math.Sqrt2, 0, 4.0 / math.Sqrt(7.0), -2.0 * math.Sqrt(3.0/7.0), math.Sqrt2},
	})
}

func TestSPDsolve01(tst *testing.T) {

	chk.PrintTitle("TestSPDsolve 01")

	a := [][]float64{
		{2, 1, 1, 3, 2},
		{1, 2, 2, 1, 1},
		{1, 2, 9, 1, 5},
		{3, 1, 1, 7, 1},
		{2, 1, 5, 1, 8},
	}
	b := []float64{-2, 4, 3, -5, 1}
	x := make([]float64, 5)
	SPDsolve(x, a, b)
	check_residR(tst, 1e-14, a, x, b)
	chk.Vector(tst, "x = inv(a) * b", 1e-13, x, []float64{
		-629.0 / 98.0,
		237.0 / 49.0,
		-53.0 / 49.0,
		62.0 / 49.0,
		23.0 / 14.0,
	})
}

func TestSPDsolve02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("TestSPDsolve 02")

	a := [][]float64{
		{2, 1, 1, 3, 2},
		{1, 2, 2, 1, 1},
		{1, 2, 9, 1, 5},
		{3, 1, 1, 7, 1},
		{2, 1, 5, 1, 8},
	}
	b := []float64{-2, 4, 3, -5, 1}
	B := []float64{24, 29, 110, 12, 102}
	x := make([]float64, 5)
	X := make([]float64, 5)
	SPDsolve2(x, X, a, b, B)
	check_residR(tst, 1e-14, a, x, b)
	check_residR(tst, 1e-14, a, X, B)
	chk.Vector(tst, "x = inv(a) * b", 1e-13, x, []float64{
		-629.0 / 98.0,
		237.0 / 49.0,
		-53.0 / 49.0,
		62.0 / 49.0,
		23.0 / 14.0,
	})
	chk.Vector(tst, "X = inv(a) * B", 1e-13, X, []float64{0, 4, 7, -1, 8})
}
