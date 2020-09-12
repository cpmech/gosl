// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func TestSpBlas01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas01. (real) matrix vector multiplication")

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

	io.Pf("a" + ad.Print("%5g") + "\n")
	io.Pf("u = %v\n", u)
	io.Pf("w = %v\n", w)
	io.Pf("r = %v\n", r)
	io.Pf("s = %v\n", s)
	io.Pf("x = %v\n", x)
	io.Pf("W = %v\n", W)

	p := make([]float64, 3)
	SpMatVecMul(p, 1, a, u) // p := 1*a*u
	chk.Array(tst, "p = a*u", 1e-17, p, []float64{5.5, 0.55, 55})

	SpMatVecMulAdd(s, 1, a, u) // s += dot(a, u)
	chk.Array(tst, "s += a*u", 1e-17, s, []float64{1005.5, 1000.55, 1055})

	SpMatVecMulAddX(x, a, 2, u, 3, W) // x += a * (2*u + 3*W)
	chk.Array(tst, "x += a * (2*u + 3*W)", 1e-17, x, []float64{1176, 2017.6, 4760})

	q := make([]float64, 5)
	SpMatTrVecMul(q, 1, a, w) // q = dot(transpose(a), w)
	chk.Array(tst, "q = trans(a)*w", 1e-17, q, []float64{312, 624, 936, 1248, 1560})

	SpMatTrVecMulAdd(r, 1, a, w) // r += dot(transpose(a), w)
	chk.Array(tst, "r += trans(a)*w", 1e-17, r, []float64{1312, 1624, 1936, 2248, 2560})
}

func TestSpBlas02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas02. (complex) matrix vector multiplication")

	var t TripletC
	t.Init(3, 5, 15)
	t.Put(0, 0, 1+0i)
	t.Put(0, 1, 2+0i)
	t.Put(0, 2, 3+0i)
	t.Put(0, 3, 4+0i)
	t.Put(0, 4, 5+0i)
	t.Put(1, 0, 0.1+0i)
	t.Put(1, 1, 0.2+0i)
	t.Put(1, 2, 0.3+0i)
	t.Put(1, 3, 0.4+0i)
	t.Put(1, 4, 0.5+0i)
	t.Put(2, 0, 10+0i)
	t.Put(2, 1, 20+0i)
	t.Put(2, 2, 30+0i)
	t.Put(2, 3, 40+0i)
	t.Put(2, 4, 50+0i)

	a := t.ToMatrix(nil)
	ad := a.ToDense()
	u := []complex128{0.1, 0.2, 0.3, 0.4, 0.5}
	w := []complex128{10.0, 20.0, 30.0}
	r := []complex128{1000, 1000, 1000, 1000, 1000}
	s := []complex128{1000, 1000, 1000}

	io.Pf("a" + ad.Print("%5g", "%+4gi") + "\n")
	io.Pf("u = %v\n", u)
	io.Pf("w = %v\n", w)
	io.Pf("r = %v\n", r)
	io.Pf("s = %v\n", s)

	p := make([]complex128, 3)
	SpMatVecMulC(p, 1, a, u) // p := 1*a*u
	chk.ArrayC(tst, "p = a*u", 1e-17, p, []complex128{5.5, 0.55, 55})

	SpMatVecMulAddC(s, 1, a, u) // s += dot(a, u)
	chk.ArrayC(tst, "s += a*u", 1e-17, s, []complex128{1005.5, 1000.55, 1055})

	q := make([]complex128, 5)
	SpMatTrVecMulC(q, 1, a, w) // q = dot(transpose(a), w)
	chk.ArrayC(tst, "q = trans(a)*w", 1e-17, q, []complex128{312, 624, 936, 1248, 1560})

	SpMatTrVecMulAddC(r, 1, a, w) // r += dot(transpose(a), w)
	chk.ArrayC(tst, "r += trans(a)*w", 1e-17, r, []complex128{1312, 1624, 1936, 2248, 2560})
}

