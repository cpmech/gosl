// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"math"
	"testing"
)

func Test_mylab01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab01")

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
	Pf("I  = %v\n", I)
	Pf("Js = %v\n", Js)
	Pf("J  = %v\n", J)
	Pf("A  = %v\n", A)
	Pf("a  = %v\n", a)
	Pf("b  = %v\n", b)
	Pf("c  = %v\n", c)
	Pf("C  = %v\n", C)
	Pf("d  = %v\n", d)
	Pf("D  = %v\n", D)
	Pf("e  = %v\n", e)
	Pf("f  = %v\n", f)
	Pf("G  = %v\n", G)
	Pf("H  = %v\n", H)
	Pf("g  = %v\n", g)
	Pf("h  = %v\n", h)
	Pf("M  = %v\n", M)
	CompareInts(tst, "I", I, []int{666, 666, 666, 666, 666})
	CompareInts(tst, "J", J, []int{666, 666, 666, 666, 666})
	CompareStrs(tst, "Js", Js, []string{"666", "666", "666", "666", "666"})
	CompareInts(tst, "A", A, []int{})
	CompareInts(tst, "a", a, []int{})
	CompareInts(tst, "b", b, []int{0})
	CompareInts(tst, "c", c, []int{0, 1, 2, 3, 4})
	CompareInts(tst, "C", C, []int{0, -1, -2, -3, -4})
	CompareInts(tst, "d", d, []int{2, 3, 4})
	CompareInts(tst, "D", D, []int{-2, -1, 0, 1, 2, 3, 4})
	CompareInts(tst, "e", e, []int{0, 1, 2, 3, 4, 5, 6})
	CompareDbls(tst, "f", f, []float64{1, 1, 1, 1, 1})
	CompareDbls(tst, "ff", ff, []float64{666, 666, 666, 666, 666})
	CompareInts(tst, "h", h, []int{0, 1, 2, 3, 4, 7, 8, 9, 11, 12, 23, 32, 37})
	CompareInts(tst, "H", H, []int{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 11, 15, 16, 38, 39})
	CheckIntMat(tst, "M", M, [][]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}})
}

func Test_mylab02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab02")

	a := []string{"66", "644", "666", "653", "10", "0", "1", "1", "1"}
	idx := StrIndexSmall(a, "666")
	Pf("a = %v\n", a)
	Pf("idx of '666' = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
	}
	idx = StrIndexSmall(a, "1")
	Pf("idx of '1'   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
	}

	b := []int{66, 644, 666, 653, 10, 0, 1, 1, 1}
	idx = IntIndexSmall(b, 666)
	Pf("b = %v\n", b)
	Pf("idx of 666 = %v\n", idx)
	if idx != 2 {
		tst.Errorf("idx is wrong")
	}
	idx = IntIndexSmall(b, 1)
	Pf("idx of 1   = %v\n", idx)
	if idx != 6 {
		tst.Errorf("idx is wrong")
	}
}

func TestMyLab03a(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab03a")

	A := []int{1, 2, 3, -1, -2, 0, 8, -3}
	mi, ma := IntMinMax(A)
	Pf("A      = %v\n", A)
	Pf("min(A) = %v\n", mi)
	Pf("max(A) = %v\n", ma)
	if mi != -3 {
		Panic("min(A) failed")
	}
	if ma != 8 {
		Panic("max(A) failed")
	}
}

func TestMyLab03b(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab03b")

	a := []int{1, 2, 3, -1, -2, 0, 8, -3}
	b := IntFilter(a, func(i int) bool {
		if a[i] < 0 {
			return true
		}
		return false
	})
	c := IntNegOut(a)
	Pf("a = %v\n", a)
	Pf("b = %v\n", b)
	Pf("c = %v\n", c)
	CompareInts(tst, "b", b, []int{1, 2, 3, 0, 8})
	CompareInts(tst, "c", c, []int{1, 2, 3, 0, 8})

	A := []float64{1, 2, 3, -1, -2, 0, 8, -3}
	s := DblSum(A)
	mi, ma := DblMinMax(A)
	Pf("A      = %v\n", A)
	Pf("sum(A) = %v\n", s)
	Pf("min(A) = %v\n", mi)
	Pf("max(A) = %v\n", ma)
	CheckScalar(tst, "sum(A)", 1e-17, s, 8)
	CheckScalar(tst, "min(A)", 1e-17, mi, -3)
	CheckScalar(tst, "max(A)", 1e-17, ma, 8)
}

