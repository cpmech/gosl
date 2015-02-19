// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func TestSparseLA01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()
	chk.PrintTitle("TestSparse LA01")

	var t Triplet
	t.Init(3, 5, 15)
	t.Put(0, 0, 1)
	t.Put(0, 1, 2)
	t.Put(0, 2, 3)
	t.Put(0, 3, 4)
	t.Put(0, 4, 5)
	t.Put(1, 0, 0.1)
	t.Put(1, 1, 0.2)
	t.Put(1, 2, 0.3)
	t.Put(1, 3, 0.4)
	t.Put(1, 4, 0.5)
	t.Put(2, 0, 10)
	t.Put(2, 1, 20)
	t.Put(2, 2, 30)
	t.Put(2, 3, 40)
	t.Put(2, 4, 50)

	a := t.ToMatrix(nil)
	ad := a.ToDense()
	u := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	w := []float64{10.0, 20.0, 30.0}
	r := []float64{1000, 1000, 1000, 1000, 1000}
	s := []float64{1000, 1000, 1000}
	x := []float64{1000, 2000, 3000}
	W := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	PrintMat("a", ad, "%5g", false)
	PrintVec("u", u, "%5g", false)
	PrintVec("w", w, "%5g", false)
	PrintVec("r", r, "%5g", false)
	PrintVec("s", s, "%5g", false)
	PrintVec("x", x, "%5g", false)
	PrintVec("W", W, "%5g", false)

	io.Pf("\nfunc SpMatVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	p := make([]float64, 3)
	SpMatVecMul(p, 1, a, u) // p := 1*a*u
	PrintVec("p = a*u", p, "%5g", false)
	chk.Vector(tst, "p = a*u", 1e-17, p, []float64{5.5, 0.55, 55})

	io.Pf("\nfunc SpMatVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatVecMulAdd(s, 1, a, u) // s += dot(a, u)
	PrintVec("s += a*u", s, "%9g", false)
	chk.Vector(tst, "s += a*u", 1e-17, s, []float64{1005.5, 1000.55, 1055})

	io.Pf("\nfunc SpMatVecMulAddX(v []float64, a *CCMatrix, α float64, u []float64, β float64, w []float64)\n")
	SpMatVecMulAddX(x, a, 2, u, 3, W) // x += a * (2*u + 3*W)
	PrintVec("x += a * (2*u + 3*W)", x, "%9g", false)
	chk.Vector(tst, "x += a * (2*u + 3*W)", 1e-17, x, []float64{1176, 2017.6, 4760})

	io.Pf("\nfunc SpMatTrVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	q := make([]float64, 5)
	SpMatTrVecMul(q, 1, a, w) // q = dot(transpose(a), w)
	PrintVec("q = trans(a)*w", q, "%5g", false)
	chk.Vector(tst, "q = trans(a)*w", 1e-17, q, []float64{312, 624, 936, 1248, 1560})

	io.Pf("\nfunc SpMatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatTrVecMulAdd(r, 1, a, w) // r += dot(transpose(a), w)
	PrintVec("r += trans(a)*w", r, "%5g", false)
	chk.Vector(tst, "r += trans(a)*w", 1e-17, r, []float64{1312, 1624, 1936, 2248, 2560})
}

func TestSparseLA02(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()
	chk.PrintTitle("TestSparse LA02")

	var t TripletC
	t.Init(3, 5, 15, false)
	t.Put(0, 0, 1, 0)
	t.Put(0, 1, 2, 0)
	t.Put(0, 2, 3, 0)
	t.Put(0, 3, 4, 0)
	t.Put(0, 4, 5, 0)
	t.Put(1, 0, 0.1, 0)
	t.Put(1, 1, 0.2, 0)
	t.Put(1, 2, 0.3, 0)
	t.Put(1, 3, 0.4, 0)
	t.Put(1, 4, 0.5, 0)
	t.Put(2, 0, 10, 0)
	t.Put(2, 1, 20, 0)
	t.Put(2, 2, 30, 0)
	t.Put(2, 3, 40, 0)
	t.Put(2, 4, 50, 0)

	a := t.ToMatrix(nil)
	ad := a.ToDense()
	u := []complex128{0.1, 0.2, 0.3, 0.4, 0.5}
	w := []complex128{10.0, 20.0, 30.0}
	r := []complex128{1000, 1000, 1000, 1000, 1000}
	s := []complex128{1000, 1000, 1000}

	PrintMatC("a", ad, "(%4g", " +%4gi)  ", false)
	PrintVecC("u", u, "(%4g", " +%4gi)  ", false)
	PrintVecC("w", w, "(%4g", " +%4gi)  ", false)
	PrintVecC("r", r, "(%4g", " +%4gi)  ", false)
	PrintVecC("s", s, "(%4g", " +%4gi)  ", false)

	io.Pf("\nfunc SpMatVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	p := make([]complex128, 3)
	SpMatVecMulC(p, 1, a, u) // p := 1*a*u
	PrintVecC("p = a*u", p, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "p = a*u", 1e-17, p, []complex128{5.5, 0.55, 55})

	io.Pf("\nfunc SpMatVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatVecMulAddC(s, 1, a, u) // s += dot(a, u)
	PrintVecC("s += a*u", s, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "s += a*u", 1e-17, s, []complex128{1005.5, 1000.55, 1055})

	io.Pf("\nfunc SpMatTrVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	q := make([]complex128, 5)
	SpMatTrVecMulC(q, 1, a, w) // q = dot(transpose(a), w)
	PrintVecC("q = trans(a)*w", q, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "q = trans(a)*w", 1e-17, q, []complex128{312, 624, 936, 1248, 1560})

	io.Pf("\nfunc SpMatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatTrVecMulAddC(r, 1, a, w) // r += dot(transpose(a), w)
	PrintVecC("r += trans(a)*w", r, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "r += trans(a)*w", 1e-17, r, []complex128{1312, 1624, 1936, 2248, 2560})
}

func TestSparseLA03(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()
	chk.PrintTitle("TestSparse LA03")

	var t TripletC
	t.Init(3, 5, 15, false)
	t.Put(0, 0, 1, 0)
	t.Put(0, 1, 2, -1)
	t.Put(0, 2, 3, 0)
	t.Put(0, 3, 4, 3)
	t.Put(0, 4, 5, 2)
	t.Put(1, 0, 0.1, 1)
	t.Put(1, 1, 0.2, 0)
	t.Put(1, 2, 0.3, -2)
	t.Put(1, 3, 0.4, 0)
	t.Put(1, 4, 0.5, -1)
	t.Put(2, 0, 10, 0)
	t.Put(2, 1, 20, 2)
	t.Put(2, 2, 30, 0)
	t.Put(2, 3, 40, -1)
	t.Put(2, 4, 50, 0)

	a := t.ToMatrix(nil)
	ad := a.ToDense()
	u := []complex128{0.1, 0.2 + 10i, 0.3, 0.4 + 3i, 0.5}
	w := []complex128{10.0 + 1i, 20.0 - 0.5i, 30.0}
	r := []complex128{1000 + 1i, 1000 + 1i, 1000 + 1i, 1000 + 1i, 1000 + 1i}
	s := []complex128{1000 - 1i, 1000 - 1i, 1000 - 1i}

	PrintMatC("a", ad, "(%4g", " +%4gi)  ", false)
	PrintVecC("u", u, "(%4g", " +%4gi)  ", false)
	PrintVecC("w", w, "(%4g", " +%4gi)  ", false)
	PrintVecC("r", r, "(%4g", " +%4gi)  ", false)
	PrintVecC("s", s, "(%4g", " +%4gi)  ", false)

	io.Pf("\nfunc SpMatVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	p := make([]complex128, 3)
	SpMatVecMulC(p, 1, a, u) // p := 1*a*u
	PrintVecC("p = a*u", p, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "p = a*u", 1e-17, p, []complex128{6.5 + 34i, 0.55 + 2.2i, 38 + 320i})

	io.Pf("\nfunc SpMatVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatVecMulAddC(s, 1, a, u) // s += dot(a, u)
	PrintVecC("s += a*u", s, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "s += a*u", 1e-15, s, []complex128{1006.5 + 33i, 1000.55 + 1.2i, 1038 + 319i})

	io.Pf("\nfunc SpMatTrVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	q := make([]complex128, 5)
	SpMatTrVecMulC(q, 1, a, w) // q = dot(transpose(a), w)
	PrintVecC("q = trans(a)*w", q, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "q = trans(a)*w", 1e-14, q, []complex128{312.5 + 20.95i, 625 + 51.9i, 935 - 37.15i, 1245 + 3.8i, 1557.5 + 4.75i})

	io.Pf("\nfunc SpMatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	SpMatTrVecMulAddC(r, 1, a, w) // r += dot(transpose(a), w)
	PrintVecC("r += trans(a)*w", r, "(%2g", " +%4gi)  ", false)
	chk.VectorC(tst, "r += trans(a)*w", 1e-14, r, []complex128{1312.5 + 21.95i, 1625 + 52.9i, 1935 - 36.15i, 2245 + 4.8i, 2557.5 + 5.75i})
}

func TestSparseLA04(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA04")

	var ta Triplet
	ta.Init(3, 3, 4)
	ta.Put(0, 0, 1.0)
	ta.Put(1, 0, 2.0)
	ta.Put(2, 1, 3.0)
	ta.Put(2, 2, 4.0)
	a := ta.ToMatrix(nil)
	io.Pf("a = %+v\n", a)
	PrintMat("a", a.ToDense(), "%2g ", false)
	chk.Vector(tst, "a.x", a.x, []float64{1, 2, 3, 4})
	chk.Ints(tst, "a.i", a.i, []int{0, 1, 2, 2})
	chk.Ints(tst, "a.p", a.p, []int{0, 2, 3, 4})

	var tb Triplet
	tb.Init(3, 3, 4)
	tb.Put(1, 0, 3.0)
	tb.Put(0, 1, 2.0)
	tb.Put(1, 2, 1.0)
	tb.Put(2, 2, 1.0)
	b := tb.ToMatrix(nil)
	io.Pf("b = %+v\n", b)
	PrintMat("b", b.ToDense(), "%2g ", false)
	chk.Vector(tst, "b.x", b.x, []float64{3, 2, 1, 1})
	chk.Ints(tst, "b.i", b.i, []int{1, 0, 1, 2})
	chk.Ints(tst, "b.p", b.p, []int{0, 1, 2, 4})

	c, a2c, b2c := SpAllocMatAddMat(a, b)
	SpMatAddMat(c, 1, a, 1, b, a2c, b2c)
	io.Pf("c = %+v\n", c)
	PrintMat("c", c.ToDense(), "%2g ", false)
	chk.Vector(tst, "c.x", c.x, []float64{1, 5, 2, 3, 1, 5})
	chk.Ints(tst, "c.i", c.i, []int{0, 1, 0, 2, 1, 2})
	chk.Ints(tst, "c.p", c.p, []int{0, 2, 4, 6})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 3, 5})
	chk.Ints(tst, "b2c", b2c, []int{1, 2, 4, 5})
	chk.Matrix(tst, "c", 1e-17, c.ToDense(), [][]float64{{1, 2, 0}, {5, 0, 1}, {0, 3, 5}})
}

func TestSparseLA05(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA05")

	var ta Triplet
	ta.Init(5, 6, 9)
	ta.Put(0, 0, 1)
	ta.Put(2, 0, 2)
	ta.Put(4, 0, 3)
	ta.Put(1, 2, 3)
	ta.Put(3, 2, 1)
	ta.Put(0, 3, 1)
	ta.Put(4, 3, 5)
	ta.Put(0, 5, 7)
	ta.Put(2, 5, 8)
	a := ta.ToMatrix(nil)
	io.Pf("a = %+v\n", a)
	PrintMat("a", a.ToDense(), "%2g ", false)
	chk.Vector(tst, "a.x", a.x, []float64{1, 2, 3, 3, 1, 1, 5, 7, 8})
	chk.Ints(tst, "a.i", a.i, []int{0, 2, 4, 1, 3, 0, 4, 0, 2})
	chk.Ints(tst, "a.p", a.p, []int{0, 3, 3, 5, 7, 7, 9})

	var tb Triplet
	tb.Init(5, 6, 12)
	tb.Put(1, 1, 1)
	tb.Put(3, 1, 8)
	tb.Put(0, 2, 1)
	tb.Put(1, 2, 2)
	tb.Put(2, 2, 3)
	tb.Put(3, 2, 4)
	tb.Put(4, 2, 5)
	tb.Put(2, 4, 5)
	tb.Put(0, 5, 1)
	tb.Put(1, 5, 4)
	tb.Put(2, 5, 2)
	tb.Put(4, 5, 1)
	b := tb.ToMatrix(nil)
	io.Pf("b = %+v\n", b)
	PrintMat("b", b.ToDense(), "%2g ", false)
	chk.Vector(tst, "b.x", b.x, []float64{1, 8, 1, 2, 3, 4, 5, 5, 1, 4, 2, 1})
	chk.Ints(tst, "b.i", b.i, []int{1, 3, 0, 1, 2, 3, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "b.p", b.p, []int{0, 0, 2, 7, 7, 8, 12})

	c, a2c, b2c := SpAllocMatAddMat(a, b)
	SpMatAddMat(c, 0.1, a, 2, b, a2c, b2c)
	io.Pf("c = %+v\n", c)
	PrintMat("c", c.ToDense(), "%10.6f ", false)
	chk.Vector(tst, "c.x", c.x, []float64{0.1, 0.2, 0.3, 2, 16, 2, 4.3, 6, 8.1, 10, 0.1, 0.5, 10, 2.7, 8, 4.8, 2})
	chk.Ints(tst, "c.i", c.i, []int{0, 2, 4, 1, 3, 0, 1, 2, 3, 4, 0, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "c.p", c.p, []int{0, 3, 5, 10, 12, 13, 17})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 2, 6, 8, 10, 11, 13, 15})
	chk.Ints(tst, "b2c", b2c, []int{3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16})
	chk.Matrix(tst, "c", 1e-16, c.ToDense(), [][]float64{
		{0.1, 0, 2.0, 0.1, 0, 2.7},
		{0, 2, 4.3, 0, 0, 8.0},
		{0.2, 0, 6.0, 0, 10, 4.8},
		{0, 16, 8.1, 0, 0, 0.0},
		{0.3, 0, 10.0, 0.5, 0, 2.0},
	})
}

func TestSparseLA06(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA06")

	var ta Triplet
	ta.Init(5, 6, 9)
	ta.Put(0, 0, 1)
	ta.Put(2, 0, 2)
	ta.Put(4, 0, 3)
	ta.Put(1, 2, 3)
	ta.Put(3, 2, 1)
	ta.Put(0, 3, 1)
	ta.Put(4, 3, 5)
	ta.Put(0, 5, 7)
	ta.Put(2, 5, 8)
	a := ta.ToMatrix(nil)
	io.Pf("a = %+v\n", a)
	PrintMat("a", a.ToDense(), "%2g ", false)
	chk.Vector(tst, "a.x", a.x, []float64{1, 2, 3, 3, 1, 1, 5, 7, 8})
	chk.Ints(tst, "a.i", a.i, []int{0, 2, 4, 1, 3, 0, 4, 0, 2})
	chk.Ints(tst, "a.p", a.p, []int{0, 3, 3, 5, 7, 7, 9})

	var tb Triplet
	tb.Init(5, 6, 12)
	tb.Put(1, 1, 1)
	tb.Put(3, 1, 8)
	tb.Put(0, 2, 1)
	tb.Put(1, 2, 2)
	tb.Put(2, 2, 3)
	tb.Put(3, 2, 4)
	tb.Put(4, 2, 5)
	tb.Put(2, 4, 5)
	tb.Put(0, 5, 1)
	tb.Put(1, 5, 4)
	tb.Put(2, 5, 2)
	tb.Put(4, 5, 1)
	b := tb.ToMatrix(nil)
	io.Pf("b = %+v\n", b)
	PrintMat("b", b.ToDense(), "%2g ", false)
	chk.Vector(tst, "b.x", b.x, []float64{1, 8, 1, 2, 3, 4, 5, 5, 1, 4, 2, 1})
	chk.Ints(tst, "b.i", b.i, []int{1, 3, 0, 1, 2, 3, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "b.p", b.p, []int{0, 0, 2, 7, 7, 8, 12})

	c, a2c, b2c := SpAllocMatAddMat(a, b)
	var d CCMatrixC
	SpInitSimilarR2C(&d, c)
	chk.Ints(tst, "c.i", c.i, []int{0, 2, 4, 1, 3, 0, 1, 2, 3, 4, 0, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "c.p", c.p, []int{0, 3, 5, 10, 12, 13, 17})
	chk.Ints(tst, "d.i", d.i, []int{0, 2, 4, 1, 3, 0, 1, 2, 3, 4, 0, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "d.p", d.p, []int{0, 3, 5, 10, 12, 13, 17})
	α, β, γ, μ := 0.1, 1.0, 0.1, 2.0
	//    c :=      γ*a + μ*b
	//    d := (α+βi)*a + μ*b
	SpMatAddMatC(&d, c, α, β, γ, a, μ, b, a2c, b2c)
	io.Pf("c = %+v\n", c)
	io.Pf("d = %+v\n", d)
	PrintMat("c", c.ToDense(), "%5.2f ", false)
	PrintMatC("d", d.ToDense(), "(%5.2f", "%5.2f) ", false)
	chk.Vector(tst, "c.x", c.x, []float64{0.1, 0.2, 0.3, 2, 16, 2, 4.3, 6, 8.1, 10, 0.1, 0.5, 10, 2.7, 8, 4.8, 2})
	chk.Vector(tst, "d.x", d.x, []float64{0.1, 0.2, 0.3, 2, 16, 2, 4.3, 6, 8.1, 10, 0.1, 0.5, 10, 2.7, 8, 4.8, 2})
	chk.Vector(tst, "d.z", d.z, []float64{1, 2, 3, 0, 0, 0, 3, 0, 1, 0, 1, 5, 0, 7, 0, 8, 0})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 2, 6, 8, 10, 11, 13, 15})
	chk.Ints(tst, "b2c", b2c, []int{3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16})
	chk.Matrix(tst, "c", 1e-16, c.ToDense(), [][]float64{
		{0.1, 0, 2.0, 0.1, 0, 2.7},
		{0, 2, 4.3, 0, 0, 8.0},
		{0.2, 0, 6.0, 0, 10, 4.8},
		{0, 16, 8.1, 0, 0, 0.0},
		{0.3, 0, 10.0, 0.5, 0, 2.0},
	})
	chk.MatrixC(tst, "d", 1e-16, d.ToDense(), [][]complex128{
		{0.1 + 1i, 0, 2.0, 0.1 + 1i, 0, 2.7 + 7i},
		{0, 2, 4.3 + 3i, 0, 0, 8.0},
		{0.2 + 2i, 0, 6.0, 0, 10, 4.8 + 8i},
		{0, 16, 8.1 + 1i, 0, 0, 0.0},
		{0.3 + 3i, 0, 10.0, 0.5 + 5i, 0, 2.0},
	})
}

func TestSparseLA07(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA07")

	var ta Triplet
	ta.Init(5, 6, 9)
	ta.Put(0, 0, 1)
	ta.Put(2, 0, 2)
	ta.Put(4, 0, 3)
	ta.Put(1, 2, 3)
	ta.Put(3, 2, 1)
	ta.Put(0, 3, 1)
	ta.Put(4, 3, 5)
	ta.Put(0, 5, 7)
	ta.Put(2, 5, 8)
	io.Pf("ta = %+v\n", ta)

	var tb Triplet
	tb.Init(5, 6, 12)
	tb.Put(1, 1, 1)
	tb.Put(3, 1, 8)
	tb.Put(0, 2, 1)
	tb.Put(1, 2, 2)
	tb.Put(2, 2, 3)
	tb.Put(3, 2, 4)
	tb.Put(4, 2, 5)
	tb.Put(2, 4, 5)
	tb.Put(0, 5, 1)
	tb.Put(1, 5, 4)
	tb.Put(2, 5, 2)
	tb.Put(4, 5, 1)
	io.Pf("tb = %+v\n", tb)

	var tc Triplet
	tc.Init(5, 6, ta.Len()+tb.Len())
	SpTriAdd(&tc, 0.1, &ta, 2, &tb)
	io.Pf("tc = %+v\n", tc)
	chk.Matrix(tst, "c", 1e-16, tc.ToMatrix(nil).ToDense(), [][]float64{
		{0.1, 0, 2.0, 0.1, 0, 2.7},
		{0, 2, 4.3, 0, 0, 8.0},
		{0.2, 0, 6.0, 0, 10, 4.8},
		{0, 16, 8.1, 0, 0, 0.0},
		{0.3, 0, 10.0, 0.5, 0, 2.0},
	})
}

func TestSparseLA08(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA08")

	var ta Triplet
	ta.Init(5, 6, 9)
	ta.Put(0, 0, 1)
	ta.Put(2, 0, 2)
	ta.Put(4, 0, 3)
	ta.Put(1, 2, 3)
	ta.Put(3, 2, 1)
	ta.Put(0, 3, 1)
	ta.Put(4, 3, 5)
	ta.Put(0, 5, 7)
	ta.Put(2, 5, 8)
	io.Pf("ta = %+v\n", ta)

	var tb Triplet
	tb.Init(5, 6, 12)
	tb.Put(1, 1, 1)
	tb.Put(3, 1, 8)
	tb.Put(0, 2, 1)
	tb.Put(1, 2, 2)
	tb.Put(2, 2, 3)
	tb.Put(3, 2, 4)
	tb.Put(4, 2, 5)
	tb.Put(2, 4, 5)
	tb.Put(0, 5, 1)
	tb.Put(1, 5, 4)
	tb.Put(2, 5, 2)
	tb.Put(4, 5, 1)
	io.Pf("tb = %+v\n", tb)

	var td TripletC
	td.Init(5, 6, ta.Len()+tb.Len(), false)
	α, β, μ := 0.1, 1.0, 2.0
	// d := (α+βi)*a + μ*b
	SpTriAddR2C(&td, α, β, &ta, μ, &tb)
	io.Pf("td = %+v\n", td)
	chk.MatrixC(tst, "d", 1e-16, td.ToMatrix(nil).ToDense(), [][]complex128{
		{0.1 + 1i, 0, 2.0, 0.1 + 1i, 0, 2.7 + 7i},
		{0, 2, 4.3 + 3i, 0, 0, 8.0},
		{0.2 + 2i, 0, 6.0, 0, 10, 4.8 + 8i},
		{0, 16, 8.1 + 1i, 0, 0, 0.0},
		{0.3 + 3i, 0, 10.0, 0.5 + 5i, 0, 2.0},
	})
}

func TestSparseLA09(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	chk.PrintTitle("TestSparse LA09")

	var a Triplet
	SpTriSetDiag(&a, 4, 666.0)
	A := a.ToMatrix(nil).ToDense()
	PrintMat("diag(4) with 666 =", A, "%8g", false)
	chk.Matrix(tst, "a", 1e-17, A, [][]float64{
		{666, 0, 0, 0},
		{0, 666, 0, 0},
		{0, 0, 666, 0},
		{0, 0, 0, 666},
	})
}

func TestSparseLA10(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("TestSparse LA10: SpTriMatTrVecMul")

	var a Triplet
	a.Init(3, 5, 15)
	a.Put(0, 0, 10.0)
	a.Put(0, 1, 20.0)
	a.Put(0, 2, 30.0)
	a.Put(0, 3, 40.0)
	a.Put(0, 4, 50.0)
	a.Put(1, 0, 1.0)
	a.Put(1, 1, 20.0)
	a.Put(1, 2, 3.0)
	a.Put(1, 3, 40.0)
	a.Put(1, 4, 5.0)
	a.Put(2, 0, 10.0)
	a.Put(2, 1, 2.0)
	a.Put(2, 2, 30.0)
	a.Put(2, 3, 4.0)
	a.Put(2, 4, 50.0)

	x := []float64{0.5, 0.4, 0.3}
	y := make([]float64, 5)
	SpTriMatTrVecMul(y, &a, x) // y := transpose(a) * x

	chk.Vector(tst, "y=Tr(a)*x", 1e-17, y, []float64{8.4, 18.6, 25.2, 37.2, 42})
}
