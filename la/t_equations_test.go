// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestEqs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs01")

	// some prescribed
	var e Equations
	e.Init(9, []int{0, 6, 3})
	e.Print(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{1, 2, 4, 5, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, 0, 1, -1, 2, 3, -1, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 3, 6})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, -1, -1, 1, -1, -1, 2, -1, -1})

	// some prescribed
	io.Pl()
	e.Init(9, []int{0, 2, 1})
	e.Print(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, -1, -1, 0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 1, 2})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, 1, 2, -1, -1, -1, -1, -1, -1})

	// none prescribed
	io.Pl()
	e.Init(5, nil)
	e.Print(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "FtoU", e.FtoU, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "KtoF", e.KtoF, nil)
	chk.Ints(tst, "FtoK", e.FtoK, []int{-1, -1, -1, -1, -1})

	// all prescribed
	io.Pl()
	e.Init(5, []int{0, 1, 2, 3, 4})
	e.Print(true)
	chk.Ints(tst, "UtoF", e.UtoF, nil)
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, -1, -1, -1, -1})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, 1, 2, 3, 4})
}

func TestEqs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs02. Put items in partitioned sparse A matrix")

	// equations classifier
	var e Equations
	e.Init(5, []int{4, 2})
	chk.Ints(tst, "UtoF", e.UtoF, []int{0, 1, 3})
	chk.Ints(tst, "KtoF", e.KtoF, []int{2, 4})

	// allocate partitioned A matrix
	e.Alloc(nil, true, true)

	// assemble A matrix
	// 11  12  13  14  15    0
	// 21  22  23  24  25    1
	// 31  32  33  34  35    2 ◀ known
	// 41  42  43  44  45    3
	// 51  52  53  54  55    4 ◀ known
	//  0   1   2   3   4
	//          ▲       ▲
	e.Start()
	e.Put(0, 0, 11) // 0
	e.Put(0, 1, 12) // 1
	e.Put(0, 2, 13) // 2
	e.Put(0, 3, 14) // 3
	e.Put(0, 4, 15) // 4
	e.Put(1, 0, 21) // 5
	e.Put(1, 1, 22) // 6
	e.Put(1, 2, 23) // 7
	e.Put(1, 3, 24) // 8
	e.Put(1, 4, 25) // 9
	e.Put(2, 0, 31) // 10
	e.Put(2, 1, 32) // 11
	e.Put(2, 2, 33) // 12
	e.Put(2, 3, 34) // 13
	e.Put(2, 4, 35) // 14
	e.Put(3, 0, 41) // 15
	e.Put(3, 1, 42) // 16
	e.Put(3, 2, 43) // 17
	e.Put(3, 3, 44) // 18
	e.Put(3, 4, 45) // 19
	e.Put(4, 0, 51) // 20
	e.Put(4, 1, 52) // 21
	e.Put(4, 2, 53) // 22
	e.Put(4, 3, 54) // 23
	e.Put(4, 4, 55) // 24

	// check
	chk.Deep2(tst, "Auu", 1e-17, e.Auu.ToDense().GetDeep2(), [][]float64{
		{11, 12, 14},
		{21, 22, 24},
		{41, 42, 44},
	})
	chk.Deep2(tst, "Auk", 1e-17, e.Auk.ToDense().GetDeep2(), [][]float64{
		{13, 15},
		{23, 25},
		{43, 45},
	})
	chk.Deep2(tst, "Aku", 1e-17, e.Aku.ToDense().GetDeep2(), [][]float64{
		{31, 32, 34},
		{51, 52, 54},
	})
	chk.Deep2(tst, "Akk", 1e-17, e.Akk.ToDense().GetDeep2(), [][]float64{
		{33, 35},
		{53, 55},
	})

	// check lengths of vectors
	chk.Int(tst, "len(Bu)", len(e.Bu), 3)
	chk.Int(tst, "len(Bk)", len(e.Bk), 2)
	chk.Int(tst, "len(Xu)", len(e.Xu), 3)
	chk.Int(tst, "len(Xk)", len(e.Xk), 2)

	// set vectors
	e.Bu[0] = 100
	e.Bu[1] = 101
	e.Bu[2] = 103
	e.Bk[0] = 102
	e.Bk[1] = 104
	b := NewVector(5)
	e.JoinVector(b, e.Bu, e.Bk)
	chk.Array(tst, "b=join(bu,bk)", 1e-17, b, []float64{100, 101, 102, 103, 104})

	vu := NewVector(3)
	vk := NewVector(2)
	e.SplitVector(vu, vk, b)
	chk.Array(tst, "vu=split(b)_u", 1e-17, vu, []float64{100, 101, 103})
	chk.Array(tst, "vk=split(b)_k", 1e-17, vk, []float64{102, 104})
}

func TestEqs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs03. Split dense A matrix")

	// A matrix
	A := NewMatrixDeep2([][]float64{
		{11, 12, 13, 14, 15}, // 0
		{21, 22, 23, 24, 25}, // 1
		{31, 32, 33, 34, 35}, // 2 ◀ known
		{41, 42, 43, 44, 45}, // 3
		{51, 52, 53, 54, 55}, // 4 ◀ known
	}) // 0   1   2   3   4
	//            ▲       ▲

	// equations classifier
	var e Equations
	e.Init(A.M, []int{4, 2})
	chk.Ints(tst, "UtoF", e.UtoF, []int{0, 1, 3})
	chk.Ints(tst, "KtoF", e.KtoF, []int{2, 4})

	// split
	e.SetDense(A, true)
	chk.Deep2(tst, "Auu", 1e-17, e.Duu.GetDeep2(), [][]float64{
		{11, 12, 14},
		{21, 22, 24},
		{41, 42, 44},
	})
	chk.Deep2(tst, "Auk", 1e-17, e.Duk.GetDeep2(), [][]float64{
		{13, 15},
		{23, 25},
		{43, 45},
	})
	chk.Deep2(tst, "Aku", 1e-17, e.Dku.GetDeep2(), [][]float64{
		{31, 32, 34},
		{51, 52, 54},
	})
	chk.Deep2(tst, "Akk", 1e-17, e.Dkk.GetDeep2(), [][]float64{
		{33, 35},
		{53, 55},
	})
}
