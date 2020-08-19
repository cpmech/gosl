// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"gosl/chk"
)

func TestSpConversion01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpConversion01. Triplet to CCMatrix")

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
	chk.Deep2(tst, "a", 1e-17, ad.GetDeep2(), [][]float64{
		{10, 11, 12},
		{20, 21, 22},
		{30, 31, 32},
	})
}

func TestSpConversion02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpConversion02. Triplet to CCMatrix")

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
	chk.Deep2(tst, "a", 1e-17, ad.GetDeep2(), [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})
}

func TestSpConversion03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpConversion03. TripletC to CCMatrixC")

	var t TripletC
	t.Init(3, 3, 10)
	t.Put(0, 0, +5.0+0i)
	t.Put(0, 0, +5.0+0i)
	t.Put(0, 1, 11.0+0i)
	t.Put(0, 2, 12.0+0i)
	t.Put(1, 0, 20.0+0i)
	t.Put(1, 1, 21.0+0i)
	t.Put(1, 2, 22.0+0i)
	t.Put(2, 0, 30.0+0i)
	t.Put(2, 1, 31.0+0i)
	t.Put(2, 2, 32.0+1i)
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	chk.Deep2c(tst, "a", 1.0e-17, ad.GetDeep2(), [][]complex128{
		{10, 11, 12},
		{20, 21, 22},
		{30, 31, 32 + 1i},
	})
}

func TestSpConversion04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpConversion04. TripletC to CCMatrixC")

	var t TripletC
	t.Init(4, 3, 4*3+2)
	t.Put(0, 0, +1.0+1.0i)
	t.Put(0, 1, +2.0+1.0i)
	t.Put(0, 2, +3.0+2.0i)
	t.Put(1, 0, +4.0+2.0i)
	t.Put(1, 1, +5.0+2.0i)
	t.Put(1, 2, +6.0+1.0i)
	t.Put(2, 0, +7.0+3.0i)
	t.Put(2, 1, +8.0+3.0i)
	t.Put(2, 2, +9.0+4.0i)
	t.Put(3, 0, +4.0+1.1i)
	t.Put(3, 1, 11.0+4.0i)
	t.Put(3, 2, 12.0+3.0i)
	t.Put(3, 0, +3.0+1.4i) // repeated
	t.Put(3, 0, +3.0+1.5i) // repeated
	a := t.ToMatrix(nil)
	ad := a.ToDense()
	chk.Deep2c(tst, "a", 1.0e-17, ad.GetDeep2(), [][]complex128{
		{1 + 1i, 2 + 1i, 3 + 2i},
		{4 + 2i, 5 + 2i, 6 + 1i},
		{7 + 3i, 8 + 3i, 9 + 4i},
		{10 + 4i, 11 + 4i, 12 + 3i},
	})
}
