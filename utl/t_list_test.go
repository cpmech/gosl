// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"reflect"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_list01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("list01. DblList.Append")

	var m DblList
	m.Append(2, 2.0)
	io.Pforan("m = %v\n", m)
	equal := reflect.DeepEqual(m.Vals, [][]float64{{}, {}, {2}})
	if !equal {
		chk.PrintFail("DblList Append")
	}

	m.Append(0, 0.0)
	m.Append(1, 1.0)
	io.Pforan("m = %v\n", m)
	equal = reflect.DeepEqual(m.Vals, [][]float64{{0}, {1}, {2}})
	if !equal {
		chk.PrintFail("DblList Append")
	}
}

func Test_list02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("list02. DblSlist.Append")

	var L DblSlist

	L.Append(true, 0.0)
	io.Pforan("L.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "0: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "0: vals", 1e-17, L.Vals, []float64{0.0})
	chk.Ints(tst, "0: ptrs", L.Ptrs, []int{0, 1})

	L.Append(true, 1.0)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "1: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "1: vals", 1e-17, L.Vals, []float64{0.0, 1.0})
	chk.Ints(tst, "1: ptrs", L.Ptrs, []int{0, 1, 2})

	L.Append(false, 1.1)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "2: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "2: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1})
	chk.Ints(tst, "2: ptrs", L.Ptrs, []int{0, 1, 3})

	L.Append(false, 1.2)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "3: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "3: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2})
	chk.Ints(tst, "3: ptrs", L.Ptrs, []int{0, 1, 4})

	L.Append(false, 1.3)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "4: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "4: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3})
	chk.Ints(tst, "4: ptrs", L.Ptrs, []int{0, 1, 5})

	L.Append(true, 2.0)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "5: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "5: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3, 2.0})
	chk.Ints(tst, "5: ptrs", L.Ptrs, []int{0, 1, 5, 6})

	L.Append(false, 2.1)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "6: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "6: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3, 2.0, 2.1})
	chk.Ints(tst, "6: ptrs", L.Ptrs, []int{0, 1, 5, 7})

	L.Append(true, 3.0)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "7: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "7: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3, 2.0, 2.1, 3.0})
	chk.Ints(tst, "7: ptrs", L.Ptrs, []int{0, 1, 5, 7, 8})

	L.Append(false, 3.1)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "8: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "8: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3, 2.0, 2.1, 3.0, 3.1})
	chk.Ints(tst, "8: ptrs", L.Ptrs, []int{0, 1, 5, 7, 9})

	L.Append(false, 3.2)
	io.Pforan("\nL.Vals = %v\n", L.Vals)
	io.Pfpink("L.Ptrs = %v\n", L.Ptrs)
	chk.Ints(tst, "9: lens", []int{L.Ptrs[len(L.Ptrs)-1]}, []int{len(L.Vals)})
	chk.Vector(tst, "9: vals", 1e-17, L.Vals, []float64{0.0, 1.0, 1.1, 1.2, 1.3, 2.0, 2.1, 3.0, 3.1, 3.2})
	chk.Ints(tst, "9: ptrs", L.Ptrs, []int{0, 1, 5, 7, 10})

	io.Pf("\n")
	L.Print("%10g")
}
