// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
	"reflect"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_deep01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("deep01")

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

func Test_deep02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("deep02")

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

func Test_deep03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("deep03. Deep2mat")

	var m Deep2mat
	m.Append(2, 2.0)
	io.Pforan("m = %v\n", m)
	equal := reflect.DeepEqual(m.Vals, [][]float64{{}, {}, {2}})
	if !equal {
		chk.PrintFail("Deep2mat Append")
	}

	m.Append(0, 0.0)
	m.Append(1, 1.0)
	io.Pforan("m = %v\n", m)
	equal = reflect.DeepEqual(m.Vals, [][]float64{{0}, {1}, {2}})
	if !equal {
		chk.PrintFail("Deep2mat Append")
	}
}
