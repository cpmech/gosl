// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_map01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("map01")

	m := map[int][]int{
		1: {100, 101},
		2: {1000},
		3: {200, 300, 400},
	}
	io.Pforan("m (before) = %v\n", m)
	IntIntsMapAppend(m, 1, 102)
	io.Pfpink("m (after) = %v\n", m)
	chk.Ints(tst, "m[1]", m[1], []int{100, 101, 102})
	chk.Ints(tst, "m[2]", m[2], []int{1000})
	chk.Ints(tst, "m[3]", m[3], []int{200, 300, 400})
	IntIntsMapAppend(m, 4, 666)
	io.Pfcyan("m (after) = %v\n", m)
	chk.Ints(tst, "m[1]", m[1], []int{100, 101, 102})
	chk.Ints(tst, "m[2]", m[2], []int{1000})
	chk.Ints(tst, "m[3]", m[3], []int{200, 300, 400})
	chk.Ints(tst, "m[4]", m[4], []int{666})
	chk.Ints(tst, "m[5]", m[5], nil)
}

func Test_map02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("map02")

	m := map[string][]float64{
		"a": {100, 101},
		"b": {1000},
		"c": {200, 300, 400},
	}
	io.Pforan("m (before) = %v\n", m)
	StrFltsMapAppend(m, "a", 102)
	io.Pfpink("m (after) = %v\n", m)
	chk.Array(tst, "m[\"a\"]", 1e-16, m["a"], []float64{100, 101, 102})
	chk.Array(tst, "m[\"b\"]", 1e-16, m["b"], []float64{1000})
	chk.Array(tst, "m[\"c\"]", 1e-16, m["c"], []float64{200, 300, 400})
	StrFltsMapAppend(m, "d", 666)
	io.Pfcyan("m (after) = %v\n", m)
	chk.Array(tst, "m[\"a\"]", 1e-16, m["a"], []float64{100, 101, 102})
	chk.Array(tst, "m[\"b\"]", 1e-16, m["b"], []float64{1000})
	chk.Array(tst, "m[\"c\"]", 1e-16, m["c"], []float64{200, 300, 400})
	chk.Array(tst, "m[\"d\"]", 1e-16, m["d"], []float64{666})
	chk.Array(tst, "m[\"e\"]", 1e-16, m["e"], nil)
}

func Test_map03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("map03")

	m := map[string][]int{
		"a": {100, 101},
		"b": {1000},
		"c": {200, 300, 400},
	}
	io.Pforan("m (before) = %v\n", m)
	StrIntsMapAppend(m, "a", 102)
	io.Pfpink("m (after) = %v\n", m)
	chk.Ints(tst, "m[\"a\"]", m["a"], []int{100, 101, 102})
	chk.Ints(tst, "m[\"b\"]", m["b"], []int{1000})
	chk.Ints(tst, "m[\"c\"]", m["c"], []int{200, 300, 400})
	StrIntsMapAppend(m, "d", 666)
	io.Pfcyan("m (after) = %v\n", m)
	chk.Ints(tst, "m[\"a\"]", m["a"], []int{100, 101, 102})
	chk.Ints(tst, "m[\"b\"]", m["b"], []int{1000})
	chk.Ints(tst, "m[\"c\"]", m["c"], []int{200, 300, 400})
	chk.Ints(tst, "m[\"d\"]", m["d"], []int{666})
	chk.Ints(tst, "m[\"e\"]", m["e"], nil)
}
