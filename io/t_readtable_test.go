// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestReadTable01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ReadTable 01")

	keys, res := ReadTableOrPanic("data/table01.dat")

	chk.Strings(tst, "keys", keys, []string{"a", "b", "c", "d"})
	chk.Vector(tst, "a", 1.0e-17, res["a"], []float64{1, 4, 7})
	chk.Vector(tst, "b", 1.0e-17, res["b"], []float64{2, 5, 8})
	chk.Vector(tst, "c", 1.0e-17, res["c"], []float64{3, 6, 9})
	chk.Vector(tst, "d", 1.0e-17, res["d"], []float64{666, 777, 641})
}

func TestReadMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ReadMatrix 01")

	res, err := ReadMatrix("data/mat01.dat")
	if err != nil {
		tst.Errorf("[1;31mfile cannot be read:[0m\n%v\n", err.Error())
	}

	chk.Matrix(tst, "mat", 1.0e-17, res, [][]float64{
		{1, 2, 3, 4},
		{10, 20, 30, 40},
		{-1, -2, -3, -4},
	})

	Pforan("res = %v\n", res)
}
