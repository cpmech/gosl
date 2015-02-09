// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "testing"
)

func Test_sort01(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("sort01")

    i := []int{33, 0, 7, 8}
    x := []float64{1000.33, 0, -77.7, 88.8}

    Pforan("by 'i'\n")
    I, X, _, _ := SortQuadruples(i, x, nil, nil, "i")
    CompareInts(tst, "i", I, []int{0, 7, 8, 33})
    CompareDbls(tst, "x", X, []float64{0, -77.7, 88.8, 1000.33})

    Pforan("by 'x'\n")
    I, X, _, _ = SortQuadruples(i, x, nil, nil, "x")
    CompareInts(tst, "i", I, []int{7, 0, 8, 33})
    CompareDbls(tst, "x", X, []float64{-77.7, 0.0, 88.8, 1000.33})
}

func Test_sort02(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("sort02")

    i := []int{33, 0, 7, 8}
    x := []float64{1000.33, 0, -77.7, 88.8}
    y := []float64{1e-5, 1e-7, 1e-2, 1e-9}
    z := []float64{-8000, -7000, 0, -1}

    Pforan("by 'i'\n")
    I, X, Y, Z := SortQuadruples(i, x, y, z, "i")
    CompareInts(tst, "i", I, []int{0, 7, 8, 33})
    CompareDbls(tst, "x", X, []float64{0, -77.7, 88.8, 1000.33})
    CompareDbls(tst, "y", Y, []float64{1e-7, 1e-2, 1e-9, 1e-5})
    CompareDbls(tst, "z", Z, []float64{-7000, 0, -1, -8000})

    Pforan("by 'x'\n")
    I, X, Y, Z = SortQuadruples(i, x, y, z, "x")
    CompareInts(tst, "i", I, []int{7, 0, 8, 33})
    CompareDbls(tst, "x", X, []float64{-77.7, 0.0, 88.8, 1000.33})
    CompareDbls(tst, "y", Y, []float64{1e-2, 1e-7, 1e-9, 1e-5})
    CompareDbls(tst, "z", Z, []float64{0, -7000, -1, -8000})

    Pforan("by 'y'\n")
    I, X, Y, Z = SortQuadruples(i, x, y, z, "y")
    CompareInts(tst, "i", I, []int{8, 0, 33, 7})
    CompareDbls(tst, "x", X, []float64{88.8, 0, 1000.33, -77.7})
    CompareDbls(tst, "y", Y, []float64{1e-9, 1e-7, 1e-5, 1e-2})
    CompareDbls(tst, "z", Z, []float64{-1, -7000, -8000, 0})

    Pforan("by 'z'\n")
    I, X, Y, Z = SortQuadruples(i, x, y, z, "z")
    CompareInts(tst, "i", I, []int{33, 0, 8, 7})
    CompareDbls(tst, "x", X, []float64{1000.33, 0, 88.8, -77.7})
    CompareDbls(tst, "y", Y, []float64{1e-5, 1e-7, 1e-9, 1e-2})
    CompareDbls(tst, "z", Z, []float64{-8000, -7000, -1, 0})
}

func Test_sort03(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("sort03")

    a, b, c := 8.0, -5.5, 4.0
    DblSort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{-5.5, 4, 8})
    DblDsort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{8, 4, -5.5})

    a, b, c = -18.0, -5.5, 4.0
    DblSort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{-18, -5.5, 4})
    DblDsort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{4, -5.5, -18})

    a, b, c = 1.0, 2.0, 3.0
    DblSort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{1, 2, 3})
    DblDsort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{3, 2, 1})

    a, b, c = 1.0, 3.0, 2.0
    DblSort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{1, 2, 3})
    DblDsort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{3, 2, 1})

    a, b, c = 3.0, 2.0, 1.0
    DblSort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{1, 2, 3})
    DblDsort3(&a, &b, &c)
    Pforan("a b c = %v %v %v\n", a, b, c)
    CompareDbls(tst, "a b c", []float64{a, b, c}, []float64{3, 2, 1})
}

func Test_sort04(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("sort04")

    a := []float64{-3, -7, 8, 11, 3, 0, -11, 8}
    b := DblGetSorted(a)
    CompareDbls(tst, "a(sorted)", b, []float64{-11, -7, -3, 0, 3, 8, 8, 11})
}

func Test_sort05(tst *testing.T) {

    prevTs := Tsilent
    defer func() {
        Tsilent = prevTs
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    //Tsilent = false
    TTitle("sort05")

    a := map[string]int{"a":1, "z":2, "c":3, "y":4, "d":5, "b":6, "x":7}
    b := map[string]float64{"a":1, "z":2, "c":3, "y":4, "d":5, "b":6, "x":7}
    c := map[string]bool{"a":false, "z":true, "c":false, "y":true, "d":true, "b":false, "x":true}
    ka := StrIntMapSort(a)
    kb := StrDblMapSort(b)
    kc := StrBoolMapSort(c)
    Pforan("sorted_keys(a) = %v\n", ka)
    Pforan("sorted_keys(b) = %v\n", kb)
    Pforan("sorted_keys(c) = %v\n", kc)
    CompareStrs(tst, "ka", ka, []string{"a", "b", "c", "d", "x", "y", "z"})
    CompareStrs(tst, "kb", kb, []string{"a", "b", "c", "d", "x", "y", "z"})
    CompareStrs(tst, "kc", kc, []string{"a", "b", "c", "d", "x", "y", "z"})

    ka, va := StrIntMapSortSplit(a)
    Pfpink("sorted_keys(a) = %v\n", ka)
    Pfpink("sorted_vals(a) = %v\n", va)
    CompareStrs(tst, "ka", ka, []string{"a", "b", "c", "d", "x", "y", "z"})
    CompareInts(tst, "va", va, []int{1, 6, 3, 5, 7, 4, 2})

    kb, vb := StrDblMapSortSplit(b)
    Pfcyan("sorted_keys(b) = %v\n", kb)
    Pfcyan("sorted_vals(b) = %v\n", vb)
    CompareStrs(tst, "kb", kb, []string{"a", "b", "c", "d", "x", "y", "z"})
    CompareDbls(tst, "vb", vb, []float64{1, 6, 3, 5, 7, 4, 2})

    kc, vc := StrBoolMapSortSplit(c)
    Pfcyan("sorted_keys(c) = %v\n", kc)
    Pfcyan("sorted_vals(c) = %v\n", vc)
    CompareStrs(tst, "kc", kc, []string{"a", "b", "c", "d", "x", "y", "z"})
    CompareBools(tst, "vc", vc, []bool{false, false, false, true, true, true, true})
}
