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
	io.Pf("I  = %v\n", I)
	io.Pf("Js = %v\n", Js)
	io.Pf("J  = %v\n", J)
	io.Pf("A  = %v\n", A)
	io.Pf("a  = %v\n", a)
	io.Pf("b  = %v\n", b)
	io.Pf("c  = %v\n", c)
	io.Pf("C  = %v\n", C)
	io.Pf("d  = %v\n", d)
	io.Pf("D  = %v\n", D)
	io.Pf("e  = %v\n", e)
	io.Pf("f  = %v\n", f)
	io.Pf("G  = %v\n", G)
	io.Pf("H  = %v\n", H)
	io.Pf("g  = %v\n", g)
	io.Pf("h  = %v\n", h)
	io.Pf("M  = %v\n", M)
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

func TestMyLab03a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab03a")

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
}

func TestMyLab03b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("mylab03b")

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

func TestMyLab04(tst *testing.T) {

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
