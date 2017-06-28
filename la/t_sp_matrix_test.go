// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestSpMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpMatrix01. Setting CCMatrix")

	//          ↓     ↓        ↓           ↓  ↓     ↓
	//          0  1  2  3  4  5  6  7  8  9 10 11 12
	Ai := []int{0, 1, 0, 2, 4, 1, 2, 3, 4, 2, 1, 4}
	Ax := []float64{2, 3, 3, -1, 4, 4, -3, 1, 2, 2, 6, 1}
	Ap := []int{0, 2, 5, 9, 10, 12}
	var A CCMatrix
	A.Set(5, 5, Ap, Ai, Ax)
	Ad := A.ToDense()

	chk.Matrix(tst, "A", 1e-17, Ad.GetSlice(), [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 4, 2, 0, 1},
	})
}

func TestSpMatrix02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpMatrix02. PutMatAndMatT, PutCCMatAndMatT")

	var K, L, A Triplet
	K.Init(6, 6, 36+2*6) // 2*6 == number of nonzeros in A
	L.Init(6, 6, 36+2*6) // 2*6 == number of nonzeros in A
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			K.Put(i, j, 1000)
			L.Put(i, j, 1000)
		}
	}
	A.Init(2, 3, 6)
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			A.Put(i, j, float64(10*(i+1)+j+1))
		}
	}
	Am := A.ToMatrix(nil)
	K.PutMatAndMatT(&A)
	L.PutCCMatAndMatT(Am)
	Kaug := K.ToMatrix(nil).ToDense()
	Laug := L.ToMatrix(nil).ToDense()
	Cor := [][]float64{
		{1000, 1000, 1000, 1011, 1021, 1000},
		{1000, 1000, 1000, 1012, 1022, 1000},
		{1000, 1000, 1000, 1013, 1023, 1000},
		{1011, 1012, 1013, 1000, 1000, 1000},
		{1021, 1022, 1023, 1000, 1000, 1000},
		{1000, 1000, 1000, 1000, 1000, 1000},
	}
	chk.Matrix(tst, "Kaug", 1.0e-17, Kaug.GetSlice(), Cor)
	chk.Matrix(tst, "Laug", 1.0e-17, Laug.GetSlice(), Cor)
}
