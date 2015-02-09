// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"
)

func Test_map01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("map01")

	m := map[int][]int{
		1: []int{100, 101},
		2: []int{1000},
		3: []int{200, 300, 400},
	}
	Pforan("m (before) = %v\n", m)
	IntIntsMapAppend(&m, 1, 102)
	Pfpink("m (after) = %v\n", m)
	CompareInts(tst, "m[1]", m[1], []int{100, 101, 102})
	CompareInts(tst, "m[2]", m[2], []int{1000})
	CompareInts(tst, "m[3]", m[3], []int{200, 300, 400})
	IntIntsMapAppend(&m, 4, 666)
	Pfcyan("m (after) = %v\n", m)
	CompareInts(tst, "m[1]", m[1], []int{100, 101, 102})
	CompareInts(tst, "m[2]", m[2], []int{1000})
	CompareInts(tst, "m[3]", m[3], []int{200, 300, 400})
	CompareInts(tst, "m[4]", m[4], []int{666})
	CompareInts(tst, "m[5]", m[5], nil)
}

func Test_map02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("map02")

	m := map[string][]float64{
		"a": []float64{100, 101},
		"b": []float64{1000},
		"c": []float64{200, 300, 400},
	}
	Pforan("m (before) = %v\n", m)
	StrDblsMapAppend(&m, "a", 102)
	Pfpink("m (after) = %v\n", m)
	CompareDbls(tst, "m[\"a\"]", m["a"], []float64{100, 101, 102})
	CompareDbls(tst, "m[\"b\"]", m["b"], []float64{1000})
	CompareDbls(tst, "m[\"c\"]", m["c"], []float64{200, 300, 400})
	StrDblsMapAppend(&m, "d", 666)
	Pfcyan("m (after) = %v\n", m)
	CompareDbls(tst, "m[\"a\"]", m["a"], []float64{100, 101, 102})
	CompareDbls(tst, "m[\"b\"]", m["b"], []float64{1000})
	CompareDbls(tst, "m[\"c\"]", m["c"], []float64{200, 300, 400})
	CompareDbls(tst, "m[\"d\"]", m["d"], []float64{666})
	CompareDbls(tst, "m[\"e\"]", m["e"], nil)
}
