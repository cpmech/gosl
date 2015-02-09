// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"
)

func TestReadTable01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	TTitle("ReadTable 01")

	keys, res, err := ReadTable("data/table01.dat")
	if err != nil {
		tst.Errorf("[1;31mfile cannot be read:[0m\n%v\n", err.Error())
	}

	CompareStrs(tst, "keys", keys, []string{"a", "b", "c", "d"})
	CheckVector(tst, "a", 1.0e-17, res["a"], []float64{1, 4, 7})
	CheckVector(tst, "b", 1.0e-17, res["b"], []float64{2, 5, 8})
	CheckVector(tst, "c", 1.0e-17, res["c"], []float64{3, 6, 9})
	CheckVector(tst, "d", 1.0e-17, res["d"], []float64{666, 777, 641})
}

func TestReadMatrix01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("ReadMatrix 01")

	res, err := ReadMatrix("data/mat01.dat")
	if err != nil {
		tst.Errorf("[1;31mfile cannot be read:[0m\n%v\n", err.Error())
	}

	CheckMatrix(tst, "mat", 1.0e-17, res, [][]float64{
		{1, 2, 3, 4},
		{10, 20, 30, 40},
		{-1, -2, -3, -4},
	})

	Pforan("res = %v\n", res)
}
