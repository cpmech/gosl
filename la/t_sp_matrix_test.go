// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestSpTriplet01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpTriplet01")

	//   0 2 0 0
	//   1 0 4 0
	//   0 0 0 5
	//   0 3 0 6
	a := new(Triplet)
	a.Init(4, 4, 6)
	a.Put(1, 0, 1)
	a.Put(0, 1, 2)
	a.Put(3, 1, 3)
	a.Put(1, 2, 4)
	a.Put(2, 3, 5)
	a.Put(3, 3, 6)

	l := a.GetDenseMatrix().Print("%2g")
	io.Pf("%v\n", l)
	chk.String(tst, l, " 0 2 0 0\n 1 0 4 0\n 0 0 0 5\n 0 3 0 6")

	a.WriteSmat("/tmp/gosl/la", "triplet01", 0)
	d, err := io.ReadFile("/tmp/gosl/la/triplet01.smat")
	status(tst, err)
	io.Pforan("d = %v\n", string(d))
	chk.String(tst, string(d), "4  4  6\n  1  0    1.000000000000000e+00\n  0  1    2.000000000000000e+00\n  3  1    3.000000000000000e+00\n  1  2    4.000000000000000e+00\n  2  3    5.000000000000000e+00\n  3  3    6.000000000000000e+00\n")
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
