// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"math"
	"testing"
)

func TestScalars(tst *testing.T) {

	//Verbose = true
	PrintTitle("Scalars")

	A := 123.123
	B := 123.123
	t0 := new(testing.T)
	Float64(t0, "|A-B|", 1e-17, A, B)
	if t0.Failed() {
		tst.Errorf("t1 should NOT have failed\n")
		return
	}

	a := 123.123 + 456.456i
	b := 123.123 + 456.456i
	t1 := new(testing.T)
	Complex128(t1, "|a-b|", 1e-17, a, b)
	if t1.Failed() {
		tst.Errorf("t1 should NOT have failed\n")
		return
	}

	b = 123.1231 + 456.456i
	t2 := new(testing.T)
	Complex128(t2, "|a-b|", 1e-17, a, b)
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}
}

func TestAnaNum(tst *testing.T) {

	//Verbose = true
	PrintTitle("AnaNum and AnaNumC")

	t1 := new(testing.T)
	AnaNum(t1, "a", 1e-17, 123, 456, false)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	AnaNum(t2, "a", 1e-17, 123, 123, false)
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}

	t3 := new(testing.T)
	AnaNumC(t3, "a", 1e-17, 123i, 456i, false)
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	AnaNumC(t4, "a", 1e-17, 123i, 123i, false)
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}
}

func TestString(tst *testing.T) {

	//Verbose = true
	PrintTitle("String")

	t1 := new(testing.T)
	String(t1, "a", "b")
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	String(t2, "a", "a")
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}
}

func TestInt(tst *testing.T) {

	//Verbose = true
	PrintTitle("Int")

	t1 := new(testing.T)
	Int(t1, "x", 1, 2)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Int(t2, "x", 1, 1)
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}
}

func TestInt32(tst *testing.T) {

	//Verbose = true
	PrintTitle("Int32")

	t1 := new(testing.T)
	Int32(t1, "x", 1, 2)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Int32(t2, "x", 1, 1)
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}
}

func TestInt64(tst *testing.T) {

	//Verbose = true
	PrintTitle("Int64")

	t1 := new(testing.T)
	Int64(t1, "x", 1, 2)
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Int64(t2, "x", 1, 1)
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}
}

func TestInts01(tst *testing.T) {

	//Verbose = true
	PrintTitle("Ints01")

	// std int -----------------------------------------

	a := []int{1, 2, 3, 4}
	t1 := new(testing.T)
	Ints(t1, "a", a, []int{1, 2, 3, 3})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Ints(t2, "a(fixed)", a, []int{1, 2, 3, 4})
	if t2.Failed() {
		tst.Errorf("t2 should NOT have failed\n")
		return
	}

	// 32 int -----------------------------------------

	b := []int32{1, 2, 3, 4}
	t3 := new(testing.T)
	Int32s(t3, "b", b, []int32{1, 2, 3, 3})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Int32s(t4, "b(fixed)", b, []int32{1, 2, 3, 4})
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	// 64 int -----------------------------------------

	c := []int64{1, 2, 3, 4}
	t5 := new(testing.T)
	Int64s(t5, "c", c, []int64{1, 2, 3, 3})
	if !t5.Failed() {
		tst.Errorf("t5 should have failed\n")
		return
	}

	t6 := new(testing.T)
	Int64s(t6, "c(fixed)", c, []int64{1, 2, 3, 4})
	if t6.Failed() {
		tst.Errorf("t6 should NOT have failed\n")
		return
	}
}

func TestInts02(tst *testing.T) {

	//Verbose = true
	PrintTitle("Ints02")

	t1 := new(testing.T)
	Ints(t1, "x", []int{1}, []int{3, 4})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Int32s(t2, "x", []int32{1}, []int32{3, 4})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Int64s(t3, "x", []int64{1}, []int64{3, 4})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}
}

func TestBools(tst *testing.T) {

	//Verbose = true
	PrintTitle("Bools")

	t1 := new(testing.T)
	Bools(t1, "x", []bool{true, false}, []bool{true, true})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Bools(t2, "x", []bool{true, false, false}, []bool{true, true})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Bools(t3, "x(fixed)", []bool{true, false, false}, []bool{true, false, false})
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}
}

