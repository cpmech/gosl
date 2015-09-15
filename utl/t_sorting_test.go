// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_sort01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort01")

	i := []int{33, 0, 7, 8}
	x := []float64{1000.33, 0, -77.7, 88.8}

	io.Pforan("by 'i'\n")
	I, X, _, _, err := SortQuadruples(i, x, nil, nil, "i")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{0, 7, 8, 33})
	chk.Vector(tst, "x", 1e-16, X, []float64{0, -77.7, 88.8, 1000.33})

	io.Pforan("by 'x'\n")
	I, X, _, _, err = SortQuadruples(i, x, nil, nil, "x")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{7, 0, 8, 33})
	chk.Vector(tst, "x", 1e-16, X, []float64{-77.7, 0.0, 88.8, 1000.33})

	x = []float64{1000.33, 0, -77.7, 88.8}
	DblSort4(&x[0], &x[1], &x[2], &x[3])
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-16, x, []float64{-77.7, 0.0, 88.8, 1000.33})

	x = []float64{1, 10.33, 0, -8.7}
	DblSort4(&x[0], &x[1], &x[2], &x[3])
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-16, x, []float64{-8.7, 0, 1, 10.33})

	x = []float64{100.33, 10, -77.7, 8.8}
	DblSort4(&x[0], &x[1], &x[2], &x[3])
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-16, x, []float64{-77.7, 8.8, 10, 100.33})

	x = []float64{-10.33, 0, 7.7, -8.8}
	DblSort4(&x[0], &x[1], &x[2], &x[3])
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-16, x, []float64{-10.33, -8.8, 0, 7.7})

	x = []float64{-1000.33, 8, -177.7, 0.8}
	DblSort4(&x[0], &x[1], &x[2], &x[3])
	io.Pforan("x = %v\n", x)
	chk.Vector(tst, "x", 1e-16, x, []float64{-1000.33, -177.7, 0.8, 8})
}

func Test_sort02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort02")

	i := []int{33, 0, 7, 8}
	x := []float64{1000.33, 0, -77.7, 88.8}
	y := []float64{1e-5, 1e-7, 1e-2, 1e-9}
	z := []float64{-8000, -7000, 0, -1}

	io.Pforan("by 'i'\n")
	I, X, Y, Z, err := SortQuadruples(i, x, y, z, "i")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{0, 7, 8, 33})
	chk.Vector(tst, "x", 1e-16, X, []float64{0, -77.7, 88.8, 1000.33})
	chk.Vector(tst, "y", 1e-16, Y, []float64{1e-7, 1e-2, 1e-9, 1e-5})
	chk.Vector(tst, "z", 1e-16, Z, []float64{-7000, 0, -1, -8000})

	io.Pforan("by 'x'\n")
	I, X, Y, Z, err = SortQuadruples(i, x, y, z, "x")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{7, 0, 8, 33})
	chk.Vector(tst, "x", 1e-16, X, []float64{-77.7, 0.0, 88.8, 1000.33})
	chk.Vector(tst, "y", 1e-16, Y, []float64{1e-2, 1e-7, 1e-9, 1e-5})
	chk.Vector(tst, "z", 1e-16, Z, []float64{0, -7000, -1, -8000})

	io.Pforan("by 'y'\n")
	I, X, Y, Z, err = SortQuadruples(i, x, y, z, "y")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{8, 0, 33, 7})
	chk.Vector(tst, "x", 1e-16, X, []float64{88.8, 0, 1000.33, -77.7})
	chk.Vector(tst, "y", 1e-16, Y, []float64{1e-9, 1e-7, 1e-5, 1e-2})
	chk.Vector(tst, "z", 1e-16, Z, []float64{-1, -7000, -8000, 0})

	io.Pforan("by 'z'\n")
	I, X, Y, Z, err = SortQuadruples(i, x, y, z, "z")
	if err != nil {
		tst.Errorf("%v\n", err)
	}
	chk.Ints(tst, "i", I, []int{33, 0, 8, 7})
	chk.Vector(tst, "x", 1e-16, X, []float64{1000.33, 0, 88.8, -77.7})
	chk.Vector(tst, "y", 1e-16, Y, []float64{1e-5, 1e-7, 1e-9, 1e-2})
	chk.Vector(tst, "z", 1e-16, Z, []float64{-8000, -7000, -1, 0})
}

