// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestEqs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs01")

	// some prescribed
	e := NewEquations(9, []int{0, 6, 3})
	e.Info(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{1, 2, 4, 5, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, 0, 1, -1, 2, 3, -1, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 3, 6})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, -1, -1, 1, -1, -1, 2, -1, -1})

	// some prescribed
	io.Pl()
	e = NewEquations(9, []int{0, 2, 1})
	e.Info(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{3, 4, 5, 6, 7, 8})
	chk.Ints(tst, "FtoU", e.FtoU, []int{-1, -1, -1, 0, 1, 2, 3, 4, 5})
	chk.Ints(tst, "KtoF", e.KtoF, []int{0, 1, 2})
	chk.Ints(tst, "FtoK", e.FtoK, []int{0, 1, 2, -1, -1, -1, -1, -1, -1})

	// none prescribed
	io.Pl()
	e = NewEquations(5, nil)
	e.Info(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "FtoU", e.FtoU, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "KtoF", e.KtoF, nil)
	chk.Ints(tst, "FtoK", e.FtoK, []int{-1, -1, -1, -1, -1})

	// most prescribed
	io.Pl()
	e = NewEquations(5, []int{1, 2, 3, 4})
	e.Info(true)
	chk.Ints(tst, "UtoF", e.UtoF, []int{0})
	chk.Ints(tst, "FtoU", e.FtoU, []int{0, -1, -1, -1, -1})
	chk.Ints(tst, "KtoF", e.KtoF, []int{1, 2, 3, 4})
	chk.Ints(tst, "FtoK", e.FtoK, []int{-1, 0, 1, 2, 3})
}

func TestEqs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs02. Put items in partitioned sparse A matrix")

	// equations classifier
	e := NewEquations(5, []int{4, 2})
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
	e := NewEquations(A.M, []int{4, 2})
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

func TestEqs04a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs04a. errors due to kx")

	defer chk.RecoverTstPanicIsOK(tst)

	e := NewEquations(0, nil)
	io.Pf("equations = %v+#\n", e)
}

func TestEqs04b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs04b. errors due to kx")

	defer chk.RecoverTstPanicIsOK(tst)

	e := NewEquations(5, []int{0, 4, 2, 1, 3, 5})
	io.Pf("equations = %v+#\n", e)
}

func TestEqs04c(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs04c. errors due to kx")

	defer chk.RecoverTstPanicIsOK(tst)

	e := NewEquations(9, []int{0, 4, 8, 3, 7, 11})
	io.Pf("equations = %v+#\n", e)
}

func TestEqs04d(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs04d. errors due to kx")

	defer chk.RecoverTstPanicIsOK(tst)

	e := NewEquations(15, []int{0, 4, 8, 12, 3, 7, 11, 15, 0, 1, 2, 3, 12, 13, 14, 15})
	io.Pf("equations = %v+#\n", e)
}

func TestEqs05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs05. Solve()")

	/*
	          Auu       Auk  ⋅  xu ?  =  bu ✓
	          Aku       Akk     xk ✓     bk ?

	    1   3  -2  │ 11   0     x0        5          x0 = -15
	    3   5   6  │  4   1     x1        7     ⇒    x1 =   8
	    2   4   3  │ -3   0  ⋅  x2    =   8          x2 =   2
	    ———————————│———————     ——       ——
	   -1   2   3  │  2   1      0       b3     ⇒    b3 =  37
	    4   1  -3  │  5   0      0       b4          b4 = -58
	*/

	// system
	n := 5                   // total number of equations
	kx := []int{3, 4}        // known equations
	nnz := []int{9, 4, 6, 3} // number of non-zeros in Auu, Auk, Aku, Akk
	kparts := true           // also allocates Aku and Akk
	vectors := true          // also allocates B and X vectors
	e := NewEquations(n, kx)
	e.Alloc(nnz, kparts, vectors)

	// Auu
	e.Start()
	e.Put(0, 0, 1.0)
	e.Put(0, 1, 3.0)
	e.Put(0, 2, -2.0)
	e.Put(1, 0, 3.0)
	e.Put(1, 1, 5.0)
	e.Put(1, 2, 6.0)
	e.Put(2, 0, 2.0)
	e.Put(2, 1, 4.0)
	e.Put(2, 2, 3.0)

	// Auk
	e.Put(0, 3, 11.0)
	e.Put(1, 3, 4.0)
	e.Put(1, 4, 1.0)
	e.Put(2, 3, -3.0)

	// Aku
	e.Put(3, 0, -1.0)
	e.Put(3, 1, 2.0)
	e.Put(3, 2, 3.0)
	e.Put(4, 0, 4.0)
	e.Put(4, 1, 1.0)
	e.Put(4, 2, -3.0)

	// Akk
	e.Put(3, 3, 2.0)
	e.Put(3, 4, 1.0)
	e.Put(4, 3, 5.0)

	// functions
	calcXk := func(I int, t float64) float64 {
		return 0.0 // returns zero for all known values: x3 and x4
	}
	calcBu := func(I int, t float64) float64 {
		switch I {
		case 0:
			return 5.0
		case 1:
			return 7.0
		case 2:
			return 8.0
		}
		chk.Panic("I=%d does not correspond to unknown value\n", I)
		return 0 // unreachable
	}

	// solve linear system
	e.SolveOnce(calcXk, calcBu)

	// check
	chk.Array(tst, "{xu}", 1e-13, e.Xu, []float64{-15, 8, 2})
	chk.Array(tst, "{xk}", 1e-13, e.Xk, []float64{0, 0})
	chk.Array(tst, "{bu}", 1e-13, e.Bu, []float64{5, 7, 8})
	chk.Array(tst, "{bk}", 1e-12, e.Bk, []float64{37, -58})

	// check full system
	b := NewVector(n)
	bRef := NewVector(n)
	buRef := NewVector(e.Nu)
	bkRef := NewVector(e.Nk)
	SpMatVecMul(buRef, 1.0, e.Auu.ToMatrix(nil), e.Xu)
	SpMatVecMulAdd(buRef, 1.0, e.Auk.ToMatrix(nil), e.Xk)
	SpMatVecMul(bkRef, 1.0, e.Aku.ToMatrix(nil), e.Xu)
	SpMatVecMulAdd(bkRef, 1.0, e.Akk.ToMatrix(nil), e.Xk)
	e.JoinVector(b, e.Bu, e.Bk)
	e.JoinVector(bRef, buRef, bkRef)
	chk.Array(tst, "{b}", 1e-12, bRef, b)
}