func TestSpBlas03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas03 (complex) matrix vector multiplication")

	var t TripletC
	t.Init(3, 5, 15)
	t.Put(0, 0, 1+0i)
	t.Put(0, 1, 2-1i)
	t.Put(0, 2, 3+0i)
	t.Put(0, 3, 4+3i)
	t.Put(0, 4, 5+2i)
	t.Put(1, 0, 0.1+1i)
	t.Put(1, 1, 0.2+0i)
	t.Put(1, 2, 0.3-2i)
	t.Put(1, 3, 0.4+0i)
	t.Put(1, 4, 0.5-1i)
	t.Put(2, 0, 10+0i)
	t.Put(2, 1, 20+2i)
	t.Put(2, 2, 30+0i)
	t.Put(2, 3, 40-1i)
	t.Put(2, 4, 50+0i)

	a := t.ToMatrix(nil)
	ad := a.ToDense()
	u := []complex128{0.1, 0.2 + 10i, 0.3, 0.4 + 3i, 0.5}
	w := []complex128{10.0 + 1i, 20.0 - 0.5i, 30.0}
	r := []complex128{1000 + 1i, 1000 + 1i, 1000 + 1i, 1000 + 1i, 1000 + 1i}
	s := []complex128{1000 - 1i, 1000 - 1i, 1000 - 1i}

	io.Pf("a" + ad.Print("%5g", "%+4gi") + "\n")
	io.Pf("u = %v\n", u)
	io.Pf("w = %v\n", w)
	io.Pf("r = %v\n", r)
	io.Pf("s = %v\n", s)

	p := make([]complex128, 3)
	SpMatVecMulC(p, 1, a, u) // p := 1*a*u
	chk.ArrayC(tst, "p = a*u", 1e-17, p, []complex128{6.5 + 34i, 0.55 + 2.2i, 38 + 320i})

	SpMatVecMulAddC(s, 1, a, u) // s += dot(a, u)
	chk.ArrayC(tst, "s += a*u", 1e-15, s, []complex128{1006.5 + 33i, 1000.55 + 1.2i, 1038 + 319i})

	q := make([]complex128, 5)
	SpMatTrVecMulC(q, 1, a, w) // q = dot(transpose(a), w)
	chk.ArrayC(tst, "q = trans(a)*w", 1e-14, q, []complex128{312.5 + 20.95i, 625 + 51.9i, 935 - 37.15i, 1245 + 3.8i, 1557.5 + 4.75i})

	SpMatTrVecMulAddC(r, 1, a, w) // r += dot(transpose(a), w)
	chk.ArrayC(tst, "r += trans(a)*w", 1e-14, r, []complex128{1312.5 + 21.95i, 1625 + 52.9i, 1935 - 36.15i, 2245 + 4.8i, 2557.5 + 5.75i})
}

func TestSpBlas04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas04. Matrix addition")

	var ta Triplet
	ta.Init(3, 3, 4)
	ta.Put(0, 0, 1.0)
	ta.Put(1, 0, 2.0)
	ta.Put(2, 1, 3.0)
	ta.Put(2, 2, 4.0)
	a := ta.ToMatrix(nil)
	io.Pf("a = %+v\n", a)
	chk.Array(tst, "a.x", 1e-17, a.x, []float64{1, 2, 3, 4})
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
	chk.Array(tst, "b.x", 1e-17, b.x, []float64{3, 2, 1, 1})
	chk.Ints(tst, "b.i", b.i, []int{1, 0, 1, 2})
	chk.Ints(tst, "b.p", b.p, []int{0, 1, 2, 4})

	c, a2c, b2c := SpAllocMatAddMat(a, b)
	SpMatAddMat(c, 1, a, 1, b, a2c, b2c)
	io.Pf("c = %+v\n", c)
	chk.Array(tst, "c.x", 1e-17, c.x, []float64{1, 5, 2, 3, 1, 5})
	chk.Ints(tst, "c.i", c.i, []int{0, 1, 0, 2, 1, 2})
	chk.Ints(tst, "c.p", c.p, []int{0, 2, 4, 6})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 3, 5})
	chk.Ints(tst, "b2c", b2c, []int{1, 2, 4, 5})
	chk.Deep2(tst, "c", 1e-17, c.ToDense().GetDeep2(), [][]float64{{1, 2, 0}, {5, 0, 1}, {0, 3, 5}})
}