func TestStrings(tst *testing.T) {

	//Verbose = true
	PrintTitle("Strings")

	t1 := new(testing.T)
	Strings(t1, "x", []string{"t", "f"}, []string{"t", "t"})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Strings(t2, "x", []string{"t", "f", "f"}, []string{"t", "t"})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Strings(t3, "x(fixed)", []string{"t", "f", "f"}, []string{"t", "f", "f"})
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}
}

func TestArray(tst *testing.T) {

	//Verbose = true
	PrintTitle("Array")

	t1 := new(testing.T)
	Array(t1, "x", 1e-17, []float64{1, 2}, []float64{1, 1})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Array(t2, "x", 1e-17, []float64{1, 2, 2}, []float64{1, 1})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Array(t3, "x", 1e-17, []float64{1, 2, 2}, []float64{1, 2, 2})
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}

	t4 := new(testing.T)
	Array(t4, "x", 1e-17, []float64{0, 0, 0}, nil)
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	t5 := new(testing.T)
	Array(t5, "x", 1e-17, []float64{math.Sqrt(-1), 0, 0}, nil)
	if !t5.Failed() {
		tst.Errorf("t5 should have failed\n")
		return
	}
}

func TestArrayC(tst *testing.T) {

	//Verbose = true
	PrintTitle("ArrayC")

	t1 := new(testing.T)
	ArrayC(t1, "x", 1e-17, []complex128{1, 2}, []complex128{1, 1})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	ArrayC(t2, "x", 1e-17, []complex128{1, 2, 2}, []complex128{1, 1})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	ArrayC(t3, "x", 1e-17, []complex128{1, 2, 2}, []complex128{1, 2, 2})
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}

	t4 := new(testing.T)
	ArrayC(t4, "x", 1e-17, []complex128{1 + 1i, 2 + 1i, 2 + 1i}, []complex128{1 + 1i, 2 + 1i, 2 + 1i})
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	t5 := new(testing.T)
	ArrayC(t5, "x", 1e-17, []complex128{1 + 1i, 2 + 1i, 2 + 1i}, []complex128{1, 2 + 1i, 2 + 1i})
	if !t5.Failed() {
		tst.Errorf("t5 should have failed\n")
		return
	}

	t6 := new(testing.T)
	ArrayC(t6, "x", 1e-17, []complex128{0, 0, 0}, nil)
	if t6.Failed() {
		tst.Errorf("t6 should NOT have failed\n")
		return
	}

	t7 := new(testing.T)
	ArrayC(t7, "x", 1e-17, []complex128{complex(math.NaN(), 0), 0, 0}, nil)
	if !t7.Failed() {
		tst.Errorf("t7 should have failed\n")
		return
	}

	t8 := new(testing.T)
	ArrayC(t8, "x", 1e-17, []complex128{complex(0, math.NaN()), 0, 0}, nil)
	if !t8.Failed() {
		tst.Errorf("t8 should have failed\n")
		return
	}
}