func TestEqs06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Eqs06. Solve()")

	/*
		  k   u  u   k   u
		k 1 │ 2  1 │ 2 │ 1     -7.5       b0          b0 =  1.0
		  ——│——————│———│——     ————       ——
		u 2 │ 4  4 │ 6 │ 1       x1        2     ⇒    x1 =  2.5
		u 3 │ 6  1 │ 4 │ 5       x2        4          x2 = -1.5
		  ——│——————│———│——  ⋅  ————   =   ——
		k 1 │ 2  3 │ 5 │ 1        2       b3     ⇒    b3 =  4.0
		  ——│——————│———│——     ————       ——
		u 1 │ 1  0 │ 3 │ 0       x4        1          x4 =  1.0
	*/

	// system
	n := 5                   // total number of equations
	kx := []int{0, 3}        // known equations
	nnz := []int{7, 6, 6, 4} // number of non-zeros in Auu, Auk, Aku, Akk
	kparts := true           // also allocates Aku and Akk
	vectors := true          // also allocates B and X vectors
	e := NewEquations(n, kx)
	e.Alloc(nnz, kparts, vectors)

	// Auu
	e.Start()
	e.Put(1, 1, 4.0)
	e.Put(1, 2, 4.0)
	e.Put(1, 4, 1.0)
	e.Put(2, 1, 6.0)
	e.Put(2, 2, 1.0)
	e.Put(2, 4, 5.0)
	e.Put(4, 1, 1.0)

	// Auk
	e.Put(1, 0, 2.0)
	e.Put(1, 3, 6.0)
	e.Put(2, 0, 3.0)
	e.Put(2, 3, 4.0)
	e.Put(4, 0, 1.0)
	e.Put(4, 3, 3.0)

	// Aku
	e.Put(0, 1, 2.0)
	e.Put(0, 2, 1.0)
	e.Put(0, 4, 1.0)
	e.Put(3, 1, 2.0)
	e.Put(3, 2, 3.0)
	e.Put(3, 4, 1.0)

	// Akk
	e.Put(0, 0, 1.0)
	e.Put(0, 3, 2.0)
	e.Put(3, 0, 1.0)
	e.Put(3, 3, 5.0)

	// check A
	A := e.GetAmat()
	D := A.ToMatrix(nil).ToDense()
	io.Pf("%v\n", D.Print("%4g"))
	chk.Deep2(tst, "D", 1e-17, D.GetDeep2(), [][]float64{
		{1, 2, 1, 2, 1},
		{2, 4, 4, 6, 1},
		{3, 6, 1, 4, 5},
		{1, 2, 3, 5, 1},
		{1, 1, 0, 3, 0},
	})

	// functions
	calcXk := func(I int, t float64) float64 {
		switch I {
		case 0:
			return -7.5
		case 3:
			return 2.0
		}
		chk.Panic("I=%d does not correspond to known value\n", I)
		return 0 // unreachable
	}
	calcBu := func(I int, t float64) float64 {
		switch I {
		case 1:
			return 2.0
		case 2:
			return 4.0
		case 4:
			return 1.0
		}
		chk.Panic("I=%d does not correspond to unknown value\n", I)
		return 0 // unreachable
	}

	// solve linear system
	e.SolveOnce(calcXk, calcBu)

	// calc {bu} because SolveOnce will return with modified {bu}
	e.Bu[0] = calcBu(1, 0)
	e.Bu[1] = calcBu(2, 0)
	e.Bu[2] = calcBu(4, 0)

	// check
	chk.Array(tst, "{xu}", 1e-13, e.Xu, []float64{2.5, -1.5, 1.0})
	chk.Array(tst, "{xk}", 1e-13, e.Xk, []float64{-7.5, 2.0})
	chk.Array(tst, "{bu}", 1e-13, e.Bu, []float64{2.0, 4.0, 1.0})
	chk.Array(tst, "{bk}", 1e-12, e.Bk, []float64{1.0, 4.0})

	// check full system
	b := NewVector(n)
	bRef := NewVector(n)
	buRef := NewVector(e.Nu)
	bkRef := NewVector(e.Nk)
	SpMatVecMul(buRef, 1.0, e.Auu.ToMatrix(nil), e.Xu)
	SpMatVecMulAdd(buRef, 1.0, e.Auk.ToMatrix(nil), e.Xk)
	SpMatVecMul(bkRef, 1.0, e.Aku.ToMatrix(nil), e.Xu)
	SpMatVecMulAdd(bkRef, 1.0, e.Akk.ToMatrix(nil), e.Xk)
	e.JoinVector(b, e.Bu, e.Bk)
	e.JoinVector(bRef, buRef, bkRef)
	chk.Array(tst, "{b}", 1e-12, bRef, b)
}
