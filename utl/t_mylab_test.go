// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_mylab01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab01")

	I := make([]int, 5)
	IntFill(I, 666)
	J := IntVals(5, 666)
	Js := StrVals(5, "666")
	M := IntsAlloc(3, 4)
	N := DblsAlloc(3, 4)
	S := StrsAlloc(2, 3)
	A := IntRange(-1)
	a := IntRange2(0, 0)
	b := IntRange2(0, 1)
	c := IntRange2(0, 5)
	C := IntRange3(0, -5, -1)
	d := IntRange2(2, 5)
	D := IntRange2(-2, 5)
	e := IntAddScalar(D, 2)
	f := DblOnes(5)
	ff := DblVals(5, 666)
	g := []int{1, 2, 3, 4, 3, 4, 2, 1, 1, 2, 3, 4, 4, 2, 3, 7, 8, 3, 8, 3, 9, 0, 11, 23, 1, 2, 32, 12, 4, 32, 4, 11, 37}
	h := IntUnique(g)
	G := []int{1, 2, 3, 38, 3, 5, 3, 1, 2, 15, 38, 1, 11}
	H := IntUnique(D, C, G, []int{16, 39})
	X, Y := MeshGrid2D(3, 6, 10, 20, 4, 3)
	P := [][]int{
		{1, 2, 3, 4, 5},
		{-1, -2, -3, -4, -5},
		{6, 7, 8, 9, 10},
	}
	Pc := IntsClone(P)
	chk.Ints(tst, "I", I, []int{666, 666, 666, 666, 666})
	chk.Ints(tst, "J", J, []int{666, 666, 666, 666, 666})
	chk.Strings(tst, "Js", Js, []string{"666", "666", "666", "666", "666"})
	chk.Ints(tst, "A", A, []int{})
	chk.Ints(tst, "a", a, []int{})
	chk.Ints(tst, "b", b, []int{0})
	chk.Ints(tst, "c", c, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "C", C, []int{0, -1, -2, -3, -4})
	chk.Ints(tst, "d", d, []int{2, 3, 4})
	chk.Ints(tst, "D", D, []int{-2, -1, 0, 1, 2, 3, 4})
	chk.Ints(tst, "e", e, []int{0, 1, 2, 3, 4, 5, 6})
	chk.Vector(tst, "f", 1e-16, f, []float64{1, 1, 1, 1, 1})
	chk.Vector(tst, "ff", 1e-16, ff, []float64{666, 666, 666, 666, 666})
	chk.Ints(tst, "h", h, []int{0, 1, 2, 3, 4, 7, 8, 9, 11, 12, 23, 32, 37})
	chk.Ints(tst, "H", H, []int{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 11, 15, 16, 38, 39})
	chk.IntMat(tst, "M", M, [][]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}})
	chk.Matrix(tst, "N", 1e-17, N, [][]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}})
	chk.Matrix(tst, "X", 1e-17, X, [][]float64{{3, 4, 5, 6}, {3, 4, 5, 6}, {3, 4, 5, 6}})
	chk.Matrix(tst, "Y", 1e-17, Y, [][]float64{{10, 10, 10, 10}, {15, 15, 15, 15}, {20, 20, 20, 20}})
	chk.StrMat(tst, "S", S, [][]string{{"", "", ""}, {"", "", ""}})
	chk.IntMat(tst, "Pc", Pc, P)
}

func Test_mylab02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab02")

	a := []string{"66", "644", "666", "653", "10", "0", "1", "1", "1"}
	idx := StrIndexSmall(a, "666")
	io.Pf("a = %v\n", a)
	io.Pf("idx of '666' = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
	}
	idx = StrIndexSmall(a, "1")
	io.Pf("idx of '1'   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
	}

	b := []int{66, 644, 666, 653, 10, 0, 1, 1, 1}
	idx = IntIndexSmall(b, 666)
	io.Pf("b = %v\n", b)
	io.Pf("idx of 666 = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
	}
	idx = IntIndexSmall(b, 1)
	io.Pf("idx of 1   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
	}
}

