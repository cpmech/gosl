// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_conv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv01")

	var t Triplet
	t.Init(3, 3, 10)
	t.Put(0, 0, 5.0)
	t.Put(0, 0, 5.0)
	t.Put(0, 1, 11.0)
	t.Put(0, 2, 12.0)
	t.Put(1, 0, 20.0)
	t.Put(1, 1, 21.0)
	t.Put(1, 2, 22.0)
	t.Put(2, 0, 30.0)
	t.Put(2, 1, 31.0)
	t.Put(2, 2, 32.0)
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	if chk.Verbose {
		PrintMat("a", ad, "%5g", false)
	}
	chk.Matrix(tst, "a", 1e-17, ad, [][]float64{
		{10, 11, 12},
		{20, 21, 22},
		{30, 31, 32},
	})
}

func Test_conv02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv02")

	var t Triplet
	t.Init(4, 3, 4*3+2)
	t.Put(0, 0, 1.0)
	t.Put(0, 1, 2.0)
	t.Put(0, 2, 3.0)
	t.Put(1, 0, 4.0)
	t.Put(1, 1, 5.0)
	t.Put(1, 2, 6.0)
	t.Put(2, 0, 7.0)
	t.Put(2, 1, 8.0)
	t.Put(2, 2, 9.0)
	t.Put(3, 0, 4.0)
	t.Put(3, 1, 11.0)
	t.Put(3, 2, 12.0)
	t.Put(3, 0, 3.0) // repeated
	t.Put(3, 0, 3.0) // repeated
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	if chk.Verbose {
		PrintMat("a", ad, "%5g", false)
	}
	chk.Matrix(tst, "a", 1e-17, ad, [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})
}

func Test_conv03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv03")

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
	Ad := Am.ToDense()
	if chk.Verbose {
		Kd := K.ToMatrix(nil).ToDense()
		Ld := L.ToMatrix(nil).ToDense()
		PrintMat("K", Kd, "%8g", false)
		PrintMat("L", Ld, "%8g", false)
		PrintMat("A", Ad, "%8g", false)
	}
	K.PutMatAndMatT(&A)
	L.PutCCMatAndMatT(Am)
	Kaug := K.ToMatrix(nil).ToDense()
	Laug := L.ToMatrix(nil).ToDense()
	if chk.Verbose {
		PrintMat("K augmented", Kaug, "%8g", false)
		PrintMat("L augmented", Laug, "%8g", false)
	}
	Cor := [][]float64{
		{1000, 1000, 1000, 1011, 1021, 1000},
		{1000, 1000, 1000, 1012, 1022, 1000},
		{1000, 1000, 1000, 1013, 1023, 1000},
		{1011, 1012, 1013, 1000, 1000, 1000},
		{1021, 1022, 1023, 1000, 1000, 1000},
		{1000, 1000, 1000, 1000, 1000, 1000},
	}
	chk.Matrix(tst, "Kaug", 1.0e-17, Kaug, Cor)
	chk.Matrix(tst, "Laug", 1.0e-17, Laug, Cor)
}

func Test_conv04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv04")

	var t TripletC
	t.Init(3, 3, 10, false)
	t.Put(0, 0, 5.0, 0)
	t.Put(0, 0, 5.0, 0)
	t.Put(0, 1, 11.0, 0)
	t.Put(0, 2, 12.0, 0)
	t.Put(1, 0, 20.0, 0)
	t.Put(1, 1, 21.0, 0)
	t.Put(1, 2, 22.0, 0)
	t.Put(2, 0, 30.0, 0)
	t.Put(2, 1, 31.0, 0)
	t.Put(2, 2, 32.0, 666.0)
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	if chk.Verbose {
		PrintMatC("a", ad, "(%2g", " +%4gi)  ", false)
	}
	chk.MatrixC(tst, "a", 1.0e-17, ad, [][]complex128{
		{10, 11, 12},
		{20, 21, 22},
		{30, 31, 32 + 666i},
	})
}

