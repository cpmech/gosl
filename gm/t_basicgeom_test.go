// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_basicgeom01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("basicgeom01. Point")

	a := &Point{0, 0, 0}
	b := a.NewCopy()
	δ := DistPointPoint(a, b)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("δ(a,b) = %v\n", δ)
	chk.Scalar(tst, "dist(a,b)", 1e-17, δ, 0.0)

	SQ3 := math.Sqrt(3.0)
	c := a.NewDisp(-1.0/SQ3, -1.0/SQ3, -1.0/SQ3)
	δ = DistPointPoint(a, c)
	io.Pforan("c = %+v\n", c)
	io.Pforan("δ(a,c) = %v\n", δ)
	chk.Scalar(tst, "dist(a,c)", 1e-17, δ, 1.0)

	b.X, b.Y, b.Z = 1, 0, 0
	δ = DistPointPoint(a, b)
	io.Pforan("b = %+v\n", b)
	io.Pforan("δ(a,b) = %v\n", δ)
	chk.Scalar(tst, "dist(a,b)", 1e-17, δ, 1)

	b.X, b.Y, b.Z = 0, 2, 0
	δ = DistPointPoint(a, b)
	io.Pforan("b = %+v\n", b)
	io.Pforan("δ(a,b) = %v\n", δ)
	chk.Scalar(tst, "dist(a,b)", 1e-17, δ, 2)

	b.X, b.Y, b.Z = 0, 3, 0
	δ = DistPointPoint(a, b)
	io.Pforan("b = %+v\n", b)
	io.Pforan("δ(a,b) = %v\n", δ)
	chk.Scalar(tst, "dist(a,b)", 1e-17, δ, 3)

	b.X, b.Y, b.Z = 1, 1, 1
	δ = DistPointPoint(a, b)
	io.Pforan("b = %+v\n", b)
	io.Pforan("δ(a,b) = %v\n", δ)
	chk.Scalar(tst, "dist(a,b)", 1e-17, δ, math.Sqrt(3.0))
}

func Test_basicgeom02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("basicgeom02. Vector")

	u := []float64{1, 2, 3}
	v := []float64{4, 5, 6}
	s := VecDot(u, v)
	r := VecNew(2, u)
	w := VecNewAdd(2, u, -3, v)
	io.Pforan("u = %v  norm = %g\n", u, VecNorm(u))
	io.Pforan("v = %v  norm = %g\n", v, VecNorm(v))
	io.Pforan("w = %v  norm = %g\n", w, VecNorm(w))
	io.Pforan("u.v = %v\n", s)
	chk.Scalar(tst, "u.v", 1e-17, s, 32.0)
	chk.Scalar(tst, "norm(u)", 1e-17, VecNorm(u), math.Sqrt(14.0))
	chk.Scalar(tst, "norm(v)", 1e-17, VecNorm(v), math.Sqrt(77.0))
	chk.Vector(tst, "r", 1e-17, r, []float64{2, 4, 6})
	chk.Vector(tst, "w", 1e-17, w, []float64{-10, -11, -12})
}

func Test_basicgeom03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("basicgeom03. Segment")

	u := &Segment{&Point{1, 7, 0}, &Point{4, 7, 0}}
	v := &Segment{&Point{2, 5, 0}, &Point{5, 9, 0}}
	w := v.New(2.0)
	U := u.Vector(1)
	V := v.Vector(1)
	W := w.Vector(1)
	io.Pforan("u = %v\n", u)
	io.Pforan("v = %v\n", v)
	io.Pforan("w = %v\n", w)
	io.Pforan("U = %v\n", U)
	io.Pforan("V = %v\n", V)
	io.Pforan("W = %v\n", W)
	chk.Scalar(tst, "len(u)", 1e-17, u.Len(), 3.0)
	chk.Scalar(tst, "len(v)", 1e-17, v.Len(), 5.0)
	chk.Scalar(tst, "len(w)", 1e-17, w.Len(), 10.0)
	chk.Scalar(tst, "w.A.X", 1e-17, w.A.X, v.A.X)
	chk.Scalar(tst, "w.A.Y", 1e-17, w.A.Y, v.A.Y)
	chk.Scalar(tst, "w.A.Z", 1e-17, w.A.Z, v.A.Z)
	chk.Vector(tst, "U", 1e-17, U, []float64{3, 0, 0})
	chk.Vector(tst, "V", 1e-17, V, []float64{3, 4, 0})
	chk.Vector(tst, "W", 1e-17, W, []float64{6, 8, 0})
	udotv := VecDot(U, V)
	udotw := VecDot(U, W)
	vdotw := VecDot(V, W)
	io.Pforan("dot(u,v) = %v\n", udotv)
	io.Pforan("dot(u,w) = %v\n", udotw)
	io.Pforan("dot(v,w) = %v\n", vdotw)
	chk.Scalar(tst, "dot(u,v)", 1e-17, udotv, 9.0)
	chk.Scalar(tst, "dot(u,w)", 1e-17, udotw, 18.0)
	chk.Scalar(tst, "dot(v,w)", 1e-17, vdotw, 50.0)
}

func Test_basicgeom04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("basicgeom04")

	zero := 1e-8

	io.PfYel("\n(1)\n")
	a := &Point{0, 0, 0}
	b := &Point{0, 0, 0}
	p := &Point{1, 1, 0}
	δ := DistPointLine(p, a, b, zero, true)
	tol := 1e-15
	SQ2 := math.Sqrt(2.0)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-SQ2) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, SQ2)
	}

	io.PfYel("\n(2)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 0, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(3)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 3, 4, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(4)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 3, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-2.4) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 2.4)
	}

	io.PfYel("\n(4)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 1.08, 1.44, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(5)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 0, 0, 0
	p.X, p.Y, p.Z = 1, 1, 1
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	SQ3 := math.Sqrt(3.0)
	if δ < 0.0 || math.Abs(δ-SQ3) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(6)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 0, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(7)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 1, 1, 1
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	tol = 1e-15
	if δ < 0.0 || math.Abs(δ-SQ2) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	io.PfYel("\n(8)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 1, 1, 0
	δ = DistPointLine(p, a, b, zero, true)
	io.Pforan("a = %+v\n", a)
	io.Pforan("b = %+v\n", b)
	io.Pforan("p = %+v\n", p)
	io.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-1.0) > tol {
		chk.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}
}

func Test_basicgeom05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("basicgeom05")

	zero := 1e-1
	told := 1e-9
	tolin := 1e-3

	o := &Point{0, 0, 0}
	a := &Point{0.5, 0.5, 0.5}
	b := &Point{0.1, 0.5, 0.8}
	c := &Point{1, 1, 1}
	cmin, cmax := PointsLims([]*Point{o, a, b, c})
	chk.Vector(tst, "cmin", 1e-17, cmin, []float64{0, 0, 0})
	chk.Vector(tst, "cmax", 1e-17, cmax, []float64{1, 1, 1})
	if !IsPointIn(a, cmin, cmax, zero) {
		chk.Panic("a=%v must be in box")
	}
	if !IsPointIn(b, cmin, cmax, zero) {
		chk.Panic("b=%v must be in box")
	}
	if !IsPointIn(c, cmin, cmax, zero) {
		chk.Panic("c=%v must be in box")
	}

	p := &Point{0.5, 0.5, 0.5}
	q := &Point{1.5, 1.5, 1.5}
	if !IsPointInLine(p, o, c, zero, told, tolin) {
		chk.Panic("p=%v must be in line")
	}
	if IsPointInLine(q, o, c, zero, told, tolin) {
		chk.Panic("q=%v must not be in line")
	}
}
