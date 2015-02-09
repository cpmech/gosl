// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"
)

func Test_bestsq01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("bestsq01")

	for i := 1; i <= 12; i++ {
		nrow, ncol := BestSquare(i)
		Pforan("nrow, ncol, nrow*ncol = %2d, %2d, %2d\n", nrow, ncol, nrow*ncol)
		if nrow*ncol != i {
			Panic("BestSquare failed")
		}
	}
}

func Test_print01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("print01")

	a := []float64{3, 33.3, 666, 2, 7, 5}
	s := DblPrint(a, "%v")
	x := "[3, 33.3, 666, 2, 7, 5]"
	Pforan("s = %v\n", s)
	if s != x {
		Panic("print(a) failed: %s != %s", s, x)
	}
	s = DblPrint(a, "%.1f")
	Pforan("s = %v\n", s)
	y := "[3.0, 33.3, 666.0, 2.0, 7.0, 5.0]"
	if s != y {
		Panic("print(a) failed: %s != %s", s, y)
	}

	b := []int{3, 33, 66, 7, 8, 9}
	B := IntPrint(b, "%3d")
	Pforan("B = %v\n", B)
	//    [..., ..., ..., ..., ..., ...]
	C := "[  3,  33,  66,   7,   8,   9]"
	if B != C {
		Panic("print(b) failed: %s != %s", B, C)
	}
}
