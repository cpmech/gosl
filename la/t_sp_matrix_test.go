// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestTriplet01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpTriplet01")

	//   0 2 0 0
	//   1 0 4 0
	//   0 0 0 5
	//   0 3 0 6
	a := NewTriplet(4, 4, 6)
	a.Put(1, 0, 1)
	a.Put(0, 1, 2)
	a.Put(3, 1, 3)
	a.Put(1, 2, 4)
	a.Put(2, 3, 5)
	a.Put(3, 3, 6)

	l := a.ToDense().Print("%2g")
	io.Pf("%v\n", l)
	chk.String(tst, l, " 0 2 0 0\n 1 0 4 0\n 0 0 0 5\n 0 3 0 6")

	a.ToMatrix(nil).WriteSmat("/tmp/gosl/la", "triplet01", 0)
	d := io.ReadFile("/tmp/gosl/la/triplet01.smat")
	io.Pforan("d = %v\n", string(d))
	smat1 := "4  4  6\n  1  0    1.000000000000000e+00\n  0  1    2.000000000000000e+00\n  3  1    3.000000000000000e+00\n  1  2    4.000000000000000e+00\n  2  3    5.000000000000000e+00\n  3  3    6.000000000000000e+00\n"
	chk.String(tst, string(d), smat1)

	b := new(Triplet)
	b.ReadSmat("/tmp/gosl/la/triplet01.smat")
	chk.Deep2(tst, "b=a", 1e-17, a.ToDense().GetDeep2(), b.ToDense().GetDeep2())

	a.WriteSmat("/tmp/gosl/la", "triplet01b", 0)
	db := io.ReadFile("/tmp/gosl/la/triplet01b.smat")
	chk.String(tst, string(db), smat1)
}

func TestTriplet02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpTriplet02. (complex version)")

	//   0    2+2i 0    0
	//   1+1i 0    4+4i 0
	//   0    0    0    5-5i
	//   0    3-3i 0    6+6i
	a := NewTripletC(4, 4, 6)
	a.Put(1, 0, 1+1i)
	a.Put(0, 1, 2+2i)
	a.Put(3, 1, 3-3i)
	a.Put(1, 2, 4+4i)
	a.Put(2, 3, 5-5i)
	a.Put(3, 3, 6+6i)

	l := a.ToDense().Print("%2g", "%+2g")
	io.Pf("%v\n", l)
	chk.String(tst, l, " 0+0i  2+2i  0+0i  0+0i\n 1+1i  0+0i  4+4i  0+0i\n 0+0i  0+0i  0+0i  5-5i\n 0+0i  3-3i  0+0i  6+6i")

	am := a.ToMatrix(nil)
	am.WriteSmat("/tmp/gosl/la", "triplet02", 0)
	d := io.ReadFile("/tmp/gosl/la/triplet02.smat")
	io.Pforan("d = %v\n", string(d))
	smat1 := "4  4  6\n  1  0    1.000000000000000e+00  +1.000000000000000e+00\n  0  1    2.000000000000000e+00  +2.000000000000000e+00\n  3  1    3.000000000000000e+00  -3.000000000000000e+00\n  1  2    4.000000000000000e+00  +4.000000000000000e+00\n  2  3    5.000000000000000e+00  -5.000000000000000e+00\n  3  3    6.000000000000000e+00  +6.000000000000000e+00\n"
	chk.String(tst, string(d), smat1)

	b := new(TripletC)
	b.ReadSmat("/tmp/gosl/la/triplet02.smat")
	chk.Deep2c(tst, "b=a", 1e-17, a.ToDense().GetDeep2(), b.ToDense().GetDeep2())

	am.WriteSmatAbs("/tmp/gosl/la", "triplet02b", 0)
	c := new(Triplet)
	c.ReadSmat("/tmp/gosl/la/triplet02b.smat")
	chk.Deep2(tst, "b=a", 1e-14, c.ToDense().GetDeep2(), [][]float64{
		{0, math.Sqrt(8), 0, 0},
		{math.Sqrt2, 0, math.Sqrt(32), 0},
		{0, 0, 0, math.Sqrt(50)},
		{0, math.Sqrt(18), 0, math.Sqrt(72)},
	})

	a.WriteSmat("/tmp/gosl/la", "triplet02b", 0)
	db := io.ReadFile("/tmp/gosl/la/triplet02b.smat")
	chk.String(tst, string(db), smat1)
}

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

	chk.Deep2(tst, "A", 1e-17, Ad.GetDeep2(), [][]float64{
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
	chk.Deep2(tst, "Kaug", 1.0e-17, Kaug.GetDeep2(), Cor)
	chk.Deep2(tst, "Laug", 1.0e-17, Laug.GetDeep2(), Cor)
}

func TestSmat01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Smat01. read/write .smat file")

	correct := [][]float64{
		{2, 3, 0, 0, 0},
		{3, 0, 4, 0, 6},
		{0, -1, -3, 2, 0},
		{0, 0, 1, 0, 0},
		{0, 4, 2, 0, 1},
	}

	var T Triplet
	T.ReadSmat("data/small-sparse-matrix.mtx", true)
	chk.Deep2(tst, "T", 1e-17, T.ToDense().GetDeep2(), correct)

	T.WriteSmat("/tmp/gosl/la", "small-test-matrix", 1e-17)
	var S Triplet
	S.ReadSmat("/tmp/gosl/la/small-test-matrix.smat")
	chk.Deep2(tst, "S", 1e-17, S.ToDense().GetDeep2(), correct)
}
