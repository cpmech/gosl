// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/utl"
)

func Test_basicgeom01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("basicgeom01. Point")

	a := &Point{0, 0, 0}
	b := a.NewCopy()
	δ := DistPointPoint(a, b)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("δ(a,b) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,b)", 1e-17, δ, 0.0)

	SQ3 := math.Sqrt(3.0)
	c := a.NewDisp(-1.0/SQ3, -1.0/SQ3, -1.0/SQ3)
	δ = DistPointPoint(a, c)
	utl.Pforan("c = %+v\n", c)
	utl.Pforan("δ(a,c) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,c)", 1e-17, δ, 1.0)

	b.X, b.Y, b.Z = 1, 0, 0
	δ = DistPointPoint(a, b)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("δ(a,b) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,b)", 1e-17, δ, 1)

	b.X, b.Y, b.Z = 0, 2, 0
	δ = DistPointPoint(a, b)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("δ(a,b) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,b)", 1e-17, δ, 2)

	b.X, b.Y, b.Z = 0, 3, 0
	δ = DistPointPoint(a, b)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("δ(a,b) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,b)", 1e-17, δ, 3)

	b.X, b.Y, b.Z = 1, 1, 1
	δ = DistPointPoint(a, b)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("δ(a,b) = %v\n", δ)
	utl.CheckScalar(tst, "dist(a,b)", 1e-17, δ, math.Sqrt(3.0))
}

func Test_basicgeom02(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("basicgeom02. Vector")

	u := []float64{1, 2, 3}
	v := []float64{4, 5, 6}
	s := VecDot(u, v)
	r := VecNew(2, u)
	w := VecNewAdd(2, u, -3, v)
	utl.Pforan("u = %v  norm = %g\n", u, VecNorm(u))
	utl.Pforan("v = %v  norm = %g\n", v, VecNorm(v))
	utl.Pforan("w = %v  norm = %g\n", w, VecNorm(w))
	utl.Pforan("u.v = %v\n", s)
	utl.CheckScalar(tst, "u.v", 1e-17, s, 32.0)
	utl.CheckScalar(tst, "norm(u)", 1e-17, VecNorm(u), math.Sqrt(14.0))
	utl.CheckScalar(tst, "norm(v)", 1e-17, VecNorm(v), math.Sqrt(77.0))
	utl.CheckVector(tst, "r", 1e-17, r, []float64{2, 4, 6})
	utl.CheckVector(tst, "w", 1e-17, w, []float64{-10, -11, -12})
}

func Test_basicgeom03(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("basicgeom03. Segment")

	u := &Segment{&Point{1, 7, 0}, &Point{4, 7, 0}}
	v := &Segment{&Point{2, 5, 0}, &Point{5, 9, 0}}
	w := v.New(2.0)
	U := u.Vector(1)
	V := v.Vector(1)
	W := w.Vector(1)
	utl.Pforan("u = %v\n", u)
	utl.Pforan("v = %v\n", v)
	utl.Pforan("w = %v\n", w)
	utl.Pforan("U = %v\n", U)
	utl.Pforan("V = %v\n", V)
	utl.Pforan("W = %v\n", W)
	utl.CheckScalar(tst, "len(u)", 1e-17, u.Len(), 3.0)
	utl.CheckScalar(tst, "len(v)", 1e-17, v.Len(), 5.0)
	utl.CheckScalar(tst, "len(w)", 1e-17, w.Len(), 10.0)
	utl.CheckScalar(tst, "w.A.X", 1e-17, w.A.X, v.A.X)
	utl.CheckScalar(tst, "w.A.Y", 1e-17, w.A.Y, v.A.Y)
	utl.CheckScalar(tst, "w.A.Z", 1e-17, w.A.Z, v.A.Z)
	utl.CheckVector(tst, "U", 1e-17, U, []float64{3, 0, 0})
	utl.CheckVector(tst, "V", 1e-17, V, []float64{3, 4, 0})
	utl.CheckVector(tst, "W", 1e-17, W, []float64{6, 8, 0})
	udotv := VecDot(U, V)
	udotw := VecDot(U, W)
	vdotw := VecDot(V, W)
	utl.Pforan("dot(u,v) = %v\n", udotv)
	utl.Pforan("dot(u,w) = %v\n", udotw)
	utl.Pforan("dot(v,w) = %v\n", vdotw)
	utl.CheckScalar(tst, "dot(u,v)", 1e-17, udotv, 9.0)
	utl.CheckScalar(tst, "dot(u,w)", 1e-17, udotw, 18.0)
	utl.CheckScalar(tst, "dot(v,w)", 1e-17, vdotw, 50.0)
}

func Test_basicgeom04(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("basicgeom04")

	zero := 1e-8

	utl.PfYel("\n(1)\n")
	a := &Point{0, 0, 0}
	b := &Point{0, 0, 0}
	p := &Point{1, 1, 0}
	δ := DistPointLine(p, a, b, zero, true)
	tol := 1e-15
	SQ2 := math.Sqrt(2.0)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-SQ2) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, SQ2)
	}

	utl.PfYel("\n(2)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 0, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(3)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 3, 4, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(4)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 3, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-2.4) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 2.4)
	}

	utl.PfYel("\n(4)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 3, 4, 0
	p.X, p.Y, p.Z = 1.08, 1.44, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(5)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 0, 0, 0
	p.X, p.Y, p.Z = 1, 1, 1
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	SQ3 := math.Sqrt(3.0)
	if δ < 0.0 || math.Abs(δ-SQ3) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(6)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 0, 0, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(7)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 1, 1, 1
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	tol = 1e-15
	if δ < 0.0 || math.Abs(δ-SQ2) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}

	utl.PfYel("\n(8)\n")
	a.X, a.Y, a.Z = 0, 0, 0
	b.X, b.Y, b.Z = 1, 0, 0
	p.X, p.Y, p.Z = 1, 1, 0
	δ = DistPointLine(p, a, b, zero, true)
	utl.Pforan("a = %+v\n", a)
	utl.Pforan("b = %+v\n", b)
	utl.Pforan("p = %+v\n", p)
	utl.Pforan("δ(p,a,b) = %v\n", δ)
	if δ < 0.0 || math.Abs(δ-1.0) > tol {
		utl.Panic("dist for p=%v failed. dist: %g != %g", p, δ, 0.0)
	}
}

func Test_basicgeom05(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("basicgeom05")

	zero := 1e-1
	told := 1e-9
	tolin := 1e-3

	o := &Point{0, 0, 0}
	a := &Point{0.5, 0.5, 0.5}
	b := &Point{0.1, 0.5, 0.8}
	c := &Point{1, 1, 1}
	cmin, cmax := PointsLims([]*Point{o, a, b, c})
	utl.CheckVector(tst, "cmin", 1e-17, cmin, []float64{0, 0, 0})
	utl.CheckVector(tst, "cmax", 1e-17, cmax, []float64{1, 1, 1})
	if !IsPointIn(a, cmin, cmax, zero) {
		utl.Panic("a=%v must be in box")
	}
	if !IsPointIn(b, cmin, cmax, zero) {
		utl.Panic("b=%v must be in box")
	}
	if !IsPointIn(c, cmin, cmax, zero) {
		utl.Panic("c=%v must be in box")
	}

	p := &Point{0.5, 0.5, 0.5}
	q := &Point{1.5, 1.5, 1.5}
	if !IsPointInLine(p, o, c, zero, told, tolin) {
		utl.Panic("p=%v must be in line")
	}
	if IsPointInLine(q, o, c, zero, told, tolin) {
		utl.Panic("q=%v must not be in line")
	}
}