func Test_mylab03a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab03a. ints: min and max. dbls: min and max")

	A := []int{1, 2, 3, -1, -2, 0, 8, -3}
	mi, ma := IntMinMax(A)
	io.Pf("A      = %v\n", A)
	io.Pf("min(A) = %v\n", mi)
	io.Pf("max(A) = %v\n", ma)
	if mi != -3 {
		chk.Panic("min(A) failed")
	}
	if ma != 8 {
		chk.Panic("max(A) failed")
	}

	if Imin(-1, 2) != -1 {
		chk.Panic("Imin(-1,2) failed")
	}
	if Imax(-1, 2) != 2 {
		chk.Panic("Imax(-1,2) failed")
	}
	if Min(-1, 2) != -1.0 {
		chk.Panic("Min(-1,2) failed")
	}
	if Max(-1, 2) != 2.0 {
		chk.Panic("Max(-1,2) failed")
	}
}

func Test_mylab03b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab03b. ints: neg out and dbls min and max")

	a := []int{1, 2, 3, -1, -2, 0, 8, -3}
	b := IntFilter(a, func(i int) bool {
		if a[i] < 0 {
			return true
		}
		return false
	})
	c := IntNegOut(a)
	io.Pf("a = %v\n", a)
	io.Pf("b = %v\n", b)
	io.Pf("c = %v\n", c)
	chk.Ints(tst, "b", b, []int{1, 2, 3, 0, 8})
	chk.Ints(tst, "c", c, []int{1, 2, 3, 0, 8})

	A := []float64{1, 2, 3, -1, -2, 0, 8, -3}
	s := DblSum(A)
	mi, ma := DblMinMax(A)
	io.Pf("A      = %v\n", A)
	io.Pf("sum(A) = %v\n", s)
	io.Pf("min(A) = %v\n", mi)
	io.Pf("max(A) = %v\n", ma)
	chk.Scalar(tst, "sum(A)", 1e-17, s, 8)
	chk.Scalar(tst, "min(A)", 1e-17, mi, -3)
	chk.Scalar(tst, "max(A)", 1e-17, ma, 8)
}

func Test_mylab04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab04")

	n := 5
	a := LinSpace(2.0, 3.0, n)
	δ := (3.0 - 2.0) / float64(n-1)
	r := make([]float64, n)
	for i := 0; i < n; i++ {
		r[i] = 2.0 + float64(i)*δ
	}
	io.Pf("δ = %v\n", δ)
	io.Pf("a = %v\n", a)
	io.Pf("r = %v\n", r)
	chk.Vector(tst, "linspace(2,3,5)", 1e-17, a, []float64{2.0, 2.25, 2.5, 2.75, 3.0})

	b := LinSpaceOpen(2.0, 3.0, n)
	Δ := (3.0 - 2.0) / float64(n)
	R := make([]float64, n)
	for i := 0; i < n; i++ {
		R[i] = 2.0 + float64(i)*Δ
	}
	io.Pf("Δ = %v\n", Δ)
	io.Pf("b = %v\n", b)
	io.Pf("R = %v\n", R)
	chk.Vector(tst, "linspace(2,3,5,open)", 1e-17, b, []float64{2.0, 2.2, 2.4, 2.6, 2.8})

	c := LinSpace(2.0, 3.0, 1)
	io.Pf("c = %v\n", c)
	chk.Vector(tst, "linspace(2,3,1)", 1e-17, c, []float64{2.0})
}