func TestMyLab04(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab04")

	n := 5
	a := LinSpace(2.0, 3.0, n)
	δ := (3.0 - 2.0) / float64(n-1)
	r := make([]float64, n)
	for i := 0; i < n; i++ {
		r[i] = 2.0 + float64(i)*δ
	}
	Pf("δ = %v\n", δ)
	Pf("a = %v\n", a)
	Pf("r = %v\n", r)
	CheckVector(tst, "linspace(2,3,5)", 1e-17, a, []float64{2.0, 2.25, 2.5, 2.75, 3.0})

	b := LinSpaceOpen(2.0, 3.0, n)
	Δ := (3.0 - 2.0) / float64(n)
	R := make([]float64, n)
	for i := 0; i < n; i++ {
		R[i] = 2.0 + float64(i)*Δ
	}
	Pf("Δ = %v\n", Δ)
	Pf("b = %v\n", b)
	Pf("R = %v\n", R)
	CheckVector(tst, "linspace(2,3,5,open)", 1e-17, b, []float64{2.0, 2.2, 2.4, 2.6, 2.8})

	c := LinSpace(2.0, 3.0, 1)
	Pf("c = %v\n", c)
	CheckVector(tst, "linspace(2,3,1)", 1e-17, c, []float64{2.0})

}

func TestMyLab05(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("mylab05")

	//                    0   1   2   3   4   5
	allkeys := []string{"a", "b", "A", "B", "α", "β"}
	excluded := []string{"A", "β"}
	g2l, l2g := LocGlobMaps(allkeys, excluded)
	Pfcyan("allkeys  = %v\n", allkeys)
	Pfcyan("excluded = %v\n", excluded)
	Pforan("g2l      = %v\n", g2l)
	Pforan("l2g      = %v\n", l2g)
	CompareInts(tst, "g2l", g2l, []int{0, 1, -1, 2, 3, -1})
	CompareInts(tst, "l2g", l2g, []int{0, 1, 3, 4})

	excludedM := map[string]bool{"A": true, "β": true}
	g2lM, l2gM := LocGlobMapsM(allkeys, excludedM)
	Pfcyan("allkeys   = %v\n", allkeys)
	Pfcyan("excludedM = %v\n", excludedM)
	Pforan("g2lM      = %v\n", g2lM)
	Pforan("l2gM      = %v\n", l2gM)
	CompareInts(tst, "g2lM", g2lM, []int{0, 1, -1, 2, 3, -1})
	CompareInts(tst, "l2gM", l2gM, []int{0, 1, 3, 4})

	//                 0    1    2    3    4    5     6     7
	names := []string{"a", "b", "c", "A", "B", "β", "x0", "y0"}
	fixed := []string{"a", "c", "A", "x0", "y0"}
	g2l, l2g = LocGlobMaps(names, fixed)
	Pfcyan("names = %v\n", names)
	Pfcyan("fixed = %v\n", fixed)
	Pforan("g2l   = %v\n", g2l)
	Pforan("l2g   = %v\n", l2g)
	CompareInts(tst, "g2l", g2l, []int{-1, 0, -1, -1, 1, 2, -1, -1})
	CompareInts(tst, "l2g", l2g, []int{1, 4, 5})

	namesN := []string{"a", "b", "c", "A", "B", "β", "x0", "y0"}
	fixedN := map[string]bool{"a": true, "c": true, "A": true, "x0": true, "y0": false}
	g2lN, l2gN := LocGlobMapsM(namesN, fixedN)
	Pfcyan("namesN = %v\n", namesN)
	Pfcyan("fixedN = %v\n", fixedN)
	Pforan("g2lN   = %v\n", g2lN)
	Pforan("l2gN   = %v\n", l2gN)
	CompareInts(tst, "g2lN", g2lN, []int{-1, 0, -1, -1, 1, 2, -1, 3})
	CompareInts(tst, "l2gN", l2gN, []int{1, 4, 5, 7})
}

func Test_conversions01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("conversions01")

	v := []float64{2.48140019424242e-08, 0.0014621532754275238, 5.558773630697262e-09, 3.0581358492226644e-08, 0.001096211253647636}
	s := Dbl2Str(v, "%.17e")
	w := Str2Dbl(s)
	CheckVector(tst, "v => s => w", 1e-17, v, w)
}

func Test_split01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("split01")

	r := DblSplit(" 1e4 1 3   8   88   ")
	Pfblue2("r = %v\n", r)
	CompareDbls(tst, "r", r, []float64{1e4, 1, 3, 8, 88})
}