func TestSpBlas05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas05 Matrix addition")

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
	chk.Array(tst, "a.x", 1e-17, a.x, []float64{1, 2, 3, 3, 1, 1, 5, 7, 8})
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
	chk.Array(tst, "b.x", 1e-16, b.x, []float64{1, 8, 1, 2, 3, 4, 5, 5, 1, 4, 2, 1})
	chk.Ints(tst, "b.i", b.i, []int{1, 3, 0, 1, 2, 3, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "b.p", b.p, []int{0, 0, 2, 7, 7, 8, 12})

	c, a2c, b2c := SpAllocMatAddMat(a, b)
	SpMatAddMat(c, 0.1, a, 2, b, a2c, b2c)
	io.Pf("c = %+v\n", c)
	chk.Array(tst, "c.x", 1e-16, c.x, []float64{0.1, 0.2, 0.3, 2, 16, 2, 4.3, 6, 8.1, 10, 0.1, 0.5, 10, 2.7, 8, 4.8, 2})
	chk.Ints(tst, "c.i", c.i, []int{0, 2, 4, 1, 3, 0, 1, 2, 3, 4, 0, 4, 2, 0, 1, 2, 4})
	chk.Ints(tst, "c.p", c.p, []int{0, 3, 5, 10, 12, 13, 17})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 2, 6, 8, 10, 11, 13, 15})
	chk.Ints(tst, "b2c", b2c, []int{3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16})
	chk.Deep2(tst, "c", 1e-16, c.ToDense().GetDeep2(), [][]float64{
		{0.1, +0, +2.0, 0.1, +0, 2.7},
		{0.0, +2, +4.3, 0.0, +0, 8.0},
		{0.2, +0, +6.0, 0.0, 10, 4.8},
		{0.0, 16, +8.1, 0.0, +0, 0.0},
		{0.3, +0, 10.0, 0.5, +0, 2.0},
	})
}

func TestSpBlas06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas06. (real and complex) Matrix addition")

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
	chk.Array(tst, "a.x", 1e-17, a.x, []float64{1, 2, 3, 3, 1, 1, 5, 7, 8})
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
	chk.Array(tst, "b.x", 1e-17, b.x, []float64{1, 8, 1, 2, 3, 4, 5, 5, 1, 4, 2, 1})
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
	chk.Array(tst, "c.x", 1e-16, c.x, []float64{0.1, 0.2, 0.3, 2, 16, 2, 4.3, 6, 8.1, 10, 0.1, 0.5, 10, 2.7, 8, 4.8, 2})
	chk.ArrayC(tst, "d.x", 1e-16, d.x, []complex128{0.1 + 1i, 0.2 + 2i, 0.3 + 3i, 2, 16, 2, 4.3 + 3i, 6, 8.1 + 1i, 10, 0.1 + 1i, 0.5 + 5i, 10, 2.7 + 7i, 8, 4.8 + 8i, 2})
	chk.Ints(tst, "a2c", a2c, []int{0, 1, 2, 6, 8, 10, 11, 13, 15})
	chk.Ints(tst, "b2c", b2c, []int{3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16})
	chk.Deep2(tst, "c", 1e-16, c.ToDense().GetDeep2(), [][]float64{
		{0.1, +0, +2.0, 0.1, +0, 2.7},
		{0.0, +2, +4.3, 0.0, +0, 8.0},
		{0.2, +0, +6.0, 0.0, 10, 4.8},
		{0.0, 16, +8.1, 0.0, +0, 0.0},
		{0.3, +0, 10.0, 0.5, +0, 2.0},
	})
	chk.Deep2c(tst, "d", 1e-16, d.ToDense().GetDeep2(), [][]complex128{
		{0.1 + 1i, +0, +2.0 + 0i, 0.1 + 1i, +0, 2.7 + 7i},
		{0.0 + 0i, +2, +4.3 + 3i, 0.0 + 0i, +0, 8.0 + 0i},
		{0.2 + 2i, +0, +6.0 + 0i, 0.0 + 0i, 10, 4.8 + 8i},
		{0.0 + 0i, 16, +8.1 + 1i, 0.0 + 0i, +0, 0.0 + 0i},
		{0.3 + 3i, +0, 10.0 + 0i, 0.5 + 5i, +0, 2.0 + 0i},
	})
}

func TestSpBlas07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas07. SpTriAdd")

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
	chk.Deep2(tst, "c", 1e-16, tc.ToDense().GetDeep2(), [][]float64{
		{0.1, +0, +2.0, 0.1, +0, 2.7},
		{0.0, +2, +4.3, 0.0, +0, 8.0},
		{0.2, +0, +6.0, 0.0, 10, 4.8},
		{0.0, 16, +8.1, 0.0, +0, 0.0},
		{0.3, +0, 10.0, 0.5, +0, 2.0},
	})
}

