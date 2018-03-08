// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestSerial01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Serial01. Serialize and deserialize Deep2")

	a := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 0, -1, -2},
	}
	v := SerializeDeep2(a)
	b := DeserializeDeep2(v, 3, 4)
	io.Pforan("a => v = %v\n", v)
	io.Pforan("v => a = %v\n", b)
	chk.Array(tst, "a => v", 1e-15, v, []float64{1, 5, 9, 2, 6, 0, 3, 7, -1, 4, 8, -2})
	chk.Deep2(tst, "v => a", 1e-15, b, a)
}

func TestSerial02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Serial02. Serialize and deserialize Deep3")

	A := [][][]float64{
		{{100, 101, 102}, {103}, {104, 105}},
		{{106}, {107}},
		{{108}, {109, 110}},
		{{111}},
		{{112, 113, 114, 115}, {116}, {117, 118}, {119, 120, 121}},
	}

	// serialize
	PrintDeep3("A", A)
	I, P, S := SerializeDeep3(A)
	Deep3GetInfo(I, P, S, true)

	// check serialization
	chk.Ints(tst, "I", I, []int{0, 0, 0, 1, 1, 2, 2, 3, 4, 4, 4, 4})
	chk.Ints(tst, "P", P, []int{0, 3, 4, 6, 7, 8, 9, 11, 12, 16, 17, 19, 22})
	Scor := LinSpace(100, 121, 22)
	io.Pf("Scor = %v\n", Scor)
	chk.Array(tst, "S", 1e-16, S, Scor)

	// deserialize
	B := DeserializeDeep3(I, P, S, false)
	PrintDeep3("B", B)

	// check deserialization
	chk.Deep3(tst, "A", 1e-16, A, B)
}

func TestSerial03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Serial03. Serialize and deserialize Deep3")

	A := [][][]float64{
		{{1, 3, 4}, {6}, {232, 23, 292, 2023}, {2, 3}},
		{{0}, {1}, {0}},
		{{0, 5, 6, 8, 3, 0}},
	}

	// serialize
	PrintDeep3("A", A)
	I, P, S := SerializeDeep3(A)
	Deep3GetInfo(I, P, S, true)

	// check serialization
	chk.Ints(tst, "I", I, []int{0, 0, 0, 0, 1, 1, 1, 2})
	chk.Ints(tst, "P", P, []int{0, 3, 4, 8, 10, 11, 12, 13, 19})
	chk.Array(tst, "S", 1e-16, S, []float64{1, 3, 4, 6, 232, 23, 292, 2023, 2, 3, 0, 1, 0, 0, 5, 6, 8, 3, 0})

	// deserialize
	B := DeserializeDeep3(I, P, S, false)
	PrintDeep3("B", B)

	// check deserialization
	chk.Deep3(tst, "A", 1e-16, A, B)
}