func Test_functions01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	TTitle("functions01")

	x := LinSpace(-2, 2, 21)
	ym := make([]float64, len(x))
	yh := make([]float64, len(x))
	ys := make([]float64, len(x))
	yAbs2Ramp := make([]float64, len(x))
	yHea2Ramp := make([]float64, len(x))
	ySig2Heav := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		ym[i] = Ramp(x[i])
		yh[i] = Heav(x[i])
		ys[i] = Sign(x[i])
		yAbs2Ramp[i] = (x[i] + math.Abs(x[i])) / 2.0
		yHea2Ramp[i] = x[i] * yh[i]
		ySig2Heav[i] = (1.0 + ys[i]) / 2.0
	}
	CheckVector(tst, "abs => ramp", 1e-17, ym, yAbs2Ramp)
	CheckVector(tst, "hea => ramp", 1e-17, ym, yHea2Ramp)
	CheckVector(tst, "sig => heav", 1e-17, yh, ySig2Heav)

	var b bytes.Buffer
	Ff(&b, "from gosl_fig import *\n")
	Gen4Arrays(&b, "x", "ym", "yh", "ys", x, ym, yh, ys)
	Ff(&b, "subplot(3,1,1)\n")
	Ff(&b, "plot(x,ym,label='Ramp/Macaulay',clip_on=0,lw=2,marker='o')\n")
	Ff(&b, "axis([axis()[0],axis()[1],-0.1,axis()[3]])\n")
	Ff(&b, "Cross()\n")
	Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	Ff(&b, "subplot(3,1,2)\n")
	Ff(&b, "plot(x,yh,label='Heaviside',clip_on=0,lw=2,marker='o')\n")
	Ff(&b, "axis([axis()[0],axis()[1],-0.1,1.1])\n")
	Ff(&b, "Cross()\n")
	Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	Ff(&b, "subplot(3,1,3)\n")
	Ff(&b, "plot(x,ys,label='Sign',clip_on=0,lw=2,marker='o')\n")
	Ff(&b, "axis([axis()[0],axis()[1],-1.1,1.1])\n")
	Ff(&b, "Cross()\n")
	Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	Ff(&b, "show()\n")
	WriteFileD("/tmp/gosl/", "functions01.py", &b)
}

// numderiv employs a 1st order forward difference to approximate the derivative of f(x) w.r.t x @ x
func numderiv(f func(x float64) float64, x float64) float64 {
	eps, cte1 := 1e-16, 1e-5
	delta := math.Sqrt(eps * max(cte1, math.Abs(x)))
	return (f(x+delta) - f(x)) / delta
}

func Test_functions02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("functions02")

	β := 6.0
	f := func(x float64) float64 { return Sramp(x, β) }
	ff := func(x float64) float64 { return SrampD1(x, β) }

	np := 401
	//x  := LinSpace(-5e5, 5e5, np)
	//x  := LinSpace(-5e2, 5e2, np)
	x := LinSpace(-5e1, 5e1, np)
	y := make([]float64, np)
	g := make([]float64, np)
	h := make([]float64, np)
	tolg, tolh := 1e-6, 1e-5
	with_err := false
	for i := 0; i < np; i++ {
		y[i] = Sramp(x[i], β)
		g[i] = SrampD1(x[i], β)
		h[i] = SrampD2(x[i], β)
		gnum := numderiv(f, x[i])
		hnum := numderiv(ff, x[i])
		errg := math.Abs(g[i] - gnum)
		errh := math.Abs(h[i] - hnum)
		clrg, clrh := "[1;32m", "[1;32m"
		if errg > tolg {
			clrg, with_err = "[1;31m", true
		}
		if errh > tolh {
			clrh, with_err = "[1;31m", true
		}
		Pf("errg = %s%23.15e   errh = %s%23.15e[0m\n", clrg, errg, clrh, errh)
	}

	var b bytes.Buffer
	Ff(&b, "from gosl_fig import *\n")
	Gen4Arrays(&b, "x", "y", "g", "h", x, y, g, h)
	Ff(&b, "subplot(3,1,1)\n")
	Ff(&b, "plot(x,y, 'b-', lw=2)\n")
	Ff(&b, "axis('equal')\n")
	Ff(&b, "Cross()\n")
	Ff(&b, "Gll('x','y',leg=0)\n")
	Ff(&b, "subplot(3,1,2)\n")
	Ff(&b, "plot(x,g, 'b-', lw=2)\n")
	Ff(&b, "Gll('x','g',leg=0)\n")
	Ff(&b, "subplot(3,1,3)\n")
	Ff(&b, "plot(x,h, 'b-', lw=2)\n")
	Ff(&b, "Gll('x','h',leg=0)\n")
	Ff(&b, "show()\n")
	WriteFileD("/tmp/gosl/", "functions02.py", &b)
	PfBlue("file <results/functions02.py> saved\n")

	if with_err {
		Panic("errors found")
	}
}