func Test_mylab05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab05. cumsum, gtpenalty")

	p := []float64{1, 2, 3, 4, 5}
	cs := make([]float64, len(p))
	CumSum(cs, p)
	io.Pforan("cs = %v\n", cs)
	chk.Vector(tst, "cumsum", 1e-17, cs, []float64{1, 3, 6, 10, 15})

	res := GtPenalty(0, 0, 1)
	io.Pforan("\n0 > 0: penalty = %v\n", res)
	chk.Scalar(tst, "0 > 0", 1e-18, res, 1e-16)

	res = GtePenalty(0, 0, 1)
	io.Pforan("\n0 ≥ 0: penalty = %v\n", res)
	chk.Scalar(tst, "0 ≥ 0", 1e-18, res, 0)

	res = GtPenalty(0, -1, 1)
	io.Pforan("\n0 > -1: penalty = %v\n", res)
	chk.Scalar(tst, "0 > -1", 1e-18, res, 0)

	res = GtPenalty(1, 0, 1)
	io.Pforan("\n1 > 0: penalty = %v\n", res)
	chk.Scalar(tst, "1 > 0", 1e-18, res, 0)

	res = GtPenalty(1, 1, 1)
	io.Pforan("\n1 > 1: penalty = %v\n", res)
	chk.Scalar(tst, "1 > 1", 1e-18, res, 1e-16)

	res = GtePenalty(1, 1, 1)
	io.Pforan("\n1 ≥ 1: penalty = %v\n", res)
	chk.Scalar(tst, "1 ≥ 1", 1e-18, res, 0)

	res = GtPenalty(23, 123, 10)
	io.Pforan("\n23 > 123: (m=10) penalty = %v\n", res)
	chk.Scalar(tst, "23 > 123 (m=10)", 1e-18, res, 1000)

	res = GtePenalty(23, 123, 10)
	io.Pforan("\n23 ≥ 123: (m=10) penalty = %v\n", res)
	chk.Scalar(tst, "23 ≥ 123 (m=10)", 1e-18, res, 1000)
}

func Test_mylab06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab06. scaling")

	// |dx|>0: increasing
	io.Pfblue2("\n|dx|>0: increasing\n")
	reverse := false
	useinds := true
	x := []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	s := make([]float64, len(x))
	Scaling(s, x, -2.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, LinSpace(-2, -1, len(x)))

	// |dx|>0: reverse
	io.Pfblue2("\n|dx|>0: reverse\n")
	reverse = true
	Scaling(s, x, -3.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, LinSpace(-2, -3, len(x)))

	// |dx|>0: increasing
	io.Pfblue2("\n|dx|>0: increasing (shuffled)\n")
	reverse = false
	x = []float64{11, 10, 12, 19, 15, 20, 17, 16, 18, 13, 14}
	Scaling(s, x, 0.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, []float64{0.1, 0.0, 0.2, 0.9, 0.5, 1.0, 0.7, 0.6, 0.8, 0.3, 0.4})

	// |dx|=0: increasing (using indices)
	io.Pfblue2("\n|dx|=0: increasing (using indices)\n")
	reverse = false
	x = []float64{123, 123, 123, 123, 123}
	s = make([]float64, len(x))
	Scaling(s, x, 10.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, []float64{10, 10.25, 10.5, 10.75, 11})

	// |dx|=0: reverse (using indices)
	io.Pfblue2("\n|dx|=0: reverse (using indices)\n")
	reverse = true
	Scaling(s, x, 10.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, []float64{11, 10.75, 10.5, 10.25, 10})

	// |dx|=0: increasing (not using indices)
	io.Pfblue2("\n|dx|=0: increasing (not using indices)\n")
	reverse = false
	useinds = false
	Scaling(s, x, 88.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, DblVals(len(x), 88))

	// |dx|=0: reverse (not using indices)
	io.Pfblue2("\n|dx|=0: reverse (not using indices)\n")
	reverse = true
	Scaling(s, x, 88.0, 1e-16, reverse, useinds)
	io.Pfpink("x = %v\n", x)
	io.Pforan("s = %v\n", s)
	chk.Vector(tst, "s", 1e-15, s, DblVals(len(x), 88))
}

func Test_conversions01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("conversions01")

	v := []float64{2.48140019424242e-08, 0.0014621532754275238, 5.558773630697262e-09, 3.0581358492226644e-08, 0.001096211253647636}
	s := Dbl2Str(v, "%.17e")
	w := Str2Dbl(s)
	chk.Vector(tst, "v => s => w", 1e-17, v, w)
}

func Test_split01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("split01")

	r := DblSplit(" 1e4 1 3   8   88   ")
	io.Pfblue2("r = %v\n", r)
	chk.Vector(tst, "r", 1e-16, r, []float64{1e4, 1, 3, 8, 88})
}

func Test_copy01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("copy01")

	v := []float64{1, 2, 3, 4, 4, 5, 5, 6, 6, 6}
	w := DblCopy(v)
	io.Pfblue2("v = %v\n", v)
	chk.Vector(tst, "w==v", 1e-16, w, v)
}
