// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"
)

func Test_serial01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("serial01")

	A := [][][]float64{
		{{100, 101, 102}, {103}, {104, 105}},
		{{106}, {107}},
		{{108}, {109, 110}},
		{{111}},
		{{112, 113, 114, 115}, {116}, {117, 118}, {119, 120, 121}},
	}

	// serialize
	Deep3Print("A", A)
	I, P, S := Deep3Serialize(A)
	Deep3GetInfo(I, P, S, true)

	// check serialization
	CompareInts(tst, "I", I, []int{0, 0, 0, 1, 1, 2, 2, 3, 4, 4, 4, 4})
	CompareInts(tst, "P", P, []int{0, 3, 4, 6, 7, 8, 9, 11, 12, 16, 17, 19, 22})
	Scor := LinSpace(100, 121, 22)
	Pf("Scor = %v\n", Scor)
	CompareDbls(tst, "S", S, Scor)

	// deserialize
	B := Deep3Deserialize(I, P, S, false)
	Deep3Print("B", B)

	// check deserialization
	CompareDeep3(tst, "A", A, B)
}

func Test_serial02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("serial02")

	A := [][][]float64{
		{{1, 3, 4}, {6}, {232, 23, 292, 2023}, {2, 3}},
		{{0}, {1}, {0}},
		{{0, 5, 6, 8, 3, 0}},
	}

	// serialize
	Deep3Print("A", A)
	I, P, S := Deep3Serialize(A)
	Deep3GetInfo(I, P, S, true)

	// check serialization
	CompareInts(tst, "I", I, []int{0, 0, 0, 0, 1, 1, 1, 2})
	CompareInts(tst, "P", P, []int{0, 3, 4, 8, 10, 11, 12, 13, 19})
	CompareDbls(tst, "S", S, []float64{1, 3, 4, 6, 232, 23, 292, 2023, 2, 3, 0, 1, 0, 0, 5, 6, 8, 3, 0})

	// deserialize
	B := Deep3Deserialize(I, P, S, false)
	Deep3Print("B", B)

	// check deserialization
	CompareDeep3(tst, "A", A, B)
}

func Test_serial03(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("serial03")

	a := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}
	v := DblMatToArray(a)
	b := DblArrayToMat(v, 3, 4)
	Pforan("a => v = %v\n", v)
	Pforan("v => a = %v\n", b)
	CheckVector(tst, "a => v", 1e-15, v, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})
	CheckMatrix(tst, "v => a", 1e-15, b, a)
}