func Test_conv05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv05")

	var t TripletC
	t.Init(4, 3, 4*3+2, false)
	t.Put(0, 0, 1.0, 1.0)
	t.Put(0, 1, 2.0, 1.0)
	t.Put(0, 2, 3.0, 2.0)
	t.Put(1, 0, 4.0, 2.0)
	t.Put(1, 1, 5.0, 2.0)
	t.Put(1, 2, 6.0, 1.0)
	t.Put(2, 0, 7.0, 3.0)
	t.Put(2, 1, 8.0, 3.0)
	t.Put(2, 2, 9.0, 4.0)
	t.Put(3, 0, 4.0, 1.1)
	t.Put(3, 1, 11.0, 4.0)
	t.Put(3, 2, 12.0, 3.0)
	t.Put(3, 0, 3.0, 1.4) // repeated
	t.Put(3, 0, 3.0, 1.5) // repeated
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	if chk.Verbose {
		PrintMatC("a", ad, "(%2g", " +%4gi)  ", false)
	}
	chk.MatrixC(tst, "a", 1.0e-17, ad, [][]complex128{
		{1 + 1i, 2 + 1i, 3 + 2i},
		{4 + 2i, 5 + 2i, 6 + 1i},
		{7 + 3i, 8 + 3i, 9 + 4i},
		{10 + 4i, 11 + 4i, 12 + 3i},
	})
}

func Test_conv06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv06")

	a := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	}
	a_ := MatAlloc(4, 3)
	am := MatToColMaj(a)
	aa := ColMajToMatNew(am, 4, 3)
	ColMajToMat(a_, am)
	io.Pforan("a  = %v\n", a)
	io.Pforan("am = %v\n", am)
	io.Pforan("aa = %v\n", aa)
	chk.Vector(tst, "a => am", 1e-17, am, []float64{1, 4, 7, 10, 2, 5, 8, 11, 3, 6, 9, 12})
	chk.Matrix(tst, "am => a", 1e-17, aa, a)
	chk.Matrix(tst, "am => a", 1e-17, a_, a)

	b := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}
	bm := MatToColMaj(b)
	bb := ColMajToMatNew(bm, 3, 4)
	io.Pforan("b  = %v\n", b)
	io.Pforan("bm = %v\n", bm)
	io.Pforan("bb = %v\n", bb)
	chk.Vector(tst, "b => bm", 1e-15, bm, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})
	chk.Matrix(tst, "bm => b", 1e-15, bb, b)
}

func Test_conv07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conv07")

	r := []float64{1, 2, 3, 4, 5, 6}
	c := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6}
	rc := RCtoComplex(r, c)
	chk.VectorC(tst, "rc", 1e-17, rc, []complex128{1 + 0.1i, 2 + 0.2i, 3 + 0.3i, 4 + 0.4i, 5 + 0.5i, 6 + 0.6i})

	R, C := ComplexToRC(rc)
	chk.Vector(tst, "r", 1e-17, R, r)
	chk.Vector(tst, "c", 1e-17, C, c)

	pa := []float64{1, 0.1, 2, 0.2, 3, 0.3, 4, 0.4, 5, 0.5, 6, 0.6}
	v := RCpairsToComplex(pa)
	chk.VectorC(tst, "p→v", 1e-17, v, []complex128{1 + 0.1i, 2 + 0.2i, 3 + 0.3i, 4 + 0.4i, 5 + 0.5i, 6 + 0.6i})

	pb := ComplexToRCpairs(v)
	chk.Vector(tst, "v→p", 1e-17, pb, pa)
}

func Test_spset01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("spset01")

	//          ↓     ↓        ↓           ↓  ↓     ↓
	//          0  1  2  3  4  5  6  7  8  9 10 11 12
	Ai := []int{0, 1, 0, 2, 4, 1, 2, 3, 4, 2, 1, 4}
	Ax := []float64{2, 3, 3, -1, 4, 4, -3, 1, 2, 2, 6, 1}
	Ap := []int{0, 2, 5, 9, 10, 12}
	var A CCMatrix
	A.Set(5, 5, Ap, Ai, Ax)
	Ad := A.ToDense()
	PrintMat("A", Ad, "%5g", false)

	chk.Matrix(tst, "A", 1e-17, Ad, [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 4, 2, 0, 1},
	})
}
