// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestDeep01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deep01")

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
	io.Pf("a = %v\n", a)

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
	io.Pf("b = %v\n", b)
}

func TestDeep02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deep02")

	a := Deep3alloc(3, 2, 4)
	Deep3set(a, 666)
	chk.Deep3(tst, "a", 1e-16, a, [][][]float64{
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
		{{666, 666, 666, 666}, {666, 666, 666, 666}},
	})
	io.Pf("a = %v\n", a)

	b := Deep4alloc(3, 2, 1, 2)
	Deep4set(b, 666)
	chk.Deep4(tst, "b", 1e-16, b, [][][][]float64{
		{{{666, 666}}, {{666, 666}}},
		{{{666, 666}}, {{666, 666}}},
		{{{666, 666}}, {{666, 666}}},
	})
	io.Pf("b = %v\n", b)
}

func TestDeep03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deep03. check sizes of Deep2")

	a := [][]float64{
		{1, 2, 3},
		{3, 4, 5},
	}
	ok := Deep2checkSize(2, 3, a)
	if !ok {
		tst.Errorf("check should have returned true")
	}

	ok = Deep2checkSize(1, 3, a)
	if ok {
		tst.Errorf("check should have returned false")
	}

	ok = Deep2checkSize(2, 2, a)
	if ok {
		tst.Errorf("check should have returned false")
	}

	b := [][]float64{}
	ok = Deep2checkSize(0, 0, b)
	if !ok {
		tst.Errorf("check should have returned true")
	}
}

func TestDeep04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deep04. check sizes of Deep3")

	a := [][][]float64{
		{
			{1, 2, 3},
			{3, 4, 5},
		},
	}
	ok := Deep3checkSize(1, 2, 3, a)
	if !ok {
		tst.Errorf("check should have returned true")
	}

	ok = Deep3checkSize(2, 2, 3, a)
	if ok {
		tst.Errorf("check should have returned false")
	}

	ok = Deep3checkSize(1, 1, 3, a)
	if ok {
		tst.Errorf("check should have returned false")
	}

	ok = Deep3checkSize(1, 2, 2, a)
	if ok {
		tst.Errorf("check should have returned false")
	}

	b := [][][]float64{}
	ok = Deep3checkSize(0, 0, 0, b)
	if !ok {
		tst.Errorf("check should have returned true")
	}

	c := [][][]float64{
		{
			{},
		},
	}
	ok = Deep3checkSize(1, 1, 0, c)
	if !ok {
		tst.Errorf("check should have returned true")
	}

	d := [][][]float64{
		{},
	}
	ok = Deep3checkSize(1, 0, 0, d)
	if !ok {
		tst.Errorf("check should have returned true")
	}
}

func TestDeep05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Deep05. transpose")

	a := [][]float64{
		{1, 2, 3},
		{3, 4, 5},
	}
	aT := Deep2transpose(a)
	chk.Deep2(tst, "trans(a)", 1e-17, aT, [][]float64{
		{1, 3},
		{2, 4},
		{3, 5},
	})
}