func Test_sort03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort03")

	a, b, c := 8.0, -5.5, 4.0
	DblSort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{-5.5, 4, 8})
	DblDsort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{8, 4, -5.5})

	a, b, c = -18.0, -5.5, 4.0
	DblSort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{-18, -5.5, 4})
	DblDsort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{4, -5.5, -18})

	a, b, c = 1.0, 2.0, 3.0
	DblSort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{1, 2, 3})
	DblDsort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{3, 2, 1})

	a, b, c = 1.0, 3.0, 2.0
	DblSort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{1, 2, 3})
	DblDsort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{3, 2, 1})

	a, b, c = 3.0, 2.0, 1.0
	DblSort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{1, 2, 3})
	DblDsort3(&a, &b, &c)
	io.Pforan("a b c = %v %v %v\n", a, b, c)
	chk.Vector(tst, "a b c", 1e-16, []float64{a, b, c}, []float64{3, 2, 1})
}

func Test_sort04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort04")

	a := []float64{-3, -7, 8, 11, 3, 0, -11, 8}
	b := DblGetSorted(a)
	chk.Vector(tst, "a(sorted)", 1e-16, b, []float64{-11, -7, -3, 0, 3, 8, 8, 11})
}

func Test_sort05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort05")

	a := map[string]int{"a": 1, "z": 2, "c": 3, "y": 4, "d": 5, "b": 6, "x": 7}
	b := map[string]float64{"a": 1, "z": 2, "c": 3, "y": 4, "d": 5, "b": 6, "x": 7}
	c := map[string]bool{"a": false, "z": true, "c": false, "y": true, "d": true, "b": false, "x": true}
	ka := StrIntMapSort(a)
	kb := StrDblMapSort(b)
	kc := StrBoolMapSort(c)
	io.Pforan("sorted_keys(a) = %v\n", ka)
	io.Pforan("sorted_keys(b) = %v\n", kb)
	io.Pforan("sorted_keys(c) = %v\n", kc)
	chk.Strings(tst, "ka", ka, []string{"a", "b", "c", "d", "x", "y", "z"})
	chk.Strings(tst, "kb", kb, []string{"a", "b", "c", "d", "x", "y", "z"})
	chk.Strings(tst, "kc", kc, []string{"a", "b", "c", "d", "x", "y", "z"})

	ka, va := StrIntMapSortSplit(a)
	io.Pfpink("sorted_keys(a) = %v\n", ka)
	io.Pfpink("sorted_vals(a) = %v\n", va)
	chk.Strings(tst, "ka", ka, []string{"a", "b", "c", "d", "x", "y", "z"})
	chk.Ints(tst, "va", va, []int{1, 6, 3, 5, 7, 4, 2})

	kb, vb := StrDblMapSortSplit(b)
	io.Pfcyan("sorted_keys(b) = %v\n", kb)
	io.Pfcyan("sorted_vals(b) = %v\n", vb)
	chk.Strings(tst, "kb", kb, []string{"a", "b", "c", "d", "x", "y", "z"})
	chk.Vector(tst, "vb", 1e-16, vb, []float64{1, 6, 3, 5, 7, 4, 2})

	kc, vc := StrBoolMapSortSplit(c)
	io.Pfcyan("sorted_keys(c) = %v\n", kc)
	io.Pfcyan("sorted_vals(c) = %v\n", vc)
	chk.Strings(tst, "kc", kc, []string{"a", "b", "c", "d", "x", "y", "z"})
	chk.Bools(tst, "vc", vc, []bool{false, false, false, true, true, true, true})
}

func Test_sort06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort06. int => ??? maps")

	a := map[int]bool{100: true, 101: false, 102: true, 10: false, 9: true, 8: false, 0: true}
	k := IntBoolMapSort(a)
	io.Pforan("sorted_keys(a) = %v\n", k)
	chk.Ints(tst, "k", k, []int{0, 8, 9, 10, 100, 101, 102})
}

func Test_sort07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sort07. sort 3 ints")

	x := []int{0, 10, 1, 3}
	IntSort3(&x[0], &x[1], &x[2])
	chk.Ints(tst, "sort3(x)", x, []int{0, 1, 10, 3})

	x = []int{0, 10, 1, 3}
	IntSort4(&x[0], &x[1], &x[2], &x[3])
	chk.Ints(tst, "sort4(x)", x, []int{0, 1, 3, 10})
}