func TestDeep2(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deep2")

	t1 := new(testing.T)
	Deep2(t1, "x", 1e-17, [][]float64{{1}, {2}}, [][]float64{{1}, {1}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Deep2(t2, "x", 1e-17, [][]float64{{1}, {2}, {2}}, [][]float64{{1}, {1}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Deep2(t3, "x", 1e-17, [][]float64{{1, 1}, {2}, {2}}, [][]float64{{1}, {2}, {2}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Deep2(t4, "x", 1e-17, [][]float64{{1}, {2}, {2}}, [][]float64{{1}, {2}, {2}})
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	t5 := new(testing.T)
	Deep2(t5, "x", 1e-17, [][]float64{{0}, {0}, {0}}, nil)
	if t5.Failed() {
		tst.Errorf("t5 should NOT have failed\n")
		return
	}

	t6 := new(testing.T)
	Deep2(t6, "x", 1e-17, [][]float64{{math.Sqrt(-1)}, {0}, {0}}, nil)
	if !t6.Failed() {
		tst.Errorf("t6 should have failed\n")
		return
	}
}

func TestDeep2c(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deep2c")

	t1 := new(testing.T)
	Deep2c(t1, "x", 1e-17, [][]complex128{{1}, {2}}, [][]complex128{{1}, {1}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Deep2c(t2, "x", 1e-17, [][]complex128{{1}, {2}, {2}}, [][]complex128{{1}, {1}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Deep2c(t3, "x", 1e-17, [][]complex128{{1, 1}, {2}, {2}}, [][]complex128{{1}, {2}, {2}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Deep2c(t4, "x", 1e-17, [][]complex128{{1 + 1i}, {2}, {2}}, [][]complex128{{1}, {2}, {2}})
	if !t4.Failed() {
		tst.Errorf("t4 should have failed\n")
		return
	}

	t5 := new(testing.T)
	Deep2c(t5, "x", 1e-17, [][]complex128{{1}, {2}, {2}}, [][]complex128{{1}, {2}, {2}})
	if t5.Failed() {
		tst.Errorf("t5 should NOT have failed\n")
		return
	}

	t6 := new(testing.T)
	Deep2c(t6, "x", 1e-17, [][]complex128{{0}, {0}, {0}}, nil)
	if t6.Failed() {
		tst.Errorf("t6 should NOT have failed\n")
		return
	}

	t7 := new(testing.T)
	Deep2c(t7, "x", 1e-17, [][]complex128{{complex(math.NaN(), 0)}, {0}, {0}}, nil)
	if !t7.Failed() {
		tst.Errorf("t7 should have failed\n")
		return
	}

	t8 := new(testing.T)
	Deep2c(t8, "x", 1e-17, [][]complex128{{complex(0, math.NaN())}, {0}, {0}}, nil)
	if !t8.Failed() {
		tst.Errorf("t8 should have failed\n")
		return
	}
}

func TestStrDeep2(tst *testing.T) {

	//Verbose = true
	PrintTitle("StrDeep2")

	t1 := new(testing.T)
	StrDeep2(t1, "x", [][]string{{"1"}, {"2"}}, [][]string{{"1"}, {"1"}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	StrDeep2(t2, "x", [][]string{{"1"}, {"2"}, {"2"}}, [][]string{{"1"}, {"1"}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	StrDeep2(t3, "x", [][]string{{"1", "1"}, {"2"}, {"2"}}, [][]string{{"1"}, {"2"}, {"2"}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	StrDeep2(t4, "x", [][]string{{"1"}, {"2"}, {"2"}}, [][]string{{"1"}, {"2"}, {"2"}})
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	t5 := new(testing.T)
	StrDeep2(t5, "x", [][]string{{""}, {""}, {""}}, nil)
	if t5.Failed() {
		tst.Errorf("t5 should NOT have failed\n")
		return
	}
}

func TestIntDeep2(tst *testing.T) {

	//Verbose = true
	PrintTitle("IntDeep2")

	t1 := new(testing.T)
	IntDeep2(t1, "x", [][]int{{1}, {2}}, [][]int{{1}, {1}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	IntDeep2(t2, "x", [][]int{{1}, {2}, {2}}, [][]int{{1}, {1}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	IntDeep2(t3, "x", [][]int{{1, 1}, {2}, {2}}, [][]int{{1}, {2}, {2}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	IntDeep2(t4, "x", [][]int{{1}, {2}, {2}}, [][]int{{1}, {2}, {2}})
	if t4.Failed() {
		tst.Errorf("t4 should NOT have failed\n")
		return
	}

	t5 := new(testing.T)
	IntDeep2(t5, "x", [][]int{{0}, {0}, {0}}, nil)
	if t5.Failed() {
		tst.Errorf("t5 should NOT have failed\n")
		return
	}
}

func TestDeep3(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deep3")

	t1 := new(testing.T)
	Deep3(t1, "x", 1e-17, [][][]float64{{{1}}, {{2}}}, [][][]float64{{{1}}, {{1}}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Deep3(t2, "x", 1e-17, [][][]float64{{{1}}, {{2}}, {{2}}}, [][][]float64{{{1}}, {{1}}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Deep3(t3, "x", 1e-17, [][][]float64{{{1}, {1}}, {{2}}, {{2}}}, [][][]float64{{{1}}, {{2}}, {{2}}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Deep3(t4, "x", 1e-17, [][][]float64{{{1, 1}}, {{2}}, {{2}}}, [][][]float64{{{1}}, {{2}}, {{2}}})
	if !t4.Failed() {
		tst.Errorf("t4 should have failed\n")
		return
	}

	t5 := new(testing.T)
	Deep3(t5, "x", 1e-17, [][][]float64{{{1}}, {{2}}, {{2}}}, [][][]float64{{{1}}, {{2}}, {{2}}})
	if t5.Failed() {
		tst.Errorf("t5 should NOT have failed\n")
		return
	}

	t6 := new(testing.T)
	Deep3(t6, "x", 1e-17, [][][]float64{{{0}}, {{0}}, {{0}}}, nil)
	if t6.Failed() {
		tst.Errorf("t6 should NOT have failed\n")
		return
	}

	t7 := new(testing.T)
	Deep3(t7, "x", 1e-17, [][][]float64{{{math.Sqrt(-1)}}, {{0}}, {{0}}}, nil)
	if !t7.Failed() {
		tst.Errorf("t7 should have failed\n")
		return
	}
}

func TestDeep4(tst *testing.T) {

	//Verbose = true
	PrintTitle("Deep4")

	t1 := new(testing.T)
	Deep4(t1, "x", 1e-17, [][][][]float64{{{{1}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{1}}}})
	if !t1.Failed() {
		tst.Errorf("t1 should have failed\n")
		return
	}

	t2 := new(testing.T)
	Deep4(t2, "x", 1e-17, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{1}}}})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Deep4(t3, "x", 1e-17, [][][][]float64{{{{1}}, {{1}}}, {{{2}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}})
	if !t3.Failed() {
		tst.Errorf("t3 should have failed\n")
		return
	}

	t4 := new(testing.T)
	Deep4(t4, "x", 1e-17, [][][][]float64{{{{1}, {1}}}, {{{2}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}})
	if !t4.Failed() {
		tst.Errorf("t4 should have failed\n")
		return
	}

	t5 := new(testing.T)
	Deep4(t5, "x", 1e-17, [][][][]float64{{{{1, 1}}}, {{{2}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}})
	if !t5.Failed() {
		tst.Errorf("t5 should have failed\n")
		return
	}

	t6 := new(testing.T)
	Deep4(t6, "x", 1e-17, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}}, [][][][]float64{{{{1}}}, {{{2}}}, {{{2}}}})
	if t6.Failed() {
		tst.Errorf("t6 should NOT have failed\n")
		return
	}

	t7 := new(testing.T)
	Deep4(t7, "x", 1e-17, [][][][]float64{{{{0}}}, {{{0}}}, {{{0}}}}, nil)
	if t7.Failed() {
		tst.Errorf("t7 should NOT have failed\n")
		return
	}

	t8 := new(testing.T)
	Deep4(t8, "x", 1e-17, [][][][]float64{{{{math.Sqrt(-1)}}}, {{{0}}}, {{{0}}}}, nil)
	if !t8.Failed() {
		tst.Errorf("t8 should have failed\n")
		return
	}
}

func TestSymmetry(tst *testing.T) {

	//Verbose = true
	PrintTitle("Symmetry")

	t1 := new(testing.T)
	Symmetry(t1, "x1", []float64{-1, -0.2, 0, +0.2, +1})
	if t1.Failed() {
		tst.Errorf("t1 should NOT have failed\n")
		return
	}

	t2 := new(testing.T)
	Symmetry(t2, "x1", []float64{-1, -0.4, 0, +0.2, +1})
	if !t2.Failed() {
		tst.Errorf("t2 should have failed\n")
		return
	}

	t3 := new(testing.T)
	Symmetry(t3, "x1", []float64{-1, -0.2, +0.2, +1})
	if t3.Failed() {
		tst.Errorf("t3 should NOT have failed\n")
		return
	}

	g4 := new(testing.T)
	Symmetry(g4, "x1", []float64{-1, -0.4, +0.2, +1})
	if !g4.Failed() {
		tst.Errorf("g4 should have failed\n")
		return
	}
}
