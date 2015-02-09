// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
	"testing"
)

func Test_deep01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("deep01")

	a := Deep3alloc(3, 2, 4)
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 4; k++ {
				if math.Abs(a[i][j][k]) > 1e-17 {
					tst.Errorf("[1;31ma[i][j][k] failed[0m")
				}
			}
		}
	}
	Pf("a = %v\n", a)

	b := Deep4alloc(3, 2, 1, 2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 1; k++ {
				for l := 0; l < 2; l++ {
					if math.Abs(b[i][j][k][l]) > 1e-17 {
						tst.Errorf("[1;31mb[i][j][k][l] failed[0m")
					}
				}
			}
		}
	}
	Pf("b = %v\n", b)
}

func Test_deep02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("deep02")

	a := Deep3alloc(3, 2, 4)
	Deep3set(a, 666)
	CompareDeep3(tst, "a", a, [][][]float64{
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
	})
	Pf("a = %v\n", a)

	b := Deep4alloc(3, 2, 1, 2)
	Deep4set(b, 666)
	CompareDeep4(tst, "b", b, [][][][]float64{
		{{{666, 666}}, {{666, 666}}},
		{{{666, 666}}, {{666, 666}}},
		{{{666, 666}}, {{666, 666}}},
	})
	Pf("b = %v\n", b)
}