func TestSpBlas08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas08. SpTriAddR2C")

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
	td.Init(5, 6, ta.Len()+tb.Len())
	α, β, μ := 0.1, 1.0, 2.0
	// d := (α+βi)*a + μ*b
	SpTriAddR2C(&td, α, β, &ta, μ, &tb)
	io.Pf("td = %+v\n", td)
	chk.Deep2c(tst, "d", 1e-16, td.ToDense().GetDeep2(), [][]complex128{
		{0.1 + 1i, +0, +2.0 + 0i, 0.1 + 1i, +0, 2.7 + 7i},
		{0.0 + 0i, +2, +4.3 + 3i, 0.0 + 0i, +0, 8.0 + 0i},
		{0.2 + 2i, +0, +6.0 + 0i, 0.0 + 0i, 10, 4.8 + 8i},
		{0.0 + 0i, 16, +8.1 + 1i, 0.0 + 0i, +0, 0.0 + 0i},
		{0.3 + 3i, +0, 10.0 + 0i, 0.5 + 5i, +0, 2.0 + 0i},
	})
}

func TestSpBlas09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas09. SpTriSetDiag")

	var a Triplet
	SpTriSetDiag(&a, 4, 666.0)
	A := a.ToMatrix(nil).ToDense()
	chk.Deep2(tst, "a", 1e-17, A.GetDeep2(), [][]float64{
		{666, 0, 0, 0},
		{0, 666, 0, 0},
		{0, 0, 666, 0},
		{0, 0, 0, 666},
	})
}

func TestSpBlas10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas10: SpTriMatTrVecMul")

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
	io.Pforan("y = %v\n", y)
	chk.Array(tst, "y=Tr(a)*x", 1e-17, y, []float64{8.4, 18.6, 25.2, 37.2, 42})

	u := []float64{8.4, 18.6, 25.2, 37.2, 42}
	z := make([]float64, 3)
	SpTriMatVecMul(z, &a, u)
	io.Pfcyan("z = %v\n", z)
	chk.Array(tst, "z=a*u", 1e-17, z, []float64{4800, 2154, 3126})
}

func TestSpBlas11(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpBlas11: SpMatMatTrMul")

	var T1 Triplet
	T1.Init(2, 3, 3)
	T1.Put(0, 0, 1)
	T1.Put(1, 1, 1)
	T1.Put(1, 2, 1)

	a1 := T1.ToMatrix(nil)

	b1 := NewMatrix(a1.m, a1.m)
	SpMatMatTrMul(b1, 1, a1)
	chk.Deep2(tst, "b1", 1e-17, b1.GetDeep2(), [][]float64{{1, 0}, {0, 2}})

	io.Pf("\n----------------------------------\n")
	var T2 Triplet
	T2.Init(3, 2, 6)
	T2.Put(0, 0, 1)
	T2.Put(0, 1, -1)
	T2.Put(1, 0, 2)
	T2.Put(1, 1, 3)
	T2.Put(2, 0, -2)
	T2.Put(2, 1, 4)

	a2 := T2.ToMatrix(nil)

	b2 := NewMatrix(a2.m, a2.m)
	SpMatMatTrMul(b2, 1, a2)
	chk.Deep2(tst, "b2", 1e-17, b2.GetDeep2(), [][]float64{{2, -1, -6}, {-1, 13, 8}, {-6, 8, 20}})

	io.Pf("\n----------------------------------\n")
	var T3 Triplet
	T3.Init(4, 1, 3)
	T3.Put(0, 0, 1)
	T3.Put(1, 0, 2)
	T3.Put(2, 0, -3)

	a3 := T3.ToMatrix(nil)

	b3 := NewMatrix(a3.m, a3.m)
	SpMatMatTrMul(b3, 1, a3)
	chk.Deep2(tst, "b3", 1e-17, b3.GetDeep2(), [][]float64{{1, 2, -3, 0}, {2, 4, -6, 0}, {-3, -6, 9, 0}, {0, 0, 0, 0}})

	io.Pf("\n----------------------------------\n")
	var T4 Triplet
	T4.Init(1, 3, 2)
	T4.Put(0, 0, 1)
	T4.Put(0, 2, 2)

	a4 := T4.ToMatrix(nil)

	b4 := NewMatrix(a4.m, a4.m)
	SpMatMatTrMul(b4, 1, a4)
	chk.Deep2(tst, "b4", 1e-17, b4.GetDeep2(), [][]float64{{5}})
}
