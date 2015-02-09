// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/utl"
)

func TestMatVec01(tst *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()
	utl.TTitle("TestMatVec 01")

	a := [][]float64{
		{1, 2, 3, 4, 5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{10, 20, 30, 40, 50},
	}
	b := [][]float64{
		{1.1, 2.2, 3.3, 4.4, 5.5},
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{1, 2, 3, 4, 5},
	}
	c := [][]float64{
		{1, 2, 3},
		{0, 4, 5},
		{1, 0, 6},
	}
	d := [][]float64{
		{0.1, 10, 0.1},
		{0.2, 20, 0.2},
		{0.3, 30, 0.3},
		{0.4, 40, 0.4},
		{0.5, 50, 0.5},
	}
	e := [][]float64{
		{0.01, 0.02, 1},
		{-0.02, 0.03, 0.01},
		{0.5, -0.01, 0.02},
	}
	u := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	v := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	w := []float64{10.0, 20.0, 30.0}
	r := []float64{1000, 1000, 1000, 1000, 1000}
	s := []float64{1000, 1000, 1000}

	PrintMat("a", a, "%5g", false)
	PrintMat("b", b, "%5g", false)
	PrintMat("c", c, "%5g", false)
	PrintMat("d", d, "%5g", false)
	PrintMat("e", e, "%7g", false)
	PrintVec("u", u, "%5g", false)
	PrintVec("v", v, "%5g", false)
	PrintVec("w", w, "%5g", false)
	PrintVec("r", r, "%5g", false)
	PrintVec("s", s, "%5g", false)

	utl.Pf("\nfunc MatVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	p := make([]float64, len(a))
	MatVecMul(p, 1, a, u) // p := 1*a*u
	PrintVec("p = a*u", p, "%5g", false)
	utl.CheckVector(tst, "p = a*u", 1e-17, p, []float64{5.5, 0.55, 55})

	utl.Pf("\nfunc MatVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	MatVecMulAdd(s, 1, a, u) // s += dot(a, u)
	PrintVec("s += a*u", s, "%9g", false)
	utl.CheckVector(tst, "s += a*u", 1e-17, s, []float64{1005.5, 1000.55, 1055})

	utl.Pf("\nfunc MatTrVecMul(v []float64, α float64, a [][]float64, u []float64)\n")
	q := make([]float64, 5)
	MatTrVecMul(q, 1, a, w) // q = dot(transpose(a), w)
	PrintVec("q = trans(a)*w", q, "%5g", false)
	utl.CheckVector(tst, "q = trans(a)*w", 1e-17, q, []float64{312, 624, 936, 1248, 1560})

	utl.Pf("\nfunc MatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64)\n")
	MatTrVecMulAdd(r, 1, a, w) // r += dot(transpose(a), w)
	PrintVec("r += trans(a)*w", r, "%5g", false)
	utl.CheckVector(tst, "r += trans(a)*w", 1e-17, r, []float64{1312, 1624, 1936, 2248, 2560})

	utl.Pf("\nfunc MatVecMulCopyAdd(w []float64, α float64, v []float64, β float64, a [][]float64, u[]float64)\n")
	MatVecMulCopyAdd(p, 0.5, w, 2.0, a, u) // p := 0.5*w + 2*a*u
	PrintVec("p = 0.5*w + 2*a*u", p, "%5g", false)
	utl.CheckVector(tst, "p = 0.5*w + 2*a*u", 1e-17, p, []float64{16, 11.1, 125})

	utl.Pf("\nfunc MatMul(c [][]float64, α float64, a, b [][]float64)\n")
	f := MatAlloc(3, 3)
	MatMul(f, 1, a, d) // f = dot(a, d)
	PrintMat("f = a*d", f, "%5g", false)
	utl.CheckMatrix(tst, "f = a*d", 1e-17, f, [][]float64{
		{5.5, 550, 5.5},
		{0.55, 55, 0.55},
		{55, 5500, 55},
	})

	utl.Pf("\nfunc MatMul3(d [][]float64, α float64, a, b, c [][]float64)\n")
	g := MatAlloc(3, 3)
	MatMul3(g, 1, c, e, f) // g = dot(c, dot(e, f))
	PrintMat("g = c*e*f", g, "%23.15e", false)
	utl.CheckMatrix(tst, "g = c*e*f", 1e-12, g, [][]float64{
		{6.751250e+01, 6.751250e+03, 6.751250e+01},
		{2.104850e+01, 2.104850e+03, 2.104850e+01},
		{7.813300e+01, 7.813300e+03, 7.813300e+01},
	})

	utl.Pf("\nfunc MatTrMul3(d [][]float64, α float64, a, b, c [][]float64)\n")
	h := MatAlloc(5, 3)
	MatTrMul3(h, 1, a, e, f) // h = dot(transpose(a), dot(e, f))
	PrintMat("h = trans(a)*e*f", h, "%23.15e", false)
	utl.CheckMatrix(tst, "h = trans(a)*e*f", 1e-13, h, [][]float64{
		{9.35566500e+01, 9.35566500e+03, 9.35566500e+01},
		{1.87113300e+02, 1.87113300e+04, 1.87113300e+02},
		{2.80669950e+02, 2.80669950e+04, 2.80669950e+02},
		{3.74226600e+02, 3.74226600e+04, 3.74226600e+02},
		{4.67783250e+02, 4.67783250e+04, 4.67783250e+02},
	})

	utl.Pf("\nfunc MatTrMulAdd3(d [][]float64, α float64, a, b, c [][]float64)\n")
	m := MatAlloc(5, 3)
	MatFill(m, 10000)
	n := MatAlloc(5, 3)
	MatCopy(n, 1, m)            // n := m
	MatTrMulAdd3(n, 1, a, e, f) // n += dot(transpose(a), dot(e, f))
	PrintMat("m", m, "%9g", false)
	PrintMat("n = m + trans(a)*e*f", n, "%23.15e", false)
	utl.CheckMatrix(tst, "n = m + trans(a)*e*f", 1e-11, n, [][]float64{
		{1.0093556650e+04, 1.93556650e+04, 1.0093556650e+04},
		{1.0187113300e+04, 2.87113300e+04, 1.0187113300e+04},
		{1.0280669950e+04, 3.80669950e+04, 1.0280669950e+04},
		{1.0374226600e+04, 4.74226600e+04, 1.0374226600e+04},
		{1.0467783250e+04, 5.67783250e+04, 1.0467783250e+04},
	})

	utl.Pf("\nfunc VecOuterAdd(a [][]float64, α float64, u, v []float64)\n")
	udyv := MatAlloc(5, 5)
	MatFill(udyv, 1000)
	VecOuterAdd(udyv, 0.5, u, v)
	utl.CheckMatrix(tst, "udyv += 0.5 * u dyad v", 1e-17, udyv, [][]float64{
		{1000.05, 1000.1, 1000.15, 1000.2, 1000.25},
		{1000.1, 1000.2, 1000.3, 1000.4, 1000.5},
		{1000.15, 1000.3, 1000.45, 1000.6, 1000.75},
		{1000.2, 1000.4, 1000.6, 1000.8, 1001.0},
		{1000.25, 1000.5, 1000.75, 1001.0, 1001.25},
	})
}
